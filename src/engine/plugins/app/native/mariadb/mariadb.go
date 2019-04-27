package main

import (
	"engine/util"
	"engine/client/k8s"
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

func (a appPlugin) SetEnv(c util.Config) util.Result {
	config = c
	var err error
	var result util.Result
	var messages []util.Message

	if config.AppPluginParameters["MysqlDumpEnable"] != "true" {
		dsn := getDSN(config)

		//reuse database connection
		if conn == nil {
			msg := util.SetMessage("INFO","Creating connection to database [" + config.AppPluginParameters["MysqlDb"] + "]")
			messages = append(messages,msg)
			conn, err = getConn(dsn)
		} else {
			msg := util.SetMessage("INFO","Reusing connection to database [" + config.AppPluginParameters["MysqlDb"] + "]")
			messages = append(messages,msg)
		}

		if conn == nil || err != nil {
			msg := util.SetMessage("ERROR", "Couldn't connect to database [" + config.AppPluginParameters["MysqlDb"] + "] " + err.Error())
			messages = append(messages,msg)

			result = util.SetResult(1,messages)
		} else {
			msg := util.SetMessage("INFO", "Connection to database [" + config.AppPluginParameters["MysqlDb"] + "] established")
			messages = append(messages,msg)
			result = util.SetResult(0,messages)
		}
	}	

	return result
}	

func (a appPlugin) Discover() util.DiscoverResult {
	var discoverResult util.DiscoverResult

	if config.AppPluginParameters["MysqlDumpEnable"] == "true" {
		discoverResult = executeDiscoverDump()
	} else {
		discoverResult = executeDiscover()	
	}

	return discoverResult
}	

func (a appPlugin) Quiesce() util.Result {	
	var result util.Result

	if config.AppPluginParameters["MysqlDumpEnable"] == "true" {
		result = executeDump()
	} else {
		result = executeQuiesce()	
	}

	return result
}

func (a appPlugin) Unquiesce() util.Result {	

	var result util.Result

	if config.AppPluginParameters["MysqlDumpEnable"] == "true" {
		result = executeDumpDelete()
	} else {
		result = executeUnquiesce()
	}

	return result
}

func (a appPlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "mariadb"
	plugin.Description = "MariaDB plugin for backing up MySql or MariaDB"
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

func executeQuiesce() util.Result {

	var result util.Result
	var messages []util.Message
	var resultCode int = 0
	msg := util.SetMessage("INFO","Flushing tables with read lock for database [" + config.AppPluginParameters["MysqlDb"] + "]")
	messages = append(messages,msg)

	_, err := conn.DB.Exec("flush tables with read lock")
	if err != nil {
		msg = util.SetMessage("ERROR","Flushing tables with read lock for database [" + config.AppPluginParameters["MysqlDb"] + "] failed! " + err.Error())
		messages = append(messages,msg)
		result = util.SetResult(1,messages)

		return result
	} else {
		msg = util.SetMessage("INFO","Flushing tables with read lock for database [" + config.AppPluginParameters["MysqlDb"] + "] successful")
		messages = append(messages,msg)
	}

	msg = util.SetMessage("INFO","Flushing logs for database [" + config.AppPluginParameters["MysqlDb"] + "]")
	messages = append(messages,msg)

	_, err = conn.DB.Exec("flush logs")
	if err != nil {
		msg = util.SetMessage("ERROR","Logs flushed for database [" + config.AppPluginParameters["MysqlDb"] + "] failed! " + err.Error())
		messages = append(messages,msg)
		result = util.SetResult(1,messages)

		return result
	} else {
		msg = util.SetMessage("INFO","Flushing logs for database [" + config.AppPluginParameters["MysqlDb"] + "] successful")
		messages = append(messages,msg)
	}

	result = util.SetResult(resultCode, messages)
	return result
}

func executeUnquiesce() util.Result {

	var result util.Result
	var messages []util.Message
	var resultCode int = 0
	msg := util.SetMessage("INFO","Unlocking tables for database [" + config.AppPluginParameters["MysqlDb"] + "]")
	messages = append(messages,msg)

	_, err := conn.DB.Exec("unlock tables")
	if err != nil {	
		msg = util.SetMessage("ERROR","Unlock tables for database[" + config.AppPluginParameters["MysqlDb"] + "] failed! " + err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1,messages)
		return result
	} else {
		msg = util.SetMessage("INFO","Unlock tables for database [" + config.AppPluginParameters["MysqlDb"] + "] successful")
		messages = append(messages,msg)
	}

	conn.DB.Close()
	conn = nil

	result = util.SetResult(resultCode, messages)
	return result
}

func executeDump() util.Result {
	var result util.Result
	var messages []util.Message

	var args []string
	var mkdirArgs []string

	podName,err := k8s.GetPod(config.AppPluginParameters["Namespace"],config.AppPluginParameters["ServiceName"],config.AppPluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1, messages)
		return result
	}

	dumpPath := config.AppPluginParameters["MysqlDumpPath"] + "/" + config.WorkflowId 

	//create tmp directory for storing dump
	mkdirArgs = append(mkdirArgs,"mkdir")
	mkdirArgs = append(mkdirArgs,"-p")
	mkdirArgs = append(mkdirArgs,dumpPath)
	
	cmdResult := k8s.ExecuteCommand(podName,config.AppPluginParameters["ContainerName"],config.AppPluginParameters["Namespace"],config.AppPluginParameters["AccessWithinCluster"],mkdirArgs...)
	
	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages,cmdResult.Messages)
	}

	//execute database dump
	args = append(args,"/opt/rh/rh-mariadb102/root/usr/bin/mysqldump")
	args = append(args,"-h")
	args = append(args,config.AppPluginParameters["MysqlHost"])
	args = append(args,"-P")
	args = append(args,config.AppPluginParameters["MysqlPort"])
	args = append(args,"-u")
	args = append(args,config.AppPluginParameters["MysqlUser"])

	if config.AppPluginParameters["MysqlPassword"] != "" {
		args = append(args,"-p" + config.AppPluginParameters["MysqlPassword"])
	} 	

	args = append(args,config.AppPluginParameters["MysqlDb"])

	args = append(args,"-T")
	args = append(args,dumpPath)

	cmdResult = k8s.ExecuteCommand(podName,config.AppPluginParameters["ContainerName"],config.AppPluginParameters["Namespace"],config.AppPluginParameters["AccessWithinCluster"],args...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages,cmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func executeDumpDelete() util.Result {
	var result util.Result
	var messages []util.Message

	dumpPath := config.AppPluginParameters["MysqlDumpPath"] + "/" + config.WorkflowId 

	var args []string
	args = append(args,"rm")
	args = append(args,"-rf")
	args = append(args,dumpPath)

	podName,err := k8s.GetPod(config.AppPluginParameters["Namespace"],config.AppPluginParameters["ServiceName"],config.AppPluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1, messages)
		return result
	}

	cmdResult := k8s.ExecuteCommand(podName,config.AppPluginParameters["ContainerName"],config.AppPluginParameters["Namespace"],config.AppPluginParameters["AccessWithinCluster"],args...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(cmdResult.Messages,messages)
	}

	result = util.SetResult(0, messages)
	return result

}

