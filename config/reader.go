package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"sqt/sqtcipher"
	"strconv"
	"strings"
)

const (
	SERVER_CONFIG_PATH   = "./.sqtconfig"
	CLIENT_CONFIG_PATH   = "./.sqtconfig"
	MAIN_CONFIG_PATH     = "./.sqt"
	MAIN_CONFIG_RAW_PATH = "./.sqt_raw"
	MAIN_CONFIG_KEY      = "1Q2W3E$r%t^y"
)

var (
	values        = make(map[string]string)
	defaultValues = map[string]string{
		"sqt_port":             "13343",
		"db_host":              "localhost",
		"db_port":              "3306",
		"cache_host":           "localhost",
		"cache_port":           "6379",
		"db_password":          "",
		"db_key_column_name":   "ID",
		"db_value_column_name": "VALUE",
	}
	Values Params
)

type Params struct {
	ConnType string
	ConnPort string

	DbHost            string
	DbPort            string
	DbLogin           string
	DbPassword        string
	DbName            string
	DbTable           string
	DbKeyColumnName   string
	DbValueColumnName string

	CacheHost string
	CachePort string

	ReadTimeInit       int
	ReadTimeStep       int
	MaxStackSize       int
	ReadTimeGrowth     string
	ReadTimeParameter1 float64
	ReadTimeParameter2 float64
	ReadTimeParameter3 float64
}

func ReadServerConfigs() {
	readConfigFile(SERVER_CONFIG_PATH)
	readMainServerConfig()
	ParseServerConfigValues()
}
func ReadClientConfigs() {
	readConfigFile(CLIENT_CONFIG_PATH)
	ParseClientConfigValues()
}

func ParseServerConfigValues() {
	Values.ConnType = "tcp"

	if val, ok := values["sqt_port"]; ok {
		Values.ConnPort = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"sqt_port\" not passed; using default value (%[1]s) instead", defaultValues["sqt_port"]))
		Values.ConnPort = defaultValues["sqt_port"]
	}

	if val, ok := values["db_host"]; ok {
		Values.DbHost = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"db_host\" not passed; using default value (%[1]s) instead", defaultValues["db_host"]))
		Values.DbHost = defaultValues["db_host"]
	}

	if val, ok := values["db_port"]; ok {
		Values.DbPort = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"db_port\" not passed; using default value (%[1]s) instead", defaultValues["db_port"]))
		Values.DbPort = defaultValues["db_port"]
	}

	if val, ok := values["db_login"]; ok {
		Values.DbLogin = val
	} else {
		fmt.Println("Missed required parameter \"db_login\"!")
		os.Exit(1)
	}

	if val, ok := values["db_password"]; ok {
		Values.DbPassword = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"db_password\" not passed; using default value (%[1]s) instead", defaultValues["db_password"]))
		Values.DbPassword = defaultValues["db_password"]
	}

	if val, ok := values["db_name"]; ok {
		Values.DbName = val
	} else {
		fmt.Println("Missed required parameter \"db_name\"!")
		os.Exit(1)
	}

	if val, ok := values["db_table"]; ok {
		Values.DbTable = val
	} else {
		fmt.Println("Missed required parameter \"db_table\"!")
		os.Exit(1)
	}

	if val, ok := values["db_key_column_name"]; ok {
		Values.DbKeyColumnName = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"db_key_column_name\" not passed; using default value (%[1]s) instead", defaultValues["db_key_column_name"]))
		Values.DbKeyColumnName = defaultValues["db_key_column_name"]
	}

	if val, ok := values["db_value_column_name"]; ok {
		Values.DbValueColumnName = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"db_value_column_name\" not passed; using default value (%[1]s) instead", defaultValues["db_value_column_name"]))
		Values.DbValueColumnName = defaultValues["db_value_column_name"]
	}

	if val, ok := values["cache_host"]; ok {
		Values.CacheHost = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"cache_host\" not passed; using default value (%[1]s) instead", defaultValues["cache_host"]))
		Values.CacheHost = defaultValues["cache_host"]
	}

	if val, ok := values["cache_port"]; ok {
		Values.CachePort = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"cache_port\" not passed; using default value (%[1]s) instead", defaultValues["cache_port"]))
		Values.CachePort = defaultValues["cache_port"]
	}

	if val, ok := values["read_time_init"]; ok {
		ival, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error reading integer parameter \"read_time_init\" (%[1]s) %[2]s", val, err.Error()))
			os.Exit(1)
		}
		Values.ReadTimeInit = ival
	} else {
		fmt.Println("Missed required parameter \"read_time_init\"!")
		os.Exit(1)
	}

	if val, ok := values["read_time_step"]; ok {
		ival, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error reading integer parameter \"read_time_step\" (%[1]s) %[2]s", val, err.Error()))
			os.Exit(1)
		}
		Values.ReadTimeStep = ival
	} else {
		fmt.Println("Missed required parameter \"read_time_step\"!")
		os.Exit(1)
	}

	if val, ok := values["max_stack_size"]; ok {
		ival, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error reading integer parameter \"max_stack_size\" (%[1]s) %[2]s", val, err.Error()))
			os.Exit(1)
		}
		Values.MaxStackSize = ival
	} else {
		fmt.Println("Missed required parameter \"max_stack_size\"!")
		os.Exit(1)
	}

	if val, ok := values["read_time_growth"]; ok {
		Values.ReadTimeGrowth = val
	} else {
		fmt.Println("Missed required parameter \"read_time_growth\"!")
		os.Exit(1)
	}

	if val, ok := values["read_time_parameter1"]; ok {
		fval, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error reading integer parameter \"read_time_parameter1\" (%[1]s) %[2]s", val, err.Error()))
			os.Exit(1)
		}
		Values.ReadTimeParameter1 = fval
	}

	if val, ok := values["read_time_parameter2"]; ok {
		fval, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error reading integer parameter \"read_time_parameter2\" (%[1]s) %[2]s", val, err.Error()))
			os.Exit(1)
		}
		Values.ReadTimeParameter2 = fval
	}

	if val, ok := values["read_time_parameter3"]; ok {
		fval, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error reading integer parameter \"read_time_parameter3\" (%[1]s) %[2]s", val, err.Error()))
			os.Exit(1)
		}
		Values.ReadTimeParameter3 = fval
	}
}

