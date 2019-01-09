package main

import (
	"bufio"
	"bytes"
	_ "context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	// _ "github.com/mongodb/mongo-go-driver/mongo"
)

func main() {
	// All fileWatcher goroutines push updates in this channel,
	// master proccess will work with them
	var logsChan chan LogLine
	logsChan = make(chan LogLine)

	for _, logFile := range parseArgs() {
		go watchLog(logFile, logsChan)
	}
	for {
		// TODO: save to persistant storage here
		line := <-logsChan
		line = line
	}
}

func watchLog(logFile LogFile, logsChan chan LogLine) {
	var lastChange time.Time
	for {
		info, err := os.Stat(logFile.fileName)
		if err != nil {
			log.Println("error while gathering file statics")
		}
		thisChange := info.ModTime()
		if thisChange.After(lastChange) {
			lastChange = thisChange
			commitChanges(&logFile, logsChan)
		}
		time.Sleep(CheckInterval * time.Second)
	}
}

func commitChanges(logFile *LogFile, logsChan chan LogLine) {
	oldLinesCount := logFile.linesCount
	newLinesCount, err := getLineCount(logFile)
	linesDiff := oldLinesCount
	if err != nil {
		log.Fatal(err)
	} else {
		linesDiff = (uint64(newLinesCount) - oldLinesCount)
	}
	if linesDiff > 0 {
		log.Printf("[%s]Found %d new lines (%d-%d)\n", logFile.fileName, linesDiff, oldLinesCount, newLinesCount)
		logFile.linesCount = uint64(newLinesCount)

		file, err := os.OpenFile(logFile.fileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		sc := bufio.NewScanner(file)
		var lastLine uint64
		for sc.Scan() {
			lastLine++
			if lastLine-1 > oldLinesCount {
				go parseAndSendLine(sc.Text(), logFile, logsChan)
			}
		}
	} else {
		log.Printf("[%s]lines count doesn't change\n", logFile.fileName)
	}
}

func parseAndSendLine(text string, logFile *LogFile, logsChan chan LogLine) {
	pos := strings.Index(text, LogDelimer)
	if pos == -1 {
		// TODO: We can't just ignore invalid lines, need to do something with them too!
		log.Printf("parse error, invalid format, delimer('%s') not found in line", LogDelimer)
		return
	}
	logTime, err := time.Parse(DateFormatMap[logFile.logType], strings.Trim(text[0:pos], " "))
	if err != nil {
		// TODO: We can't just ignore invalid lines, need to do something with them too!
		log.Println(fmt.Sprintf("parse error, can't parse %s using format %s", text[0:pos], DateFormatMap[logFile.logType]))
	}
	logMsg := text[pos:]
	fileName := logFile.fileName
	logFormat := logFile.logType
	logsChan <- LogLine{logTime, logMsg, fileName, logFormat}
}

func getLineCount(logFile *LogFile) (int, error) {
	file, err := os.OpenFile(logFile.fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	//last line not always end with \n and cause missing one line in result, pain
	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count + 1, nil

		case err != nil:
			return count, err
		}
	}
}
