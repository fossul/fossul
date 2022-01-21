/*
Copyright 2019 The Fossul Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/fossul/fossul/src/engine/server/docs"
	"github.com/fossul/fossul/src/engine/util"
	httpSwagger "github.com/swaggo/http-swagger"
)

const version = "1.0.0"

var port string = os.Getenv("FOSSUL_SERVER_SERVICE_PORT")
var configDir string = os.Getenv("FOSSUL_SERVER_CONFIG_DIR")
var dataDir string = os.Getenv("FOSSUL_SERVER_DATA_DIR")
var myUser string = os.Getenv("FOSSUL_USERNAME")
var myPass string = os.Getenv("FOSSUL_PASSWORD")
var serverHostname string = os.Getenv("FOSSUL_SERVER_CLIENT_HOSTNAME")
var serverPort string = os.Getenv("FOSSUL_SERVER_CLIENT_PORT")
var appHostname string = os.Getenv("FOSSUL_APP_CLIENT_HOSTNAME")
var appPort string = os.Getenv("FOSSUL_APP_CLIENT_PORT")
var storageHostname string = os.Getenv("FOSSUL_STORAGE_CLIENT_HOSTNAME")
var storagePort string = os.Getenv("FOSSUL_STORAGE_CLIENT_PORT")
var debug string = os.Getenv("FOSSUL_SERVER_DEBUG")

var runningWorkflowMap sync.Map

//var runningWorkflowMapMutex = sync.RWMutex{}
//var runningWorkflowMap map[string]string = make(map[string]string)
// @title Fossul Framework Server API
// @version 1.0
// @description APIs for managing Fossul workflows, jobs, profile, and configs
// @description JSON API definition can be retrieved at <a href="/api/v1/swagger/doc.json">/api/v1/swagger/doc.json</a>
// @termsOfService http://swagger.io/terms/

// @contact.name Keith Tenzer
// @contact.url http://www.keithtenzer.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {

	log.Println("Configs directory [" + configDir + "] Data directory [" + dataDir + "]")
	err := util.CreateDir(configDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = util.CreateDir(dataDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	router := NewRouter()
	router.PathPrefix("/api/v1").Handler(httpSwagger.WrapHandler)

	StartCron()
	err = LoadCronSchedules()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// Catch Terminal (ctr-c) and SigTerm
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// Signal recieved, shutdown
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println("Storage service shutdown failed! %v", err)
		}

		log.Println("Stopping server service on port [" + port + "]")
		close(idleConnsClosed)
	}()

	log.Println("Starting server service on port [" + port + "]")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Println("Server service start failed! %v", err)
	}

	<-idleConnsClosed
}

func printConfigDebug(config util.Config) {
	if debug == "true" {
		log.Println("[DEBUG]", config)
	}
}

func printConfigMapDebug(configMap map[string]string) {
	if debug == "true" {
		log.Println("[DEBUG]", configMap)
	}
}
