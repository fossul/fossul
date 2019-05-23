package main

import (
    "context"
	"fmt"
	"time"
	"fossil/src/engine/util"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type appPlugin string
var AppPlugin appPlugin

type key string
const (
	hostKey     = key("hostKey")
	usernameKey = key("usernameKey")
	passwordKey = key("passwordKey")
	databaseKey = key("databaseKey")
)

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
		msg := util.SetMessage("ERROR", "Couldn't connect to database [" + config.AppPluginParameters["MongoDb"] + "] " + err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1,messages)
		discoverResult.Result = result
		return discoverResult
	} else {
		msg := util.SetMessage("INFO", "Connection to database [" + config.AppPluginParameters["MongoDb"] + "] established")
		messages = append(messages,msg)
		result = util.SetResult(0,messages)
	}

	discover.Instance = config.AppPluginParameters["MongoDb"]
	serverOptions, err := conn.Database("admin").RunCommand(
		context.Background(),
		bsonx.Doc{{"getCmdLineOpts", bsonx.Int32(1)}},
	).DecodeBytes()
	if err != nil {
		msg := util.SetMessage("ERROR","Discovery for database [" + config.AppPluginParameters["MongoDb"] + "] failed! " + err.Error())
		messages = append(messages,msg)
		result = util.SetResult(1,messages)

		discoverResult.Result = result
		return discoverResult
	}
	
	var serverOptionsResult map[string]interface{}
	bson.Unmarshal(serverOptions, &serverOptionsResult)

	parsed := serverOptionsResult["parsed"].(map[string]interface{})
	storage := parsed["storage"].(map[string]interface{})
	dataDir := fmt.Sprintf("%s",storage["dbPath"])

	var dataFilePaths []string
	dataFilePaths = append(dataFilePaths,dataDir)
	discover.DataFilePaths = dataFilePaths

	msg := util.SetMessage("INFO", "Data Directory is [" + dataDir + "]")
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
		msg := util.SetMessage("ERROR", "Couldn't connect to database [" + config.AppPluginParameters["MongoDb"] + "] " + err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1,messages)
		return result
	} else {
		msg := util.SetMessage("INFO", "Connection to database [" + config.AppPluginParameters["MongoDb"] + "] established")
		messages = append(messages,msg)
		result = util.SetResult(0,messages)
	}

	msg := util.SetMessage("INFO","Flushing writes and locking for database [" + config.AppPluginParameters["MongoDb"] + "]")
	messages = append(messages,msg)

	lockResults, err := conn.Database("admin").RunCommand(
		context.Background(),
		bsonx.Doc{{"fsync",bsonx.Int32(1)},{"lock",bsonx.Boolean(true)}},
	).DecodeBytes()
	if err != nil {
		msg := util.SetMessage("ERROR","Flushing writes and locking for database [" + config.AppPluginParameters["MongoDb"] + "] failed! " + err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1,messages)
		return result
	} else {
		msg := util.SetMessage("INFO","Flushing writes and locking for database [" + config.AppPluginParameters["MongoDb"] + "] successful")
		messages = append(messages,msg)

		msg = util.SetMessage("INFO",fmt.Sprintf("%s",lockResults))
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
		msg := util.SetMessage("ERROR", "Couldn't connect to database [" + config.AppPluginParameters["MongoDb"] + "] " + err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1,messages)
		return result
	} else {
		msg := util.SetMessage("INFO", "Connection to database [" + config.AppPluginParameters["MongoDb"] + "] established")
		messages = append(messages,msg)
		result = util.SetResult(0,messages)
	}

	msg := util.SetMessage("INFO","Unlocking database [" + config.AppPluginParameters["MongoDb"] + "]")
	messages = append(messages,msg)

	unlockResults, err := conn.Database("admin").RunCommand(
		context.Background(),
		bsonx.Doc{{"fsyncUnlock", bsonx.Int32(1)}},
	).DecodeBytes()
	if err != nil {
		msg := util.SetMessage("ERROR","Unlocking database [" + config.AppPluginParameters["MongoDb"] + "] failed! " + err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1,messages)
		return result
	} else {
		msg := util.SetMessage("INFO","Unlocking database [" + config.AppPluginParameters["MongoDb"] + "] successful")
		messages = append(messages,msg)

		msg = util.SetMessage("INFO",fmt.Sprintf("%s",unlockResults))
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
	plugin.Name = "mongo"
	plugin.Description = "Mongo plugin for backing up Mongo databases"
	plugin.Version = "1.0.0"
	plugin.Type = "app"

	var capabilities []util.Capability
	var discoverCap util.Capability
	discoverCap.Name = "discover"

	var quiesceCap util.Capability
	quiesceCap.Name = "quiesce"

	var unquiesceCap util.Capability
	unquiesceCap.Name = "unquiesce"

	var preRestoreCap util.Capability
	preRestoreCap.Name = "preRestore"

	var postRestoreCap util.Capability
	postRestoreCap.Name = "postRestore"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities,discoverCap,quiesceCap,unquiesceCap,preRestoreCap,postRestoreCap,infoCap)

	plugin.Capabilities = capabilities
	
	return plugin
}

func checkErr(err error) {
	fmt.Println("error handling")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getDSN(c util.Config) string { //(context.Context,error) {
	dsn := fmt.Sprintf("mongodb://%s:%s@%s/%s",
	c.AppPluginParameters["MongoUser"],c.AppPluginParameters["MongoPassword"],c.AppPluginParameters["MongoHost"],
	c.AppPluginParameters["MongoDb"])

	return dsn
}

func getConn(dsn string) (*mongo.Client, error) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))

	if err != nil {
		return nil, err
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = conn.Ping(ctx, readpref.Primary())

	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s (%s)\n", err, dsn))
	}

	return conn, nil
}

func main() {}