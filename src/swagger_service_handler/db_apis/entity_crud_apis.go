/*
* This file contains implementation of entity struct(currently: medicines, tests and adivces)
entities are all data containing name and description only.*
*/
package apis

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/read"
	"telehealers.in/router/restapi/operations/register"
	"telehealers.in/router/restapi/operations/remove"
	"telehealers.in/router/restapi/operations/update"
)

type insertEntityReq struct {
	req       *models.Entity
	createdBy int64
	table     *string
}

func (param *insertEntityReq) makeQuery() (sqlReq sqlExeParams, err error) {
	if param.req.Name == "" {
		err = newQueryError("non-empty name is required")
		return
	}
	columns, values := "name", "?"
	sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.req.Name)
	if param.req.Description != "" {
		columns += ",description"
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.req.Description)
		values += ",?"
	}
	sqlReq.Query = fmt.Sprintf(generalInsertQuery, *param.table, columns) + " (" + values + ")"
	logger.Printf("[Checkpoint] Insert Entity Query:%v [ARGS:%v]", sqlReq.Query, sqlReq.QueryArgs)
	return
}

type updateEntityReq struct {
	req   *models.Entity
	table *string
}

func (param *updateEntityReq) makeQuery() (sqlReq sqlExeParams, err error) {
	if param.req.ID == 0 {
		err = newQueryError("positive id required")
		return
	}
	if (param.req.Name == "") && (param.req.Description == "") {
		err = newQueryError("one of name or description is required")
		return
	}
	setter := ""
	if param.req.Name != "" {
		updateQueryListString(&setter, "name", ",")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.req.Name)
	}
	if param.req.Description != "" {
		updateQueryListString(&setter, "description", ",")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.req.Description)
	}
	sqlReq.Query = fmt.Sprintf(generalUpdateQuery, *param.table, setter, "id = ?")
	sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.req.ID)
	return
}

type deleteEntityReq struct {
	deleteID int64
	table    *string
}

func (param *deleteEntityReq) makeQuery() (sqlReq sqlExeParams, err error) {
	if param.deleteID == 0 {
		err = newQueryError("non zero id is required")
		return
	}
	sqlReq.Query = fmt.Sprintf(generalDeleteQuery, *param.table, "id = ?")
	sqlReq.QueryArgs = append(sqlReq.QueryArgs, param.deleteID)
	return
}

type findEntityReq struct {
	id           *int64
	substrInName *string
	page         int64
	pageSize     int64
	table        *string
}

func (req *findEntityReq) makeQuery() (sqlReq sqlExeParams, err error) {
	if (req.page <= 0) || (req.pageSize <= 0) {
		err = newQueryError("positive(non-zero) page & page_size is needed")
		return
	}
	queryFilter := ""
	if (req.id != nil) && (*req.id != 0) {
		updateQueryListString(&queryFilter, "id", " AND ")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, req.id)
	}
	if (req.substrInName != nil) && (*req.substrInName != "") {
		updateQueryListStringWithOperation(&queryFilter,
			"LOWER(name) LIKE CONCAT('%', CONCAT(LOWER(?),'%'))", " AND ")
		sqlReq.QueryArgs = append(sqlReq.QueryArgs, *req.substrInName)
	}
	if queryFilter != "" {
		queryFilter = " WHERE " + queryFilter
	}
	sqlReq.Query = fmt.Sprintf(generalFetchQuery, "id,name,description",
		*req.table+queryFilter,
		(req.page-1)*req.pageSize,
		req.pageSize)
	logger.Printf("[Checkpoint]Find entity query:%v [ARGS:%v]", sqlReq.Query, sqlReq.QueryArgs)
	return
}

/*** Medicine Cruds ***/
//Insert
type insertMedicineReq struct {
	insertEntityReq
}

func (param *insertMedicineReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Register medicine API]Error:code:%v|error:%v",
		httpCode, err)
	return register.NewPutMedicineRegisterDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (*insertMedicineReq) okResponse(id int64, affectedRows int64) middleware.Responder {
	return register.NewPutMedicineRegisterOK().WithPayload(&models.PassedRegInfo{ID: id})
}

func RegisterMedicineAPI(param register.PutMedicineRegisterParams, p *models.Principal) middleware.Responder {
	req := &insertMedicineReq{insertEntityReq{req: param.Info.Data, createdBy: param.Info.CreatedBy, table: &medicineTbl}}
	return UpdateAndRespond(req)
}

// Update
type updateMedicineReq struct {
	updateEntityReq
}

func (param *updateMedicineReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Update medicine API]Error:code:%v|error:%v",
		httpCode, err)
	return update.NewPostMedicineUpdateDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (*updateMedicineReq) okResponse(int64, int64) middleware.Responder {
	return update.NewPostMedicineUpdateOK()
}

