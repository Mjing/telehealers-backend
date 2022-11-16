package apis

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/appointment"
)

var (
	regAptQuery = "INSERT INTO " + aptTbl + " (%v)"
)

/** Register appointment implementation **/
type RegisterAppointment appointment.PutAppointmentRegisterParams

/** inv: len(columns) == len(valsAndPlaceHolders) AND len(new(queryArgs)) == len(PlaceHolders) **/
func parseAptInfoTimeAttrs(apt *models.AppointmentInfo, queryArgs []any) (columns, valsAndPlaceHolders []string) {
	/** Flags to indicate which columns should be initialized to now**/
	var dateNow, strNow, stNow, endNow, endRNow bool
	/**Parse init array and toggle found flags**/
	for _, v := range apt.InitializeToNow {
		switch v {
		case "date":
			dateNow = true
		case "start_time_requested":
			strNow = true
		case "start_time":
			stNow = true
		case "end_time_requested":
			endRNow = true
		case "end_time":
			endNow = true
		}
	}

	queryArgs = append(queryArgs, apt.DoctorID, apt.PatientID)
	if dateNow {
		columns = append(columns, "date")
		valsAndPlaceHolders = append(valsAndPlaceHolders, "DATE(NOW())")
	} else if apt.Date != "" {
		columns = append(columns, "date")
		valsAndPlaceHolders = append(valsAndPlaceHolders, "?")
		queryArgs = append(queryArgs, apt.Date)
	}

	if strNow {
		columns = append(columns, "requested_start_time")
		valsAndPlaceHolders = append(valsAndPlaceHolders, "TIME(NOW())")
	} else if apt.StartTimeRequested != "" {
		columns = append(columns, "requested_start_time")
		valsAndPlaceHolders = append(valsAndPlaceHolders, "?")
		queryArgs = append(queryArgs, apt.StartTimeRequested)
	}

	if stNow {
		columns = append(columns, "start_time")
		valsAndPlaceHolders = append(valsAndPlaceHolders, "TIME(NOW())")
	} else if apt.StartTime != "" {
		columns = append(columns, "start_time")
		valsAndPlaceHolders = append(valsAndPlaceHolders, "?")
		queryArgs = append(queryArgs, apt.StartTime)
	}

	if endRNow {
		columns = append(columns, "requested_end_time")
		valsAndPlaceHolders = append(valsAndPlaceHolders, "TIME(NOW())")
	} else if apt.EndTimeRquested != "" {
		columns = append(columns, "requested_end_time")
		valsAndPlaceHolders = append(valsAndPlaceHolders, "?")
		queryArgs = append(queryArgs, apt.EndTimeRquested)
	}

	if endNow {
		columns = append(columns, "end_time")
		valsAndPlaceHolders = append(valsAndPlaceHolders, "TIME(NOW())")
	} else if apt.EndTime != "" {
		columns = append(columns, "end_time")
		valsAndPlaceHolders = append(valsAndPlaceHolders, "?")
		queryArgs = append(queryArgs, apt.EndTime)
	}
	return
}

func (param *RegisterAppointment) makeQuery() (sqlReq sqlExeParams, err error) {
	columns := "doctor_id, patient_id"
	if (param.Appointment.DoctorID == 0) || (param.Appointment.PatientID == 0) {
		err = newQueryError("doctor_id and patient_id can't be empty")
		return
	}
	sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Appointment.DoctorID, param.Appointment.PatientID)
	valueSet := "?,?"
	timeColumns, timeVorPs := parseAptInfoTimeAttrs(param.Appointment, sqlReq.QueryArgs)
	if len(timeColumns) > 0 {
		columns += "," + strings.Join(timeColumns, ",")
		valueSet += "," + strings.Join(timeVorPs, ",")
	}
	sqlReq.Query = fmt.Sprintf(regAptQuery, columns) + " VALUES (" + valueSet + ")"
	logger.Printf("[Appointment QUery]:%v", sqlReq)
	return
}

