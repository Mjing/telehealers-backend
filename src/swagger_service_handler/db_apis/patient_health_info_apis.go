/** patient health info API implementation **/
package apis

import (
	"database/sql"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"telehealers.in/router/models"
	phi "telehealers.in/router/restapi/operations/patient_health_info"
)

// Insert health info implementation
type InsertHealthInfo phi.PutPatientHealthInfoAddParams

func (param *InsertHealthInfo) makeQuery() (sqlReq sqlExeParams,
	err error) {
	if (param.Info.PatientID == 0) || (param.Info.Date == "") || (param.Info.Time == "") {
		err = newQueryError("patient_id, date, time can't be empty")
		return
	}
	columns, values := "patient_id, date, time", "?,?,?"
	sqlReq.QueryArgs = []any{param.Info.PatientID, param.Info.Date, param.Info.Time}
	if param.Info.BloodPressure != "" {
		columns += ",blood_pressure"
		values += ",?"
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.BloodPressure)
	}
	if param.Info.Complaint != "" {
		columns += ",health_complaints"
		values += ",?"
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.Complaint)
	}
	if param.Info.Height != "" {
		columns += ",height"
		values += ",?"
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.Height)
	}
	if param.Info.Weight != "" {
		columns += ",weight"
		values += ",?"
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.Info.Weight)
	}

	sqlReq.Query = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)",
		patienHealthInfoTbl, columns, values)
	return
}

func (param *InsertHealthInfo) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Insert health-info]API Error:http-code:%v|:%v", httpCode, err)
	return phi.NewPutPatientHealthInfoAddDefault(httpCode).WithPayload(models.Error(err.Error()))
}

func (param *InsertHealthInfo) okResponse(int64, int64) middleware.Responder {
	logger.Printf("[Insert health-info]finished successfuly")
	return phi.NewPutPatientHealthInfoAddOK()
}

func AddHealthInfoAPI(param phi.PutPatientHealthInfoAddParams,
	p *models.Principal) middleware.Responder {
	req := InsertHealthInfo(param)
	return UpdateAndRespond(&req)
}

//Find health info implementation

type FindHealthInfo struct {
	req         phi.GetPatientHealthInfoFindParams
	fetchedData models.HealthInfo
}

func (param *FindHealthInfo) makeQuery() (sqlReq sqlExeParams, err error) {
	columns := "id,gender,height,weight,bp,health_complaints,patient_id,DATE(created_on) as date, TIME(created_on) as time"
	if (param.req.AppointmentID != nil) && (*param.req.AppointmentID != 0) {
		sqlReq.Query = fmt.Sprintf("SELECT %v FROM %v infoTbl, %v aptTbl WHERE aptTbl.id = ? AND "+
			"infoTbl.id = aptTbl.patient_health_info_id", columns, patienHealthInfoTbl, aptTbl)
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.req.AppointmentID)
	} else if (param.req.PatientID != nil) && (*param.req.PatientID != 0) {
		sqlReq.Query = fmt.Sprintf("SELECT %v FROM %v info, %v patient WHERE patient.id = ? AND "+
			"info.id = patient.health_info_id", columns, patienHealthInfoTbl, patientTbl)
	} else {
		err = newQueryError("one of appointment_id or patient_id is required")
	}
	return
}

func (param *FindHealthInfo) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Find health-info]Error: status:%v|err:%v", httpCode, err)
	return phi.NewGetPatientHealthInfoFindDefault(httpCode).WithPayload(models.Error(err.Error()))
}

func (param *FindHealthInfo) scanRows(rows *sql.Rows) (scanErr error) {
	scanErr = rows.Scan(&param.fetchedData.ID, &param.fetchedData.Gender, &param.fetchedData.Height,
		&param.fetchedData.Weight, &param.fetchedData.BloodPressure, &param.fetchedData.Complaint,
		&param.fetchedData.PatientID, &param.fetchedData.Date, &param.fetchedData.Time)
	if scanErr != nil {
		logger.Printf("[Find health-info]row-scan error:%v", scanErr)
	}
	return
}

func (param *FindHealthInfo) okResponse(int64, int64) middleware.Responder {
	logger.Printf("[Find health-info]Successful")
	return phi.NewGetPatientHealthInfoFindOK().WithPayload(&param.fetchedData)
}

func FindHealthInfoAPI(param phi.GetPatientHealthInfoFindParams, p *models.Principal) middleware.Responder {
	req := &FindHealthInfo{req: param}
	return FetchAndRespond(req)
}
