package main

import (
	"github.com/gorilla/mux"
    "net/http"
    "engine/util"
)
 
type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}
 
type Routes []Route

func NewRouter() *mux.Router {
 
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler
 
        handler = route.HandlerFunc
        handler = util.LogApi(handler, route.Name)
 
        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }
 
    return router
}

var routes = Routes{
    Route{
        "GetStatus",
        "GET",
        "/status",
        GetStatus,
    },
    Route{
        "StartBackupWorkflow",
        "POST",
        "/startBackupWorkflow",
        StartBackupWorkflow,
    },   
    Route{
        "SendTrapSuccessCmd",
        "POST",
        "/sendTrapSuccessCmd",
        SendTrapSuccessCmd,
    },   
    Route{
        "SendTrapErrorCmd",
        "POST",
        "/sendTrapErrorCmd",
        SendTrapErrorCmd,
    },
    Route{
        "GetConfig",
        "GET",
        "/getConfig/{profileName}/{configName}",
        GetConfig,
    },
    Route{
        "GetDefaultConfig",
        "GET",
        "/getDefaultConfig",
        GetDefaultConfig,
    },
    Route{
        "GetDefaultPluginConfig",
        "GET",
        "/getDefaultPluginConfig/{pluginName}",
        GetDefaultPluginConfig,
    },
    Route{
        "GetWorkflowStepResults",
        "GET",
        "/getWorkflowStepResults/{profileName}/{configName}/{workflowId}/{stepId}",
        GetWorkflowStepResults,
    },
    Route{
        "GetWorkflowStatus",
        "GET",
        "/getWorkflowStatus/{profileName}/{configName}/{id}",
        GetWorkflowStatus,
    }, 
    Route{
        "DeleteWorkflowResults",
        "POST",
        "/deleteWorkflowResults/{profileName}/{configName}/{id}",
        DeleteWorkflowResults,
    },                                         
}