package helper

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"utopia/internal/cmc"
	"utopia/internal/excel"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/urfave/cli.v1"
)

var (
	TRANSFER_SHEET_NAME  = "transfer"
	TRANSFER_LIST_HEADER = []string{"index", "from", "to", "value", "notes"}

	CURRENCY_SHEET_NAME  = "currency"
	CURRENCY_LIST_HEADER = []string{"id", "symbol", "rank", "current", "total", "pairs", "platform", "address", "price", "marketcap", "lastupdated"}
)

// transfer information
type TransferInfo struct {
	From  string
	To    string
	Value string
	Notes string
}

// create new app and init properties
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

// change unit form wei to ether
func WeiToEth(wei *big.Int) float32 {
	eth, _ := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(big.NewInt(1e+18))).Float32()
	return eth
}

// change unit form ether to wei
func EthToWei(eth float32) *big.Int {
	wei := new(big.Int)
	new(big.Float).Mul(big.NewFloat(float64(eth)), new(big.Float).SetInt(big.NewInt(1e+18))).Int(wei)
	return wei
}

func DefaultVlue(value string, def string) string {
	if value == "" {
		return def
	}

	return value
}

// parse function sig and call data
func ParseParams(params string) (string, []string, error) {
	params = strings.Trim(params, " ")

	// parse method
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

	// parse parameters
	params = params[:index]
	args := strings.Split(params, ",")

	return method, args, nil
}

// change call string value to dst type
func Str2Type(input string, totype reflect.Type) (interface{}, error) {
	input = strings.Trim(input, "\"")

	switch totype {
	case reflect.TypeOf(uint8(0)):
	case reflect.TypeOf(uint16(0)):
	case reflect.TypeOf(uint32(0)):
	case reflect.TypeOf(uint64(0)):
		// string to decimal in base 10
		return strconv.ParseUint(strings.Trim(input, " "), 10, 64)
	case reflect.TypeOf(int8(0)):
	case reflect.TypeOf(int16(0)):
	case reflect.TypeOf(int32(0)):
	case reflect.TypeOf(int64(0)):
		// string to decimal in base 10
		return strconv.ParseInt(strings.Trim(input, " "), 10, 64)
	case reflect.TypeOf(&big.Int{}):
		// string to big.int in base 10
		result, ok := new(big.Int).SetString(strings.Trim(input, " "), 10)
		if !ok {
			return nil, errors.New("Convert big.int failed")
		}
		return result, nil
	case reflect.TypeOf(false):
		// string to bool in base 10
		return strconv.ParseBool(strings.Trim(input, " "))
	case reflect.TypeOf(""):
		// string return and not trim space
		return input, nil
	case reflect.TypeOf(common.Address{}):
		// hex string to address
		return common.HexToAddress(strings.Trim(input, " ")), nil
	case reflect.ArrayOf(32, reflect.TypeOf(byte(0))):
		// hex string to byte32
		return common.FromHex(strings.Trim(input, " ")), nil
	case reflect.SliceOf(reflect.TypeOf(byte(0))):
		// hex string to []byte
		return common.FromHex(strings.Trim(input, " ")), nil
	default:
		return nil, errors.New("Not support type")
	}

	return nil, nil
}

// change result data to string value
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
		// decimal to string
		return fmt.Sprintf("%d", input), nil
	case reflect.TypeOf(&big.Int{}):
		// big.int to string
		return input.(*big.Int).String(), nil
	case reflect.TypeOf(false):
		// bool to string
		return fmt.Sprintf("%b", input), nil
	case reflect.TypeOf(""):
		// return string
		return input.(string), nil
	case reflect.TypeOf(common.Address{}):
		// address to string
		return input.(common.Address).Hex(), nil
	case reflect.ArrayOf(32, reflect.TypeOf(byte(0))):
	case reflect.SliceOf(reflect.TypeOf(byte(0))):
		// []byte to string
		return fmt.Sprintf("0x%s", common.Bytes2Hex(input.([]byte))), nil
	default:
		return "", errors.New("Not support type")
	}

	return "", nil
}

