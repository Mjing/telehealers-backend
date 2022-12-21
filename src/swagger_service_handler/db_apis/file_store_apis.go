package apis

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-openapi/runtime/middleware"
	"telehealers.in/router/models"
	fileStoreAPI "telehealers.in/router/restapi/operations/file"
)

func UploadAPI(param fileStoreAPI.PostFileUploadParams, p *models.Principal) middleware.Responder {
	var savePath string
	var fileID int64
	var newFile *os.File
	_, mpfh, mperr := param.HTTPRequest.FormFile("file")
	if mperr != nil {
		logger.Printf("[Error]In reading req file-name:%v", mperr)
		return fileStoreAPI.NewPostFileUploadDefault(500).WithPayload("internal server error:in reading file")
	}

	if sessionIDCookie, cookieErr := param.HTTPRequest.Cookie("th-ssid"); cookieErr != nil {
		logger.Printf("[Error]In reading cookie:%v", cookieErr)
		return fileStoreAPI.NewPostFileUploadDefault(400).WithPayload(
			"user login required:sessionid cookie absent in request")
	} else {
		if userID, userType, respCode, err := getLoginData(sessionIDCookie.Value); err != nil {
			logger.Printf("[Error]In fetching session-id:%v", err)
			return fileStoreAPI.NewPostFileUploadDefault(respCode).WithPayload(models.Error(err.Error()))
		} else {
			savePath = makeFileStoreSavePath(userType, userID)
			if mkdirErr := os.MkdirAll(savePath, os.ModePerm); mkdirErr != nil {
				logger.Printf("[Error]Unable to create dir[path:%v]:%v", savePath, mkdirErr)
				return fileStoreAPI.NewPostFileUploadDefault(500).WithPayload("internal server error")
			}
			var err error
			newFile, err = os.CreateTemp(savePath, "*"+mpfh.Filename)
			if err != nil {
				logger.Printf("[Error]Open image file:%v", err)
				return fileStoreAPI.NewPostFileUploadDefault(500).WithPayload("internal server error")
			}
			savePath = newFile.Name()
			if fileID, err = addNewFileToStoreDB(userID, userType, savePath); err != nil {
				logger.Printf("[Error]Insert file entry to db:%v", err)
				return fileStoreAPI.NewPostFileUploadDefault(500).WithPayload("internal db error")
			}
		}
	}
	imgBuffer := make([]byte, 1024)
	logger.Printf("[Checkpoint]Opened file:%v", savePath)
	for {
		bytesRead, readErr := param.File.Read(imgBuffer)
		if readErr == io.EOF {
			break
		} else if readErr != nil {
			logger.Printf("[Error]Reading image API req: %v", readErr)
			return fileStoreAPI.NewPostFileUploadDefault(403).WithPayload("Bad Image in request")
		}
		_, writeErr := newFile.Write(imgBuffer[:bytesRead])
		if writeErr != nil {
			logger.Printf("[Error]Writing img:%v", writeErr)
			return fileStoreAPI.NewPostFileUploadDefault(500).WithPayload("internal server error")
		}
	}
	if _, readErr := param.File.Read(imgBuffer); readErr != nil && readErr != io.EOF {
		logger.Printf("[Error]Reading image API req:%v", readErr)
		return fileStoreAPI.NewPostFileUploadDefault(403).WithPayload("Bad Image in request")
	}
	if closeErr := param.File.Close(); closeErr != nil {
		logger.Printf("[Error]In closing image:%v", closeErr)
		return fileStoreAPI.NewPostFileUploadDefault(500).WithPayload("internal server error")
	}
	logger.Printf("[Checkpoint]Image successfuly writen")
	return fileStoreAPI.NewPostFileUploadOK().WithPayload(&models.PassedRegInfo{ID: fileID})
}

