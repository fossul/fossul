package main

import (
	"github.com/gorilla/mux"
	"net/http"
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
        handler = Logger(handler, route.Name)
 
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
        "PreQuiesceCmd",
        "GET",
        "/preQuiesceCmd",
        PreQuiesceCmd,
    },
    Route{
        "QuiesceCmd",
        "GET",
        "/quiesceCmd",
        QuiesceCmd,
    },
    Route{
        "Quiesce",
        "GET",
        "/quiesce",
        Quiesce,
    },
    Route{
        "PostQuiesceCmd",
        "GET",
        "/postQuiesceCmd",
        PostQuiesceCmd,
    },
    Route{
        "PreUnquiesceCmd",
        "GET",
        "/preUnquiesceCmd",
        PreUnquiesceCmd,
    },
    Route{
        "UnquiesceCmd",
        "GET",
        "/unquiesceCmd",
        UnquiesceCmd,
    },
    Route{
        "Unquiesce",
        "GET",
        "/unquiesce",
        Unquiesce,
    },
    Route{
        "PostUnquiesceCmd",
        "GET",
        "/postUnquiesceCmd",
        PostUnquiesceCmd,
    },    
}