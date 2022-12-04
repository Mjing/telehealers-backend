package apis

import (
	"errors"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/appointment"
	"telehealers.in/router/restapi/operations/prescription"
)

/** Returns respondable errors**/
func insertPrescription(req *models.Prescription) (id int64, httpStatusCode int, err error) {
	var sqlQuery sqlExeParams
	sqlQuery.Query = "INSERT INTO " + prescriptionTbl + " (%v) VALUES (%v)"
	onDuplicate := " ON DUPLICATE KEY UPDATE %v"
	columns, values, updateString := "", "", ""
	updateArgs := make([]any, 1)
	httpStatusCode = 200
	if req.CreatedBy != 0 {
		updateQueryListStringWithOperation(&columns, "created_by", ",")
		updateQueryListStringWithOperation(&values, "?", ",")
		updateQueryListString(&updateString, "created_by", ",")
		updateArgs = append(updateArgs, req.CreatedBy)
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.CreatedBy)
	} else {
		err = newQueryError("created_by required")
		httpStatusCode = 403
		return
	}
	if req.Name != "" {
		updateQueryListStringWithOperation(&columns, "name", ",")
		updateQueryListStringWithOperation(&values, "?", ",")
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.Name)
	}
	if req.CommentsOnMeds != "" {
		updateQueryListStringWithOperation(&columns, "comment_on_medicines", ",")
		updateQueryListStringWithOperation(&values, "?", ",")
		updateQueryListString(&updateString, "comment_on_medicines", ",")
		updateArgs = append(updateArgs, req.CommentsOnMeds)
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.CommentsOnMeds)
	}
	if req.CommentsOnTests != "" {
		updateQueryListStringWithOperation(&columns, "comment_on_tests", ",")
		updateQueryListStringWithOperation(&values, "?", ",")
		updateQueryListString(&updateString, "comment_on_tests", ",")
		updateArgs = append(updateArgs, req.CommentsOnTests)
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.CommentsOnTests)
	}
	if req.OverallAdvice != "" {
		updateQueryListStringWithOperation(&columns, "comment_on_advices", ",")
		updateQueryListStringWithOperation(&values, "?", ",")
		updateQueryListString(&updateString, "comment_on_advices", ",")
		updateArgs = append(updateArgs, req.OverallAdvice)
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, req.OverallAdvice)
	}

	if columns == "" {
		err = errors.New("one of name, comments_on_meds, comments_on_tests or overall_advice is required")
		httpStatusCode = 403
	} else {
		sqlQuery.Query = fmt.Sprintf(sqlQuery.Query, columns, values)
	}
	if updateString != "" {
		sqlQuery.Query = sqlQuery.Query + fmt.Sprint(onDuplicate, updateString)
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, updateArgs...)
	}
	if prescID, _, insErr := ExecDataUpdateQuery(sqlQuery.Query, sqlQuery.QueryArgs...); insErr != nil {
		logger.Printf("[Query Error]In Prescription insertion:%v:{%v[Args:%v]}", insErr, sqlQuery.Query, sqlQuery.QueryArgs)
		err = errors.New("internal db error")
		httpStatusCode = 500
	} else {
		id = prescID
	}
	logger.Printf("[Success]Prescription Insertion")
	return
}

func addToPrescMap(tbl, mapToColumn string, prescID int64, mapToElems []*models.MapObject) (err error) {
	if prescID == 0 {
		err = newQueryError("Zero prescription ID not allowed")
	} else {
		var sqlQuery sqlExeParams
		if len(mapToElems) == 0 {
			return
		}
		sqlQuery.Query = "INSERT IGNORE INTO " + tbl + " (prescription_id, " + mapToColumn + ", description) VALUES "
		for _, elem := range mapToElems {
			if elem.ID == 0 {
				err = newQueryError("Zero " + mapToColumn + " not allowed")
			}
			updateQueryListStringWithOperation(&sqlQuery.Query, fmt.Sprintf("(%v,%v,?)",
				prescID, elem.ID), ",")
			sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, elem.Description)
		}
		if _, _, updateErr := ExecDataUpdateQuery(sqlQuery.Query, sqlQuery.QueryArgs...); updateErr != nil {
			logger.Printf("[Query Error]Query:%v [err:%v]", sqlQuery.Query, updateErr)
		} else {
			logger.Printf("[Insertion complete ]Added %v[%v] to %v", mapToElems, prescID, tbl)
		}
	}
	return
}

