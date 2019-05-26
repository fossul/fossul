package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"fossil/src/engine/util"
	"net/http"
	"strings"
)

// GetStatus godoc
// @Description Status and version information for the service
// @Accept  json
// @Produce  json
// @Success 200 {string} string 
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /status [get]
func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status util.Status
	status.Msg = "OK"
	status.Version = version
	
	json.NewEncoder(w).Encode(status)
}

func SendTrapSuccessCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	if config.SendTrapSuccessCmd != "" {
		args := strings.Split(config.SendTrapSuccessCmd, ",")
		message := util.SetMessage("INFO", "Performing send trap success command [" + config.SendTrapSuccessCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func SendTrapErrorCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	if config.SendTrapSuccessCmd != "" {
		args := strings.Split(config.SendTrapErrorCmd, ",")
		message := util.SetMessage("INFO", "Performing send trap error command [" + config.SendTrapSuccessCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func GetJobs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	jobsDir := dataDir + "/" + profileName + "/" + configName
	
	var result util.Result
	var messages []util.Message

	var jobs util.Jobs
	var jobList []util.Job
	jobList,err := util.ListJobs(jobsDir)

	if err != nil {
		msg := util.SetMessage("ERROR","Job list failed! " + err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		jobs.Result = result

		_ = json.NewDecoder(r.Body).Decode(&jobs)
		json.NewEncoder(w).Encode(jobs)
	}

	jobs.Result.Code = 0
	jobs.Jobs = jobList

	_ = json.NewDecoder(r.Body).Decode(&jobs)
	json.NewEncoder(w).Encode(jobs)
}