func UpdateMedicineAPI(param update.PostMedicineUpdateParams, p *models.Principal) middleware.Responder {
	req := &updateMedicineReq{updateEntityReq{req: param.Info, table: &medicineTbl}}
	return UpdateAndRespond(req)
}

//Remove

type removeMedReq struct {
	deleteEntityReq
}

func (param *removeMedReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Remove medicine API]Error:code:%v|error:%v",
		httpCode, err)
	return remove.NewDeleteMedicineRemoveDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (*removeMedReq) okResponse(int64, int64) middleware.Responder {
	return remove.NewDeleteMedicineRemoveOK()
}

func RemoveMedicineAPI(param remove.DeleteMedicineRemoveParams, p *models.Principal) middleware.Responder {
	req := &removeMedReq{deleteEntityReq{deleteID: param.ID, table: &medicineTbl}}
	return UpdateAndRespond(req)
}

//Find

type findMedReq struct {
	findEntityReq
	fetchedData *read.GetMedicineFindOKBody
}

func (req *findMedReq) scanRows(rows *sql.Rows) (err error) {
	req.fetchedData = &read.GetMedicineFindOKBody{}
	for rows.Next() {
		row := &models.Entity{}
		if scanErr := rows.Scan(&row.ID, &row.Name, &row.Description); scanErr != nil {
			logger.Printf("[Find medicine api]Error:%v", scanErr)
			err = errors.New("internal db read error")
			break
		}
		req.fetchedData.Data = append(req.fetchedData.Data, row)
	}
	return
}

func (param *findMedReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Find medicine API]Error:code:%v|error:%v",
		httpCode, err)
	return read.NewGetMedicineFindDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (req *findMedReq) okResponse(int64, int64) middleware.Responder {
	return read.NewGetMedicineFindOK().WithPayload(req.fetchedData)
}

func FindMedicineAPI(param read.GetMedicineFindParams, p *models.Principal) middleware.Responder {
	logger.Printf("asd %v", *param.NameContaining)
	req := &findMedReq{findEntityReq: findEntityReq{table: &medicineTbl,
		id: param.ID, substrInName: param.NameContaining, page: *param.Page, pageSize: *param.PageSize}}
	return FetchAndRespond(req)
}

/***************************END MEDICINE APIs*******************************/

/*** Medical advices CRUDs ***/
//Insert
type insertAdviceReq struct {
	insertEntityReq
}

func (param *insertAdviceReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Register advice API]Error:code:%v|error:%v",
		httpCode, err)
	return register.NewPutMedicalAdviceRegisterDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (*insertAdviceReq) okResponse(id int64, ar int64) middleware.Responder {
	return register.NewPutMedicalAdviceRegisterOK().WithPayload(&models.PassedRegInfo{ID: id})
}

func RegisterAdviceAPI(param register.PutMedicalAdviceRegisterParams, p *models.Principal) middleware.Responder {
	req := &insertAdviceReq{insertEntityReq{req: param.Info.Data, createdBy: param.Info.CreatedBy, table: &adviceTbl}}
	return UpdateAndRespond(req)
}

// Update
type updateAdviceReq struct {
	updateEntityReq
}

func (param *updateAdviceReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Update advice API]Error:code:%v|error:%v",
		httpCode, err)
	return update.NewPostMedicalAdviceUpdateDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (*updateAdviceReq) okResponse(int64, int64) middleware.Responder {
	return update.NewPostMedicalAdviceUpdateOK()
}

func UpdateAdviceAPI(param update.PostMedicalAdviceUpdateParams, p *models.Principal) middleware.Responder {
	req := &updateAdviceReq{updateEntityReq{req: param.Info, table: &adviceTbl}}
	return UpdateAndRespond(req)
}

//Remove

type removeAdviceReq struct {
	deleteEntityReq
}

func (param *removeAdviceReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Remove advice API]Error:code:%v|error:%v",
		httpCode, err)
	return remove.NewDeleteMedicalAdviceRemoveDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (*removeAdviceReq) okResponse(int64, int64) middleware.Responder {
	return remove.NewDeleteMedicalAdviceRemoveOK()
}

func RemoveAdviceAPI(param remove.DeleteMedicalAdviceRemoveParams, p *models.Principal) middleware.Responder {
	req := &removeAdviceReq{deleteEntityReq{deleteID: param.ID, table: &adviceTbl}}
	return UpdateAndRespond(req)
}

//Find

type findAdviceReq struct {
	findEntityReq
	fetchedData *read.GetMedicalAdviceFindOKBody
}

