// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SimpleMetaData contains all meta data concerning the Simple contract.
var SimpleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"msg\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"GetMessage\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"msg\",\"type\":\"string\"}],\"name\":\"SetMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"data\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"03e33b53": "GetMessage()",
		"88fabb3a": "SetMessage(string)",
		"73d4a13a": "data()",
	},
	Bin: "0x608060405234801561001057600080fd5b506040516105ed3803806105ed83398101604081905261002f916100f8565b8051610042906000906020840190610049565b5050610201565b828054610055906101c7565b90600052602060002090601f01602090048101928261007757600085556100bd565b82601f1061009057805160ff19168380011785556100bd565b828001600101855582156100bd579182015b828111156100bd5782518255916020019190600101906100a2565b506100c99291506100cd565b5090565b5b808211156100c957600081556001016100ce565b634e487b7160e01b600052604160045260246000fd5b6000602080838503121561010b57600080fd5b82516001600160401b038082111561012257600080fd5b818501915085601f83011261013657600080fd5b815181811115610148576101486100e2565b604051601f8201601f19908116603f01168101908382118183101715610170576101706100e2565b81604052828152888684870101111561018857600080fd5b600093505b828410156101aa578484018601518185018701529285019261018d565b828411156101bb5760008684830101525b98975050505050505050565b600181811c908216806101db57607f821691505b6020821081036101fb57634e487b7160e01b600052602260045260246000fd5b50919050565b6103dd806102106000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806303e33b531461004657806373d4a13a1461006457806388fabb3a1461006c575b600080fd5b61004e610081565b60405161005b9190610251565b60405180910390f35b61004e610113565b61007f61007a3660046102bc565b6101a1565b005b6060600080546100909061036d565b80601f01602080910402602001604051908101604052809291908181526020018280546100bc9061036d565b80156101095780601f106100de57610100808354040283529160200191610109565b820191906000526020600020905b8154815290600101906020018083116100ec57829003601f168201915b5050505050905090565b600080546101209061036d565b80601f016020809104026020016040519081016040528092919081815260200182805461014c9061036d565b80156101995780601f1061016e57610100808354040283529160200191610199565b820191906000526020600020905b81548152906001019060200180831161017c57829003601f168201915b505050505081565b80516101b49060009060208401906101b8565b5050565b8280546101c49061036d565b90600052602060002090601f0160209004810192826101e6576000855561022c565b82601f106101ff57805160ff191683800117855561022c565b8280016001018555821561022c579182015b8281111561022c578251825591602001919060010190610211565b5061023892915061023c565b5090565b5b80821115610238576000815560010161023d565b600060208083528351808285015260005b8181101561027e57858101830151858201604001528201610262565b81811115610290576000604083870101525b50601f01601f1916929092016040019392505050565b634e487b7160e01b600052604160045260246000fd5b6000602082840312156102ce57600080fd5b813567ffffffffffffffff808211156102e657600080fd5b818401915084601f8301126102fa57600080fd5b81358181111561030c5761030c6102a6565b604051601f8201601f19908116603f01168101908382118183101715610334576103346102a6565b8160405282815287602084870101111561034d57600080fd5b826020860160208301376000928101602001929092525095945050505050565b600181811c9082168061038157607f821691505b6020821081036103a157634e487b7160e01b600052602260045260246000fd5b5091905056fea2646970667358221220d50e24aa6785ef6b515861544b1bbe550d2c762753f139480aa0a0c474353e3564736f6c634300080d0033",
}

// SimpleABI is the input ABI used to generate the binding from.
// Deprecated: Use SimpleMetaData.ABI instead.
var SimpleABI = SimpleMetaData.ABI

// Deprecated: Use SimpleMetaData.Sigs instead.
// SimpleFuncSigs maps the 4-byte function signature to its string representation.
var SimpleFuncSigs = SimpleMetaData.Sigs

// SimpleBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SimpleMetaData.Bin instead.
var SimpleBin = SimpleMetaData.Bin

// DeploySimple deploys a new Ethereum contract, binding an instance of Simple to it.
func DeploySimple(auth *bind.TransactOpts, backend bind.ContractBackend, msg string) (common.Address, *types.Transaction, *Simple, error) {
	parsed, err := SimpleMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimpleBin), backend, msg)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Simple{SimpleCaller: SimpleCaller{contract: contract}, SimpleTransactor: SimpleTransactor{contract: contract}, SimpleFilterer: SimpleFilterer{contract: contract}}, nil
}

