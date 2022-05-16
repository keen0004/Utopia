package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"utopia/internal/helper"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

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

func main() {
	// fmt.Printf(ParseParams("transfer(0x1234, 0x5678, 100)"))

	data := "0x01234567890abcdef"

	fmt.Printf("0x%s\n", hex.EncodeToString(crypto.Keccak256(helper.Str2bytes(data))))

	var buf []byte
	hash := sha3.NewLegacyKeccak256()
	hash.Write(helper.Str2bytes(data))
	buf = hash.Sum(buf)

	fmt.Printf("0x%s\n", hex.EncodeToString(buf))
}
