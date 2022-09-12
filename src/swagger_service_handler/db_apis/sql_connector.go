package apis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	_ "github.com/go-sql-driver/mysql"
)

/**
Function to setup DB connection variables at start
dbNameV: Name of SQL database
dbUserV: SQL username
dbPassV: Password for SQL username
dbAddrV: Server address in host:ip format e.g. localhost:3306
**/
func SetConnectionVars(dbNameV, dbUserV, dbPassV, dbAddrV string) {
	dbName = dbNameV
	dbUser = dbUserV
	dbPass = dbPassV
	dbAddr = dbAddrV
}

/** To be called at the start of the application,
for setup and initialization of the package**/
func InitConnection() (err error) {
	if pool == nil {
		pool, err = sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v)/%v",
			dbUser, dbPass, dbAddr, dbName))
		if err != nil {
			logger.Fatalf("[%v]Unable to use source name:%v", logIDFlag, err)
			return err
		}
	}
	ctx, cancel := getTimeOutContext()
	defer cancel()
	err = pool.PingContext(ctx)
	if err != nil {
		logger.Fatalf("[%v]Unable to connect to database: %v", logIDFlag, err)
	}
	logger.Printf("[%v]DB Connected", logIDFlag)
	return
}

/** DB API Helper
	Both functions create context locally if need be add another paramerter
		to take context variable
**/

func UpdateAndRespond(data UpdateAPIs) middleware.Responder {
	queryParams, queryErr := data.makeQuery()
	if queryErr != nil {
		logger.Printf("[Query Error]query:%v|error:%v", queryParams.Query, queryErr)
		return data.errResponse(400, queryErr)
	}
	lastId, rowId, updateErr := ExecDataUpdateQuery(queryParams.Query, queryParams.QueryArgs...)
	if updateErr != nil {
		logger.Printf("[DB Update Error]query:%v|err:%v", queryParams.Query, updateErr)
		if duplicateEntryError(updateErr) {
			return data.errResponse(400, newQueryError("entry with same name already present"))
		}
		return data.errResponse(500, errors.New("internal db error"))
	}
	return data.okResponse(lastId, rowId)
}

func FetchAndRespond(data ReadAPIs) middleware.Responder {
	queryParams, queryErr := data.makeQuery()
	if queryErr != nil {
		logger.Printf("[Query Error]query:%v|error:%v", queryParams.Query, queryErr)
		return data.errResponse(400, queryErr)
	}
	ctx, cancel := getTimeOutContext()
	defer cancel()
	foundRows, fetchErr := ExecDataFetchQuery(ctx, queryParams.Query, queryParams.QueryArgs...)
	if fetchErr != nil {
		logger.Printf("Query[Fetch Error]query:%v|error:%v", fetchErr, queryParams.Query)
		return data.errResponse(500, errors.New("internal db read error"))
	}
	data.scanRows(foundRows)
	return data.okResponse(0, 0)
}

/** Execute Insert queries.
returns lastId inserted, rows affected, error
Data update here implies: Creation, Updation and Deletion
**/
func ExecDataUpdateQuery(query string, queryParams ...any) (int64, int64, error) {
	ctx, cancel := getTimeOutContext()
	defer cancel()
	result, err := pool.ExecContext(ctx, query, queryParams...)
	if err != nil {
		logger.Printf("[Error]In data Updater:%v", err)
		return 0, 0, err
	} else {
		lastId, _ := result.LastInsertId()
		rowsAffected, _ := result.RowsAffected()
		logger.Printf("[Success]data updated:%v:%v", lastId, rowsAffected)
		return lastId, rowsAffected, nil
	}
}

func ExecDataFetchQuery(ctx context.Context, query string, queryParams ...any) (*sql.Rows, error) {
	rows, err := pool.QueryContext(ctx, query, queryParams...)
	if err != nil {
		logger.Printf("[Error]In data fetch:%v", err)
		return nil, err
	} else {
		return rows, nil
	}
}

/*** SQL HELPER FUNCTION ***/
/*** Append base string with 'columnName = ?' with joinString as conjunction.
Provided to create SQL operation lists of form
	"col1 = ?, col2 = ?,..." or "col1 = ? AND col2 = ? AND..."
Some simple usage
("", "col", ",") == "col = ?"
("col1 = ?", "col2", ",") == "col1 = ? , col2 = ?"
("col1 > 2", "col2", "AND") == "col1 > 2 AND col2 = ?"
***/
func updateQueryListString(base *string, columnName, joinString string) {
	if *base != "" {
		*base += " " + joinString + " "
	}
	*base += columnName + " = ?"
}

/*** Append base string with operation with joinString as conjunction.
Some simple usage
("", "col", ",") == "col"
("col1 = ?", "col2 = ?", ",") == "col1 = ?, col2 = ?"
("col1 > 2", "col1 < 20", "AND") == "col1 > 2 AND col1 < 20"
***/
func updateQueryListStringWithOperation(base *string, operation, joinString string) {
	if (*base != "") && (operation != "") {
		*base += " " + joinString + " "
	}
	*base += operation
}
