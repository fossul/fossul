package util

import (
	"log"
	"os"
	"time"
	"io"
)

func logMessageToConsole (l *log.Logger, message Message) {
	l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " [" + message.Level + "] ")
    l.Print(message.Message)
}

func logMessagesToConsole (l *log.Logger, messages []Message) {
	for _, message := range messages {
		l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " [" + message.Level + "] ")
		l.Print(message.Message)
	}	
}

func logMessage (l *log.Logger, message Message) {
	l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " [" + message.Level + "] ")

	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

	mw := io.MultiWriter(os.Stdout, file)
	l.SetOutput(mw)

    l.Print(message.Message)
}

func logMessages (l *log.Logger, messages []Message) {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }

    defer file.Close()

	mw := io.MultiWriter(os.Stdout, file)
	l.SetOutput(mw)

	for _, message := range messages {
		l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " [" + message.Level + "] ")
		l.Print(message.Message)
	}	
}