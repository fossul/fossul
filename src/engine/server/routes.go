package main

import (
	"github.com/gorilla/mux"
    "net/http"
    "fossul/src/engine/util"
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
 
        handler = basicAuth(handler)

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }
 
    return router
}

func basicAuth(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user, pass, _ := r.BasicAuth()

        if myUser != user || myPass != pass {
            w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
            http.Error(w, "Unauthorized.", http.StatusUnauthorized)
            return
        }

        h.ServeHTTP(w, r)
    })
}

var routes = Routes{
    Route{
        "GetStatus",
        "GET",
        "/status",
        GetStatus,
    },
    Route{
        "StartBackupWorkflowLocalConfig",
        "POST",
        "/startBackupWorkflowLocalConfig",
        StartBackupWorkflowLocalConfig,
    },       
    Route{
        "StartBackupWorkflow",
        "POST",
        "/startBackupWorkflow/{profileName}/{configName}/{policy}",
        StartBackupWorkflow,
    },
    Route{
        "StartRestoreWorkflowLocalConfig",
        "POST",
        "/startRestoreWorkflowLocalConfig/{workflowId}",
        StartRestoreWorkflowLocalConfig,
    },       
    Route{
        "StartRestoreWorkflow",
        "POST",
        "/startRestoreWorkflow/{profileName}/{configName}/{policy}/{workflowId}",
        StartRestoreWorkflow,
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
        "GetPluginConfig",
        "GET",
        "/getPluginConfig/{profileName}/{configName}/{pluginName}",
        GetPluginConfig,
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
        "GET",
        "/deleteWorkflowResults/{profileName}/{configName}/{id}",
        DeleteWorkflowResults,
    },
    Route{
        "GetJobs",
        "GET",
        "/getJobs/{profileName}/{configName}",
        GetJobs,
    },
    Route{
        "AddConfig",
        "POST",
        "/addConfig/{profileName}/{configName}",
        AddConfig,
    },
    Route{
        "AddPluginConfig",
        "POST",
        "/addPluginConfig/{profileName}/{configName}/{pluginName}",
        AddPluginConfig,
    },    
    Route{
        "DeleteConfig",
        "GET",
        "/deleteConfig/{profileName}/{configName}",
        DeleteConfig,
    },
    Route{
        "DeleteConfigDir",
        "GET",
        "/deleteConfigDir/{profileName}/{configName}",
        DeleteConfigDir,
    },    
    Route{
        "DeletePluginConfig",
        "GET",
        "/deletePluginConfig/{profileName}/{configName}/{pluginName}",
        DeletePluginConfig,
    },    
    Route{
        "AddProfile",
        "GET",
        "/addProfile/{profileName}",
        AddProfile,
    },          
    Route{
        "DeleteProfile",
        "GET",
        "/deleteProfile/{profileName}",
        DeleteProfile,
    },
    Route{
        "ListProfiles",
        "GET",
        "/listProfiles",
        ListProfiles,
    },
    Route{
        "ListConfigs",
        "GET",
        "/listConfigs/{profileName}",
        ListConfigs,
    },
    Route{
        "ListPluginConfigs",
        "GET",
        "/listPluginConfigs/{profileName}/{configName}",
        ListPluginConfigs,
    },
    Route{
        "AddSchedule",
        "POST",
        "/addSchedule/{profileName}/{configName}/{policy}",
        AddSchedule,
    },
    Route{
        "DeleteSchedule",
        "GET",
        "/deleteSchedule/{profileName}/{configName}/{policy}",
        DeleteSchedule,
    },
    Route{
        "ListSchedules",
        "GET",
        "/listSchedules",
        ListSchedules,
    },                                                                                             
}