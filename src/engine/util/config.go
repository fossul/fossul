package util

import (
	"io/ioutil"
	"log"
	"github.com/BurntSushi/toml"
)

type Config struct {
	BackupName string `json:"backupName"`
	PluginDir string `json:"pluginDir"`
	AppPlugin string `json:"appPlugin"`
	StoragePlugin string `json:"storagePlugin"`
	BackupRetentions []BackupRetention `json:"backupRetentions"`
	PreAppQuiesceCmd string `json:"preAppQuiesceCmd,omitempty"`
	AppQuiesceCmd string `json:"appQuiesceCmd,omitempty"`
	PostAppQuiesceCmd string `json:"postAppQuiesceCmd,omitempty"`
	BackupCreateCmd string `json:"backupCreateCmd,omitempty"`
	PreAppUnquiesceCmd string `json:"preAppUnquiesceCmd,omitempty"`
	AppUnquiesceCmd string `json:"appUnquiesceCmd,omitempty"`
	PostAppUnquiesceCmd string `json:"postAppUnquiesceCmd,omitempty"`
	SendTrapErrorCmd string `json:"sendTrapErrorCmd,omitempty"`
	SendTrapSuccessCmd string `json:"sendTrapSuccessCmd,omitempty"`
}

type BackupRetention struct {
	Policy string `json:"policy"`
	RetentionDays int `json:"retentionDays"`	
}

func ReadConfig (filename string) Config {
    b, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Println(err)
    }

	str := string(b)
	var config Config = decodeConfig(str)
	
	return config
}

func decodeConfig (blob string) Config {
	var config Config
	if _, err := toml.Decode(blob, &config); err != nil {
  		log.Println(err)
	}

	return config
}