package main
 
import (
    "log"
    "net/http"
    "os"
    "fossil/src/engine/util"
)

const pluginDir = "plugins"
const myUser = "admin"
const myPass = "redhat123"
 
func main() {
    err := util.CreateDir(pluginDir,0755)
    if err != nil {
        log.Fatal(err)
    }

    os.Setenv("MyUser", myUser)
    os.Setenv("MyPass", myPass)
    
    router := NewRouter()
 
    log.Fatal(http.ListenAndServe(":8002", router))
}