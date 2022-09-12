package apis

import (
	"database/sql"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/appointment"
)

var (
	regAptQuery = "INSERT INTO " + aptTbl + " (%v)"
)

/** Register appointment implementation **/
type RegisterAppointment appointment.PutAppointmentRegisterParams

func (param *RegisterAppointment) makeQuery() (sqlReq sqlExeParams, err error) {
	columns := "doctor_id, patient_id, date, start_time_requested, end_time_requested"
	if (param.Appointment.Date == "") || (param.Appointment.DoctorID == 0) || (param.Appointment.PatientID == 0) ||
		(param.Appointment.StartTimeRequested == "") || (param.Appointment.EndTimeRquested == "") {
		err = newQueryError(columns + " can't be empty")
		return
	}
	sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Appointment.DoctorID, param.Appointment.PatientID, param.Appointment.Date, param.Appointment.StartTimeRequested,
		param.Appointment.EndTimeRquested)
	sqlReq.Query = fmt.Sprintf(regAptQuery, columns)
	return
}

func (param *RegisterAppointment) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Register Appointment]Error:code:%v|error:%v", httpCode, err)
	return appointment.NewPutAppointmentRegisterDefault(httpCode).WithPayload(models.Error(err.Error()))
}

func (param *RegisterAppointment) okResponse(lastId, recordsAffected int64) middleware.Responder {
	logger.Printf("[Register Appointment]Successful")
	return appointment.NewPutAppointmentRegisterOK()
}

func RegisterppointmentAPI(param appointment.PutAppointmentRegisterParams,
	p *models.Principal) middleware.Responder {
	req := RegisterAppointment(param)
	return UpdateAndRespond(&req)
}

/** Update appointment implementation **/
type UpdateAppointment appointment.PostAppointmentUpdateParams

func (param *UpdateAppointment) makeQuery() (sqlReq sqlExeParams, err error) {
	if param.Info.ID != 0 {
		err = newQueryError("empty id not allowed")
		return
	}
	sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.ID)
	var setter string
	if param.Info.DoctorID != 0 {
		updateQueryListString(&setter, "doctor_id", ",")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.DoctorID)
	}
	if param.Info.PatientID != 0 {
		updateQueryListString(&setter, "patient_id", ",")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.PatientID)
	}
	if param.Info.StartTimeRequested == "" {
		updateQueryListString(&setter, "start_time_requested", ",")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.StartTimeRequested)
	}
	if param.Info.EndTimeRquested == "" {
		updateQueryListString(&setter, "end_time_requested", ",")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.EndTimeRquested)
	}
	if param.Info.Date == "" {
		updateQueryListString(&setter, "date", ",")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.Date)
	}
	sqlReq.Query = fmt.Sprintf(generalUpdateQuery, aptTbl, setter, "id = ?")
	return
}

func (*UpdateAppointment) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Update appointment API]Error:code:%v|error:%v", httpCode, err)
	return appointment.NewPostAppointmentUpdateDefault(httpCode).WithPayload(models.Error(err.Error()))
}

func (*UpdateAppointment) okResponse(int64, int64) middleware.Responder {
	logger.Printf("[Update appointment API]Successful")
	return appointment.NewPostAppointmentUpdateOK()
}

func UpdateAppointmentAPI(param appointment.PostAppointmentUpdateParams,
	p *models.Principal) middleware.Responder {
	req := UpdateAppointment(param)
	return UpdateAndRespond(&req)
}

/** Remove appointment implementaiton **/
type RemoveAppointment appointment.DeleteAppointmentRemoveParams

func (param *RemoveAppointment) makeQuery() (sqlReq sqlExeParams, err error) {
	if param.ID != 0 {
		err = newQueryError("id should be positive(non-zero)")
		return
	}
	sqlReq.Query = fmt.Sprintf(generalDeleteQuery, aptTbl, "id = ?")
	sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.ID)
	return
}
func (*RemoveAppointment) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Remove appointment API]Error:code:%v|error:%v", httpCode, err)
	return appointment.NewDeleteAppointmentRemoveDefault(httpCode).WithPayload(models.Error(err.Error()))
}
func (*RemoveAppointment) okResponse(lastId, affectedRows int64) middleware.Responder {
	logger.Printf("[Remove appointment API]Successful")
	return appointment.NewDeleteAppointmentRemoveOK()
}

func RemoveAppointmentAPI(param appointment.DeleteAppointmentRemoveParams,
	p *models.Principal) middleware.Responder {
	req := RemoveAppointment(param)
	return UpdateAndRespond(&req)
}

/** Find appointment implementation **/
type FindAppointment struct {
	param       appointment.GetAppointmentFindParams
	fetchedData *appointment.GetAppointmentFindOKBody
}

func (param *FindAppointment) makeQuery() (sqlReq sqlExeParams, err error) {
	req := &param.param
	if req.OnDate == "" {
		err = newQueryError("invalid on_date value")
		return
	}
	condition := " WHERE (on_date = ?)"
	orderBy := ""
	sqlReq.QueryArgs = append(sqlReq.QueryArgs, req.OnDate)
	if *req.DoctorID != 0 {
		condition += " AND (doctor_id = ?)"
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, *req.DoctorID)
	}
	if *req.PatientID != 0 {
		condition += " AND (patient_id = ?)"
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, *req.PatientID)
	}
	orderBy = " ORDER BY date, requested_start_time ASC"
	sqlReq.Query = fmt.Sprintf(generalFetchQuery, aptFetchColumns, aptTbl+condition, req.Size*(req.Page-1), req.Size) + orderBy
	return
}

func (param *FindAppointment) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Find-appointment API]Error:code:%v|error:%v", httpCode,
		err)
	return appointment.NewGetAppointmentFindDefault(httpCode).WithPayload(models.Error(err.Error()))
}

func (param *FindAppointment) okResponse(int64, int64) middleware.Responder {
	logger.Printf("[Find-appointment API]Successful")
	return appointment.NewGetAppointmentFindOK().WithPayload(param.fetchedData)
}

func (param *FindAppointment) scanRows(rows *sql.Rows) error {
	data := &appointment.GetAppointmentFindOKBody{}
	for rows.Next() {
		aptData := &models.AppointmentInfo{}
		if scanErr := rows.Scan(&aptData.ID, &aptData.DoctorID, &aptData.PatientID, &aptData.PatientHealthInfoID,
			&aptData.PrescriptionID, &aptData.Date, &aptData.StartTime, &aptData.EndTime,
			&aptData.StartTimeRequested, &aptData.EndTimeRquested); scanErr != nil {
			logger.Printf("[Find appointment API]Scan error:%v", scanErr)
			return scanErr
		}
		data.Appointments = append(data.Appointments, aptData)
	}
	param.fetchedData = data
	return nil
}

func FindAppointmentAPI(param appointment.GetAppointmentFindParams,
	p *models.Principal) middleware.Responder {
	req := &FindAppointment{param: param}
	return FetchAndRespond(req)
}