func (param *RegisterAppointment) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Register Appointment]Error:code:%v|error:%v", httpCode, err)
	return appointment.NewPutAppointmentRegisterDefault(httpCode).WithPayload(models.Error(err.Error()))
}

func (param *RegisterAppointment) okResponse(lastId, recordsAffected int64) middleware.Responder {
	logger.Printf("[Register Appointment]Successful")
	return appointment.NewPutAppointmentRegisterOK().WithPayload(lastId)
}

func RegisterppointmentAPI(param appointment.PutAppointmentRegisterParams,
	p *models.Principal) middleware.Responder {
	req := RegisterAppointment(param)
	return UpdateAndRespond(&req)
}

/** Update appointment implementation **/
type UpdateAppointment appointment.PostAppointmentUpdateParams

func (param *UpdateAppointment) makeQuery() (sqlReq sqlExeParams, err error) {
	if param.Info.ID == 0 {
		err = newQueryError("empty id not allowed")
		return
	}
	var setter string
	if param.Info.DoctorID != 0 {
		updateQueryListString(&setter, "doctor_id", ",")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.DoctorID)
	}
	if param.Info.PatientID != 0 {
		updateQueryListString(&setter, "patient_id", ",")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.PatientID)
	}
	timeColumns, timeVorPs := parseAptInfoTimeAttrs(param.Info, sqlReq.QueryArgs)
	for i, col := range timeColumns {
		updateQueryListStringWithOperation(&setter, fmt.Sprintf("%v = %v", col, timeVorPs[i]), ",")
	}
	sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.ID)
	sqlReq.Query = fmt.Sprintf(generalUpdateQuery, aptTbl, setter, "id = ?")
	logger.Printf("[Update apt query]:%v", sqlReq)
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
	if req.OnlyPending == nil {
		if *req.OnlyPending {
			condition += " AND (prescription_id == NULL)"
		} else {
			condition += " AND (prescription_id != NULL)"
		}
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

/** /appointment/count: Find appointment implementation **/
type CountAppointment struct {
	param       appointment.GetAppointmentCountParams
	fetchedData *appointment.GetAppointmentCountOKBody
}

func (param *CountAppointment) makeQuery() (sqlReq sqlExeParams, err error) {
	req := &param.param
	if req.OnDate == "" {
		err = newQueryError("invalid on_date value")
		return
	}
	condition := " WHERE (on_date = ?)"
	sqlReq.QueryArgs = append(sqlReq.QueryArgs, req.OnDate)
	if req.DoctorID == nil && *req.DoctorID != 0 {
		condition += " AND (doctor_id = ?)"
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, *req.DoctorID)
	}
	if req.PatientID == nil && *req.PatientID != 0 {
		condition += " AND (patient_id = ?)"
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, *req.PatientID)
	}
	if req.OnlyPending == nil {
		if *req.OnlyPending {
			condition += " AND (prescription_id == NULL)"
		} else {
			condition += " AND (prescription_id != NULL)"
		}
	}
	sqlReq.Query = "SELECT COUNT(*) FROM " + aptTbl + condition
	return
}

func (param *CountAppointment) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Find-appointment API]Error:code:%v|error:%v", httpCode,
		err)
	return appointment.NewGetAppointmentCountDefault(httpCode).WithPayload(models.Error(err.Error()))
}

func (param *CountAppointment) okResponse(int64, int64) middleware.Responder {
	logger.Printf("[Find-appointment API]Successful")
	return appointment.NewGetAppointmentCountOK().WithPayload(param.fetchedData)
}

func (param *CountAppointment) scanRows(rows *sql.Rows) error {
	data := &appointment.GetAppointmentCountOKBody{}
	for rows.Next() {
		if scanErr := rows.Scan(&data.Count); scanErr != nil {
			logger.Printf("[Find appointment API]Scan error:%v", scanErr)
			return scanErr
		}
	}
	param.fetchedData = data
	return nil
}

func CountAppointmentAPI(param appointment.GetAppointmentCountParams,
	p *models.Principal) middleware.Responder {
	req := &CountAppointment{param: param}
	return FetchAndRespond(req)
}