// Simple is an auto generated Go binding around an Ethereum contract.
type Simple struct {
	SimpleCaller     // Read-only binding to the contract
	SimpleTransactor // Write-only binding to the contract
	SimpleFilterer   // Log filterer for contract events
}

// SimpleCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleSession struct {
	Contract     *Simple           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SimpleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleCallerSession struct {
	Contract *SimpleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SimpleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleTransactorSession struct {
	Contract     *SimpleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SimpleRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleRaw struct {
	Contract *Simple // Generic contract binding to access the raw methods on
}

// SimpleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleCallerRaw struct {
	Contract *SimpleCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleTransactorRaw struct {
	Contract *SimpleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimple creates a new instance of Simple, bound to a specific deployed contract.
func NewSimple(address common.Address, backend bind.ContractBackend) (*Simple, error) {
	contract, err := bindSimple(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Simple{SimpleCaller: SimpleCaller{contract: contract}, SimpleTransactor: SimpleTransactor{contract: contract}, SimpleFilterer: SimpleFilterer{contract: contract}}, nil
}

// NewSimpleCaller creates a new read-only instance of Simple, bound to a specific deployed contract.
func NewSimpleCaller(address common.Address, caller bind.ContractCaller) (*SimpleCaller, error) {
	contract, err := bindSimple(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleCaller{contract: contract}, nil
}

// NewSimpleTransactor creates a new write-only instance of Simple, bound to a specific deployed contract.
func NewSimpleTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleTransactor, error) {
	contract, err := bindSimple(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleTransactor{contract: contract}, nil
}

// NewSimpleFilterer creates a new log filterer instance of Simple, bound to a specific deployed contract.
func NewSimpleFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleFilterer, error) {
	contract, err := bindSimple(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleFilterer{contract: contract}, nil
}

// bindSimple binds a generic wrapper to an already deployed contract.
func bindSimple(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Simple *SimpleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Simple.Contract.SimpleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Simple *SimpleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Simple.Contract.SimpleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Simple *SimpleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Simple.Contract.SimpleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Simple *SimpleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Simple.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Simple *SimpleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Simple.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Simple *SimpleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Simple.Contract.contract.Transact(opts, method, params...)
}

// GetMessage is a free data retrieval call binding the contract method 0x03e33b53.
//
// Solidity: function GetMessage() view returns(string)
func (_Simple *SimpleCaller) GetMessage(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Simple.contract.Call(opts, &out, "GetMessage")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetMessage is a free data retrieval call binding the contract method 0x03e33b53.
//
// Solidity: function GetMessage() view returns(string)
func (_Simple *SimpleSession) GetMessage() (string, error) {
	return _Simple.Contract.GetMessage(&_Simple.CallOpts)
}

// GetMessage is a free data retrieval call binding the contract method 0x03e33b53.
//
// Solidity: function GetMessage() view returns(string)
func (_Simple *SimpleCallerSession) GetMessage() (string, error) {
	return _Simple.Contract.GetMessage(&_Simple.CallOpts)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(string)
func (_Simple *SimpleCaller) Data(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Simple.contract.Call(opts, &out, "data")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(string)
func (_Simple *SimpleSession) Data() (string, error) {
	return _Simple.Contract.Data(&_Simple.CallOpts)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(string)
func (_Simple *SimpleCallerSession) Data() (string, error) {
	return _Simple.Contract.Data(&_Simple.CallOpts)
}

// SetMessage is a paid mutator transaction binding the contract method 0x88fabb3a.
//
// Solidity: function SetMessage(string msg) returns()
func (_Simple *SimpleTransactor) SetMessage(opts *bind.TransactOpts, msg string) (*types.Transaction, error) {
	return _Simple.contract.Transact(opts, "SetMessage", msg)
}

// SetMessage is a paid mutator transaction binding the contract method 0x88fabb3a.
//
// Solidity: function SetMessage(string msg) returns()
func (_Simple *SimpleSession) SetMessage(msg string) (*types.Transaction, error) {
	return _Simple.Contract.SetMessage(&_Simple.TransactOpts, msg)
}

// SetMessage is a paid mutator transaction binding the contract method 0x88fabb3a.
//
// Solidity: function SetMessage(string msg) returns()
func (_Simple *SimpleTransactorSession) SetMessage(msg string) (*types.Transaction, error) {
	return _Simple.Contract.SetMessage(&_Simple.TransactOpts, msg)
}
