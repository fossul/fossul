package util

import (
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"os"
	"bufio"
	"regexp"
	"encoding/json"
	"net/http"
	"bytes"
)

type Config struct {
	ProfileName string `json:"profileName,omitempty"`
	ConfigName string `json:"configName,omitempty"`
	WorkflowId string `json:"workflowId,omitempty"`
	AppPlugin string `json:"appPlugin"`
	StoragePlugin string `json:"storagePlugin"`
	ArchivePlugin string `json:"archivePlugin"`
	AutoDiscovery bool `json:"autoDiscovery"`
	JobRetention int `json:"jobRetention"`
	BackupRetentions []BackupRetention `json:"backupRetentions"`
	SelectedBackupPolicy string `json:"backupPolicy,omitmepty"`
	SelectedBackupRetention int `json:"backupRetention,omitmepty"`
	SelectedWorkflowId int `json:"selectedWorkflowId,omitmepty"`
	PreAppQuiesceCmd string `json:"preAppQuiesceCmd,omitempty"`
	AppQuiesceCmd string `json:"appQuiesceCmd,omitempty"`
	PostAppQuiesceCmd string `json:"postAppQuiesceCmd,omitempty"`
	BackupCreateCmd string `json:"backupCreateCmd,omitempty"`
	BackupDeleteCmd string `json:"backupDeleteCmd,omitempty"`
	ArchiveCreateCmd string `json:"archiveCreateCmd,omitempty"`
	ArchiveDeleteCmd string `json:"archiveDeleteCmd,omitempty"`	
	PreAppUnquiesceCmd string `json:"preAppUnquiesceCmd,omitempty"`
	AppUnquiesceCmd string `json:"appUnquiesceCmd,omitempty"`
	PostAppUnquiesceCmd string `json:"postAppUnquiesceCmd,omitempty"`
	PreAppRestoreCmd string `json:"preAppRestoreCmd,omitempty"`
	RestoreCmd string `json:"restoreCmd,omitempty"`
	PostAppRestoreCmd string `json:"postAppRestoreCmd,omitempty"`
	SendTrapErrorCmd string `json:"sendTrapErrorCmd,omitempty"`
	SendTrapSuccessCmd string `json:"sendTrapSuccessCmd,omitempty"`
	AppPluginParameters map[string]string `json:"appPluginParameters,omitempty"`
	StoragePluginParameters map[string]string `json:"storagePluginParameters,omitempty"`
}

type BackupRetention struct {
	Policy string `json:"policy"`
	RetentionDays int `json:"retentionDays"`	
}

func ReadConfig(filePath string) (Config,error) {
	var config Config
	b, err := ioutil.ReadFile(filePath)
	
  if err != nil {
    return config,err
  } else {
		str := string(b)
		config,err = decodeConfig(str)

		if err != nil {
			return config,err
		}

		return config,nil
	}
}

func WriteConfig(filePath string,config Config) error {
	buf,err := EncodeConfig(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func WritePluginConfig(filePath string,configMap map[string]string) error {
	buf,err := EncodePluginConfig(configMap)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func EncodePluginConfig(configMap map[string]string) (*bytes.Buffer,error) {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(configMap); err != nil {
		return buf,err
	}

	return buf,nil
}

func EncodeConfig(config Config) (*bytes.Buffer,error) {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		return buf,err
	}

	return buf,nil
}

func decodeConfig(blob string) (Config,error) {
	var config Config
	if _, err := toml.Decode(blob, &config); err != nil {
		return config,err
	}

	return config,nil
}

func ReadConfigToMap(filePath string) (map[string]string,error) {
	var configMap = map[string]string{}

	file, err := os.Open(filePath)
    if err != nil {
			return configMap,err
    }
    defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		return configMap,err
	}

  for scanner.Scan() {
		re := regexp.MustCompile(`(\S+)\s*=\s*\"(\S+)\"`)
		match := re.FindStringSubmatch(scanner.Text())

		if len(match) != 0 {
			configMap[match[1]] = match[2]
		}
  }
	
	return configMap,nil
}

func SetAppPluginParameters(appConfigPath string, config Config) (Config,error) {
	var err error
	configAppMap := make(map[string]string)

	if len(config.AppPlugin) != 0 {
		configAppMap,err = ReadConfigToMap(appConfigPath)
		if err != nil {
			return config,err
		}
	}
	config.AppPluginParameters = configAppMap

	return config,nil
}

func SetStoragePluginParameters(storageConfigPath string, config Config) (Config,error) {
	var err error
	configStorageMap := make(map[string]string)

	if len(config.StoragePlugin) != 0 {
		configStorageMap,err = ReadConfigToMap(storageConfigPath)
		if err != nil {
			return config,err
		}
	}
	config.StoragePluginParameters = configStorageMap

	return config,nil
}

func GetConfig(w http.ResponseWriter, r *http.Request) (Config,error) {

	var config Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		return config,err
	}
	defer r.Body.Close()
 
	_,err := json.Marshal(&config)
	if err != nil {
		return config,err
	}

	return config,nil
}

func GetPluginConfig(w http.ResponseWriter, r *http.Request) (map[string]string,error) {

	var configMap map[string]string
	if err := json.NewDecoder(r.Body).Decode(&configMap); err != nil {
		return configMap,err
	}
	defer r.Body.Close()
 
	_,err := json.Marshal(&configMap)
	if err != nil {
		return configMap,err
	}

	return configMap,nil
}

func ConfigMapToJson(configMap map[string]string) (string,error) {
	jsonString, err := json.Marshal(configMap)
	if err != nil {
		return " ",err
	}

	return string(jsonString),nil
}

func ExistsBackupRetention(policy string, retentions []BackupRetention) bool {
	for _, retention := range retentions {
		if retention.Policy == policy {
			return true
		}
	}
	return false
}

func GetBackupRetention(policy string, retentions []BackupRetention) int {
	for _, retention := range retentions {
		if retention.Policy == policy {
			return retention.RetentionDays
		}
	}
	return -1
}