package apis

import (
	"errors"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"telehealers.in/router/models"
	hd "telehealers.in/router/restapi/operations/helpdesk"
)

/** Handler for /helpdesk/ticket/open **/
func makeOpenTicketQuery(param hd.PutHelpdeskTicketOpenParams, userID int64, userType string) (sqlQuery sqlExeParams) {
	sqlQuery.Query = "INSERT INTO " + helpdeskTbl + " (type, status, description, created_by, creator_type) " +
		"VALUES (?,?,?,?,?)"
	sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, param.Ticket.TicketType, param.Ticket.Status,
		param.Ticket.Description, userID, userType)
	return
}

func HelpdeskTicketOpenAPI(param hd.PutHelpdeskTicketOpenParams, p *models.Principal) middleware.Responder {
	userID, userType, respLoginCode, respLoginErr := getLoginDataFromSsidCookie(param.HTTPRequest)
	if respLoginErr != nil {
		logger.Printf("[Error]In fetching session-id:%v", respLoginErr)
		return hd.NewGetHelpdeskTicketFindCountDefault(respLoginCode).WithPayload(models.Error(respLoginErr.Error()))
	}
	sqlQuery := makeOpenTicketQuery(param, userID, userType)
	ticketID, _, exeErr := ExecDataUpdateQuery(sqlQuery.Query, sqlQuery.QueryArgs...)
	if exeErr != nil {
		logger.Printf("[Error] In executing ticket query:%v", exeErr)
		return hd.NewPutHelpdeskTicketOpenDefault(500).WithPayload("internal server error")
	}
	return hd.NewPutHelpdeskTicketOpenOK().WithPayload(&models.PassedRegInfo{ID: ticketID})
}

/** Handler for /helpdesk/ticket/update **/
func makeUpdateTicketQuery(param hd.PostHelpdeskTicketUpdateParams) (sqlQuery sqlExeParams, err error) {
	sqlQuery.Query = "UPDATE " + helpdeskTbl + " SET "
	setCommands := ""
	if param.Ticket.Status != "" {
		setCommands = "status = ?"
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, param.Ticket.Status)
	}
	if param.Ticket.Description != "" {
		updateQueryListString(&setCommands, "description", ",")
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, param.Ticket.Description)
	}
	if setCommands == "" {
		logger.Printf("[Query Error] No set fields provided")
		err = errors.New("non-empty ticket data(other that id) required")
	} else {
		sqlQuery.Query += setCommands + " WHERE id = ?"
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, param.Ticket.ID)
	}
	return
}

func HelpdeskTicketUpdateAPI(param hd.PostHelpdeskTicketUpdateParams, p *models.Principal) middleware.Responder {
	if param.Ticket.ID != 0 {
		logger.Printf("[Error]Empty id in update-ticket")
		return hd.NewPostHelpdeskTicketUpdateDefault(400).WithPayload("id required")
	}
	sqlQuery, queryErr := makeUpdateTicketQuery(param)
	if queryErr != nil {
		logger.Printf("[Error]Bad request:%v", queryErr)
		return hd.NewPostHelpdeskTicketUpdateDefault(400).WithPayload(models.Error(queryErr.Error()))
	}
	_, _, exeErr := ExecDataUpdateQuery(sqlQuery.Query, sqlQuery.QueryArgs...)
	if exeErr != nil {
		logger.Printf("[Error]DB err in update ticket:%v", exeErr)
		return hd.NewPutHelpdeskTicketOpenDefault(500).WithPayload("internal db error")
	}
	return hd.NewPostHelpdeskTicketUpdateOK()
}

/** /helpdesk/ticket/find Handler **/
func makeHelpDeskTableFilter(userID int64, userType string, reqUserId *int64, reqUserType, fromDate, toDate *string) (sqlQuery sqlExeParams) {
	if userType != "admin" {
		sqlQuery.Query = fmt.Sprintf("created_by = ? AND creator_type = ? ")
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, userID, userType)
	}
	if (reqUserId != nil) && (*reqUserId != 0) {
		updateQueryListString(&sqlQuery.Query, "created_by", "AND")
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, *reqUserId)
	}
	if (reqUserType != nil) && (*reqUserType != "") {
		updateQueryListString(&sqlQuery.Query, "creator_type", "AND")
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, *reqUserType)
	}
	if (fromDate != nil) && (*fromDate != "") {
		updateQueryListStringWithOperation(&sqlQuery.Query, "last_updated >= ?", "AND")
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, *fromDate)
	}
	if (toDate != nil) && (*toDate != "") {
		updateQueryListStringWithOperation(&sqlQuery.Query, "last_updated <= ?", "AND")
		sqlQuery.QueryArgs = append(sqlQuery.QueryArgs, *toDate)
	}
	return
}

