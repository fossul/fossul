package util

import (
	"regexp"
)

type Result struct {
	Code int `json:"code,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type ResultSimple struct {
	Code int `json:"code,omitempty"`
	Messages []string `json:"messages,omitempty"`
}

type Message struct {
	Timestamp int64 `json:"time,omitempty"`
	Level string `json:"level,omitempty"`
	Message string `json:"message,omitempty"`
}

func SetMessage(level string, msg string) Message {
	time := GetTimestamp()

	var message Message
	message.Timestamp = time
	message.Level = level
	message.Message = msg

	return message
}

func SetMessages(inputMessages []string) []Message {
	var messages []Message
	for _, msg := range inputMessages {
		re := regexp.MustCompile(`(INFO|WARN|ERROR|DEBUG|CMD)\s+(.*)`)
		match := re.FindStringSubmatch(msg)

		if len(match) != 0 {
			message := SetMessage(match[1],match[2])
			messages = append(messages,message)
		} else {
			if msg != "" {
				message := SetMessage("UNKOWN",msg)
				messages = append(messages,message)			
			}	
		}	
	}

	return messages
}

func SetResultMessage(code int, level, msg string) Result {

	var messages []Message
	var message Message
	message.Level = level
	message.Message = msg
	messages = append(messages,message)

	result := SetResult(code, messages)

	return result
}


func SetResult(code int, messages []Message) Result {
	var result Result
	result.Code = code
	result.Messages = messages

	return result
}

func SetResultSimple(code int, messages []string) ResultSimple {
	var result ResultSimple
	result.Code = code
	result.Messages = messages

	return result
}
