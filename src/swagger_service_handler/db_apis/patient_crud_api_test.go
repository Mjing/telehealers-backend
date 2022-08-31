package apis

import (
	"testing"

	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/patient"
)

var (
	//Test patients
	ed = &models.PatientInfo{Name: "Ed", Email: "ed@ededdneddy",
		Phone: "8888666666"}
	edd = &models.PatientInfo{Name: "Edd", Email: "edd@ededdneddy",
		Phone: "01019124334532"}
	eddy     = &models.PatientInfo{Name: "Eddy", Email: "eddy@ededdneddy", Phone: "6666666666"}
	patients = []*models.PatientInfo{ed, edd, eddy}
)

func TestRegisterPatient(t *testing.T) {
	TestSetupDbConnection(t)
	for _, pat := range patients {
		req := patient.PutPatientRegisterParams{Info: pat}
		responder := RegisterPatient(req)
		switch resp := responder.(type) {
		case *patient.PutPatientRegisterDefault:
			if handledErrorMessage(string(resp.Payload)) {
				t.Log("[Handled fail] duplicate entry")
				continue
			}
			t.Errorf("Unhandled error:%v", resp.Payload)
		case *patient.PutPatientRegisterOK:
			t.Log("ok response")
			continue
		default:
			t.Errorf("Unhandle error for pat:%v| response type:%T | resp:%v", pat, resp, resp)
		}
	}
}

func TestUpdatePatient(t *testing.T) {
	TestSetupDbConnection(t)
	for _, pat := range patients {
		req := patient.PostPatientUpdateParams{Info: pat}
		switch resp := UpdatePatient(req).(type) {
		case *patient.PostPatientUpdateDefault:
			if handledErrorMessage(string(resp.Payload)) {
				t.Log("[Handled fail] Duplicate entry")
				continue
			}
			t.Errorf("Unhandled error:%v", resp.Payload)
		case *patient.PostPatientUpdateOK:
			continue
		default:
			t.Errorf("Unhandle error for pat:%v| response type:%T | resp:%v", pat, resp, resp)
		}
	}
}
func TestRemovePatient(t *testing.T) {
	TestSetupDbConnection(t)
	for _, pat := range patients {
		req := patient.DeletePatientRemoveParams{ID: pat.ID}
		switch resp := RemovePatient(req).(type) {
		case *patient.DeletePatientRemoveDefault:
			if handledErrorMessage(string(resp.Payload)) {
				t.Log("[Handled fail] Duplicate entry")
				continue
			}
			t.Errorf("Unhandled error:%v", resp.Payload)
		case *patient.DeletePatientRemoveOK:
			continue
		default:
			t.Errorf("Unhandle error for pat:%v| response type:%T | resp:%v", pat, resp, resp)
		}
	}
}

func TestFindPatient(t *testing.T) {
	TestSetupDbConnection(t)
	for _, pat := range patients {
		req := patient.NewGetPatientFindParams()
		req.NameContaining = &pat.Name
		var p, s = int64(1), int64(10)
		req.Page = &p
		req.Size = &s
		switch resp := FindPatient(req).(type) {
		case *patient.GetPatientFindDefault:
			if handledErrorMessage(string(resp.Payload)) {
				t.Log("[Handled fail]")
				continue
			}
		case *patient.GetPatientFindOK:
			continue
		default:
			t.Errorf("Unhandle error for pat:%v| response type:%T | resp:%v", pat, resp, resp)
		}
	}
}
