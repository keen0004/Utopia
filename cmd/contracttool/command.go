package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"utopia/contracts/token"
	"utopia/internal/chain"
	"utopia/internal/contract"
	"utopia/internal/helper"
	"utopia/internal/logger"
	utopia_network "utopia/internal/network"
	"utopia/internal/wallet"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"
)

var (
	ChainFlag = cli.StringFlag{
		Name:  "chain",
		Usage: "Chain name (ie: btc, eth, bsc)",
		Value: "",
	}
	KeyFlag = cli.StringFlag{
		Name:  "key",
		Usage: "The file path of key store ",
		Value: "",
	}
	KeyDirFlag = cli.StringFlag{
		Name:  "keydir",
		Usage: "The directory of key store",
		Value: "",
	}
	PasswordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "The password of key store",
		Value: "",
	}
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

	cmdDeploy = cli.Command{
		Name:   "deploy",
		Usage:  "Deploy smart contract on chain",
		Action: DeployContract,
		Flags: []cli.Flag{
			ChainFlag,
			KeyFlag,
			PasswordFlag,
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
			ChainFlag,
			KeyFlag,
			PasswordFlag,
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
			ChainFlag,
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
					ChainFlag,
					ContractFlag,
					AccountFlag,
				},
			},
			{
				Name:   "transfer",
				Usage:  "Transfer erc20 balance of account",
				Action: TransferERC20,
				Flags: []cli.Flag{
					ChainFlag,
					KeyFlag,
					PasswordFlag,
					ContractFlag,
					ToFlag,
					ValueFlag,
					FileFlag,
				},
			},
			{
				Name:   "merge",
				Usage:  "Merge erc20 balance of account",
				Action: MergeERC20,
				Flags: []cli.Flag{
					ChainFlag,
					KeyDirFlag,
					PasswordFlag,
					ContractFlag,
					ToFlag,
				},
			},
			{
				Name:   "approve",
				Usage:  "Approve erc20 balance of account",
				Action: ApproveERC20,
				Flags: []cli.Flag{
					ChainFlag,
					KeyFlag,
					PasswordFlag,
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
					ChainFlag,
					ContractFlag,
					AccountFlag,
				},
			},
			{
				Name:   "transfer",
				Usage:  "Transfer erc721 balance of account",
				Action: TransferERC721,
				Flags: []cli.Flag{
					ChainFlag,
					KeyFlag,
					PasswordFlag,
					ContractFlag,
					ToFlag,
					ValueFlag,
				},
			},
			{
				Name:   "merge",
				Usage:  "Merge erc721 balance of account",
				Action: MergeERC721,
				Flags: []cli.Flag{
					ChainFlag,
					KeyDirFlag,
					PasswordFlag,
					ContractFlag,
					ToFlag,
				},
			},
			{
				Name:   "approve",
				Usage:  "Approve erc721 balance of account",
				Action: ApproveERC721,
				Flags: []cli.Flag{
					ChainFlag,
					KeyFlag,
					PasswordFlag,
					ContractFlag,
					ToFlag,
					ValueFlag,
				},
			},
			{
				Name:   "property",
				Usage:  "Query properties of erc721 nft",
				Action: PropertyQuery,
				Flags: []cli.Flag{
					ChainFlag,
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
	chainName := ctx.String(ChainFlag.Name)
	key := ctx.String(KeyFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	code := ctx.String(CodeFlag.Name)
	abi := ctx.String(ABIFlag.Name)
	params := ctx.String(ParamFlag.Name)
	ivalue := ctx.String(ValueFlag.Name)

	value := common.Big0
	fv, err := strconv.ParseFloat(ivalue, 64)
	if err == nil {
		value = helper.EthToWei(float32(fv))
	}

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	wallet := wallet.NewWallet(wallet.WALLET_ETH, key, password)
	err = wallet.LoadKey()
	if err != nil {
		return err
	}

	chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if chain == nil {
		return errors.New("Connect chain failed")
	}
	defer chain.DisConnect()

	contract := contract.NewContract(chain, "", contract.COMMON_CRONTACT)
	err = contract.SetABI(abi)
	if err != nil {
		return err
	}

	bin, err := ioutil.ReadFile(code)
	if err != nil {
		return err
	}

	result, err := contract.Deploy(common.FromHex(string(bin)), params, wallet, value)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Deploy contract address %s\n", result)
	return nil
}

func CallContract(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	key := ctx.String(KeyFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	address := ctx.String(ContractFlag.Name)
	abi := ctx.String(ABIFlag.Name)
	params := ctx.String(ParamFlag.Name)
	ivalue := ctx.String(ValueFlag.Name)

	value := common.Big0
	fv, err := strconv.ParseFloat(ivalue, 64)
	if err == nil {
		value = helper.EthToWei(float32(fv))
	}

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	wallet := wallet.NewWallet(wallet.WALLET_ETH, key, password)
	err = wallet.LoadKey()
	if err != nil {
		return err
	}

	chain := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if chain == nil {
		return errors.New("Connect chain failed")
	}
	defer chain.DisConnect()

	contract := contract.NewContract(chain, address, contract.COMMON_CRONTACT)
	err = contract.SetABI(abi)
	if err != nil {
		return err
	}

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
	chainName := ctx.String(ChainFlag.Name)
	address := ctx.String(ContractFlag.Name)
	account := ctx.String(AccountFlag.Name)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	erc20, err := token.NewERC20(common.HexToAddress(address), c.(*chain.EthChain).Client)
	if err != nil {
		return err
	}

	balance, err := erc20.BalanceOf(nil, common.HexToAddress(account))
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Balance: %f\n", helper.WeiToEth(balance))
	return nil
}

func TransferERC20(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	key := ctx.String(KeyFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	contract := ctx.String(ContractFlag.Name)
	to := ctx.String(ToFlag.Name)
	value := ctx.String(ValueFlag.Name)
	file := ctx.String(FileFlag.Name)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	addresss := strings.Split(to, ",")
	values := strings.Split(value, ",")
	if len(addresss) != len(values) {
		return errors.New("Not match the address and value list")
	}

	if file != "" {
		alist, vlist, err := helper.ReadTransferFile(file)
		if err != nil {
			return err
		}

		addresss = append(addresss, alist...)
		values = append(values, vlist...)
	}

	total := new(big.Int)
	valueList := make([]*big.Int, 0)
	for _, v := range values {
		fv, _ := strconv.ParseFloat(v, 64)
		valueList = append(valueList, helper.EthToWei(float32(fv)))
		total = new(big.Int).Add(total, helper.EthToWei(float32(fv)))
	}

	logger.Debug("Total address: %d, total balance: %f", len(addresss), helper.WeiToEth(total))

	wallet := wallet.NewWallet(wallet.WALLET_ETH, key, password)
	err = wallet.LoadKey()
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	erc20, err := token.NewERC20(common.HexToAddress(contract), c.(*chain.EthChain).Client)
	if err != nil {
		return err
	}

	opts, err := c.(*chain.EthChain).GenTransOpts(wallet, nil)
	if err != nil {
		return err
	}

	for index, v := range addresss {
		tx, err := erc20.Transfer(opts, common.HexToAddress(v), valueList[index])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Tranfer %f to %s failed with err: %v\n", helper.WeiToEth(valueList[index]), v, err)
		} else {
			fmt.Fprintf(os.Stderr, "Tranfer %f to %s with transaction %s\n", helper.WeiToEth(valueList[index]), v, tx.Hash().Hex())
		}
	}

	return nil
}

func MergeERC20(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	keydir := ctx.String(KeyDirFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	contract := ctx.String(ContractFlag.Name)
	to := ctx.String(ToFlag.Name)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	wallet, err := wallet.ListWallet(wallet.WALLET_ETH, keydir, password)
	if err != nil {
		return err
	}

	logger.Debug("Total merge number is %d", len(wallet))

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	erc20, err := token.NewERC20(common.HexToAddress(contract), c.(*chain.EthChain).Client)
	if err != nil {
		return err
	}

	fromList := make([]string, 0)
	valueList := make([]string, 0)

	total := new(big.Int)
	for _, w := range wallet {
		if w.Address() == to {
			continue
		}

		opts, err := c.(*chain.EthChain).GenTransOpts(w, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gen transaction options for address %s failed with err: %v\n", w.Address(), err)
			continue
		}

		value, err := erc20.BalanceOf(nil, common.HexToAddress(w.Address()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "query balance for address %s failed with err: %v\n", w.Address(), err)
			continue
		}

		if value.Cmp(big.NewInt(0)) > 0 {
			tx, err := erc20.Transfer(opts, common.HexToAddress(to), value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "transfer balance for address %s failed with err: %v\n", w.Address(), err)
				continue
			} else {
				fmt.Fprintf(os.Stderr, "transfer balance for address %s success with transaction: %s\n", w.Address(), tx.Hash().Hex())
			}
		}

		total = new(big.Int).Add(total, value)
		fromList = append(fromList, w.Address())
		valueList = append(valueList, strconv.FormatFloat(float64(helper.WeiToEth(value)), 'f', 5, 32))
	}

	err = helper.SaveTransferFile(fromList, valueList, "./transfer.xlsx")
	if err != nil {
		log.Warn("Write transfer log failed with %v", err)
	}

	fmt.Printf("Total merge balance %f\n", helper.WeiToEth(total))
	return nil
}

func ApproveERC20(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	key := ctx.String(KeyFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	contract := ctx.String(ContractFlag.Name)
	to := ctx.String(ToFlag.Name)
	value := ctx.String(ValueFlag.Name)
	fv, _ := strconv.ParseFloat(value, 64)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	wallet := wallet.NewWallet(wallet.WALLET_ETH, key, password)
	err = wallet.LoadKey()
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	erc20, err := token.NewERC20(common.HexToAddress(contract), c.(*chain.EthChain).Client)
	if err != nil {
		return err
	}

	opts, err := c.(*chain.EthChain).GenTransOpts(wallet, nil)
	if err != nil {
		return err
	}

	allowance, err := erc20.Allowance(nil, common.HexToAddress(wallet.Address()), common.HexToAddress(to))
	if err != nil {
		return err
	}

	if allowance.Cmp(helper.EthToWei(float32(fv))) == 0 {
		fmt.Fprintf(os.Stderr, "No need to modify approve info")
		return nil
	}

	tx, err := erc20.Approve(opts, common.HexToAddress(to), helper.EthToWei(float32(fv)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Approve %s to %s failed with err: %v\n", value, to, err)
	} else {
		fmt.Fprintf(os.Stderr, "Approve %s to %s with transaction %s\n", value, to, tx.Hash().Hex())
	}

	return nil
}

func QueryERC721(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	address := ctx.String(ContractFlag.Name)
	account := ctx.String(AccountFlag.Name)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	erc721, err := token.NewERC721(common.HexToAddress(address), c.(*chain.EthChain).Client)
	if err != nil {
		return err
	}

	balance, err := erc721.BalanceOf(nil, common.HexToAddress(account))
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Balance: %f\n", helper.WeiToEth(balance))
	return nil
}

func TransferERC721(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	key := ctx.String(KeyFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	contract := ctx.String(ContractFlag.Name)
	to := ctx.String(ToFlag.Name)
	value := ctx.String(ValueFlag.Name)
	tokenId, _ := new(big.Int).SetString(value, 10)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	wallet := wallet.NewWallet(wallet.WALLET_ETH, key, password)
	err = wallet.LoadKey()
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	erc721, err := token.NewERC721(common.HexToAddress(contract), c.(*chain.EthChain).Client)
	if err != nil {
		return err
	}

	owner, err := erc721.OwnerOf(nil, tokenId)
	if err != nil {
		return err
	}

	approve, err := erc721.GetApproved(nil, tokenId)
	if err != nil {
		return err
	}

	if wallet.Address() != owner.Hex() && wallet.Address() != approve.Hex() {
		return errors.New("Not owner or approver for token")
	}

	opts, err := c.(*chain.EthChain).GenTransOpts(wallet, nil)
	if err != nil {
		return err
	}

	tx, err := erc721.SafeTransferFrom(opts, common.HexToAddress(wallet.Address()), common.HexToAddress(to), tokenId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Tranfer %s to %s failed with err: %v\n", value, to, err)
	} else {
		fmt.Fprintf(os.Stderr, "Tranfer %s to %s with transaction %s\n", value, to, tx.Hash().Hex())
	}

	return nil
}

func MergeERC721(ctx *cli.Context) error {
	return nil
}

func ApproveERC721(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	key := ctx.String(KeyFlag.Name)
	password := ctx.String(PasswordFlag.Name)
	contract := ctx.String(ContractFlag.Name)
	to := ctx.String(ToFlag.Name)
	value := ctx.String(ValueFlag.Name)
	tokenId, _ := new(big.Int).SetString(value, 10)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	wallet := wallet.NewWallet(wallet.WALLET_ETH, key, password)
	err = wallet.LoadKey()
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	erc721, err := token.NewERC721(common.HexToAddress(contract), c.(*chain.EthChain).Client)
	if err != nil {
		return err
	}

	opts, err := c.(*chain.EthChain).GenTransOpts(wallet, nil)
	if err != nil {
		return err
	}

	if tokenId != nil {
		allowance, err := erc721.GetApproved(nil, tokenId)
		if err != nil {
			return err
		}

		if allowance == common.HexToAddress(to) {
			fmt.Fprintf(os.Stderr, "No need to modify approve info")
			return nil
		}

		tx, err := erc721.Approve(opts, common.HexToAddress(to), tokenId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Approve %s to %s failed with err: %v\n", value, to, err)
		} else {
			fmt.Fprintf(os.Stderr, "Approve %s to %s with transaction %s\n", value, to, tx.Hash().Hex())
		}
	} else {
		all, err := erc721.IsApprovedForAll(nil, common.HexToAddress(wallet.Address()), common.HexToAddress(to))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Get approve all information failed with error: %v\n", err)
			return nil
		}

		if strings.ToLower(value) == "true" && all == false {
			tx, err := erc721.SetApprovalForAll(opts, common.HexToAddress(to), true)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Approve all to %s failed with err: %v\n", to, err)
			} else {
				fmt.Fprintf(os.Stderr, "Approve all to %s with transaction %s\n", to, tx.Hash().Hex())
			}
		} else if strings.ToLower(value) == "false" && all == true {
			tx, err := erc721.SetApprovalForAll(opts, common.HexToAddress(to), false)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Revoke all to %s failed with err: %v\n", to, err)
			} else {
				fmt.Fprintf(os.Stderr, "Revoke all to %s with transaction %s\n", to, tx.Hash().Hex())
			}
		} else {
			fmt.Fprintf(os.Stderr, "No need to set approve information\n")
		}
	}

	return nil
}

func PropertyQuery(ctx *cli.Context) error {
	chainName := ctx.String(ChainFlag.Name)
	contract := ctx.String(ContractFlag.Name)
	value := ctx.String(ValueFlag.Name)
	tokenId, _ := new(big.Int).SetString(value, 10)

	meta, err := chain.ChainMetaByName(chainName)
	if err != nil {
		return err
	}

	c := chain.NewChain(meta.Id, meta.Currency, meta.Name)
	if c == nil {
		return errors.New("Connect chain failed")
	}
	defer c.DisConnect()

	erc721, err := token.NewERC721(common.HexToAddress(contract), c.(*chain.EthChain).Client)
	if err != nil {
		return err
	}

	url, err := erc721.TokenURI(nil, tokenId)
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

	fmt.Fprintf(os.Stderr, "Result: 0x%s\n", hex.EncodeToString(result))
	return nil
}

func DecodeABI(ctx *cli.Context) error {
	method := ctx.String(FuncFlag.Name)
	data := ctx.String(DataFlag.Name)
	sign := ctx.Bool(SignFlag.Name)

	contract := contract.NewContract(nil, "", contract.COMMON_CRONTACT)
	result, err := contract.DecodeABI(method, common.FromHex(data), sign)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Result: %s\n", result)
	return nil
}
