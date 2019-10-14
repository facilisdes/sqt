package main

import (
	"fmt"
	"os"
	"sqt/config"
)

func main() {
	//config.GenerateExampleConfig()
	args := os.Args[1:]

	if len(args) == 0 {
		printEmptyParams()
		os.Exit(1)
	}

	command := args[0]
	switch command {
	case "generate":
		GenerateExampleConfig()
	case "pack":
		EncryptConfig()
	case "unpack":
		DecryptConfig()
	case "help":
		printHelp()
	default:
		printWrongParams()
		os.Exit(1)
	}
}

func printEmptyParams(){
	fmt.Println("No parameters passed.")
	printHelp()
}

func printWrongParams(){
	fmt.Println("Parameter(s) unrecognized.")
	printHelp()
}


func printHelp() {
	fmt.Println("Execute me in form \"sqtConfig {param}\". Allowed params:")
	fmt.Println("- \"generate\" to create example config file (.sqt_raw).")
	fmt.Println("- \"pack\" to encrypt example config file (.sqt_raw) into server-ready config file (.sqt)")
	fmt.Println("- \"unpack\" to decrypt server-ready config file (.sqt) into raw config file (.sqt_raw)")
}

func GenerateExampleConfig() {
	config.GenerateExampleConfig()
}

func EncryptConfig() {
	config.EncryptConfig()
}
func DecryptConfig() {
	config.DecryptConfig()
}