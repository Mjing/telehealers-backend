/** APIs related to doctor data **/
package apis

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/doctor"
)

var (
	insertDocQuery = "INSERT INTO " + doctorTbl + " (%v) VALUES (%v)"
	updateDocQuery = "UPDATE " + doctorTbl + " SET %v WHERE %v"
	deleteDocQuery = "DELETE FROM " + doctorTbl + " WHERE %v"
	findDocQuery   = "SELECT id, name, email, phone, about, profile_picture FROM " +
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
	query = fmt.Sprintf(insertDocQuery, columns, values)
	return
}

/** Main function to register doctor into the system.
TODO: Create register process: Apply->Verify->Approve
**/
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

/** Function to create doctor update query.
query,queryArgs are to be used together in exec-function **/
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
		logger.Printf("[Error]db error:%v", err)
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
			&docData.Phone, &docData.About, &docData.ProfilePictureID); scanErr != nil {
			logger.Printf("[Error]doctor data scan error:%v", scanErr)
			return doctor.NewGetDoctorFindDefault(500).WithPayload("internal db error in data read")
		}
		foundDocs.Doctors = append(foundDocs.Doctors, docData)
	}
	logger.Printf("[Success]find game api")
	return doctor.NewGetDoctorFindOK().WithPayload(foundDocs)
}
