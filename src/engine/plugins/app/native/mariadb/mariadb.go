package main

import (
	"database/sql"
	"errors"
	"fmt"
	"fossul/src/engine/util"
	_ "github.com/go-sql-driver/mysql"
)

type appPlugin string

var conn *MySQL
var AppPlugin appPlugin

type MySQL struct {
	DB *sql.DB
}

func (a appPlugin) SetEnv(config util.Config) util.Result {
	var err error
	var result util.Result
	var messages []util.Message

	dsn := getDSN(config)

	//reuse database connection
	if conn == nil {
		msg := util.SetMessage("INFO", "Creating connection to database ["+config.AppPluginParameters["MysqlDb"]+"]")
		messages = append(messages, msg)
		conn, err = getConn(dsn)
	} else {
		msg := util.SetMessage("INFO", "Reusing connection to database ["+config.AppPluginParameters["MysqlDb"]+"]")
		messages = append(messages, msg)
	}

	if conn == nil || err != nil {
		msg := util.SetMessage("ERROR", "Couldn't connect to database ["+config.AppPluginParameters["MysqlDb"]+"] "+err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
	} else {
		msg := util.SetMessage("INFO", "Connection to database ["+config.AppPluginParameters["MysqlDb"]+"] established")
		messages = append(messages, msg)
		result = util.SetResult(0, messages)
	}

	return result
}

func (a appPlugin) Discover(config util.Config) util.DiscoverResult {
	var discoverResult util.DiscoverResult
	var discoverList []util.Discover
	var discover util.Discover
	var result util.Result
	var messages []util.Message

	discover.Instance = config.AppPluginParameters["MysqlDb"]

	var (
		name  string
		value string
	)

	err := conn.DB.QueryRow("show global variables like 'datadir'").Scan(&name, &value)
	if err != nil {
		msg := util.SetMessage("ERROR", "Discovery for database ["+config.AppPluginParameters["MysqlDb"]+"] failed! "+err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)

		discoverResult.Result = result
		return discoverResult
	}
	var dataFilePaths []string
	dataDir := value + config.AppPluginParameters["MysqlDb"]
	dataFilePaths = append(dataFilePaths, dataDir)
	discover.DataFilePaths = dataFilePaths

	msg := util.SetMessage("INFO", "Data Directory is ["+value+config.AppPluginParameters["MysqlDb"]+"]")
	messages = append(messages, msg)

	discoverList = append(discoverList, discover)

	result = util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}

func (a appPlugin) Quiesce(config util.Config) util.Result {

	var result util.Result
	var messages []util.Message
	var resultCode int = 0
	msg := util.SetMessage("INFO", "Flushing tables with read lock for database ["+config.AppPluginParameters["MysqlDb"]+"]")
	messages = append(messages, msg)

	_, err := conn.DB.Exec("flush tables with read lock")
	if err != nil {
		msg = util.SetMessage("ERROR", "Flushing tables with read lock for database ["+config.AppPluginParameters["MysqlDb"]+"] failed! "+err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)

		return result
	} else {
		msg = util.SetMessage("INFO", "Flushing tables with read lock for database ["+config.AppPluginParameters["MysqlDb"]+"] successful")
		messages = append(messages, msg)
	}

	msg = util.SetMessage("INFO", "Flushing logs for database ["+config.AppPluginParameters["MysqlDb"]+"]")
	messages = append(messages, msg)

	_, err = conn.DB.Exec("flush logs")
	if err != nil {
		msg = util.SetMessage("ERROR", "Logs flushed for database ["+config.AppPluginParameters["MysqlDb"]+"] failed! "+err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)

		return result
	} else {
		msg = util.SetMessage("INFO", "Flushing logs for database ["+config.AppPluginParameters["MysqlDb"]+"] successful")
		messages = append(messages, msg)
	}

	result = util.SetResult(resultCode, messages)
	return result

}

func (a appPlugin) Unquiesce(config util.Config) util.Result {

	var result util.Result
	var messages []util.Message
	var resultCode int = 0
	msg := util.SetMessage("INFO", "Unlocking tables for database ["+config.AppPluginParameters["MysqlDb"]+"]")
	messages = append(messages, msg)

	_, err := conn.DB.Exec("unlock tables")
	if err != nil {
		msg = util.SetMessage("ERROR", "Unlock tables for database["+config.AppPluginParameters["MysqlDb"]+"] failed! "+err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	} else {
		msg = util.SetMessage("INFO", "Unlock tables for database ["+config.AppPluginParameters["MysqlDb"]+"] successful")
		messages = append(messages, msg)
	}

	conn.DB.Close()
	conn = nil

	result = util.SetResult(resultCode, messages)
	return result

}

func (a appPlugin) PreRestore(config util.Config) util.Result {

	var result util.Result
	var messages []util.Message

	msg := util.SetMessage("INFO", "PreRestore Not implemented")
	messages = append(messages, msg)

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) PostRestore(config util.Config) util.Result {

	var result util.Result
	var messages []util.Message

	msg := util.SetMessage("INFO", "PostRestore Not implemented")
	messages = append(messages, msg)

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "mariadb"
	plugin.Description = "MariaDB plugin for backing up MySql or MariaDB"
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

	capabilities = append(capabilities, discoverCap, quiesceCap, unquiesceCap, infoCap)

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
	var dsn string

	if c.AppPluginParameters["MysqlPassword"] == "" {
		dsn = c.AppPluginParameters["MysqlUser"] + "@" + c.AppPluginParameters["MysqlProto"] + "(" +
			c.AppPluginParameters["MysqlHost"] + ":" + c.AppPluginParameters["MysqlPort"] + ")/" +
			c.AppPluginParameters["MYSQL_DB"]
	} else {
		dsn = c.AppPluginParameters["MysqlUser"] + ":" + c.AppPluginParameters["MysqlPassword"] + "@" +
			c.AppPluginParameters["MysqlProto"] + "(" + c.AppPluginParameters["MysqlHost"] + ":" +
			c.AppPluginParameters["MysqlPort"] + ")/" + c.AppPluginParameters["MysqlDb"]
	}

	return dsn
}

func getConn(dsn string) (*MySQL, error) {

	var (
		m   = new(MySQL)
		err error
	)

	m.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = m.DB.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s (%s)\n", err, dsn))
	}

	return m, nil
}

func main() {}
