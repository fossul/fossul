package main
 
import (
    "log"
    "net/http"
    "os"
    "fossil/src/engine/util"
)

var port string = os.Getenv("FOSSIL_STORAGE_SERVICE_PORT")
var pluginDir string = os.Getenv("FOSSIL_STORAGE_PLUGIN_DIR")
var myUser string = os.Getenv("FOSSIL_USERNAME")
var myPass string = os.Getenv("FOSSIL_PASSWORD")
 
func main() {
    log.Println("Plugin directory [" + pluginDir + "]")
    err := util.CreateDir(pluginDir,0755)
    if err != nil {
        log.Fatal(err)
    }

    os.Setenv("MyUser", myUser)
    os.Setenv("MyPass", myPass)
    
    router := NewRouter()
 
    log.Println("Starting storage service on port [" + port + "]")
    log.Fatal(http.ListenAndServe(":" + port, router))
}