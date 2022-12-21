/** APIs related to doctor data **/
package apis

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/doctor"
)

var (
	insertDocQuery = "INSERT INTO " + doctorTbl + " (%v) VALUES (%v)"
	updateDocQuery = "UPDATE " + doctorTbl + " SET %v WHERE %v"
	deleteDocQuery = "DELETE FROM " + doctorTbl + " WHERE %v"
	findDocQuery   = "SELECT id, name, email, phone, about, profile_picture, sign_pic FROM " +
		doctorTbl + " WHERE "
)

const queryErrorTag = "[Query Error]"

func makeInsertDocQuery(doctor *models.DoctorInfo) (query string, queryArgs []any, err error) {
	if doctor.Name == "" || doctor.Email == "" || doctor.Phone == "" {
		err = newQueryError("doctor name, email or phone can't be empty")
		return
	}
	columns := "name,email,phone"
	values := "?,?,?"
	queryArgs = append(queryArgs, doctor.Name)
	queryArgs = append(queryArgs, doctor.Email)
	queryArgs = append(queryArgs, doctor.Phone)
	if doctor.About != "" {
		columns += ",about"
		values += ",?"
		queryArgs = append(queryArgs, doctor.About)
	}
	if doctor.ProfilePictureID != 0 {
		columns += ",profile_picture"
		values += ",?"
		queryArgs = append(queryArgs, doctor.ProfilePictureID)
	}
	if doctor.SignPictureID != 0 {
		columns += ",sign_pic"
		values += ",?"
		queryArgs = append(queryArgs, doctor.SignPictureID)
	}
	query = fmt.Sprintf(insertDocQuery, columns, values)
	return
}

/*
* Main function to register doctor into the system.
TODO: Create register process: Apply->Verify->Approve
*
*/
func RegisterDoctor(param doctor.PutDoctorRegisterParams) middleware.Responder {
	query, queryArgs, queryErr := makeInsertDocQuery(param.Info)
	if queryErr != nil {
		logger.Printf("[Error]bad input:%v", queryErr.Error())
		return doctor.NewPutDoctorRegisterDefault(400).WithPayload(models.Error(queryErr.Error()))
	}
	if _, _, execErr := ExecDataUpdateQuery(query, queryArgs...); execErr != nil {
		logger.Printf("[Error]doctor register db functionality:%v", execErr)
		if duplicateEntryError(execErr) {
			return doctor.NewPutDoctorRegisterDefault(400).WithPayload("Requested doctor already present")
		}
		return doctor.NewPutDoctorRegisterDefault(500).WithPayload("Internal error")
	}
	return doctor.NewPutDoctorRegisterOK()
}

/** /doctor/register/apply **/
type doctorRegistrationApplication struct {
	doctor.PostDoctorRegisterApplyParams
}

func (req *doctorRegistrationApplication) makeQuery() (sqlQuery sqlExeParams, err error) {
	if req.Application.DoctorInfo.RegistrationID == "" || req.Application.DoctorInfo.Name == "" || req.Application.DoctorInfo.Email == "" ||
		req.Application.DoctorInfo.Password == "" {
		err = newQueryError("registration_id, name, email and password can't be empty")
		return
	}
	columns := "registration_number, name, email, password, applied_on"
	values := "?,?,?,?,NOW()"
	sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.Application.DoctorInfo.RegistrationID,
		req.Application.DoctorInfo.Name, req.Application.DoctorInfo.Email, req.Application.DoctorInfo.Password)
	if req.Application.AdditionalInfo != "" {
		columns += ",comments"
		values += ",?"
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.Application.AdditionalInfo)
	}
	sqlQuery.Query = "INSERT INTO " + doctorRegistrationApplicationTbl + fmt.Sprintf(" (%v) VALUES (%v)", columns, values) +
		" ON DUPLICATE KEY UPDATE applied_on = NOW(), name = ?, email = ?, password = ?"
	sqlQuery.QueryArgs = append(sqlQuery.QueryArgs,
		req.Application.DoctorInfo.Name, req.Application.DoctorInfo.Email, req.Application.DoctorInfo.Password)
	return
}

func (*doctorRegistrationApplication) errResponse(httpStatusCode int, err error) middleware.Responder {
	return doctor.NewPostDoctorRegisterApplyDefault(httpStatusCode).WithPayload(models.Error(err.Error()))
}

