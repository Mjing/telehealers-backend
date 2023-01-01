package apis

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/patient"
)

var (
	insertPatQuery = "INSERT INTO " + patientTbl + " (%v) VALUES (%v)"
	updatePatQuery = "UPDATE " + patientTbl + " SET %v WHERE %v"
	deletePatQuery = "DELETE FROM " + patientTbl + " WHERE %v"
	findPatQuery   = "SELECT id, name, email, IFNULL(phone, '') as phone, " +
		"IFNULL(profile_picture, 0) as profile_picture FROM " + patientTbl + " WHERE "
)

func makeInsertPatQuery(patient *models.PatientInfo) (query string, queryArgs []any, err error) {
	if patient.Name == "" || patient.Email == "" || patient.Password == "" {
		err = newQueryError("patient name, email and password can't be empty")
		return
	}
	columns := "name,email,password"
	values := "?,?,?"
	queryArgs = append(queryArgs, patient.Name)
	queryArgs = append(queryArgs, patient.Email)
	queryArgs = append(queryArgs, patient.Password)
	if patient.Phone != "" {
		columns += ",phone"
		values += ",?"
		queryArgs = append(queryArgs, patient.Phone)
	}
	if patient.ProfilePictureID != 0 {
		columns += ",profile_picture"
		values += ",?"
		queryArgs = append(queryArgs, patient.ProfilePictureID)
	}
	query = fmt.Sprintf(insertPatQuery, columns, values)
	return
}

/*
* Main function to register patient into the system.
TODO: Test it
*
*/
func RegisterPatient(param patient.PutPatientRegisterParams) middleware.Responder {
	query, queryArgs, queryErr := makeInsertPatQuery(param.Info)
	if queryErr != nil {
		logger.Printf("[Error]bad input:%v", queryErr.Error())
		return patient.NewPutPatientRegisterDefault(400).WithPayload(models.Error(queryErr.Error()))
	}
	if lastID, _, execErr := ExecDataUpdateQuery(query, queryArgs...); execErr != nil {
		logger.Printf("[Error]patient register db functionality:%v", execErr)
		if duplicateEntryError(execErr) {
			return patient.NewPutPatientRegisterDefault(400).WithPayload("Requested patient already present")
		}
		return patient.NewPutPatientRegisterDefault(500).WithPayload("Internal error")
	} else {
		resp := &patient.PutPatientRegisterOKBody{RegisteredID: lastID}
		return patient.NewPutPatientRegisterOK().WithPayload(resp)
	}
}

/*
* Function to create patient update query.
query,queryArgs are to be used together in exec-function *
*/
func makeUpdatePatQuery(pat *models.PatientInfo) (query string, queryArgs []any, err error) {
	var set, cond string
	if pat.ID == 0 {
		return "", queryArgs, newQueryError("non-zero id required for update")
	}
	if pat.Email != "" {
		return "", queryArgs, newQueryError("email can't be updated")
	}
	if (pat.Name == "") && (pat.Phone == "") && (pat.ProfilePictureID == 0) {
		return "", queryArgs, newQueryError("one of name, phone, about, profile_picture is needed")
	}
	updateQueryListString(&cond, "id", ",")

	if pat.Name != "" {
		updateQueryListString(&set, "name", ",")
		queryArgs = append(queryArgs, pat.Name)
	}
	if pat.Phone != "" {
		updateQueryListString(&set, "phone", ",")
		queryArgs = append(queryArgs, pat.Phone)
	}
	if pat.ProfilePictureID != 0 {
		updateQueryListString(&set, "profile_picture", ",")
		queryArgs = append(queryArgs, pat.ProfilePictureID)
	}
	queryArgs = append(queryArgs, pat.ID)
	query = fmt.Sprintf(updatePatQuery, set, cond)
	return
}

/** Main function to update patient data in patientTbl **/
func UpdatePatient(param patient.PostPatientUpdateParams) middleware.Responder {
	query, queryArgs, queryErr := makeUpdatePatQuery(param.Info)
	if queryErr != nil {
		logger.Printf("[Error]Bad request:%v", queryErr)
		return patient.NewPostPatientUpdateDefault(400).WithPayload(models.Error(queryErr.Error()))
	}
	if _, _, err := ExecDataUpdateQuery(query, queryArgs...); err != nil {
		logger.Printf("[Error]db error:%v", err)
		return patient.NewPostPatientUpdateDefault(500).WithPayload("internale error")
	}
	return patient.NewPostPatientUpdateOK()
}

/** query,queryArgs are to be used together in exec-function **/
func makeDeletePatQuery(patID int64) (query string, queryArgs []any, err error) {
	if patID == 0 {
		err = newQueryError("non-zero id required")
		return
	}
	queryArgs = append(queryArgs, patID)
	query = fmt.Sprintf(deletePatQuery, " id = ? ")
	return
}

/** Main function to delete patient from patientTbl**/
func RemovePatient(param patient.DeletePatientRemoveParams) middleware.Responder {
	query, queryArgs, queryErr := makeDeletePatQuery(param.ID)
	if queryErr != nil {
		logger.Printf("[Error]Bad request:%v", queryErr)
		return patient.NewDeletePatientRemoveDefault(400).WithPayload(models.Error(queryErr.Error()))
	}
	if _, _, err := ExecDataUpdateQuery(query, queryArgs...); err != nil {
		logger.Printf("[Error]db error:%v", err)
		return patient.NewDeletePatientRemoveDefault(500).WithPayload("internale error")
	}
	return patient.NewDeletePatientRemoveOK()
}

