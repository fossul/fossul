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
