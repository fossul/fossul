package util

import (
	"encoding/json"
	"net/http"
	"log"
)

func GetConfig (w http.ResponseWriter, r *http.Request) Config {

	var config Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
 
	res,err := json.Marshal(&config)
	if err != nil {
        log.Println(err)
	}

	log.Println("DEBUG", string(res))

	return config
}