package contract

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"utopia/internal/chain"
	"utopia/internal/helper"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type EthContract struct {
	chain   chain.Chain         // Chain id which contract deployed
	address common.Address      // Contract address
	abi     string              // Contract ABI
	client  *bind.BoundContract // Client of contract
}

func NewEthContract(chain chain.Chain, address string) Contract {
	return &EthContract{
		chain:   chain,
		address: common.HexToAddress(address),
		abi:     "",
		client:  nil,
	}
}

func (c *EthContract) Address() string {
	return c.address.Hex()
}

func (c *EthContract) Code() (string, error) {
	// read code from chain
	return c.chain.Code(c.address.Hex())
}

func (c *EthContract) ABI() string {
	return c.abi
}

func (c *EthContract) SetABI(path string) error {
	// read abi file content and parse to abi
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	parsed, err := abi.JSON(strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	c.abi = string(data)
	c.client = bind.NewBoundContract(c.address, parsed, c.chain.(*chain.EthChain).Client, c.chain.(*chain.EthChain).Client, c.chain.(*chain.EthChain).Client)

	return nil
}

func (c *EthContract) Deploy(code string, params string, wallet wallet.Wallet, value *big.Int) (string, error) {
	// parse the constructor params
	method, args, err := helper.ParseParams(params)
	if err != nil {
		return "", err
	}

	if method != "" {
		return "", errors.New("method must be empty for constructor")
	}

	// parse abi for get constructor method
	parsed, err := abi.JSON(strings.NewReader(string(c.abi)))
	if err != nil {
		return "", err
	}

	data := make([]interface{}, 0)
	if len(parsed.Constructor.Inputs) > 0 {
		index := 0
		for _, p := range parsed.Constructor.Inputs {
			if len(args) <= index {
				return "", errors.New("Not enough parameters")
			}

			if p.Type.Elem != nil {
				var subdata interface{}
				subdata, index, err = helper.Str2Array(args, index, p.Type)
				if err != nil {
					return "", err
				}

				data = append(data, subdata)
			} else {
				v, err := helper.Str2Type(args[index], p.Type.GetType())
				if err != nil {
					return "", err
				}

				data = append(data, v)
				index++
			}
		}
	}

	// get transaciton options for sign tx and set value
	opts, err := c.chain.(*chain.EthChain).GenTransOpts(wallet, value)
	if err != nil {
		return "", err
	}

	// send deploy transaction
	address, _, _, err := bind.DeployContract(opts, parsed, common.Hex2Bytes(code), c.chain.(*chain.EthChain).Client, data...)
	if err != nil {
		return "", err
	}

	// get contract address
	c.address = address
	return address.Hex(), nil
}

func (c *EthContract) Call(params string, wallet wallet.Wallet, value *big.Int) ([]interface{}, error) {
	// parse the constructor params
	method, args, err := helper.ParseParams(params)
	if err != nil {
		return nil, err
	}

	// parse abi for get call method
	parsed, err := abi.JSON(strings.NewReader(string(c.abi)))
	if err != nil {
		return nil, err
	}

	m, ok := parsed.Methods[method]
	if !ok {
		return nil, errors.New("Can not found methon in abi")
	}

	index := 0
	data := make([]interface{}, 0)
	for _, p := range m.Inputs {
		if len(args) <= index {
			return nil, errors.New("Not enough parameters")
		}

		if p.Type.Elem != nil {
			var subdata interface{}
			subdata, index, err = helper.Str2Array(args, index, p.Type)
			if err != nil {
				return nil, err
			}

			data = append(data, subdata)
		} else {
			v, err := helper.Str2Type(args[index], p.Type.GetType())
			if err != nil {
				return nil, err
			}

			data = append(data, v)
			index++
		}
	}

	// call contract if method is read-only, otherwise send transaction
	var result []interface{}
	if m.IsConstant() {
		err = c.client.Call(nil, &result, method, data...)
		if err != nil {
			return nil, err
		}
	} else {
		// get transaciton options for sign tx and set value
		opts, err := c.chain.(*chain.EthChain).GenTransOpts(wallet, value)
		if err != nil {
			return nil, err
		}

		tx, err := c.client.Transact(opts, method, data...)
		if err != nil {
			return nil, err
		}

		// return transaction hash
		result = make([]interface{}, 0)
		result = append(result, tx.Hash().Hex())
	}

	// todo: parse the call result or wait transction receipt
	return result, nil
}

func (c *EthContract) EncodeABI(method string, data string, withfunc bool) (string, error) {
	funcname, argtypes, err := helper.ParseParams(method)
	if err != nil {
		return "", err
	}

	callMethod, args, err := helper.ParseParams(data)
	if err != nil {
		return "", err
	}

	if strings.ToLower(funcname) != strings.ToLower(callMethod) {
		return "", errors.New("Not match call funcion name")
	}

	// trim space and change to lower case
	arguments := make(abi.Arguments, 0)
	for i, arg := range argtypes {
		argtypes[i] = strings.ToLower(strings.Trim(arg, " "))
		if argtypes[i] == "uint" || argtypes[i] == "int" {
			argtypes[i] = argtypes[i] + "256"
		} else if argtypes[i] == "uint[]" {
			argtypes[i] = "uint256[]"
		} else if argtypes[i] == "int[]" {
			argtypes[i] = "int256[]"
		}

		argtype, _ := abi.NewType(argtypes[i], "", nil)
		arguments = append(arguments, abi.Argument{Type: argtype})
	}

	// get the function signature
	result := make([]byte, 0)
	if withfunc {
		sig := fmt.Sprintf("%v(%v)", funcname, strings.Join(argtypes, ","))
		result = append(result, crypto.Keccak256([]byte(sig))[:4]...)
	}

	// change string data to dst type
	index := 0
	argData := make([]interface{}, 0)
	for _, p := range arguments {
		if len(args) <= index {
			return "", errors.New("Not enough parameters")
		}

		if p.Type.Elem != nil {
			var subdata interface{}
			subdata, index, err = helper.Str2Array(args, index, p.Type)
			if err != nil {
				return "", err
			}

			argData = append(argData, subdata)
		} else {
			v, err := helper.Str2Type(args[index], p.Type.GetType())
			if err != nil {
				return "", err
			}

			argData = append(argData, v)
			index++
		}
	}

	// pack all parameters
	pack, err := arguments.Pack(argData...)
	if err != nil {
		return "", err
	}

	result = append(result, pack...)
	return "0x" + common.Bytes2Hex(result), nil
}

func (c *EthContract) DecodeABI(method string, data string, withfunc bool) (string, error) {
	rdata := common.Hex2Bytes(data)

	// parse the function sig
	funcname, argtypes, err := helper.ParseParams(method)
	if err != nil {
		return "", err
	}

	// trim space and change to lower case
	arguments := make(abi.Arguments, 0)
	for i, arg := range argtypes {
		argtypes[i] = strings.ToLower(strings.Trim(arg, " "))
		if argtypes[i] == "uint" || argtypes[i] == "int" {
			argtypes[i] = argtypes[i] + "256"
		} else if argtypes[i] == "uint[]" {
			argtypes[i] = "uint256[]"
		} else if argtypes[i] == "int[]" {
			argtypes[i] = "int256[]"
		}

		argtype, _ := abi.NewType(argtypes[i], "", nil)
		arguments = append(arguments, abi.Argument{Type: argtype})
	}

	// check the method signature
	if withfunc {
		sig := fmt.Sprintf("%v(%v)", funcname, strings.Join(argtypes, ","))
		funcsig := crypto.Keccak256([]byte(sig))[:4]
		if hex.EncodeToString(funcsig) != hex.EncodeToString(rdata[:4]) {
			return "", errors.New("Not match function type")
		}

		rdata = rdata[4:]
	}

	// unpack all parameters
	parsed, err := arguments.Unpack(rdata)
	if err != nil {
		return "", err
	}

	// change result to string
	var builder strings.Builder
	builder.WriteString(funcname)
	builder.WriteString("(")

	for index, p := range arguments {
		if p.Type.Elem != nil {
			subdata, err := helper.Array2Str(parsed[index], p.Type.Elem.GetType())
			if err != nil {
				return "", err
			}

			builder.WriteString(subdata)
		} else {
			v, err := helper.Type2Str(parsed[index], p.Type.GetType())
			if err != nil {
				return "", err
			}

			builder.WriteString(v)
		}

		if index != len(arguments)-1 {
			builder.WriteString(",")
		}
	}

	builder.WriteString(")")
	return builder.String(), nil
}
