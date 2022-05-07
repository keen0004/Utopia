package helper

import (
	"encoding/hex"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/urfave/cli.v1"
)

func NewApp(version string, usage string) *cli.App {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = filepath.Base(os.Args[0])
	app.Author = "fuxideng"
	app.Email = "fuxideng@gmail.com"
	app.Version = version
	app.Usage = usage

	return app
}

func Str2bytes(data string) []byte {
	if strings.HasPrefix(data, "0x") || strings.HasPrefix(data, "0X") {
		data = data[2:]

		bytes, err := hex.DecodeString(data)
		if err != nil {
			return []byte("")
		}

		return bytes
	}

	return []byte(data)
}

func WeiToEth(wei *big.Int) float32 {
	eth, _ := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(big.NewInt(1e+18))).Float32()
	return eth
}

func EthToWei(eth float32) *big.Int {
	wei := new(big.Int)
	new(big.Float).Mul(big.NewFloat(float64(eth)), new(big.Float).SetInt(big.NewInt(1e+18))).Int(wei)
	return wei
}

func ParseParams(params string) (string, []string, error) {
	params = strings.Trim(params, " ")

	index := strings.Index(params, "(")
	if index == -1 {
		return "", []string{}, errors.New("Invalid parameters")
	}

	method := params[:index]
	params = params[index+1:]

	index = strings.Index(params, ")")
	if index == -1 {
		return "", []string{}, errors.New("Invalid parameters")
	}

	params = params[:index]
	args := strings.Split(params, ",")

	return method, args, nil
}

func Str2Type(input string, totype reflect.Type) (interface{}, error) {
	switch totype {
	case reflect.TypeOf(uint8(0)):
	case reflect.TypeOf(uint16(0)):
	case reflect.TypeOf(uint32(0)):
	case reflect.TypeOf(uint64(0)):
		return strconv.ParseUint(input, 10, 64)
	case reflect.TypeOf(int8(0)):
	case reflect.TypeOf(int16(0)):
	case reflect.TypeOf(int32(0)):
	case reflect.TypeOf(int64(0)):
		return strconv.ParseInt(input, 10, 64)
	case reflect.TypeOf(&big.Int{}):
		result, ok := new(big.Int).SetString(input, 10)
		if !ok {
			return nil, errors.New("Convert big.int failed")
		}
		return result, nil
	case reflect.TypeOf(false):
		return strconv.ParseBool(input)
	case reflect.TypeOf(""):
		return input, nil
	case reflect.TypeOf(common.Address{}):
		return common.HexToAddress(input), nil
	case reflect.ArrayOf(32, reflect.TypeOf(byte(0))):
		return Str2bytes(input), nil
	case reflect.SliceOf(reflect.TypeOf(byte(0))):
		return Str2bytes(input), nil
	default:
		return nil, errors.New("Not support type")
	}

	return nil, nil
}
