package main

import (
	"os"
	"fmt"
	"github.com/pborman/getopt/v2"
	"engine/util"
	"encoding/json"
)

func main() {
	optAction := getopt.StringLong("action",'a',"","quiesce|unquiesce|list")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	if getopt.IsSet("action") != true {
		fmt.Println("ERROR: incorrect parameter")
		getopt.Usage()
		os.Exit(1)
	}

	if *optAction == "quiesce" {
		quiesce()
	} else if *optAction == "unquiesce" {
		unquiesce()
	} else if *optAction == "info" {
		info()		
	}
}	

func quiesce () {
	fmt.Println("Performing application quiesce")
}

func unquiesce () {
	fmt.Println("Performing application unquiesce")
}

func info () {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
    if err != nil {
        fmt.Println(err)
        return
	}
	
	fmt.Println(string(b))
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