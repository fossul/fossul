package util

import (
	"io/ioutil"
	"log"
	"github.com/BurntSushi/toml"
	"os"
	"bufio"
	"regexp"
	"encoding/json"
	"net/http"
)

type Config struct {
	ProfileName string `json:"profileName,omitempty"`
	ConfigName string `json:"configName,omitempty"`
	WorkflowId string `json:"workflowId,omitempty"`
	PluginDir string `json:"pluginDir"`
	AppPlugin string `json:"appPlugin"`
	StoragePlugin string `json:"storagePlugin"`
	ArchivePlugin string `json:"archivePlugin"`
	AutoDiscovery bool `json:"autoDiscovery"`
	BackupRetentions []BackupRetention `json:"backupRetentions"`
	SelectedBackupPolicy string `json:"backupPolicy,omitmepty"`
	SelectedBackupRetention int `json:"backupRetention,omitmepty"`
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
	SendTrapErrorCmd string `json:"sendTrapErrorCmd,omitempty"`
	SendTrapSuccessCmd string `json:"sendTrapSuccessCmd,omitempty"`
	AppPluginParameters map[string]string `json:"appPluginParameters,omitempty"`
	StoragePluginParameters map[string]string `json:"storagePluginParameters,omitempty"`
}

type BackupRetention struct {
	Policy string `json:"policy"`
	RetentionDays int `json:"retentionDays"`	
}

func ReadConfig(filename string) Config {
	var config Config
    b, err := ioutil.ReadFile(filename)
    if err != nil {
		log.Println("here",err)
        return config
    } else {
		str := string(b)
		config = decodeConfig(str)

		return config
	}
}

func decodeConfig(blob string) Config {
	var config Config
	if _, err := toml.Decode(blob, &config); err != nil {
		log.Println(err)
		return config
	}

	return config
}

func ReadConfigToMap(filename string) map[string]string {
	var configMap = map[string]string{}

	file, err := os.Open(filename)
    if err != nil {
		log.Println(err)
		return configMap
    }
    defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Println(err)
		return configMap
	}

    for scanner.Scan() {
		re := regexp.MustCompile(`(\S+)\s*=\s*\"(\S+)\"`)
		match := re.FindStringSubmatch(scanner.Text())

		if len(match) != 0 {
			configMap[match[1]] = match[2]
		}
    }
	
	return configMap
}

func SetAppPluginParameters(appConfigPath string, config Config) Config {
	configAppMap := make(map[string]string)

	if len(config.AppPlugin) != 0 {
		configAppMap = ReadConfigToMap(appConfigPath)
	}
	config.AppPluginParameters = configAppMap

	return config
}

func SetStoragePluginParameters(storageConfigPath string, config Config) Config {
	configStorageMap := make(map[string]string)

	if len(config.StoragePlugin) != 0 {
		configStorageMap = ReadConfigToMap(storageConfigPath)
	}
	config.StoragePluginParameters = configStorageMap

	return config
}

func GetConfig(w http.ResponseWriter, r *http.Request) Config {

	var config Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		log.Println(err)
		return config
	}
	defer r.Body.Close()
 
	res,err := json.Marshal(&config)
	if err != nil {
		log.Println(err)
		return config
	}

	log.Println("DEBUG", string(res))

	return config
}

func ConfigMapToJson(configMap map[string]string) string {
	jsonString, err := json.Marshal(configMap)
	if err != nil {
		log.Println(err)
		return ""	
	}

	return string(jsonString)
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