func ParseClientConfigValues() {
	Values.ConnType = "tcp"

	if val, ok := values["sqt_port"]; ok {
		Values.ConnPort = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"sqt_port\" not passed; using default value (%[1]s) instead", defaultValues["sqt_port"]))
		Values.ConnPort = defaultValues["sqt_port"]
	}

	if val, ok := values["db_host"]; ok {
		Values.DbHost = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"db_host\" not passed; using default value (%[1]s) instead", defaultValues["db_host"]))
		Values.DbHost = defaultValues["db_host"]
	}

	if val, ok := values["db_port"]; ok {
		Values.DbPort = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"db_port\" not passed; using default value (%[1]s) instead", defaultValues["db_port"]))
		Values.DbPort = defaultValues["db_port"]
	}

	if val, ok := values["db_login"]; ok {
		Values.DbLogin = val
	} else {
		fmt.Println("Missed required parameter \"db_login\"!")
		os.Exit(1)
	}

	if val, ok := values["db_password"]; ok {
		Values.DbPassword = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"db_password\" not passed; using default value (%[1]s) instead", defaultValues["db_password"]))
		Values.DbPassword = defaultValues["db_password"]
	}

	if val, ok := values["db_name"]; ok {
		Values.DbName = val
	} else {
		fmt.Println("Missed required parameter \"db_name\"!")
		os.Exit(1)
	}

	if val, ok := values["db_table"]; ok {
		Values.DbTable = val
	} else {
		fmt.Println("Missed required parameter \"db_table\"!")
		os.Exit(1)
	}

	if val, ok := values["db_key_column_name"]; ok {
		Values.DbKeyColumnName = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"db_key_column_name\" not passed; using default value (%[1]s) instead", defaultValues["db_key_column_name"]))
		Values.DbKeyColumnName = defaultValues["db_key_column_name"]
	}

	if val, ok := values["db_value_column_name"]; ok {
		Values.DbValueColumnName = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"db_value_column_name\" not passed; using default value (%[1]s) instead", defaultValues["db_value_column_name"]))
		Values.DbValueColumnName = defaultValues["db_value_column_name"]
	}

	if val, ok := values["cache_host"]; ok {
		Values.CacheHost = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"cache_host\" not passed; using default value (%[1]s) instead", defaultValues["cache_host"]))
		Values.CacheHost = defaultValues["cache_host"]
	}

	if val, ok := values["cache_port"]; ok {
		Values.CachePort = val
	} else {
		fmt.Println(fmt.Sprintf("Parameter \"cache_port\" not passed; using default value (%[1]s) instead", defaultValues["cache_port"]))
		Values.CachePort = defaultValues["cache_port"]
	}
}

func readConfigFromFile(fileContent string) {
	lines := strings.Split(fileContent, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		keyValuePair := strings.Split(line, "=")
		values[keyValuePair[0]] = keyValuePair[1]
	}
}

func readConfigFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("Error reading server config:", err.Error())
		os.Exit(1)
	}
	readConfigFromFile(readFile(filePath))
}

func readMainServerConfig() {
	if _, err := os.Stat(MAIN_CONFIG_PATH); os.IsNotExist(err) {
		fmt.Println("Error reading main config:", err.Error())
		os.Exit(1)
	}
	fileEncContent := readFile(MAIN_CONFIG_PATH)
	fileDecContent := sqtcipher.Decrypt(fileEncContent, MAIN_CONFIG_KEY)
	readConfigFromFile(fileDecContent)
}

func readFile(filePath string) string {
	dat, err := ioutil.ReadFile(filePath)
	check(err)
	return string(dat)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
