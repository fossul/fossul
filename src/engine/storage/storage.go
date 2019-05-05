package main
 
import (
    "log"
    "net/http"
    "os"
)

const myUser = "admin"
const myPass = "redhat123"
 
func main() {
    os.Setenv("MyUser", myUser)
    os.Setenv("MyPass", myPass)
    
    router := NewRouter()
 
    log.Fatal(http.ListenAndServe(":8002", router))
}