// Returns respondable error, userID == 0 if data not found
func getLoginData(ssid string) (userID int64, userType string, httpStatusCode int, err error) {
	query := "SELECT user_type, user_id, status FROM " + sessionTbl + " WHERE session_id = ? ORDER BY last_login DESC LIMIT 1"
	ctx, cancel := getTimeOutContext()
	defer cancel()
	if rows, dbErr := ExecDataFetchQuery(ctx, query, ssid); dbErr != nil {
		logger.Printf("[Error]In fetching session data[ssid:%v]:%v", ssid, dbErr)
		err = dbErr
		httpStatusCode = 500
	} else {
		defer rows.Close()
		for rows.Next() {
			var userStatus string
			if scanErr := rows.Scan(&userType, &userID, &userStatus); scanErr != nil {
				logger.Printf("[Error] in scanning data for [ssid:%v]:%v", ssid, scanErr)
				err = errors.New("internal server error: in accessing session data")
				httpStatusCode = 500
			} else {
				httpStatusCode = 200
				// TODO: Correct
				// if userStatus != "online" {
				// 	logger.Printf("[Checkpoint]User not online[ssid:%v] status:%v", ssid, userStatus)
				// 	err = fmt.Errorf("user should be online")
				// 	httpStatusCode = 400
				// }
			}
		}
	}
	return
}

// returns respondable error
func getLoginDataFromSsidCookie(req *http.Request) (userID int64, userType string, httpRespCode int, err error) {
	if sessionIDCookie, cookieErr := req.Cookie("th-ssid"); cookieErr != nil {
		logger.Printf("[Error]In reading cookie:%v", cookieErr)
		return 0, "", 400, errors.New("user login required:sessionid cookie absent in request")
	} else {
		return getLoginData(sessionIDCookie.Value)
	}
}

// Returns save path of directory, ends with / character
func makeFileStoreSavePath(userType string, userID int64) string {
	return fmt.Sprintf("%v/%v/%v/", fileStoreRoot, userType, userID)
}

func addNewFileToStoreDB(userID int64, userType, filePath string) (int64, error) {
	query := "INSERT IGNORE INTO " + fileStoreTbl + " (user_id, user_type, path) VALUES (?,?,?) ON DUPLICATE KEY UPDATE user_id = ?, user_type = ?"
	if storeID, _, insertErr := ExecDataUpdateQuery(query, userID, userType, filePath, userID, userType); insertErr != nil {
		logger.Printf("[Error]In inserting file[data:uid:%v, utype:%v, path:%v]:%v", userID, userType,
			filePath, insertErr)
		return 0, errors.New("internal server error:uploading file data")
	} else {
		return storeID, nil
	}
}

/**** END OF /file/upload *****/

/**** handler for /file/download ****/
func getFilePathFromDB(fileID int64) (path string, err error) {
	query := "SELECT path FROM " + fileStoreTbl + " WHERE id = ?"
	ctx, cancel := getTimeOutContext()
	defer cancel()
	if rows, fetchErr := ExecDataFetchQuery(ctx, query, fileID); fetchErr != nil {
		logger.Printf("[Error]Error in fetching data[id:%v]:%v", fileID, fetchErr)
		err = errors.New("internal server error:in finding file")
	} else {
		rows.Next()
		if scanErr := rows.Scan(&path); scanErr != nil {
			logger.Printf("[Error]Error in scanning rows:%v", scanErr)
			err = errors.New("internal server error:in finding file")
		}
		rows.Close()
	}
	return
}

func DownloadAPI(param fileStoreAPI.GetFileDownloadParams, p *models.Principal) middleware.Responder {
	if path, err := getFilePathFromDB(param.ID); err != nil {
		logger.Printf("[Error]In download api:%v", err)
		return fileStoreAPI.NewGetFileDownloadDefault(500).WithPayload(models.Error(err.Error()))
	} else {
		if file, openErr := os.Open(path); openErr != nil {
			logger.Printf("[Error] In download api openning file:%v", openErr)
			return fileStoreAPI.NewGetFileDownloadDefault(500).WithPayload("internal server error:in reading fiel")
		} else {
			logger.Printf("[Checkpoint]File requested:%v", path)
			return fileStoreAPI.NewGetFileDownloadOK().WithPayload(file)
		}
	}
}

/**** END OF /file/download ****/
