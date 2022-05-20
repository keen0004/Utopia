package helper

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

	method := strings.Trim(params[:index], " ")
	params = params[index+1:]

	index = strings.Index(params, ")")
	if index == -1 {
		return "", []string{}, errors.New("Invalid parameters")
	}

	params = params[:index]
	args := strings.Split(params, ",")

	return method, args, nil
}

// todo: not support array now
func Str2Type(input string, totype reflect.Type) (interface{}, error) {
	input = strings.Trim(input, "\"")

	switch totype {
	case reflect.TypeOf(uint8(0)):
	case reflect.TypeOf(uint16(0)):
	case reflect.TypeOf(uint32(0)):
	case reflect.TypeOf(uint64(0)):
		return strconv.ParseUint(strings.Trim(input, " "), 10, 64)
	case reflect.TypeOf(int8(0)):
	case reflect.TypeOf(int16(0)):
	case reflect.TypeOf(int32(0)):
	case reflect.TypeOf(int64(0)):
		return strconv.ParseInt(strings.Trim(input, " "), 10, 64)
	case reflect.TypeOf(&big.Int{}):
		result, ok := new(big.Int).SetString(strings.Trim(input, " "), 10)
		if !ok {
			return nil, errors.New("Convert big.int failed")
		}
		return result, nil
	case reflect.TypeOf(false):
		return strconv.ParseBool(strings.Trim(input, " "))
	case reflect.TypeOf(""):
		return input, nil
	case reflect.TypeOf(common.Address{}):
		return common.HexToAddress(strings.Trim(input, " ")), nil
	case reflect.ArrayOf(32, reflect.TypeOf(byte(0))):
		return common.FromHex(strings.Trim(input, " ")), nil
	case reflect.SliceOf(reflect.TypeOf(byte(0))):
		return common.FromHex(strings.Trim(input, " ")), nil
	default:
		return nil, errors.New("Not support type")
	}

	return nil, nil
}

// todo: not support array now
func Type2Str(input interface{}, itype reflect.Type) (string, error) {
	switch itype {
	case reflect.TypeOf(uint8(0)):
	case reflect.TypeOf(uint16(0)):
	case reflect.TypeOf(uint32(0)):
	case reflect.TypeOf(uint64(0)):
	case reflect.TypeOf(int8(0)):
	case reflect.TypeOf(int16(0)):
	case reflect.TypeOf(int32(0)):
	case reflect.TypeOf(int64(0)):
		return fmt.Sprintf("%d", input), nil
	case reflect.TypeOf(&big.Int{}):
		return input.(*big.Int).String(), nil
	case reflect.TypeOf(false):
		return fmt.Sprintf("%b", input), nil
	case reflect.TypeOf(""):
		return input.(string), nil
	case reflect.TypeOf(common.Address{}):
		return input.(common.Address).Hex(), nil
	case reflect.ArrayOf(32, reflect.TypeOf(byte(0))):
	case reflect.SliceOf(reflect.TypeOf(byte(0))):
		return fmt.Sprintf("0x%s", hex.EncodeToString(input.([]byte))), nil
	default:
		return "", errors.New("Not support type")
	}

	return "", nil
}

func Str2Array(args []string, index int, totype abi.Type) (interface{}, int, error) {
	if !strings.HasPrefix(args[index], "[") {
		return nil, index, errors.New("Need array paramter but not found")
	}

	result := reflect.MakeSlice(totype.GetType(), 0, 256)
	inarray := true

	// for empty slice
	if args[index] == "[]" {
		return result.Interface(), index + 1, nil
	}

	args[index] = args[index][1:]
	if strings.HasSuffix(args[index], "]") {
		args[index] = args[index][:len(args[index])-1]
		inarray = false
	}

	v, err := Str2Type(args[index], totype.Elem.GetType())
	if err != nil {
		return nil, index, err
	}

	result = reflect.Append(result, reflect.ValueOf(v))
	for {
		index++

		if !inarray {
			break
		}

		if len(args) <= index {
			break
		}

		if strings.HasSuffix(args[index], "]") {
			args[index] = args[index][:len(args[index])-1]
			inarray = false
		}

		v, err = Str2Type(args[index], totype.Elem.GetType())
		if err != nil {
			return nil, index, err
		}

		result = reflect.Append(result, reflect.ValueOf(v))
	}

	if inarray {
		return nil, index, errors.New("Invalid array values")
	}

	return result.Interface(), index, nil
}

func Array2Str(input interface{}, totype reflect.Type) (string, error) {
	var builder strings.Builder
	builder.WriteString("[")

	size := reflect.ValueOf(input).Len()
	for i := 0; i < size; i++ {
		output, err := Type2Str(reflect.ValueOf(input).Index(i).Interface(), totype)
		if err != nil {
			return "", err
		}

		builder.WriteString(output)
		if i != size-1 {
			builder.WriteString(",")
		}
	}

	builder.WriteString("]")
	return builder.String(), nil
}