func (req *findAdviceReq) scanRows(rows *sql.Rows) (err error) {
	req.fetchedData = &read.GetMedicalAdviceFindOKBody{}
	for rows.Next() {
		row := &models.Entity{}
		if scanErr := rows.Scan(&row.ID, &row.Name, &row.Description); scanErr != nil {
			logger.Printf("[Find advice api]Error:%v", scanErr)
			err = errors.New("internal db read error")
			break
		}
		req.fetchedData.Data = append(req.fetchedData.Data, row)
	}
	return
}

func (param *findAdviceReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Find advice API]Error:code:%v|error:%v",
		httpCode, err)
	return read.NewGetMedicalAdviceFindDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (req *findAdviceReq) okResponse(int64, int64) middleware.Responder {
	return read.NewGetMedicalAdviceFindOK().WithPayload(req.fetchedData)
}

func FindAdviceAPI(param read.GetMedicalAdviceFindParams, p *models.Principal) middleware.Responder {
	req := &findAdviceReq{findEntityReq: findEntityReq{table: &adviceTbl,
		id: param.ID, substrInName: param.NameContaining, page: *param.Page, pageSize: *param.PageSize}}
	return FetchAndRespond(req)
}

/***************************END MEDICAL Advice APIs*******************************/

/*** Medical test CRUDs ***/
//Insert
type insertMedTestReq struct {
	insertEntityReq
}

func (param *insertMedTestReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Register test API]Error:code:%v|error:%v",
		httpCode, err)
	return register.NewPutMedicalTestRegisterDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (*insertMedTestReq) okResponse(id int64, ar int64) middleware.Responder {
	return register.NewPutMedicalTestRegisterOK().WithPayload(&models.PassedRegInfo{ID: id})
}

func RegisterTestAPI(param register.PutMedicalTestRegisterParams, p *models.Principal) middleware.Responder {
	req := &insertMedTestReq{insertEntityReq{req: param.Info.Data, createdBy: param.Info.CreatedBy, table: &testTbl}}
	return UpdateAndRespond(req)
}

// Update
type updateMedTestReq struct {
	updateEntityReq
}

func (param *updateMedTestReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Update test API]Error:code:%v|error:%v",
		httpCode, err)
	return update.NewPostMedicalTestUpdateDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (*updateMedTestReq) okResponse(int64, int64) middleware.Responder {
	return update.NewPostMedicalTestUpdateOK()
}

func UpdateMedTestAPI(param update.PostMedicalTestUpdateParams, p *models.Principal) middleware.Responder {
	req := &updateMedTestReq{updateEntityReq{req: param.Info, table: &testTbl}}
	return UpdateAndRespond(req)
}

//Remove

type removeMedTestReq struct {
	deleteEntityReq
}

func (param *removeMedTestReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Remove test API]Error:code:%v|error:%v",
		httpCode, err)
	return remove.NewDeleteMedicalTestRemoveDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (*removeMedTestReq) okResponse(int64, int64) middleware.Responder {
	return remove.NewDeleteMedicalTestRemoveOK()
}

func RemoveMedTestAPI(param remove.DeleteMedicalTestRemoveParams, p *models.Principal) middleware.Responder {
	req := &removeMedTestReq{deleteEntityReq{deleteID: param.ID, table: &testTbl}}
	return UpdateAndRespond(req)
}

//Find

type findMedTestReq struct {
	findEntityReq
	fetchedData *read.GetMedicalTestFindOKBody
}

func (req *findMedTestReq) scanRows(rows *sql.Rows) (err error) {
	req.fetchedData = &read.GetMedicalTestFindOKBody{}
	for rows.Next() {
		row := &models.Entity{}
		if scanErr := rows.Scan(&row.ID, &row.Name, &row.Description); scanErr != nil {
			logger.Printf("[Find medical test api]Error:%v", scanErr)
			err = errors.New("internal db read error")
			break
		}
		req.fetchedData.Data = append(req.fetchedData.Data, row)
	}
	return
}

func (param *findMedTestReq) errResponse(httpCode int, err error) middleware.Responder {
	logger.Printf("[Find medical test API]Error:code:%v|error:%v",
		httpCode, err)
	return read.NewGetMedicalTestFindDefault(httpCode).WithPayload(
		models.Error(err.Error()))
}

func (req *findMedTestReq) okResponse(int64, int64) middleware.Responder {
	return read.NewGetMedicalTestFindOK().WithPayload(req.fetchedData)
}

func FindMedTestAPI(param read.GetMedicalTestFindParams, p *models.Principal) middleware.Responder {
	req := &findMedTestReq{findEntityReq: findEntityReq{table: &testTbl,
		id: param.ID, substrInName: param.NameContaining, page: *param.Page, pageSize: *param.PageSize}}
	return FetchAndRespond(req)
}

/***************************END MEDICINE Test APIs*******************************/
