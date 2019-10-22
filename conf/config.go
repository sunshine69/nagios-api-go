package conf

import (
	"io/ioutil"
	"encoding/json"
)

//Configuration - Application global config
type Configuration struct {
	Port string
	NagiosStatusFilePath string
}

//Config -
var Config Configuration

//LoadConfig - load config
func init() {
		if Config.Port == "" {
			jsonByte, err := ioutil.ReadFile("config.json")
			if err != nil {
				panic(err)
			}
			json.Unmarshal(jsonByte, &Config)
		}
}