package main
 
import (
    "log"
    "net/http"
    "fossil/src/engine/util"
)

const configDir = "configs/"
const dataDir = "data/"
const myUser = "admin"
const myPass = "redhat123"
var runningWorkflowMap map[string]string = make(map[string]string)
 
func main() {
    util.CreateDir(configDir,0755)
    util.CreateDir(dataDir,0755)

    router := NewRouter()
 
    log.Fatal(http.ListenAndServe(":8000", router))
}