func (*doctorRegistrationApplication) okResponse(int64, int64) middleware.Responder {
	return doctor.NewPostDoctorRegisterApplyOK()
}

func DoctorRegistrationApplicationAPI(_req doctor.PostDoctorRegisterApplyParams) middleware.Responder {
	req := &doctorRegistrationApplication{_req}
	return UpdateAndRespond(req)
}

/* /doctor/register/review */
type doctorRegistrationReview struct {
	doctor.PostDoctorRegisterReviewParams
}

func (*doctorRegistrationReview) errResponse(httpStatusCode int, err error) middleware.Responder {
	return doctor.NewPostDoctorRegisterReviewDefault(httpStatusCode).WithPayload(models.Error(err.Error()))
}

func (*doctorRegistrationReview) okResponse(int64, int64) middleware.Responder {
	return doctor.NewPostDoctorRegisterReviewOK()
}

func (req *doctorRegistrationReview) makeApproveQuery() (sqlQuery sqlExeParams, err error) {
	/** inv: req is well formed **/
	sqlQuery.Query = "INSERT INTO " + doctorTbl + " (name, email, registration_number, password) SELECT name, email, registration_number, password FROM " +
		doctorRegistrationApplicationTbl + " WHERE id = ?"
	sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.Review.ApplicationID)
	return
}

func (req *doctorRegistrationReview) makeDenyQuery() (sqlQuery sqlExeParams, err error) {
	/** inv: req is well formed **/
	if req.Review.ReviewerComments == "" {
		err = newQueryError("for denied request, reviewer comments are necessary")
		return
	}
	sqlQuery.Query = "UPDATE " + doctorRegistrationApplicationTbl + " SET approved = ?, reviewer_comments = ? WHERE id = ?"
	sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.Review.Approve, req.Review.ReviewerComments, req.Review.ApplicationID)
	return
}

func (req *doctorRegistrationReview) makeQuery() (sqlQuery sqlExeParams, err error) {
	if req.Review.ApplicationID == nil || *req.Review.ApplicationID == 0 || req.Review.Approve == nil {
		err = newQueryError("application_id, approve params can't be empty")
		return
	}
	if *req.Review.Approve {
		return req.makeApproveQuery()
	} else {
		return req.makeDenyQuery()
	}
}

func DoctorRegistrationApplicationReviewAPI(_req doctor.PostDoctorRegisterReviewParams) middleware.Responder {
	req := &doctorRegistrationReview{_req}
	return UpdateAndRespond(req)
}

/** End /doctor/register/review **/

/*
* Function to create doctor update query.
query,queryArgs are to be used together in exec-function *
*/
func makeUpdateDoctorQuery(doc *models.DoctorInfo) (query string, queryArgs []any, err error) {
	var set, cond string
	if doc.ID == 0 {
		return "", queryArgs, newQueryError("non-zero id required for update")
	}
	if doc.Email != "" {
		return "", queryArgs, newQueryError("email can't be updated")
	}
	if doc.RegistrationID != "" {
		return "", queryArgs, newQueryError("registration_id can't be updated")
	}
	if (doc.Name == "") && (doc.Phone == "") && (doc.About == "") && (doc.ProfilePictureID == 0) {
		return "", queryArgs, newQueryError("one of name, phone, about, profile_picture is needed")
	}
	updateQueryListString(&cond, "id", ",")

	if doc.Name != "" {
		updateQueryListString(&set, "name", ",")
		queryArgs = append(queryArgs, doc.Name)
	}
	if doc.Phone != "" {
		updateQueryListString(&set, "phone", ",")
		queryArgs = append(queryArgs, doc.Phone)
	}
	if doc.About != "" {
		updateQueryListString(&set, "about", ",")
		queryArgs = append(queryArgs, doc.About)
	}
	if doc.ProfilePictureID != 0 {
		updateQueryListString(&set, "profile_picture", ",")
		queryArgs = append(queryArgs, doc.ProfilePictureID)
	}
	if doc.SignPictureID != 0 {
		updateQueryListString(&set, "sign_pic", ",")
		queryArgs = append(queryArgs, doc.SignPictureID)
	}
	queryArgs = append(queryArgs, doc.ID)
	return fmt.Sprintf(updateDocQuery, set, cond), queryArgs, nil
}

