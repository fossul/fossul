package main

import (
	"fossil/src/engine/util"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"errors"
)

type appPlugin string
var AppPlugin appPlugin

func (a appPlugin) SetEnv(config util.Config) util.Result {
	var result util.Result

	return result
}	

func (a appPlugin) Discover(config util.Config) util.DiscoverResult {
	var discoverResult util.DiscoverResult
	var discoverList []util.Discover
	var discover util.Discover
	var result util.Result
	var messages []util.Message

	dsn := getDSN(config)

	conn,err := getConn(dsn)

	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't connect to database [" + config.AppPluginParameters["PqDb"] + "] " + err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1,messages)
		discoverResult.Result = result
		return discoverResult
	} else {
		defer conn.Close()
		msg := util.SetMessage("INFO", "Connection to database [" + config.AppPluginParameters["PqDb"] + "] established")
		messages = append(messages,msg)
		result = util.SetResult(0,messages)
	}

	discover.Instance = config.AppPluginParameters["PqDb"]

	var value string

	err = conn.QueryRow("show data_directory").Scan(&value)
	if err != nil {
		msg := util.SetMessage("ERROR","Discovery for database [" + config.AppPluginParameters["PqDb"] + "] failed! " + err.Error())
		messages = append(messages,msg)
		result = util.SetResult(1,messages)

		discoverResult.Result = result
		return discoverResult
	}
	var dataFilePaths []string
	dataDir := value
	dataFilePaths = append(dataFilePaths,dataDir)
	discover.DataFilePaths = dataFilePaths

	msg := util.SetMessage("INFO", "Data Directory is [" + value + "]")
	messages = append(messages,msg)

	discoverList = append(discoverList,discover)

	result = util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}	

func (a appPlugin) Quiesce(config util.Config) util.Result {	

	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	dsn := getDSN(config)
	conn,err := getConn(dsn)

	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't connect to database [" + config.AppPluginParameters["PqDb"] + "] " + err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1,messages)
		return result
	} else {
		defer conn.Close()
		msg := util.SetMessage("INFO", "Connection to database [" + config.AppPluginParameters["PqDb"] + "] established")
		messages = append(messages,msg)
		result = util.SetResult(0,messages)
	}

	backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"],config.SelectedBackupPolicy,config.WorkflowId)

	msg := util.SetMessage("INFO","Entering backup mode using label " + backupName + " for database [" + config.AppPluginParameters["PqDb"] + "]")
	messages = append(messages,msg)

	_, err = conn.Exec("SELECT pg_start_backup('" + backupName + "')")
	if err != nil {
		msg = util.SetMessage("ERROR","Entering backup mode using label " + backupName + " for database [" + config.AppPluginParameters["PqDb"] + "] failed! " + err.Error())
		messages = append(messages,msg)
		result = util.SetResult(1,messages)

		return result
	} else {
		msg = util.SetMessage("INFO","Entering backup mode using label " + backupName + " for database [" + config.AppPluginParameters["PqDb"] + "] successful")
		messages = append(messages,msg)
	}

	result = util.SetResult(resultCode, messages)
	return result

}

func (a appPlugin) Unquiesce(config util.Config) util.Result {	

	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	dsn := getDSN(config)
	conn,err := getConn(dsn)

	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't connect to database [" + config.AppPluginParameters["PqDb"] + "] " + err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1,messages)
		return result
	} else {
		defer conn.Close()
		msg := util.SetMessage("INFO", "Connection to database [" + config.AppPluginParameters["PqDb"] + "] established")
		messages = append(messages,msg)
		result = util.SetResult(0,messages)
	}

	msg := util.SetMessage("INFO","Exiting backup mode for database [" + config.AppPluginParameters["PqDb"] + "]")
	messages = append(messages,msg)

	_, err = conn.Exec("SELECT pg_stop_backup()")
	if err != nil {
		msg = util.SetMessage("ERROR","Exiting backup mode for database [" + config.AppPluginParameters["PqDb"] + "] failed! " + err.Error())
		messages = append(messages,msg)
		result = util.SetResult(1,messages)

		return result
	} else {
		msg = util.SetMessage("INFO","Exiting backup mode for database [" + config.AppPluginParameters["PqDb"] + "] successful")
		messages = append(messages,msg)
	}

	result = util.SetResult(resultCode, messages)
	return result

}

func (a appPlugin) PreRestore(config util.Config) util.Result {	

	var result util.Result
	var messages []util.Message

	msg := util.SetMessage("INFO","PreRestore Not implemented")
	messages = append(messages,msg)

	result = util.SetResult(0, messages)
	return result
}	

func (a appPlugin) PostRestore(config util.Config) util.Result {	

	var result util.Result
	var messages []util.Message

	msg := util.SetMessage("INFO","PostRestore Not implemented")
	messages = append(messages,msg)

	result = util.SetResult(0, messages)
	return result
}	

func (a appPlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "postgres"
	plugin.Description = "Postgres plugin for backing up PostgreSQL databases"
	plugin.Version = "1.0.0"
	plugin.Type = "app"

	var capabilities []util.Capability
	var discoverCap util.Capability
	discoverCap.Name = "discover"

	var quiesceCap util.Capability
	quiesceCap.Name = "quiesce"

	var unquiesceCap util.Capability
	unquiesceCap.Name = "unquiesce"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities,discoverCap,quiesceCap,unquiesceCap,infoCap)

	plugin.Capabilities = capabilities
	
	return plugin
}

func checkErr(err error) {
	fmt.Println("error handling")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getDSN(c util.Config) string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	c.AppPluginParameters["PqHost"],c.AppPluginParameters["PqPort"],c.AppPluginParameters["PqUser"],
	c.AppPluginParameters["PqPassword"], c.AppPluginParameters["PqDb"],c.AppPluginParameters["PqSslMode"])
	
	return dsn
}

func getConn(dsn string) (*sql.DB, error) {

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s (%s)\n", err, dsn))
	}

	return conn, nil
}

func main() {}