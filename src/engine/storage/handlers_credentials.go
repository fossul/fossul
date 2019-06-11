/*
Copyright 2019 The Fossul Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"fossul/src/engine/util"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
)

// AddAwsCredentials godoc
// @Description Add AWS Credentials to storage service
// @Param config body util.AwsCredentials true "aws credentials struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /addAwsCredentials [post]
func AddAwsCredentials(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var messages []util.Message

	awsCredentials, err := util.GetAwsCredentials(w, r)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't get aws credentials! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	homeDir, err := getUserHomeDir()
	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't get user home directory! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	baseDir := homeDir + "/.aws"
	filePath := baseDir + "/credentials"

	util.CreateDir(baseDir, 0755)
	err = writeAwsCredentials(awsCredentials.AwsKey, awsCredentials.AwsSecret, filePath)
	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't get aws credentials! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	msg := util.SetMessage("INFO", "Add AWS credentials completed successfully")
	messages = append(messages, msg)

	result.Code = 0

	result.Messages = messages

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)

}

// DeleteAwsCredentials godoc
// @Description Delete AWS Credentials from storage service
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /deleteAwsCredentials [post]
func DeleteAwsCredentials(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var messages []util.Message

	homeDir, err := getUserHomeDir()
	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't get user home directory! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	baseDir := homeDir + "/.aws"
	filePath := baseDir + "/credentials"

	err = os.Remove(filePath)
	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't delete aws credentials! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	msg := util.SetMessage("INFO", "Delete AWS credentials completed successfully")
	messages = append(messages, msg)

	result.Code = 0

	result.Messages = messages

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)

}

func writeAwsCredentials(key, secret, filePath string) error {
	bytes := []byte("[default]\naws_access_key_id = " + key + "\naws_secret_access_key = " + secret + "\n")

	err := ioutil.WriteFile(filePath, bytes, 0444)
	if err != nil {
		return err
	}

	return nil
}

func getUserHomeDir() (string, error) {
	var homeDir string
	usr, err := user.Current()
	if err != nil {
		return homeDir, err
	}

	homeDir = usr.HomeDir
	return homeDir, nil
}
