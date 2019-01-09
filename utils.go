package main

import (
	"log"
	"os"
	"strings"
	"time"
)

// LogFile conrains information of log file to watch for, log type and etc
type LogFile struct {
	fileName   string
	logType    string
	linesCount uint64
}

// LogLine contains all data for storing in persistant storage
type LogLine struct {
	Time      time.Time `json:"log_time"`
	Text      string    `json:"log_msg"`
	FileName  string    `json:"file_name"`
	LogFormat string    `json:"log_format"`
}

const (
	// WriteLogInterval represent interval of adding lines to log
	WriteLogInterval = 2
	// LinesPerTime number of lines adding per 1 time of writing to log
	LinesPerTime = 500
	// LogDelimer string in between log time stamp and log text
	LogDelimer = " | "
	// CheckInterval interval in seconds between checking log file
	CheckInterval = 5
)

// DateFormatMap key is log type, value is date layout ready for parsing
var DateFormatMap = map[string]string{
	"first_format":  "Jan _2, 2006 at 3:04:05pm (MST)",
	"second_format": "2006-01-02T15:04:05Z",
}

func parseArgs() (logs []LogFile) {
	args := os.Args[1:]
	if len(args)%2 != 0 {
		log.Fatal("invalid input data!")
		return
	}
	for i := 0; i < len(args)/2; i++ {
		logName := args[i*2]
		logType := strings.ToLower(args[i*2+1])
		log.Printf("%s type %s \n", logName, logType)
		if _, ok := DateFormatMap[logType]; ok {
			logs = append(logs, LogFile{logName, logType, 0})
		} else {
			log.Fatalf("format of file %s does not corresponding any of allowed formats! (%s)", logName, logType)
		}
	}
	return
}
