package main
 
import (
    "log"
    "net/http"
)

const configDir = "configs/"
 
func main() {
    router := NewRouter()
 
    log.Fatal(http.ListenAndServe(":8000", router))
}