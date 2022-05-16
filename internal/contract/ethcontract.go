package contract

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"utopia/internal/chain"
	"utopia/internal/helper"

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

func (c *EthContract) Code() ([]byte, error) {
	return c.chain.Code(c.address.Hex())
}

func (c *EthContract) ABI() string {
	return c.abi
}

func (c *EthContract) SetABI(path string) error {
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

func (c *EthContract) Call(params string) ([]interface{}, error) {
	method, args, err := helper.ParseParams(params)
	if err != nil {
		return nil, err
	}

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
			subdata, err := helper.Str2Array(args, index, p.Type.Elem.GetType())
			if err != nil {
				return nil, err
			}

			data = append(data, subdata)
			index = index + len(subdata)
		} else {
			v, err := helper.Str2Type(args[index], p.Type.GetType())
			if err != nil {
				return nil, err
			}

			data = append(data, v)
			index++
		}
	}

	var result []interface{}
	err = c.client.Call(nil, &result, method, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *EthContract) EncodeABI(method string, data string, withfunc bool) ([]byte, error) {
	funcname, argtypes, err := helper.ParseParams(method)
	if err != nil {
		return nil, err
	}

	callMethod, args, err := helper.ParseParams(data)
	if err != nil {
		return nil, err
	}

	if strings.ToLower(funcname) != strings.ToLower(callMethod) {
		return nil, errors.New("Not match call funcion name")
	}

	// trim space and change to lower case
	arguments := make(abi.Arguments, 0)
	for i, arg := range argtypes {
		argtypes[i] = strings.ToLower(strings.Trim(arg, " "))
		if argtypes[i] == "uint" || argtypes[i] == "int" {
			argtypes[i] = argtypes[i] + "256"
		}

		argtype, _ := abi.NewType(argtypes[i], "", nil)
		arguments = append(arguments, abi.Argument{Type: argtype})
	}

	result := make([]byte, 0)
	if withfunc {
		sig := fmt.Sprintf("%v(%v)", funcname, strings.Join(argtypes, ","))
		result = append(result, crypto.Keccak256([]byte(sig))[:4]...)
	}

	index := 0
	argData := make([]interface{}, 0)
	for _, p := range arguments {
		if len(args) <= index {
			return nil, errors.New("Not enough parameters")
		}

		if p.Type.Elem != nil {
			subdata, err := helper.Str2Array(args, index, p.Type.Elem.GetType())
			if err != nil {
				return nil, err
			}

			argData = append(argData, subdata)
			index = index + len(subdata)
		} else {
			v, err := helper.Str2Type(args[index], p.Type.GetType())
			if err != nil {
				return nil, err
			}

			argData = append(argData, v)
			index++
		}
	}

	pack, _ := arguments.Pack(argData...)
	result = append(result, pack...)

	return result, nil
}

func (c *EthContract) DecodeABI(method string, data []byte, withfunc bool) (string, error) {
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
		}

		argtype, _ := abi.NewType(argtypes[i], "", nil)
		arguments = append(arguments, abi.Argument{Type: argtype})
	}

	if withfunc {
		sig := fmt.Sprintf("%v(%v)", funcname, strings.Join(argtypes, ","))
		funcsig := crypto.Keccak256([]byte(sig))[:4]
		if hex.EncodeToString(funcsig) != hex.EncodeToString(data[:4]) {
			return "", errors.New("Not match function type")
		}

		data = data[4:]
	}

	parsed, err := arguments.Unpack(data)
	if err != nil {
		return "", err
	}

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
