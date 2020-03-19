package abi

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/sha3"
	"strings"
)

type ABI struct {
	jsonABI   []byte
	abi       []jsonItem
	functions []Function
}

func New(jsonABI []byte) (*ABI, error) {
	parsedABI, err := unmarshal(jsonABI)

	return &ABI{
		jsonABI: jsonABI,
		abi:     parsedABI,
	}, err
}

func (a *ABI) GetFunctionByName(name string) (*Function, error) {
	functions, err := a.GetFunctions()
	if err != nil {
		return nil, err
	}

	for _, f := range functions {
		if f.Name == name {
			return &f, nil
		}
	}

	return nil, nil
}

func (a *ABI) GetFunctions() ([]Function, error) {
	if len(a.functions) != 0 {
		return a.functions, nil
	}

	var f []Function

	items, err := unmarshal(a.jsonABI)
	if err != nil {
		return f, err
	}

	for _, item := range items {
		if item.Type != tFunction {
			continue
		}

		signature, err := getSignature(item.Name, item.Inputs, item.Outputs)
		if err != nil {
			return f, err
		}

		f = append(f, Function{
			Name:        item.Name,
			Signature:   signature,
			Hash:        getFunctionHash(item.Name, item.Inputs),
			ParamsCount: len(item.Inputs),
		})
	}

	a.functions = f

	return a.functions, nil
}

func getFunctionHash(name string, arguments []jsonParam) struct {
	Full      string
	Short     string
	Signature string
} {
	argumentTypes := make([]string, 0, len(arguments))
	for _, argument := range arguments {
		argumentTypes = append(argumentTypes, string(argument.Type))
	}

	signature := []byte(name + "(" + strings.Join(argumentTypes, ",") + ")")

	hash := sha3.NewLegacyKeccak256()
	hash.Write(signature)

	hashStr := "0x" + hex.EncodeToString(hash.Sum(nil))

	return struct {
		Full      string
		Short     string
		Signature string
	}{Full: hashStr, Short: hashStr[0:10], Signature: string(signature)}
}

func getSignature(name string, arguments []jsonParam, returnValues []jsonParam) (string, error) {
	var signature string

	args := make([]string, 0, len(arguments))
	for _, argument := range arguments {
		if len(argument.Name) > 0 {
			args = append(args, argument.Name+" "+string(argument.Type))
		} else {
			args = append(args, string(argument.Type))
		}
	}

	var returnValue = ""
	if len(returnValues) > 1 {
		return signature, errors.New("unexpected returnValues count: > 1")
	} else if len(returnValues) == 1 {
		if len(returnValues[0].Name) > 0 {
			returnValue = returnValues[0].Name + " " + string(returnValues[0].Type)
		} else {
			returnValue = string(returnValues[0].Type)
		}
	}

	signature = name + "(" + strings.Join(args, ", ") + ") " + returnValue

	return signature, nil
}

func unmarshal(jsonABI []byte) ([]jsonItem, error) {
	var items []jsonItem

	err := json.Unmarshal(jsonABI, &items)

	return items, err
}
