/**This file contains constants and global declarations.
**/

package apis

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"time"

	"log"
)

var (
	//DB network variables
	pool                 *sql.DB
	queryTimeOutDuration = time.Second * 5
	//logging constants
	logger = log.Default()
)

//DB constants
const (
	//Table names
	//columns: name, email, phone, about, profile_picture
	//Constraint: Unique email
	doctorTbl = "doctors"
)

const (
	logIDFlag = "|API-HANDLER|"
)

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
