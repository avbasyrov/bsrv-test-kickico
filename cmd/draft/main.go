package main

import (
	"bsrv-test-kickico/internal/pkg/abi"
	"bsrv-test-kickico/internal/pkg/config"
	"bsrv-test-kickico/internal/pkg/infura"
	"bsrv-test-kickico/internal/pkg/typecast"
	"encoding/hex"
	"fmt"
	"os"
)

const configPath = "configs"
const configFile = "config.yml"

const cmdList = "list"
const cmdCall = "call"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg, err := config.Read(configPath, configFile)
	if err != nil {
		panic(err)
	}

	ABI, err := abi.New(cfg.ABIContent)
	//ABIContent, err := ethABI.JSON(strings.NewReader(string(cfg.ABIContent)))
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case cmdList:
		functions, err := ABI.GetFunctions()
		if err != nil {
			panic(err)
		}
		for _, f := range functions {
			fmt.Println(f.Hash.Short, " ", f.Signature)
		}
	case cmdCall:
		call(os.Args, cfg)
	default:
		fmt.Printf("unexpected subcommand '%s'\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println(os.Args[0], cmdList, "- lists all public methods of Ethereum contract")
	fmt.Println(os.Args[0], cmdCall, "param1 param2 ... - call public methods of Ethereum contract with params (if given)")
}

func call(osArgs []string, cfg config.Config) {
	if len(os.Args) < 3 {
		fmt.Println("You must provide function name to be called")
		printUsage()
		os.Exit(1)
	}

	ABI, err := abi.New(cfg.ABIContent)
	if err != nil {
		panic(err)
	}

	callFunctionName := os.Args[2]
	callFunction, err := ABI.GetFunctionByName(callFunctionName)
	if err != nil {
		panic(err)
	}
	if _, exists := cfg.ABI.Methods[callFunctionName]; !exists {
		fmt.Printf("Method %s not exists!\n", callFunctionName)
		os.Exit(1)
	}
	argumentsCount := len(cfg.ABI.Methods[callFunctionName].Inputs)
	if len(os.Args) != (1 + 1 + 1 + argumentsCount) {
		fmt.Printf("Method %s expect %d argument(s)!\n", callFunctionName, argumentsCount)
		os.Exit(1)
	}

	fmt.Println(callFunction.Hash.Short, " ", callFunction.Signature)

	args := make([]interface{}, argumentsCount, argumentsCount)
	for i := 0; i < argumentsCount; i = i + 1 {
		args[i], err = typecast.FromString(cfg.ABI.Methods[callFunctionName].Inputs[i].Type, os.Args[3+i])
		if err != nil {
			panic(err)
		}
	}

	res, err := cfg.ABI.Pack(callFunctionName, args...)
	if err != nil {
		panic(err)
	}

	requestData := "0x" + hex.EncodeToString(res)
	fmt.Println(requestData)

	// {"jsonrpc":"2.0","method":"eth_call","params": [{"from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155","to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567","gas": "0x76c0","gasPrice": "0x9184e72a000","value": "0x9184e72a","data": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"}, "latest"],"id":1}

	const pattern string = `{"jsonrpc":"2.0","method":"eth_call","params": [{"to": "%s","data": "%s"}, "latest"],"id":%d}`
	request := fmt.Sprintf(pattern, cfg.Address, requestData, 1)
	fmt.Println(request)

	infuraClient := infura.New(cfg.Infura.ProjectID, cfg.Infura.PrivateKey)
	response, err := infuraClient.Request(request)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}
