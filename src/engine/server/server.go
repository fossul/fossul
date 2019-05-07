package main
 
import (
    "log"
    "net/http"
    "fossil/src/engine/util"
)

const configDir = "configs"
const dataDir = "data"
const myUser = "admin"
const myPass = "redhat123"
var runningWorkflowMap map[string]string = make(map[string]string)
 
func main() {
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

    log.Fatal(http.ListenAndServe(":8000", router))
}