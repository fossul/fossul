package main

import (
	"engine/util"
	"engine/plugins/pluginUtil"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"errors"
)

type appPlugin string

var conn *MySQL
var config util.Config
var AppPlugin appPlugin

type MySQL struct {
		DB *sql.DB
}

func (a appPlugin) SetEnv(c util.Config) {
	config = c
	var err error

	dsn := getDSN(config)
	fmt.Println("HERE12121",dsn)
	conn, err = getConn(dsn)
	checkErr(err)
}	

func (a appPlugin) Quiesce() util.Result {	

	var result util.Result
	var messages []util.Message
	var resultCode int = 0
	msg := util.SetMessage("INFO","Flushing tables with read lock for database [" + config.AppPluginParameters["MYSQL_DB"] + "]")
	messages = append(messages,msg)

	_, err := conn.DB.Exec("flush tables with read lock")
	if err != nil {
		msg = util.SetMessage("ERROR","Flushing tables with read lock for database [" + config.AppPluginParameters["MYSQL_DB"] + "] failed! " + err.Error())
		messages = append(messages,msg)
		resultCode = 1
	} else {
		msg = util.SetMessage("INFO","Flushing tables with read lock for database [" + config.AppPluginParameters["MYSQL_DB"] + "] successful")
		messages = append(messages,msg)
	}

	msg = util.SetMessage("INFO","Flushing logs for database [" + config.AppPluginParameters["MYSQL_DB"] + "]")
	messages = append(messages,msg)

	_, err = conn.DB.Exec("flush logs")
	if err != nil {
		msg = util.SetMessage("ERROR","Logs flushed for database [" + config.AppPluginParameters["MYSQL_DB"] + "] failed! " + err.Error())
		messages = append(messages,msg)
		resultCode = 1
	} else {
		msg = util.SetMessage("INFO","Flushing logs for database [" + config.AppPluginParameters["MYSQL_DB"] + "] successful")
		messages = append(messages,msg)
	}

	result = util.SetResult(resultCode, messages)
	return result

}

func (a appPlugin) Unquiesce() util.Result {	

	var result util.Result
	var messages []util.Message
	var resultCode int = 0
	msg := util.SetMessage("INFO","Unlocking tables for database [" + config.AppPluginParameters["MYSQL_DB"] + "]")
	messages = append(messages,msg)

	_, err := conn.DB.Exec("unlock tables")
	if err != nil {	
		msg = util.SetMessage("ERROR","Unlock tables for database[" + config.AppPluginParameters["MYSQL_DB"] + "] failed! " + err.Error())
		messages = append(messages,msg)
	} else {
		msg = util.SetMessage("INFO","Unlock tables for database [" + config.AppPluginParameters["MYSQL_DB"] + "] successful")
		messages = append(messages,msg)
		resultCode = 1
	}

	result = util.SetResult(resultCode, messages)
	return result

}

func (a appPlugin) Info() {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
    if err != nil {
        pluginUtil.LogErrorMessage(err.Error())
	} else {
		pluginUtil.PrintMessage(string(b))
	}
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "mariadb"
	plugin.Description = "MariaDB plugin for backing up MySql or MariaDB"
	plugin.Type = "app"

	var capabilities []util.Capability
	var quiesceCap util.Capability
	quiesceCap.Name = "quiesce"

	var unquiesceCap util.Capability
	unquiesceCap.Name = "unquiesce"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities,quiesceCap,unquiesceCap,infoCap)

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
	
	if c.AppPluginParameters["MYSQL_PASSWORD"] == "" {
		dsn = c.AppPluginParameters["MYSQL_USER"] + "@" + c.AppPluginParameters["MYSQL_PROTO"] + "(" + 
		c.AppPluginParameters["MYSQL_HOST"] + ":" + c.AppPluginParameters["MYSQL_PORT"] + ")/" + 
		c.AppPluginParameters["MYSQL_DB"]
	} else {
		dsn = c.AppPluginParameters["MYSQL_USER"] + ":" + c.AppPluginParameters["MYSQL_PASSWORD"] + "@" + 
		c.AppPluginParameters["MYSQL_PROTO"] + "(" + c.AppPluginParameters["MYSQL_HOST"] + ":" + 
		c.AppPluginParameters["MYSQL_PORT"] + ")/" + c.AppPluginParameters["MYSQL_DB"]
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