/** Main function to update doctor data in doctorTbl **/
func UpdateDoctor(param doctor.PostDoctorUpdateParams) middleware.Responder {
	query, queryArgs, queryErr := makeUpdateDoctorQuery(param.Info)
	if queryErr != nil {
		logger.Printf("[Error]Bad request:%v", queryErr)
		return doctor.NewPostDoctorUpdateDefault(400).WithPayload(models.Error(queryErr.Error()))
	}
	if _, _, err := ExecDataUpdateQuery(query, queryArgs...); err != nil {
		logger.Printf("[Error]db error:%v|query:%v[args:%v]", err, query, queryArgs)
		return doctor.NewPostDoctorUpdateDefault(500).WithPayload("internale error")
	}
	return doctor.NewPostDoctorUpdateOK()
}

/** query,queryArgs are to be used together in exec-function **/
func makeDeleteDoctorQuery(doc int64) (query string, queryArgs []any, err error) {
	if doc == 0 {
		err = newQueryError("non-zero id required")
		return
	}
	queryArgs = append(queryArgs, doc)
	return fmt.Sprintf(deleteDocQuery, " id = ? "), queryArgs, nil
}

/** Main function to delete doctor from doctorTbl**/
func RemoveDoctor(param doctor.DeleteDoctorRemoveParams) middleware.Responder {
	query, queryArgs, queryErr := makeDeleteDoctorQuery(param.ID)
	if queryErr != nil {
		logger.Printf("[Error]Bad request:%v", queryErr)
		return doctor.NewDeleteDoctorRemoveDefault(400).WithPayload(models.Error(queryErr.Error()))
	}
	if _, _, err := ExecDataUpdateQuery(query, queryArgs...); err != nil {
		logger.Printf("[Error]db error:%v", err)
		return doctor.NewPostDoctorUpdateDefault(500).WithPayload("internale error")
	}
	return doctor.NewPostDoctorUpdateOK()
}

/** query & queryArgs, should be used together**/
func makeFindDoctorQuery(param doctor.GetDoctorFindParams) (query string, queryArgs []any, err error) {
	if param.ID != nil && *param.ID != 0 {
		query = findDocQuery + "id = ?"
		queryArgs = append(queryArgs, *param.ID)
	} else if param.NameContaining != nil && *param.NameContaining != "" {
		query = findDocQuery + "LOWER(name) LIKE CONCAT('%',LOWER(?),'%')"
		queryArgs = append(queryArgs, *param.NameContaining)
	} else {
		err = newQueryError("one of id or name_containing param needed")
	}
	logger.Printf("[INFO]query:%v args:%v err:%v", query, queryArgs, err)
	return
}

/** Main function to search doctor in doctorTbl **/
func FindDoctor(param doctor.GetDoctorFindParams) middleware.Responder {
	query, queryArgs, queryErr := makeFindDoctorQuery(param)
	if queryErr != nil {
		logger.Printf("[Error]bad query queryErr:%v", queryErr)
		return doctor.NewGetDoctorFindDefault(400).WithPayload(models.Error(queryErr.Error()))
	}
	ctx, cancel := getTimeOutContext()
	defer cancel()
	rows, err := ExecDataFetchQuery(ctx, query, queryArgs...)
	if err != nil {
		logger.Printf("[Error]in find doc:query:%v::params:%v::error:%v",
			query, param, err)
		return doctor.NewGetDoctorFindDefault(500).WithPayload(models.Error("internal db error"))
	}
	defer rows.Close()
	foundDocs := &doctor.GetDoctorFindOKBody{}
	for rows.Next() {
		docData := &models.DoctorInfo{}
		if scanErr := rows.Scan(&docData.ID, &docData.Name, &docData.Email,
			&docData.Phone, &docData.About, &docData.ProfilePictureID, &docData.SignPictureID); scanErr != nil {
			logger.Printf("[Error]doctor data scan error:%v", scanErr)
			return doctor.NewGetDoctorFindDefault(500).WithPayload("internal db error in data read")
		}
		foundDocs.Doctors = append(foundDocs.Doctors, docData)
	}
	logger.Printf("[Success]find game api")
	return doctor.NewGetDoctorFindOK().WithPayload(foundDocs)
}

/** /doctor/login **/
type doctorLogin struct {
	doctor.GetDoctorLoginParams
	info      doctor.GetDoctorLoginOK
	dataError error
}

