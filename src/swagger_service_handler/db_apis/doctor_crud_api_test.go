package apis

import (
	"fmt"
	"strings"
	"testing"

	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/doctor"
)

var (
	//Test Doctors
	docOc = &models.DoctorInfo{Name: "Dr. Oto Octavius", Email: "8@sinister6",
		Phone: "8888666666"}
	drFanzWorth = &models.DoctorInfo{Name: "Dr. Huebert Farnzworth", Email: "Fonzi@planex",
		Phone: "01019124334532"}
	drDOOM  = &models.DoctorInfo{Name: "MF DOOM", Email: "DOOM@allcaps", Phone: "6666666666"}
	doctors = []*models.DoctorInfo{docOc, drFanzWorth, drDOOM}
)

func TestSetupDbConnection(t *testing.T) {
	dbUser = "root"
	dbName = "telehealers"
	dbAddr = "localhost:3306"
	connErr := InitConnection()
	if connErr != nil {
		t.Errorf("Check test-env vars.")
	}
}

func handledErrorMessage(errMsg string) bool {
	return strings.Contains(errMsg, "already") ||
		strings.Contains(errMsg, queryErrorTag)
}

func TestRegisterDoctor(t *testing.T) {
	TestSetupDbConnection(t)
	for _, doc := range doctors {
		t.Logf("Adding:%v", doc)
		fmt.Printf("xxxxxxxxx")
		req := doctor.PutDoctorRegisterParams{Info: doc}
		responder := RegisterDoctor(req)
		switch resp := responder.(type) {
		case *doctor.PutDoctorRegisterDefault:
			if handledErrorMessage(string(resp.Payload)) {
				t.Log("[Handled fail] duplicate entry")
				continue
			}
			t.Errorf("Unhandled error:%v", resp.Payload)
		case *doctor.PutDoctorRegisterOK:
			t.Log("ok response")
			continue
		default:
			t.Errorf("Unhandle error for doc:%v| response type:%T | resp:%v", doc, resp, resp)
		}
	}
}

func TestUpdateDoctor(t *testing.T) {
	TestSetupDbConnection(t)
	for _, doc := range doctors {
		req := doctor.PostDoctorUpdateParams{Info: doc}
		switch resp := UpdateDoctor(req).(type) {
		case *doctor.PostDoctorUpdateDefault:
			if handledErrorMessage(string(resp.Payload)) {
				t.Log("[Handled fail] Duplicate entry")
				continue
			}
			t.Errorf("Unhandled error:%v", resp.Payload)
		case *doctor.PostDoctorUpdateOK:
			continue
		default:
			t.Errorf("Unhandle error for doc:%v| response type:%T | resp:%v", doc, resp, resp)
		}
	}
}
func TestRemoveDoctor(t *testing.T) {
	TestSetupDbConnection(t)
	for _, doc := range doctors {
		req := doctor.DeleteDoctorRemoveParams{ID: doc.ID}
		switch resp := RemoveDoctor(req).(type) {
		case *doctor.DeleteDoctorRemoveDefault:
			if handledErrorMessage(string(resp.Payload)) {
				t.Log("[Handled fail] Duplicate entry")
				continue
			}
			t.Errorf("Unhandled error:%v", resp.Payload)
		case *doctor.DeleteDoctorRemoveOK:
			continue
		default:
			t.Errorf("Unhandle error for doc:%v| response type:%T | resp:%v", doc, resp, resp)
		}
	}
}

func TestFindDoctor(t *testing.T) {
	TestSetupDbConnection(t)
	for _, doc := range doctors {
		req := doctor.NewGetDoctorFindParams()
		req.NameContaining = &doc.Name
		var p, s = int64(1), int64(10)
		req.Page = &p
		req.Size = &s
		switch resp := FindDoctor(req).(type) {
		case *doctor.GetDoctorFindDefault:
			if handledErrorMessage(string(resp.Payload)) {
				t.Log("[Handled fail]")
				continue
			}
		case *doctor.GetDoctorFindOK:
			continue
		default:
			t.Errorf("Unhandle error for doc:%v| response type:%T | resp:%v", doc, resp, resp)
		}
	}
}
