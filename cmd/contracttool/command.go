package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"utopia/internal/chain"
	"utopia/internal/config"
	"utopia/internal/contract"
	"utopia/internal/helper"
	utopia_network "utopia/internal/network"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/urfave/cli.v1"
)

var (
	AccountFlag = cli.StringFlag{
		Name:  "account",
		Usage: "Account address in hex mode",
		Value: "",
	}
	ToFlag = cli.StringFlag{
		Name:  "to",
		Usage: "The dest address list in hex mode and seperate by ,",
		Value: "",
	}
	ValueFlag = cli.StringFlag{
		Name:  "value",
		Usage: "The value list in ether unit",
		Value: "",
	}
	ContractFlag = cli.StringFlag{
		Name:  "contract",
		Usage: "Contract address in hex mode",
		Value: "",
	}
	CodeFlag = cli.StringFlag{
		Name:  "code",
		Usage: "The contract code file path",
		Value: "",
	}
	ABIFlag = cli.StringFlag{
		Name:  "abi",
		Usage: "The contract abi file path",
		Value: "",
	}
	ParamFlag = cli.StringFlag{
		Name:  "param",
		Usage: "The parameters for call contract",
		Value: "",
	}
	DataFlag = cli.StringFlag{
		Name:  "data",
		Usage: "The abi argumetns",
		Value: "",
	}
	FuncFlag = cli.StringFlag{
		Name:  "func",
		Usage: "The function protype",
		Value: "",
	}
	SignFlag = cli.BoolFlag{
		Name:  "sign",
		Usage: "True or False to indicate function sign",
	}
	FileFlag = cli.StringFlag{
		Name:  "file",
		Usage: "The excel file path with address list",
		Value: "",
	}
	EnableFlag = cli.BoolFlag{
		Name:  "enable",
		Usage: "Enable approve or Disable appreove",
	}

	cmdDeploy = cli.Command{
		Name:   "deploy",
		Usage:  "Deploy smart contract on chain",
		Action: DeployContract,
		Flags: []cli.Flag{
			CodeFlag,
			ABIFlag,
			ParamFlag,
			ValueFlag,
		},
	}
	cmdCall = cli.Command{
		Name:   "call",
		Usage:  "Call smart contract on chain",
		Action: CallContract,
		Flags: []cli.Flag{
			ContractFlag,
			ABIFlag,
			ParamFlag,
			ValueFlag,
		},
	}
	cmdList = cli.Command{
		Name:   "list",
		Usage:  "List all contract with account trsacted",
		Action: ListContract,
		Flags: []cli.Flag{
			AccountFlag,
		},
	}
	cmdERC20 = cli.Command{
		Name:  "erc20",
		Usage: "ERC20 operations on chain",
		Subcommands: []cli.Command{
			{
				Name:   "balance",
				Usage:  "Query erc20 contract balance of account",
				Action: QueryERC20,
				Flags: []cli.Flag{
					ContractFlag,
					AccountFlag,
				},
			},
			{
				Name:   "transfer",
				Usage:  "Transfer erc20 balance of account",
				Action: TransferERC20,
				Flags: []cli.Flag{
					ContractFlag,
					ToFlag,
					ValueFlag,
					FileFlag,
				},
			},
			{
				Name:   "approve",
				Usage:  "Approve erc20 balance of account",
				Action: ApproveERC20,
				Flags: []cli.Flag{
					ContractFlag,
					ToFlag,
					ValueFlag,
				},
			},
		},
	}
	cmdERC721 = cli.Command{
		Name:  "erc721",
		Usage: "ERC721 operations on chain",
		Subcommands: []cli.Command{
			{
				Name:   "balance",
				Usage:  "Query erc721 contract balance of account",
				Action: QueryERC721,
				Flags: []cli.Flag{
					ContractFlag,
					AccountFlag,
				},
			},
			{
				Name:   "transfer",
				Usage:  "Transfer erc721 balance of account",
				Action: TransferERC721,
				Flags: []cli.Flag{
					ContractFlag,
					ToFlag,
					ValueFlag,
					FileFlag,
				},
			},
			{
				Name:   "approve",
				Usage:  "Approve erc721 balance of account",
				Action: ApproveERC721,
				Flags: []cli.Flag{
					ContractFlag,
					ToFlag,
					ValueFlag,
					EnableFlag,
				},
			},
			{
				Name:   "property",
				Usage:  "Query properties of erc721 nft",
				Action: PropertyQuery,
				Flags: []cli.Flag{
					ContractFlag,
					ValueFlag,
				},
			},
		},
	}
	cmdAbi = cli.Command{
		Name:  "abi",
		Usage: "ABI encode and decode",
		Subcommands: []cli.Command{
			{
				Name:   "encode",
				Usage:  "Encode abi with arguments",
				Action: EncodeABI,
				Flags: []cli.Flag{
					FuncFlag,
					DataFlag,
					SignFlag,
				},
			},
			{
				Name:   "decode",
				Usage:  "Decode abi with arguments",
				Action: DecodeABI,
				Flags: []cli.Flag{
					FuncFlag,
					DataFlag,
					SignFlag,
				},
			},
		},
	}
)

