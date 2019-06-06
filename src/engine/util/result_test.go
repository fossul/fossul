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
package util

import (
	"log"
	"testing"
)

func TestResult(t *testing.T) {
	message1 := SetMessage("INFO", "Testing1")

	var messages []string
	message2 := "INFO Testing2"

	messages = append(messages, message2)

	messageList := SetMessages(messages)
	messageList = append(messageList, message1)

	result := SetResult(0, messageList)

	log.Println(result)

	if result.Code != 0 {
		t.Fail()
	}

	if len(result.Messages) != 2 {
		t.Fail()
	}

}
