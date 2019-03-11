package main

import (
	"os"
	"fmt"
	"github.com/pborman/getopt/v2"
	"engine/util"
	"encoding/json"
)

func main() {
	optAction := getopt.StringLong("action",'a',"","quiesce|unquiesce|info")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	if getopt.IsSet("action") != true {
		fmt.Printf("ERROR incorrect parameter\n")
		getopt.Usage()
		os.Exit(1)
	}

		//load env parameters
		configMap := getEnvParams()

	if *optAction == "quiesce" {
		printEnv(configMap)
		quiesce(configMap)
	} else if *optAction == "unquiesce" {
		printEnv(configMap)
		unquiesce(configMap)
	} else if *optAction == "info" {
		info()		
	} else {
		fmt.Printf("ERROR incorrect parameter" + *optAction + "\n")
		getopt.Usage()
		os.Exit(1)
	}
}	

func quiesce (configMap map[string]string) {
	printEnv(configMap)
	fmt.Printf("INFO Performing application quiesce")
}

func unquiesce (configMap map[string]string) {
	printEnv(configMap)
	fmt.Printf("INFO Performing application unquiesce")
}

func info () {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
    if err != nil {
        fmt.Println(err)
        return
	}
	
	fmt.Printf(string(b)+ "\n")
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample"
	plugin.Description = "A sample plugin"
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

func printEnv(configMap map[string]string) {
	fmt.Printf("DEBUG Config Parameters" + "\n")
	fmt.Printf("DEBUG SampleAppVar1=" + configMap["SampleAppVar1"] + "\n")
	fmt.Printf("DEBUG SampleAppVar2=" + configMap["SampleAppVar1"] + "\n")
}

func getEnvParams() map[string]string {
	configMap := map[string]string{}

	configMap["SampleAppVar1"] = os.Getenv("SampleAppVar1")
	configMap["SampleAppVar2"] = os.Getenv("SampleAppVar2")

	return configMap
}