func DeployContract(ctx *cli.Context) error {
	code := ctx.String(CodeFlag.Name)
	abi := ctx.String(ABIFlag.Name)
	params := ctx.String(ParamFlag.Name)
	ivalue := ctx.String(ValueFlag.Name)

	value := common.Big0
	fv, err := strconv.ParseFloat(ivalue, 64)
	if err == nil {
		value = helper.EthToWei(float32(fv))
	}

	// get wallet for sign transaction
	wallet, err := wallet.GetWallet(config.Config.Chain.From)
	if err != nil {
		return err
	}

	// get and connect chain
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if chain == nil {
		return errors.New("Connect chain failed")
	}
	defer chain.DisConnect()

	// create contract
	contract := contract.NewContract(chain, "", contract.COMMON_CRONTACT)
	err = contract.SetABI(abi)
	if err != nil {
		return err
	}

	bin, err := ioutil.ReadFile(code)
	if err != nil {
		return err
	}

	// deploy contract
	result, err := contract.Deploy(string(bin), params, wallet, value)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Deploy contract address %s\n", result)
	return nil
}

func CallContract(ctx *cli.Context) error {
	address := ctx.String(ContractFlag.Name)
	abi := ctx.String(ABIFlag.Name)
	params := ctx.String(ParamFlag.Name)
	ivalue := ctx.String(ValueFlag.Name)

	value := common.Big0
	fv, err := strconv.ParseFloat(ivalue, 64)
	if err == nil {
		value = helper.EthToWei(float32(fv))
	}

	// get wallet for sign transaction
	wallet, err := wallet.GetWallet(config.Config.Chain.From)
	if err != nil {
		return err
	}

	// get chain meta and connect
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if chain == nil {
		return errors.New("Connect chain failed")
	}
	defer chain.DisConnect()

	// create contract
	contract := contract.NewContract(chain, address, contract.COMMON_CRONTACT)
	err = contract.SetABI(abi)
	if err != nil {
		return err
	}

	// call contract
	result, err := contract.Call(params, wallet, value)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Call contract result: %v\n", result)
	return nil
}

func ListContract(ctx *cli.Context) error {
	return nil
}

func QueryERC20(ctx *cli.Context) error {
	address := ctx.String(ContractFlag.Name)
	account := ctx.String(AccountFlag.Name)

	// get chain meta and connect
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	// query erc20 balance
	erc20 := contract.NewContract(c, address, contract.ERC20_CONTRACT)
	if erc20 == nil {
		return errors.New("Create erc20 contract failed")
	}

	balance, err := erc20.(*contract.ERC20Contract).Balance(account)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Balance: %f\n", helper.WeiToEth(balance))
	return nil
}

func TransferERC20(ctx *cli.Context) error {
	address := ctx.String(ContractFlag.Name)
	to := ctx.String(ToFlag.Name)
	value := ctx.String(ValueFlag.Name)
	file := ctx.String(FileFlag.Name)

	translist := make([]helper.TransferInfo, 0)
	if to != "" && value != "" {
		translist = append(translist, helper.TransferInfo{
			From:  config.Config.Chain.From,
			To:    to,
			Value: value,
		})
	}

	// read file transfer list
	if file != "" {
		list, err := helper.ReadTransferFile(file)
		if err != nil {
			return err
		}

		translist = append(translist, list...)
	}

	// get chain meta and connect it
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	// transfer erc20 balance
	erc20 := contract.NewContract(c, address, contract.ERC20_CONTRACT)
	if erc20 == nil {
		return errors.New("Create erc20 contract failed")
	}

	for _, info := range translist {
		// get wallet for sign transaction
		wallet, err := wallet.GetWallet(info.From)
		if err != nil {
			return err
		}

		fv, _ := strconv.ParseFloat(info.Value, 64)
		tx, err := erc20.(*contract.ERC20Contract).Transfer(info.To, helper.EthToWei(float32(fv)), wallet)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Tranfer %s to %s failed with err: %v\n", info.Value, info.To, err)
		} else {
			fmt.Fprintf(os.Stderr, "Tranfer %s to %s with transaction %s\n", info.Value, info.To, tx)
		}
	}

	return nil
}

func ApproveERC20(ctx *cli.Context) error {
	address := ctx.String(ContractFlag.Name)
	to := ctx.String(ToFlag.Name)
	value := ctx.String(ValueFlag.Name)
	fv, _ := strconv.ParseFloat(value, 64)

	// get wallet for sign transaction
	wallet, err := wallet.GetWallet(config.Config.Chain.From)
	if err != nil {
		return err
	}

	// get chain meta and connect it
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	// transfer erc20 balance
	erc20 := contract.NewContract(c, address, contract.ERC20_CONTRACT)
	if erc20 == nil {
		return errors.New("Create erc20 contract failed")
	}

	tx, err := erc20.(*contract.ERC20Contract).Approve(to, helper.EthToWei(float32(fv)), wallet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Approve %s to %s failed with err: %v\n", value, to, err)
	} else {
		fmt.Fprintf(os.Stderr, "Approve %s to %s with transaction %s\n", value, to, tx)
	}

	return nil
}

