package main

import (
	"fmt"
	"log"
	"os"
	"time"
	// "math/rand"
)

func main() {
	logFiles := parseArgs()
	for _, f := range logFiles {
		fmt.Println(f)
		l := f
		go eagerBeaverLogger(&f)
	}
	for {
		time.Sleep(1 * time.Second)
	}
}

func eagerBeaverLogger(logFile *LogFile) {
	log.Println("EagerBeaverLogger has started logging very buzzy in file", logFile.fileName)
	i := 0
	for {
		time.Sleep(WriteLogInterval * time.Second)
		j := 0
		newLine := ""
		for j < LinesPerTime {
			j++
			i++
			timeStamp := time.Now().Format(DateFormatMap[logFile.logType])
			message := fmt.Sprintf("This is log message № %d\n", i)
			newLine = newLine + fmt.Sprintf("%s %s %s", timeStamp, LogDelimer, message)
		}
		file, err := os.OpenFile(logFile.fileName, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Println("eagerBeaverLogger error opening file")
		}
		defer file.Close()
		log.Printf("write more %d log lines to %s \n", LinesPerTime, logFile.fileName)
		if _, err = file.WriteString(newLine); err != nil {
			log.Println("eagerBeaverLogger error writing to file", err)
		}
		file.Sync()
	}
}