func (okResp *doctorLogin) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	cookie := &http.Cookie{Name: "th-ssid", Value: okResp.info.Payload.SessionID, Path: "/",
		SameSite: http.SameSiteNoneMode, Secure: true}
	http.SetCookie(rw, cookie)
	rw.WriteHeader(200)
	if respErr := producer.Produce(rw, okResp.info.Payload); respErr != nil {
		logger.Fatalf("[CRITICAL ERROR] Unable to write response:%v", respErr)
	}
}

func (*doctorLogin) errResponse(httpStatusCode int, err error) middleware.Responder {
	return doctor.NewGetDoctorLoginDefault(httpStatusCode).WithPayload(models.Error(err.Error()))
}

func (data *doctorLogin) okResponse(int64, int64) middleware.Responder {
	if data.dataError != nil {
		logger.Printf("[Query Error] Query or DB error:%v", data.dataError)
		return doctor.NewGetDoctorLoginDefault(400).WithPayload(models.Error(data.dataError.Error()))
	}
	return data
}

func (req *doctorLogin) makeQuery() (sqlQuery sqlExeParams, err error) {
	if req.Email == "" {
		err = newQueryError("email required")
		return
	}
	sqlQuery.Query = "SELECT id, UUID() as session_id, name, email, phone, about, profile_picture, registration_number FROM " + doctorTbl +
		" WHERE email = ? AND password = ?"
	sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.Email, req.Password)
	return
}

func (resp *doctorLogin) scanRows(rows *sql.Rows) error {
	scannedRows := 0
	for ; rows.Next(); scannedRows++ {
		docData := &models.DoctorInfo{}
		respBody := &doctor.GetDoctorLoginOKBody{}
		if err := rows.Scan(&docData.ID, &respBody.SessionID, &docData.Name, &docData.Email, &docData.Phone, &docData.About, &docData.ProfilePictureID, &docData.RegistrationID); err != nil {
			logger.Printf("[Scan Error]:%v", err)
			return err
		}
		respBody.Doctor = docData
		resp.info.WithPayload(respBody)
	}
	switch scannedRows {
	case 0:
		resp.dataError = newQueryError("no doctor found with given email-id")
	case 1:
	default:
		resp.dataError = errors.New("internal db error: multiple doctors with same email id")
	}
	if err := updateLoginSession(resp.info.Payload.Doctor.ID, resp.info.Payload.SessionID, "offline", "doctor"); err != nil {
		logger.Printf("[Error] In doc login session-id updation:%v", err)
		resp.dataError = errors.New("internal db error: In creating login data")
	}
	return nil
}

func DoctorLoginAPI(_req doctor.GetDoctorLoginParams) middleware.Responder {
	req := &doctorLogin{GetDoctorLoginParams: _req}
	return FetchAndRespond(req)
}

/** End of /doctor/login **/

/** /doctor/register/pending_applications **/
type getPendingDocApp struct {
	doctor.GetDoctorRegisterPendingApplicationsParams
	data doctor.GetDoctorRegisterPendingApplicationsOKBody
}

func (*getPendingDocApp) errResponse(httpStatusCode int, err error) middleware.Responder {
	return doctor.NewGetDoctorRegisterPendingApplicationsDefault(httpStatusCode).WithPayload(models.Error(err.Error()))
}

func (data *getPendingDocApp) okResponse(int64, int64) middleware.Responder {
	return doctor.NewGetDoctorRegisterPendingApplicationsOK().WithPayload(&data.data)
}

func (req *getPendingDocApp) makeQuery() (sqlQuery sqlExeParams, err error) {
	filter := " WHERE approved = FALSE"
	if req.AppliedAfter != nil {
		updateQueryListStringWithOperation(&filter, "applied_on >= ?", " AND ")
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, *req.AppliedAfter)
	}
	if req.AppliedBefore != nil {
		updateQueryListStringWithOperation(&filter, "applied_on <= ?", " AND ")
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, *req.AppliedBefore)
	}
	if req.NameLike != nil {
		updateQueryListStringWithOperation(&filter, "name LIKE CONCAT('%', CONCAT(?, '%'))", " AND ")
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, *req.NameLike)
	}
	sqlQuery.Query = "SELECT id, name, email, registration_number, applied_on, comments, reviewer_comments FROM " + doctorRegistrationApplicationTbl + filter +
		" ORDER BY applied_on " + req.Sort + " LIMIT ?, ?"
	sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.PageSize*(req.Page-1), req.PageSize)
	return
}

