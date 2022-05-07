package contract

import (
	"errors"
	"io/ioutil"
	"reflect"
	"strings"
	"utopia/internal/chain"
	"utopia/internal/helper"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type EthContract struct {
	chain   chain.Chain         // Chain id which contract deployed
	address common.Address      // Contract address
	abi     string              // Contract ABI
	client  *bind.BoundContract // Client of contract
}

func NewEthContract(chain chain.Chain, address string) *EthContract {
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
			subdata, err := c.parseArray(args, index, p.Type.Elem.GetType())
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

func (c *EthContract) parseArray(args []string, index int, totype reflect.Type) ([]interface{}, error) {
	if !strings.HasPrefix(args[index], "[") {
		return nil, errors.New("Need array paramter but not found")
	}

	result := make([]interface{}, 0)
	inarray := true

	args[index] = args[index][1:]
	if strings.HasSuffix(args[index], "]") {
		args[index] = args[index][:len(args[index])-1]
		inarray = false
	}

	v, err := helper.Str2Type(args[index], totype)
	if err != nil {
		return nil, err
	}

	result = append(result, v)
	for {
		if !inarray {
			break
		}

		index++
		if len(args) <= index {
			break
		}

		if strings.HasSuffix(args[index], "]") {
			args[index] = args[index][:len(args[index])-1]
			inarray = false
		}

		v, err = helper.Str2Type(args[index], totype)
		if err != nil {
			return nil, err
		}

		result = append(result, v)
	}

	if inarray {
		return nil, errors.New("Invalid array values")
	}

	return result, nil
}