/** Return value is not nil only in case of error in API **/
func addMedDataToPresc(param prescription.PutPrescriptionRegisterParams) middleware.Responder {
	err := addToPrescMap(prescToMedMap, "medicine_id", param.Prescription.ID, param.Prescription.Medicines)
	if err != nil {
		logger.Printf("[Error]In adding medicines to prescription:%v", err)
		return prescription.NewPutPrescriptionRegisterDefault(500).WithPayload(
			"internal server error:In adding medicines to prescription")
	}
	err = addToPrescMap(prescToAdvMap, "advice_id", param.Prescription.ID, param.Prescription.Advices)
	if err != nil {
		logger.Printf("[Error]In adding advices to prescription:%v", err)
		return prescription.NewPutPrescriptionRegisterDefault(500).WithPayload(
			"internal server error:In adding medical-advices to prescription")
	}
	err = addToPrescMap(prescToTestMap, "test_id", param.Prescription.ID, param.Prescription.MedicalTests)
	if err != nil {
		logger.Printf("[Error]In adding medical tests to prescription:%v", err)
		return prescription.NewPutPrescriptionRegisterDefault(500).WithPayload(
			"internal server error:In adding medical-advices to prescription")
	}
	return nil
}

/** /prescription/register handler **/
func RegisterPrescriptionAPI(param prescription.PutPrescriptionRegisterParams,
	p *models.Principal) middleware.Responder {
	var err error
	var httpStatusCode int
	if param.Prescription.ID == 0 {
		param.Prescription.ID, httpStatusCode, err = insertPrescription(param.Prescription)
		if err != nil {
			logger.Printf("[Error]In inserting prescription:%v", err)
			return prescription.NewPutPrescriptionRegisterDefault(httpStatusCode).WithPayload(models.Error(err.Error()))
		}
		if errResp := addMedDataToPresc(param); errResp != nil {
			logger.Printf("[Error]In adding med-data to prescription. [req:%v]", *param.Prescription)
			return errResp
		}
		logger.Printf("[Checkpoint] Prescription[%v] saved", param.Prescription.ID)
	}
	if param.Prescription.AppointmentID != 0 {
		if param.Prescription.Name != "" {
			param.Prescription.Name = ""
			param.Prescription.ID = 0
			param.Prescription.ID, httpStatusCode, err = insertPrescription(param.Prescription)
			if err != nil {
				logger.Printf("[Error]In inserting name-less prescription:%v", err)
				return prescription.NewPutPrescriptionRegisterDefault(httpStatusCode).WithPayload(models.Error(err.Error()))
			}
			if errResp := addMedDataToPresc(param); errResp != nil {
				logger.Printf("[Error]In adding name-less med-data to prescription. [req:%v]", *param.Prescription)
				return errResp
			}
			logger.Printf("[Checkpoint] Nameless prescription[%v] saved", param.Prescription.ID)
		}
		/**inv: presc.name == "": Only name less prescriptions should be added to appointments
			So that prescription history can't be altered
		**/
		updateAptReq := appointment.NewPostAppointmentUpdateParams()
		info := &models.AppointmentInfo{PrescriptionID: param.Prescription.ID,
			ID: param.Prescription.AppointmentID}
		updateAptReq.Info = info
		logger.Printf("[Checkpoint]Saving prescription")
		responder := UpdateAppointmentAPI(updateAptReq, nil)
		switch resp := responder.(type) {
		case *appointment.PostAppointmentUpdateDefault:
			logger.Printf("[Error]In attaching prescription[%v] to appointment[%v]",
				param.Prescription.ID, param.Prescription.AppointmentID)
			return prescription.NewPutPrescriptionRegisterDefault(500).WithPayload(
				"internal server error:in adding prescription to appointment")
		case *appointment.PostAppointmentUpdateOK:
			logger.Printf("[Checkpoint] Prescription[%v] attached to appointment[%v]",
				param.Prescription.ID, param.Prescription.AppointmentID)
		default:
			logger.Printf("[Error] unknow response type %T", resp)
			return prescription.NewPutPrescriptionRegisterDefault(500).WithPayload(
				"internal server error:try again later")
		}
	}
	return prescription.NewPutPrescriptionRegisterOK().WithPayload(&models.PassedRegInfo{
		ID: param.Prescription.ID})
}