func executeDiscover() util.DiscoverResult {
	var discoverResult util.DiscoverResult
	var discoverList []util.Discover
	var discover util.Discover
	var result util.Result
	var messages []util.Message

	discover.Instance = config.AppPluginParameters["MysqlDb"]

	var (
		name string
		value string
	)

	err := conn.DB.QueryRow("show global variables like 'datadir'").Scan(&name,&value)
	if err != nil {
		msg := util.SetMessage("ERROR","Discovery for database [" + config.AppPluginParameters["MysqlDb"] + "] failed! " + err.Error())
		messages = append(messages,msg)
		result = util.SetResult(1,messages)

		discoverResult.Result = result
		return discoverResult
	}
	var dataFilePaths []string
	dataDir := value + config.AppPluginParameters["MysqlDb"]
	dataFilePaths = append(dataFilePaths,dataDir)
	discover.DataFilePaths = dataFilePaths

	msg := util.SetMessage("INFO", "Data Directory is [" + value + config.AppPluginParameters["MysqlDb"] + "]")
	messages = append(messages,msg)

	discoverList = append(discoverList,discover)

	result = util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}

func executeDiscoverDump() util.DiscoverResult {
	var discoverResult util.DiscoverResult
	var discoverList []util.Discover
	var discover util.Discover
	var result util.Result
	var messages []util.Message

	discover.Instance = config.AppPluginParameters["MysqlDb"]

	var dataFilePaths []string
	dumpPath := config.AppPluginParameters["MysqlDumpPath"] + "/" + config.WorkflowId 
	dataFilePaths = append(dataFilePaths,dumpPath)
	discover.DataFilePaths = dataFilePaths

	msg := util.SetMessage("INFO", "Data Directory is [" + config.AppPluginParameters["MysqlDumpPath"] + "]")
	messages = append(messages,msg)

	discoverList = append(discoverList,discover)

	result = util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}

func main() {}