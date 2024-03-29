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
	"regexp"
)

type Result struct {
	Code     int       `json:"code,omitempty"`
	Messages []Message `json:"messages,omitempty"`
	Data     []string  `json:"data,omitempty"`
	Backup   Backup    `json:"backup,omitempty"`
}

type ResultSimple struct {
	Code     int      `json:"code,omitempty"`
	Messages []string `json:"messages,omitempty"`
	Data     []string `json:"data,omitempty"`
}

type Message struct {
	Timestamp int64  `json:"time,omitempty"`
	Level     string `json:"level,omitempty"`
	Message   string `json:"message,omitempty"`
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
			message := SetMessage(match[1], match[2])
			messages = append(messages, message)
		} else {
			if msg != "" {
				message := SetMessage("UNKOWN", msg)
				messages = append(messages, message)
			}
		}
	}

	return messages
}

func SetResultMessage(code int, level, msg string) Result {
	time := GetTimestamp()

	var messages []Message
	var message Message
	message.Level = level
	message.Timestamp = time
	message.Message = msg
	messages = append(messages, message)

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

func PrependMessage(message Message, messages []Message) []Message {
	var prependedMessages []Message
	prependedMessages = append(prependedMessages, message)

	for _, msg := range messages {
		prependedMessages = append(prependedMessages, msg)
	}

	return prependedMessages
}

func PrependMessages(prependedMessages, postpendedMessages []Message) []Message {
	var messages []Message

	for _, msg := range prependedMessages {
		messages = append(messages, msg)
	}

	for _, msg := range postpendedMessages {
		messages = append(messages, msg)
	}

	return messages
}