/** /prescription/find handler **/
func fetchFindPrescData(param prescription.GetPrescriptionFindParams) (prescriptions []*models.Prescription, suggestions []*models.MapObject, err error) {
	query := ""
	queryFilter := ""
	queryFilterArgs := []any{}
	if (param.ID != nil) && (*param.ID != 0) {
		updateQueryListString(&queryFilter, "id", "AND")
		queryFilterArgs = append(queryFilterArgs, *param.ID)
	}
	if (param.CreatedBy != nil) && (*param.CreatedBy != 0) {
		updateQueryListString(&queryFilter, "created_by", "AND")
		queryFilterArgs = append(queryFilterArgs, *param.CreatedBy)
	}
	if (param.Search != nil) && (*param.Search != "") {
		updateQueryListStringWithOperation(&queryFilter,
			"MATCH(name, description, comment_on_medicines, comment_on_tests, comment_on_advices) AGAINST (?)", "AND")
		queryFilterArgs = append(queryFilterArgs, *param.Search)
	}
	if queryFilter != "" {
		queryFilter = " WHERE " + queryFilter
	}
	if (param.ProvideSuggestion != nil) && (*param.ProvideSuggestion) {
		query = "SELECT id, IFNULL(name, CONCAT('Presc. from:', DATE(created_on))), description FROM " + prescriptionTbl + queryFilter

	} else {
		query = "SELECT id, comment_on_medicines, comment_on_tests, comment_on_advices, DATE(last_updated) AS DATE," +
			"IFNULL(name, ''), description FROM " + prescriptionTbl + queryFilter
	}
	ctx, cancel := getTimeOutContext()
	defer cancel()
	if rows, fetchErr := ExecDataFetchQuery(ctx, query, queryFilterArgs...); fetchErr != nil {
		logger.Printf("[Error: Fetching] In find prescription.err:%v{Query:%v[args:%v]}", fetchErr, query, queryFilterArgs)
		return prescriptions, suggestions, fetchErr
	} else {
		for rows.Next() {
			if (param.ProvideSuggestion != nil) && (*param.ProvideSuggestion) {
				suggestion := &models.MapObject{}
				if err = rows.Scan(&suggestion.ID, &suggestion.Name, &suggestion.Description); err != nil {
					logger.Printf("[Error: Fetching] In find suggestion, scan err:%v{Query:%v[args:%v]}", err, query, queryFilterArgs)
					return
				}
				suggestions = append(suggestions, suggestion)
			} else {
				prescription := &models.Prescription{}
				if err = rows.Scan(&prescription.ID, &prescription.CommentsOnMeds, &prescription.CommentsOnTests, &prescription.OverallAdvice,
					&prescription.Date, &prescription.Name, &prescription.Description); err != nil {
					logger.Printf("[Error: Fetching] In find prescription, scan err:%v{Query:%v[args:%v]}", err, query, queryFilterArgs)
					return
				}
				prescriptions = append(prescriptions, prescription)
			}
		}
		return
	}
}

func FetchPrescriptionRelatedData(mapFromTbl, mapToTbl, mapFromColumn string, id int64) ([]*models.MapObject, error) {
	query := fmt.Sprintf("SELECT mapTo.id as id, mapTo.name as name, mapFrom.description as description FROM %v mapFrom, %v mapTo "+
		"WHERE (mapFrom.prescription_id = %v) AND  (mapFrom.%v = mapTo.id)", mapFromTbl, mapToTbl, id, mapFromColumn)
	ctx, onCancel := getTimeOutContext()
	defer onCancel()
	resp := []*models.MapObject{}
	if rows, exeErr := ExecDataFetchQuery(ctx, query); exeErr != nil {
		logger.Printf("[Error]Fetching prescription related data:%v[query:%v]| %v-%v[%v-%v]", exeErr, query, mapFromColumn, mapToTbl, id, mapFromColumn)
		return resp, exeErr
	} else {
		for rows.Next() {
			data := &models.MapObject{}
			if scanErr := rows.Scan(&data.ID, &data.Name, &data.Description); scanErr != nil {
				logger.Printf("[Error] Error in scanning:%v,query:%v", scanErr, query)
				return resp, scanErr
			} else {
				resp = append(resp, data)
			}
		}
	}
	return resp, nil
}

func FindPrescriptionAPI(param prescription.GetPrescriptionFindParams, p *models.Principal) middleware.Responder {
	if prescriptions, suggestions, err := fetchFindPrescData(param); err != nil {
		logger.Printf("[Error]In find prescription data API:%v", err)
		return prescription.NewGetPrescriptionFindDefault(500).WithPayload("Internal server error:In reading prescription")
	} else {
		for _, resp := range prescriptions {
			resp.Medicines, err = FetchPrescriptionRelatedData(prescToMedMap, medicineTbl, "medicine_id", resp.ID)
			if err != nil {
				logger.Printf("[Error]In find prescriptions, fetching medicines:%v", err)
				return prescription.NewGetPrescriptionFindDefault(500).WithPayload("Internal server error:In reading prescription medicines")
			}
			resp.MedicalTests, err = FetchPrescriptionRelatedData(prescToTestMap, testTbl, "test_id", resp.ID)
			if err != nil {
				logger.Printf("[Error]In find prescriptions, fetching medical-tests:%v", err)
				return prescription.NewGetPrescriptionFindDefault(500).WithPayload("Internal server error:In reading prescription medical-tests")
			}
			resp.Advices, err = FetchPrescriptionRelatedData(prescToAdvMap, adviceTbl, "advice_id", resp.ID)
			if err != nil {
				logger.Printf("[Error]In find prescriptions, fetching advices:%v", err)
				return prescription.NewGetPrescriptionFindDefault(500).WithPayload("Internal server error:In reading prescription advices")
			}
		}
		logger.Printf("[Checkpoint]Success find prescription")
		return prescription.NewGetPrescriptionFindOK().WithPayload(&prescription.GetPrescriptionFindOKBody{
			Prescriptions: prescriptions, Suggestions: suggestions})
	}
}
