package util

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var logger *log.Logger
var once sync.Once

func GetLoggerInstance() *log.Logger {
	once.Do(func() {
		logger = createLogger()
	})
	return logger
}

func createLogger() *log.Logger {
	mw := io.MultiWriter(os.Stdout)

	return log.New(mw, "", 0)
}

func createMultiLogger() *log.Logger {
	err := CreateDir(os.Getenv("GOBIN")+"/logs", 0755)
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now().Unix()
	logfile := fmt.Sprintf("%s/logs/fossil_%d.conf", os.Getenv("GOBIN"), t)
	file, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	return log.New(mw, "", 0)
}

func LogMessageToConsole(l *log.Logger, message Message) {
	//l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " [" + message.Level + "] ")
	l.SetPrefix(time.Now().Format(time.RFC3339) + " [" + message.Level + "] ")
	l.Print(message.Message)
}

func LogInfoMessage(l *log.Logger, message string) {
	l.SetPrefix(time.Now().Format(time.RFC3339) + "[INFO]")
	l.Print(message)
}

func LogWarnMessage(l *log.Logger, message string) {
	l.SetPrefix(time.Now().Format(time.RFC3339) + "[WARN]")
	l.Print(message)
}

func LogErrorMessage(l *log.Logger, message string) {
	l.SetPrefix(time.Now().Format(time.RFC3339) + "[ERROR]")
	l.Print(message)
}

func LogDebugMessage(l *log.Logger, message string) {
	l.SetPrefix(time.Now().Format(time.RFC3339) + "[DEBUG]")
	l.Print(message)
}

func LogCmdMessage(l *log.Logger, message string) {
	l.SetPrefix(time.Now().Format(time.RFC3339) + "[CMD]")
	l.Print(message)
}

func LogCommentMessage(l *log.Logger, message string) {
	l.SetPrefix("########## ")
	l.Print(message + " ##########")
}

func LogResults(l *log.Logger, result []Result) {
	for _, item := range result {
		for _, line := range item.Messages {
			if line.Level == "INFO" {
				t := time.Unix(line.Timestamp, 0)
				l.SetPrefix(t.String() + " [INFO] ")
				l.Print(line.Message)
			} else if line.Level == "WARN" {
				t := time.Unix(line.Timestamp, 0)
				l.SetPrefix(t.String() + " [WARN] ")
				l.Print(line.Message)
			} else if line.Level == "ERROR" {
				t := time.Unix(line.Timestamp, 0)
				l.SetPrefix(t.String() + " [ERROR] ")
				l.Print(line.Message)
			} else if line.Level == "DEBUG" {
				t := time.Unix(line.Timestamp, 0)
				l.SetPrefix(t.String() + " [DEBUG] ")
				l.Print(line.Message)
			} else if line.Level == "CMD" {
				t := time.Unix(line.Timestamp, 0)
				l.SetPrefix(t.String() + " [CMD] ")
				l.Print(line.Message)
			} else if line.Level == "COMMENT" {
				l.SetPrefix("########## ")
				l.Print(line.Message + " ##########")
			} else {
				t := time.Unix(line.Timestamp, 0)
				l.SetPrefix(t.String() + " [UNKNOWN] ")
				l.Print(line.Message)
			}
		}
	}
}

func LogResult(l *log.Logger, result Result) {
	for _, line := range result.Messages {
		if line.Level == "INFO" {
			t := time.Unix(line.Timestamp, 0)
			l.SetPrefix(t.String() + " [INFO] ")
			l.Print(line.Message)
		} else if line.Level == "WARN" {
			t := time.Unix(line.Timestamp, 0)
			l.SetPrefix(t.String() + " [WARN] ")
			l.Print(line.Message)
		} else if line.Level == "ERROR" {
			t := time.Unix(line.Timestamp, 0)
			l.SetPrefix(t.String() + " [ERROR] ")
			l.Print(line.Message)
		} else if line.Level == "DEBUG" {
			t := time.Unix(line.Timestamp, 0)
			l.SetPrefix(t.String() + " [DEBUG] ")
			l.Print(line.Message)
		} else if line.Level == "CMD" {
			t := time.Unix(line.Timestamp, 0)
			l.SetPrefix(t.String() + " [CMD] ")
			l.Print(line.Message)
		} else if line.Level == "COMMENT" {
			l.SetPrefix("########## ")
			l.Print(line.Message + " ##########")
		} else {
			t := time.Unix(line.Timestamp, 0)
			l.SetPrefix(t.String() + " [UNKNOWN] ")
			l.Print(line.Message)
		}
	}
}

func LogApi(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
