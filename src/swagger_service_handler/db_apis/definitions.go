/**This file contains constants and global declarations.
**/

package apis

import (
	"database/sql"
	"os"
	"time"

	"log"

	"github.com/go-openapi/runtime/middleware"
)

var (
	dbName = os.Getenv("DB_NAME")
	dbUser = os.Getenv("DB_USER")
	dbPass = os.Getenv("DB_PASS")
	dbAddr = os.Getenv("DB_ADDR")
	//DB network variables
	pool                 *sql.DB
	queryTimeOutDuration = time.Second * 5
	//logging constants
	logger = log.Default()
)

// DB constants
const (
	//General update query
	generalInsertQuery = "INSERT INTO %v (%v) VALUES " //...followed by values in ()
	generalUpdateQuery = "UPDATE %v SET %v WHERE %v"
	generalDeleteQuery = "DELETE FROM %v WHERE %v"
	generalFetchQuery  = "SELECT %v FROM %v LIMIT %v, %v" //In case of conditions extend 2 position parameter

	//Table names
	//columns: name, email, phone, about, profile_picture, registration_number, password
	//Constraint: Unique email
	doctorTbl = "doctors"
	//columns:id,name,email,registration_number,approved,comments,reviewer_comments
	doctorRegistrationApplicationTbl = "doctor_admission_applications"
	//columns: name, email, phone, profile_picture, password, profile_picture, about
	patientTbl = "patients"
	//columns: date, requested_start_time, requested_end_time, start_time,
	//end_time, doctor_id, patient_id, prescription_id,patient_health_info_id
	aptTbl          = "appointments"
	aptFetchColumns = "id, doctor_id, patient_id, patient_health_info_id, prescription_id, date, start_time, end_time, requested_start_time, requested_end_time"
	//columns:gender,height,weight,bp,health_complaints,patient_id,created_on
	patienHealthInfoTbl = "patient_health_info"
)

// General Entities
// Columns: name, description
var (
	medicineTbl = "medicines"
	testTbl     = "med_tests"
	adviceTbl   = "advices"
)

const (
	logIDFlag = "|API-HANDLER|"
)

// Struct to work with github.com/go-sql-driver/mysql sql function
// attributes should be used together
type sqlExeParams struct {
	Query     string //MySQL query with positional parameters
	QueryArgs []any  //Value of positional parameters
}

type UpdateAPIs interface {
	makeQuery() (sqlExeParams, error)
	//code 2xx not allowed for error response
	errResponse(httpStatusCode int, err error) middleware.Responder
	//Only status code 200 will be responded
	okResponse(lastId, rowAffected int64) middleware.Responder
}
type ReadAPIs interface {
	UpdateAPIs
	//In okResponse, inputs will be ignored
	scanRows(*sql.Rows) error
}