// change string to array for call function
func Str2Array(args []string, index int, totype abi.Type) (interface{}, int, error) {
	if !strings.HasPrefix(args[index], "[") {
		return nil, index, errors.New("Need array paramter but not found")
	}

	// max slice size is 256
	result := reflect.MakeSlice(totype.GetType(), 0, 256)
	inarray := true

	// for empty slice
	if args[index] == "[]" {
		return result.Interface(), index + 1, nil
	}

	// delete [] for one element
	args[index] = args[index][1:]
	if strings.HasSuffix(args[index], "]") {
		args[index] = args[index][:len(args[index])-1]
		inarray = false
	}

	// change to dst type
	v, err := Str2Type(args[index], totype.Elem.GetType())
	if err != nil {
		return nil, index, err
	}

	// add first element and interate the left
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

	// check array is completed
	if inarray {
		return nil, index, errors.New("Invalid array values")
	}

	return result.Interface(), index, nil
}

// change array to string for call result
func Array2Str(input interface{}, totype reflect.Type) (string, error) {
	var builder strings.Builder
	builder.WriteString("[")

	// iterate times to array size
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

func ReadTransferFile(path string) ([]TransferInfo, error) {
	// open excel file to read list
	file, err := excel.NewExcel(path)
	if err != nil {
		return nil, err
	}

	err = file.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close(false)

	data, err := file.ReadAll(TRANSFER_SHEET_NAME)
	if err != nil {
		return nil, err
	}

	result := make([]TransferInfo, 0)
	for index, row := range data {
		// skip the header
		if index == 0 {
			continue
		}

		if len(row) < len(TRANSFER_LIST_HEADER) {
			return nil, errors.New("Invalid file format")
		}

		// {"index", "from", "to", "value", "notes"}
		result = append(result, TransferInfo{
			From:  row[1],
			To:    row[2],
			Value: row[3],
			Notes: row[4],
		})
	}

	return result, nil
}

func SaveTransferFile(info []TransferInfo, path string) error {
	// open excel file to write list
	file, err := excel.NewExcel(path)
	if err != nil {
		return err
	}

	err = file.Open()
	if err != nil {
		return err
	}
	defer file.Close(true)

	data := make([][]string, 0)
	data = append(data, TRANSFER_LIST_HEADER)

	// {"index", "from", "to", "value", "notes"}
	for i, key := range info {
		row := make([]string, 0, len(TRANSFER_LIST_HEADER))
		row = append(row, strconv.Itoa(i+1))
		row = append(row, DefaultVlue(key.From, "0x"))
		row = append(row, DefaultVlue(key.To, "0x"))
		row = append(row, DefaultVlue(key.Value, "0"))
		row = append(row, DefaultVlue(key.Notes, "x"))

		data = append(data, row)
	}

	return file.WriteAll(TRANSFER_SHEET_NAME, data)
}

func WriteCurrencyFile(path string, money string, list []*cmc.Listing) error {
	// open excel file to write list
	file, err := excel.NewExcel(path)
	if err != nil {
		return err
	}

	err = file.Open()
	if err != nil {
		return err
	}
	defer file.Close(true)

	data := make([][]string, 0)
	data = append(data, CURRENCY_LIST_HEADER)

	// {"id", "symbol", "rank", "current", "total", "pairs", "platform", "address", "price", "marketcap", "lastupdated"}
	for _, info := range list {
		row := make([]string, 0, len(CURRENCY_LIST_HEADER))
		row = append(row, strconv.Itoa(int(info.ID)))
		row = append(row, info.Symbol)
		row = append(row, strconv.Itoa(int(info.CMCRank)))
		row = append(row, strconv.Itoa(int(info.CirculatingSupply)))
		row = append(row, strconv.Itoa(int(info.TotalSupply)))
		row = append(row, strconv.Itoa(int(info.NumMarketPairs)))
		row = append(row, DefaultVlue(info.Platform.Symbol, "x"))
		row = append(row, DefaultVlue(info.Platform.TokenAddress, "0x"))
		quote, ok := info.Quote[money]
		if ok {
			row = append(row, fmt.Sprintf("%.5f", quote.Price))
			row = append(row, fmt.Sprintf("%.2f", quote.MarketCap))
			row = append(row, quote.LastUpdated)
		} else {
			row = append(row, "0")
			row = append(row, "0")
			row = append(row, "x")
		}

		data = append(data, row)
	}

	return file.WriteAll(CURRENCY_SHEET_NAME, data)
}
