package main
 
import (
    "log"
    "net/http"
    "fossil/src/engine/util"
    "os"
)

var port string = os.Getenv("FOSSIL_SERVER_PORT")
var configDir string = os.Getenv("FOSSIL_SERVER_CONFIG_DIR")
var dataDir string = os.Getenv("FOSSIL_SERVER_DATA_DIR")
var myUser string = os.Getenv("FOSSIL_USERNAME")
var myPass string = os.Getenv("FOSSIL_PASSWORD")

var runningWorkflowMap map[string]string = make(map[string]string)
 
func main() {
    log.Println("Configs directory [" + configDir + "] Data directory [" + dataDir + "]")
    err := util.CreateDir(configDir,0755)
    if err != nil {
        log.Fatal(err)
    }

    err = util.CreateDir(dataDir,0755)
    if err != nil {
        log.Fatal(err)
    }

    router := NewRouter()

    StartCron()
    err = LoadCronSchedules()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Starting server service on port [" + port + "]")
    //log.Fatal(http.ListenAndServe(":" + port, router))
    log.Fatal(http.ListenAndServe(":8000", router))
}