func HelpdeskTicketFindAPI(param hd.GetHelpdeskTicketFindParams, p *models.Principal) middleware.Responder {
	userID, userType, respLoginCode, respLoginErr := getLoginDataFromSsidCookie(param.HTTPRequest)
	if respLoginErr != nil {
		logger.Printf("[Error]fetching session data err:%v", respLoginErr)
		return hd.NewGetHelpdeskTicketFindDefault(respLoginCode).WithPayload(
			models.Error(respLoginErr.Error()))
	}
	sqlQuery := makeHelpDeskTableFilter(userID, userType, param.UserID, param.UserType, param.FromDate, param.ToDate)

	queryTemplate := "SELECT id, type, status, description FROM " + helpdeskTbl + " %v LIMIT %v,%v"
	if sqlQuery.Query != "" {
		sqlQuery.Query = "WHERE " + sqlQuery.Query
	}
	sqlQuery.Query = fmt.Sprintf(queryTemplate, sqlQuery.Query, (*param.Page-1)*(*param.PageSize), *param.PageSize)
	ctx, cancel := getTimeOutContext()
	defer cancel()
	rows, dbErr := ExecDataFetchQuery(ctx, sqlQuery.Query, sqlQuery.QueryArgs...)
	if dbErr != nil {
		logger.Printf("[Error]Helpdesk ticket find api, db error:%v", dbErr)
		return hd.NewGetHelpdeskTicketFindDefault(500).WithPayload("internal db error")
	}
	resp := &hd.GetHelpdeskTicketFindOKBody{}
	for rows.Next() {
		rowData := &models.Ticket{}
		if scanErr := rows.Scan(&rowData.ID, &rowData.TicketType, &rowData.Status, &rowData.Description); scanErr != nil {
			logger.Printf("[Error]Ticket scanning error:%v", scanErr)
			return hd.NewGetHelpdeskTicketFindDefault(500).WithPayload("internal server error")
		}
		resp.Tickets = append(resp.Tickets, rowData)
	}
	return hd.NewGetHelpdeskTicketFindOK().WithPayload(resp)
}

func HelpdeskTicketFindCountAPI(param hd.GetHelpdeskTicketFindCountParams, p *models.Principal) middleware.Responder {
	userID, userType, respLoginCode, respLoginErr := getLoginDataFromSsidCookie(param.HTTPRequest)
	if respLoginErr != nil {
		logger.Printf("[Error]fetching session data err:%v", respLoginErr)
		return hd.NewGetHelpdeskTicketFindCountDefault(respLoginCode).WithPayload(
			models.Error(respLoginErr.Error()))
	}
	sqlQuery := makeHelpDeskTableFilter(userID, userType, param.UserID, param.UserType, param.FromDate, param.ToDate)

	query := "SELECT COUNT(*) FROM " + helpdeskTbl
	if sqlQuery.Query != "" {
		sqlQuery.Query = query + " WHERE " + sqlQuery.Query
	}
	sqlQuery.Query = query
	ctx, cancel := getTimeOutContext()
	defer cancel()
	rows, dbErr := ExecDataFetchQuery(ctx, sqlQuery.Query, sqlQuery.QueryArgs...)
	if dbErr != nil {
		logger.Printf("[Error]Helpdesk ticket find api, db error:%v", dbErr)
		return hd.NewGetHelpdeskTicketFindDefault(500).WithPayload("internal db error")
	}
	resp := &hd.GetHelpdeskTicketFindCountOKBody{}
	for rows.Next() {
		if scanErr := rows.Scan(&resp.Count); scanErr != nil {
			logger.Printf("[Error]Ticket scanning error:%v", scanErr)
			return hd.NewGetHelpdeskTicketFindCountDefault(500).WithPayload("internal server error")
		}
	}
	return hd.NewGetHelpdeskTicketFindCountOK().WithPayload(resp)
}
