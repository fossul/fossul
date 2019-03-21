package main
 
import (
    "log"
    "net/http"
    "engine/util"
)

const configDir = "configs/"
const dataDir = "data/"
 
func main() {
    util.CreateDir(configDir,0755)
    util.CreateDir(dataDir,0755)

    router := NewRouter()
 
    log.Fatal(http.ListenAndServe(":8000", router))
}