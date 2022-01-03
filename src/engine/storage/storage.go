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
	"syscall"

	_ "github.com/fossul/fossul/src/engine/storage/docs"
	"github.com/fossul/fossul/src/engine/util"
	httpSwagger "github.com/swaggo/http-swagger"
)

const version = "1.0.0"

var port string = os.Getenv("FOSSUL_STORAGE_SERVICE_PORT")
var pluginDir string = os.Getenv("FOSSUL_STORAGE_PLUGIN_DIR")
var myUser string = os.Getenv("FOSSUL_USERNAME")
var myPass string = os.Getenv("FOSSUL_PASSWORD")
var debug string = os.Getenv("FOSSUL_STORAGE_DEBUG")

// @title Fossul Framework Storage API
// @version 1.0
// @description APIs for managing Fossul storage plugins
// @description JSON API definition can be retrieved at <a href="/api/v1/swagger/doc.json">/api/v1/swagger/doc.json</a>
// @termsOfService http://swagger.io/terms/

// @contact.name Keith Tenzer
// @contact.url http://www.keithtenzer.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {
	log.Println("Plugin directory [" + pluginDir + "]")
	err := util.CreateDir(pluginDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("MyUser", myUser)
	os.Setenv("MyPass", myPass)

	router := NewRouter()
	router.PathPrefix("/api/v1").Handler(httpSwagger.WrapHandler)

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

		log.Println("Stopping storage service on port [" + port + "]")
		close(idleConnsClosed)
	}()

	log.Println("Starting storage service on port [" + port + "]")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Println("App service start failed! %v", err)
	}

	<-idleConnsClosed
}

func printConfigDebug(config util.Config) {
	if debug == "true" {
		log.Println("[DEBUG]", config)
	}
}
