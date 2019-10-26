package config

import (
	"io/ioutil"
	"sqt/sqtcipher"
)

var defaultMainConfig = "read_time_init=500\n" +
	"read_time_step=1000\n" +
	"max_stack_size=10\n" +
	"read_time_growth=sum\n" +
	"read_time_parameter=0\n"

func EncryptConfig() {
	configFile := readFile(MAIN_CONFIG_RAW_PATH)
	configFile = sqtcipher.Encrypt(configFile, MAIN_CONFIG_KEY)
	saveFile(configFile, MAIN_CONFIG_PATH)
}

func GenerateExampleConfig() {
	fileContent := defaultMainConfig
	saveFile(fileContent, MAIN_CONFIG_RAW_PATH)
}

func DecryptConfig() {
	configFile := readFile(MAIN_CONFIG_PATH)
	configFile = sqtcipher.Decrypt(configFile, MAIN_CONFIG_KEY)
	saveFile(configFile, MAIN_CONFIG_RAW_PATH)
}

func saveFile(fileContent string, filePath string) {
	d1 := []byte(fileContent)
	err := ioutil.WriteFile(filePath, d1, 0644)
	check(err)
}
