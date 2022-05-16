package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"utopia/internal/contract"
	"utopia/internal/helper"

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
					AccountFlag,
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
	return nil
}

func CallContract(ctx *cli.Context) error {
	return nil
}

func ListContract(ctx *cli.Context) error {
	return nil
}

func QueryERC20(ctx *cli.Context) error {
	return nil
}

func TransferERC20(ctx *cli.Context) error {
	return nil
}

func MergeERC20(ctx *cli.Context) error {
	return nil
}

func ApproveERC20(ctx *cli.Context) error {
	return nil
}

func QueryERC721(ctx *cli.Context) error {
	return nil
}

func TransferERC721(ctx *cli.Context) error {
	return nil
}

func MergeERC721(ctx *cli.Context) error {
	return nil
}

func ApproveERC721(ctx *cli.Context) error {
	return nil
}

func PropertyQuery(ctx *cli.Context) error {
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
	result, err := contract.DecodeABI(method, helper.Str2bytes(data), sign)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Result: %s\n", result)
	return nil
}
