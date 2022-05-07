package main

import (
	"errors"
	"fmt"
	"strings"
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
	fmt.Printf(ParseParams("transfer(0x1234, 0x5678, 100)"))
}