func QueryERC721(ctx *cli.Context) error {
	address := ctx.String(ContractFlag.Name)
	account := ctx.String(AccountFlag.Name)

	// get chain meta and connect
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	// query erc721 balance
	erc721 := contract.NewContract(c, address, contract.ERC721_CONTRACT)
	if erc721 == nil {
		return errors.New("Create erc20 contract failed")
	}

	balance, err := erc721.(*contract.ERC721Contract).Balance(account)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Balance: %d\n", balance)
	return nil
}

func TransferERC721(ctx *cli.Context) error {
	address := ctx.String(ContractFlag.Name)
	to := ctx.String(ToFlag.Name)
	value := ctx.String(ValueFlag.Name)
	file := ctx.String(FileFlag.Name)

	translist := make([]helper.TransferInfo, 0)
	if to != "" && value != "" {
		translist = append(translist, helper.TransferInfo{
			From:  config.Config.Chain.From,
			To:    to,
			Value: value,
		})
	}

	// read file transfer list
	if file != "" {
		list, err := helper.ReadTransferFile(file)
		if err != nil {
			return err
		}

		translist = append(translist, list...)
	}

	// get chain meta and connect it
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	// transfer erc20 balance
	erc721 := contract.NewContract(c, address, contract.ERC721_CONTRACT)
	if erc721 == nil {
		return errors.New("Create erc20 contract failed")
	}

	for _, info := range translist {
		// get wallet for sign transaction
		wallet, err := wallet.GetWallet(info.From)
		if err != nil {
			return err
		}

		tv, _ := strconv.ParseUint(info.Value, 10, 64)
		tx, err := erc721.(*contract.ERC721Contract).Transfer(info.To, tv, wallet)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Tranfer %s to %s failed with err: %v\n", info.Value, info.To, err)
		} else {
			fmt.Fprintf(os.Stderr, "Tranfer %s to %s with transaction %s\n", info.Value, info.To, tx)
		}
	}

	return nil
}

func ApproveERC721(ctx *cli.Context) error {
	address := ctx.String(ContractFlag.Name)
	to := ctx.String(ToFlag.Name)
	value := ctx.String(ValueFlag.Name)
	enable := ctx.Bool(EnableFlag.Name)
	tokenId, _ := strconv.ParseUint(value, 10, 64)

	// get wallet for sign transaction
	wallet, err := wallet.GetWallet(config.Config.Chain.From)
	if err != nil {
		return err
	}

	// get chain meta and connect it
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	// transfer erc20 balance
	erc721 := contract.NewContract(c, address, contract.ERC721_CONTRACT)
	if erc721 == nil {
		return errors.New("Create erc20 contract failed")
	}

	tx, err := erc721.(*contract.ERC721Contract).Approve(to, tokenId, enable, wallet)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Approve %s to %s failed with err: %v\n", value, to, err)
	} else {
		fmt.Fprintf(os.Stderr, "Approve %s to %s with transaction %s\n", value, to, tx)
	}

	return nil
}

func PropertyQuery(ctx *cli.Context) error {
	address := ctx.String(ContractFlag.Name)
	value := ctx.String(ValueFlag.Name)
	tokenId, _ := strconv.ParseUint(value, 10, 64)

	// get chain meta and connect it
	meta, err := chain.ChainMetaByName(config.Config.Chain.Network)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	// transfer erc20 balance
	erc721 := contract.NewContract(c, address, contract.ERC721_CONTRACT)
	if erc721 == nil {
		return errors.New("Create erc20 contract failed")
	}

	url, err := erc721.(*contract.ERC721Contract).TokenUrl(tokenId)
	if err != nil {
		return err
	}

	url = strings.Replace(url, "ipfs://", "https://ipfs.io/ipfs/", 1)
	result, err := utopia_network.HttpGet(url, nil, nil)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s\n", string(result))
	return nil
}

func EncodeABI(ctx *cli.Context) error {
	method := ctx.String(FuncFlag.Name)
	data := ctx.String(DataFlag.Name)
	sign := ctx.Bool(SignFlag.Name)

	contract := contract.NewContract(nil, "", contract.COMMON_CRONTACT)
	result, err := contract.EncodeABI(method, data, sign)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Result: %s\n", result)
	return nil
}

func DecodeABI(ctx *cli.Context) error {
	method := ctx.String(FuncFlag.Name)
	data := ctx.String(DataFlag.Name)
	sign := ctx.Bool(SignFlag.Name)

	contract := contract.NewContract(nil, "", contract.COMMON_CRONTACT)
	result, err := contract.DecodeABI(method, data, sign)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Result: %s\n", result)
	return nil
}
