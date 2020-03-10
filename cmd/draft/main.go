package main

import (
	"fmt"
	"kickico/internal/pkg/abi"
	"kickico/internal/pkg/config"
	"os"
)

const configPath = "configs"
const configFile = "config.yml"

const cmdList = "list"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg, err := config.Read(configPath, configFile)
	if err != nil {
		panic(err)
	}

	functions, err := abi.GetFunctions(cfg.ABI)
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case cmdList:
		for _, f := range functions {
			fmt.Println(f.Hash.Short, " ", f.Signature)
		}
	default:
		fmt.Printf("unexpected subcommand '%s'\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println(os.Args[0], cmdList, "- lists all public methods of Ethereum contact")
}
