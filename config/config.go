package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	defaultConfigFileName = ".ccavenue.json"
)

type Config struct {
	Host       string `json:"host,omitempty"`
	MerchantId int    `json:"merchant_id,omitempty"`
	AccessCode string `json:"access_code,omitempty"`
	WorkingKey string `json:"working_key,omitempty"`
	IVStr      string `json:"iv_str,omitempty"`
}

var c = Config{}

func Configuration(configFileName ...string) (Config, error) {

	if (c == Config{}) {

		var cfname string
		switch len(configFileName) {
		case 0:
			dirname, err := os.UserHomeDir()
			if err != nil {
				return c, err
			}
			cfname = fmt.Sprintf("%s/%s", dirname, defaultConfigFileName)
		case 1:
			cfname = configFileName[0]
		default:
			return c, fmt.Errorf("incorrect arguments for configuration file name")
		}

		configFile, err := os.Open(cfname)
		if err != nil {
			return c, err
		}
		defer configFile.Close()

		decoder := json.NewDecoder(configFile)
		err = decoder.Decode(&c)
		if err != nil {
			return c, err
		}
	}

	return c, nil
}