func (resp *getPendingDocApp) scanRows(rows *sql.Rows) (err error) {
	for rows.Next() {
		rowData := &models.RegistrationApplication{}
		if err = rows.Scan(&rowData.ID, &rowData.Name, &rowData.Email, &rowData.RegistrationNumber,
			&rowData.AppliedOn, &rowData.Comments, &rowData.ReviewerComments); err != nil {
			logger.Printf("[Scan Error]In /doctor/register/pending_application:%v", err)
			return
		}
		resp.data.Applications = append(resp.data.Applications, rowData)
	}
	return
}

func GetDoctorRegistrationApplicationAPI(_req doctor.GetDoctorRegisterPendingApplicationsParams) middleware.Responder {
	req := &getPendingDocApp{GetDoctorRegisterPendingApplicationsParams: _req}
	return FetchAndRespond(req)
}

/** End of /doctor/register/pending_applications **/

/** /doctor/patients API **/
type doctorTreatedPatients struct {
	doctor.GetDoctorPatientsParams
	resp doctor.GetDoctorPatientsOKBody
}

func (*doctorTreatedPatients) errResponse(httpStatus int, err error) middleware.Responder {
	return doctor.NewGetDoctorPatientsDefault(httpStatus).WithPayload(models.Error(err.Error()))
}

func (data *doctorTreatedPatients) okResponse(int64, int64) middleware.Responder {
	return doctor.NewGetDoctorPatientsOK().WithPayload(&data.resp)
}

func (req *doctorTreatedPatients) makeQuery() (sqlQuery sqlExeParams, err error) {
	if req.DoctorID == 0 {
		err = newQueryError("non-zero doctor_id required")
		return
	}
	sqlQuery.Query = "SELECT p.id, p.name, p.email, p.phone, p.profile_picture FROM " + patientTbl + " as p, " +
		aptTbl + " as apt WHERE apt.patient_id = p.id AND apt.doctor_id = ? LIMIT " +
		fmt.Sprintf(" %v, %v", (req.Page-1)*req.PageSize, req.PageSize)
	sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.DoctorID)
	return
}

func (data *doctorTreatedPatients) scanRows(rows *sql.Rows) (err error) {
	for rows.Next() {
		patientData := &models.PatientInfo{}
		if err = rows.Scan(&patientData.ID, &patientData.Name, &patientData.Email,
			&patientData.Phone, &patientData.ProfilePictureID); err != nil {
			logger.Printf("[Scan Error]fetching doctor's patients:%v", err)
			return
		}
		data.resp.Patients = append(data.resp.Patients, patientData)
	}
	return
}

func DoctorRelatedPatientsAPI(_req doctor.GetDoctorPatientsParams, p *models.Principal) middleware.Responder {
	req := &doctorTreatedPatients{GetDoctorPatientsParams: _req}
	return FetchAndRespond(req)
}

/** API /doctor/online **/

/*
* Login session maintanance helpers
status: 'ONLINE', 'OFFLINE', ..
*
*/
func updateLoginSession(userID int64, sessionId, status, userType string) (err error) {
	if (userType != "patient") && (userType != "doctor") {
		return fmt.Errorf("bad user-type:%v", userType)
	}

	query := "INSERT INTO " + sessionTbl +
		" (user_id, user_type, session_id, status) VALUES (?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE session_id = ?, status = ?"
	if _, _, err = ExecDataUpdateQuery(query, userID, userType, sessionId, status, sessionId, status); err != nil {
		logger.Printf("[Session ID Error]:%v:%v", query, err)
	}
	return
}

func DoctorOnlineAPI(req doctor.PostDoctorOnlineParams, p *models.Principal) middleware.Responder {
	status := "online"
	if *req.Req.Status == "OFFLINE" {
		status = "offline"
	}
	if err := updateLoginSession(*req.Req.DoctorID, *req.Req.SessionID, status, "doctor"); err != nil {
		logger.Printf("[Error]in updating login session:%v", err)
		return doctor.NewGetDoctorLoginDefault(500).WithPayload("Internal db error")
	}
	return doctor.NewPostDoctorOnlineOK()
}
