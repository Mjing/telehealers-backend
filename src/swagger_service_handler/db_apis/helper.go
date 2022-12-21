package apis

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
)

func newQueryError(errMsg string) error {
	return errors.New(queryErrorTag + errMsg)
}

func SetupLogFile(fileName string) *log.Logger {
	file, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[Error]|%v|In openning log file:%v", logIDFlag, err)
		log.Print("[INFO]Directing logs to stdout")
	} else {
		log.Printf("Setting log to %v", fileName)
		logger.SetOutput(file)
	}
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
	logger.SetPrefix(logIDFlag)
	return logger
}

func getTimeOutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.TODO(), queryTimeOutDuration)
}

func duplicateEntryError(err error) bool {
	return strings.Contains(err.Error(), "uplicate")
}

func SetDataRootDir(newRootPath string) {
	fileStoreRoot = newRootPath
}
