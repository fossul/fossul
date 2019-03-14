package pluginUtil

import (
    "fmt"
    "engine/util"
)

func LogInfoMessage (message string) {
    fmt.Println("INFO " + message)
}

func LogWarnMessage (message string) {
    fmt.Println("WARN " + message)
}

func LogErrorMessage (message string) {
    fmt.Println("ERROR " + message)
}

func LogDebugMessage (message string) {
    fmt.Println("DEBUG " + message)
}
func LogCmdMessage (message string) {
    fmt.Println("CMD " + message)
}

func LogCommentMessage (message string) {
    fmt.Println("########## " + message + " ##########")
}

func PrintMessage (message string) {
    fmt.Println(message)
}

func LogResultMessages (result util.Result) {
	for _, line := range result.Messages {
		fmt.Println(line.Level, line.Message)
	}
}