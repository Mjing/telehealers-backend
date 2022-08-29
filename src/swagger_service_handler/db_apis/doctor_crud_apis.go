/** APIs related to doctor data **/
package apis

import (
	"errors"
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/doctor"
)

var insertDocQuery = "INSERT INTO " + doctorTbl + " (%v) VALUES (%v)"

func makeInsertDocQuery(docInfo *models.DoctorInfo) (string, error) {
	if docInfo.Name == "" || docInfo.Email == "" || docInfo.Phone == "" {
		return "", errors.New("doctor name, email or phone can't be empty")
	}
	columns := "name, email, phone"
	values := fmt.Sprintf("'%v','%v','%v'", docInfo.Name,
		docInfo.Email, docInfo.Phone)
	if docInfo.About != "" {
		columns += ", about"
		values += fmt.Sprintf(",'%v'", docInfo.About)
	}
	if docInfo.ProfilePicture != 0 {
		columns += ", profile_picture"
		values += fmt.Sprintf(", %v", docInfo.ProfilePicture)
	}
	return fmt.Sprintf(insertDocQuery, columns, values), nil
}

/** Main function to register doctor into the system.
TODO: Create register process: Apply->Verify->Approve
**/
func RegisterDoctor(param doctor.PutDoctorRegisterParams) middleware.Responder {
	query, queryErr := makeInsertDocQuery(param.Info)
	if queryErr != nil {
		logger.Printf("[Error]bad input:%v", queryErr.Error())
		return doctor.NewPutDoctorRegisterDefault(400).WithPayload(models.Error(queryErr.Error()))
	}
	if _, _, execErr := ExecDataUpdateQuery(query); execErr != nil {
		logger.Printf("[Error]doctor register db functionality:%v", execErr)
		if duplicateEntryError(execErr) {
			return doctor.NewPostDoctorUpdateDefault(400).WithPayload("Requested doctor already present")
		}
		return doctor.NewPutDoctorRegisterDefault(500).WithPayload("Internal error")
	}
	return doctor.NewPutDoctorRegisterOK()
}