/** query & queryArgs, should be used together**/
func makeFindPatQuery(param patient.GetPatientFindParams) (query string, queryArgs []any, err error) {
	if param.ID != nil && *param.ID != 0 {
		query = "id = ?"
		queryArgs = append(queryArgs, *param.ID)
	}
	if param.NameContaining != nil && *param.NameContaining != "" {
		updateQueryListStringWithOperation(&query, "LOWER(name) LIKE CONCAT('%',LOWER(?),'%')", "AND")
		queryArgs = append(queryArgs, *param.NameContaining)
	}
	if len(param.Ids) != 0 {
		placeHolders := ""
		for _, id := range param.Ids {
			updateQueryListStringWithOperation(&placeHolders, "?", ",")
			queryArgs = append(queryArgs, id)
		}
		updateQueryListStringWithOperation(&query, "id in ("+placeHolders+")", "AND")
	}
	if param.OfDoctor != nil {
		updateQueryListStringWithOperation(&query,
			fmt.Sprintf("id IN (SELECT patient_id FROM %v WHERE doctor_id = %v)",
				aptTbl, *param.OfDoctor), "AND")
	}
	if query == "" {
		err = newQueryError("one of query-param is needed")
	}

	query = findPatQuery + query
	logger.Printf("[INFO]query:%v args:%v err:%v", query, queryArgs, err)
	return
}

/** Main function to search patient in patientTbl **/
func FindPatient(param patient.GetPatientFindParams) middleware.Responder {
	query, queryArgs, queryErr := makeFindPatQuery(param)
	if queryErr != nil {
		logger.Printf("[Error]bad query queryErr:%v", queryErr)
		return patient.NewGetPatientFindDefault(400).WithPayload(models.Error(queryErr.Error()))
	}
	ctx, cancel := getTimeOutContext()
	defer cancel()
	rows, err := ExecDataFetchQuery(ctx, query, queryArgs...)
	if err != nil {
		logger.Printf("[Error]in find pat:query:%v::params:%v::error:%v",
			query, param, err)
		return patient.NewGetPatientFindDefault(500).WithPayload(models.Error("internal db error"))
	}
	defer rows.Close()
	fountPats := &patient.GetPatientFindOKBody{}
	for rows.Next() {
		patData := &models.PatientInfo{}
		if scanErr := rows.Scan(&patData.ID, &patData.Name, &patData.Email,
			&patData.Phone, &patData.ProfilePictureID); scanErr != nil {
			logger.Printf("[Error]patient data scan error:%v", scanErr)
			return patient.NewGetPatientFindDefault(500).WithPayload("internal db error in data read")
		}
		fountPats.Patients = append(fountPats.Patients, patData)
	}
	logger.Printf("[Success]find patient API")
	return patient.NewGetPatientFindOK().WithPayload(fountPats)
}

/** /patient/login **/
type patientLogin struct {
	patient.GetPatientLoginParams
	info      patient.GetPatientLoginOK
	dataError error
}

func (okResp *patientLogin) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	cookie := &http.Cookie{Name: "th-ssid", Value: okResp.info.Payload.SessionID, Path: "/",
		SameSite: http.SameSiteNoneMode, Secure: true}
	http.SetCookie(rw, cookie)
	rw.WriteHeader(200)
	if respErr := producer.Produce(rw, okResp.info.Payload); respErr != nil {
		logger.Fatalf("[CRITICAL ERROR] Unable to write response:%v", respErr)
	}
}

func (*patientLogin) errResponse(httpStatusCode int, err error) middleware.Responder {
	return patient.NewGetPatientLoginDefault(httpStatusCode).WithPayload(models.Error(err.Error()))
}

func (data *patientLogin) okResponse(int64, int64) middleware.Responder {
	if data.dataError != nil {
		logger.Printf("[Query Error] Query or DB error:%v", data.dataError)
		return patient.NewGetPatientLoginDefault(400).WithPayload(models.Error(data.dataError.Error()))
	}
	return data
}

func (req *patientLogin) makeQuery() (sqlQuery sqlExeParams, err error) {
	if req.Email == "" {
		err = newQueryError("email required")
		return
	}
	sqlQuery.Query = "SELECT id, UUID() as session_id, name, email, IFNULL(phone, ''), IFNULL(about, ''), IFNULL(profile_picture, 0) FROM " + patientTbl +
		" WHERE email = ? AND password = ?"
	sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.Email, req.Password)
	return
}

func (resp *patientLogin) scanRows(rows *sql.Rows) error {
	scannedRows := 0
	for ; rows.Next(); scannedRows++ {
		patientData := &models.PatientInfo{}
		respBody := &patient.GetPatientLoginOKBody{}
		if err := rows.Scan(&patientData.ID, &respBody.SessionID, &patientData.Name, &patientData.Email, &patientData.Phone, &patientData.About, &patientData.ProfilePictureID); err != nil {
			logger.Printf("[Scan Error]:%v", err)
			return err
		}
		respBody.Patient = patientData
		resp.info.WithPayload(respBody)
	}
	switch scannedRows {
	case 0:
		resp.dataError = newQueryError("no doctor found with given email-id")
	case 1:
	default:
		resp.dataError = errors.New("internal db error: multiple doctors with same email id")
	}
	if err := updateLoginSession(resp.info.Payload.Patient.ID, resp.info.Payload.SessionID, "offline", "patient"); err != nil {
		logger.Printf("[Error] In doc login session-id updation:%v", err)
		resp.dataError = errors.New("internal db error: In creating login data")
	}
	return nil
}

func PatientLoginAPI(_req patient.GetPatientLoginParams) middleware.Responder {
	req := &patientLogin{GetPatientLoginParams: _req}
	return FetchAndRespond(req)
}

/** End of /doctor/login **/
