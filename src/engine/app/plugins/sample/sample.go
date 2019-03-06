package main

import (
	"os"
	"fmt"
	"github.com/pborman/getopt/v2"
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
	} else if *optAction == "list" {
		list()		
	}
}	

func quiesce () {
	fmt.Println("Performing application quiesce")
}

func unquiesce () {
	fmt.Println("Performing application unquiesce")
}

func list () {
	fmt.Println("quiesce unquiesce")
}