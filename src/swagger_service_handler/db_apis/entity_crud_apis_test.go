package apis

import (
	"testing"

	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/read"
	"telehealers.in/router/restapi/operations/register"
	"telehealers.in/router/restapi/operations/remove"
	"telehealers.in/router/restapi/operations/update"
)

var (
	testcase1  = models.Entity{Name: "test-1", Description: "testcase 1"}
	testcase2  = models.Entity{Name: "test-2", Description: "testcase 2"}
	testcase12 = models.Entity{Name: "test-12", Description: "testcase 12"}
	tests      = [](models.Entity){testcase1, testcase12, testcase2}
)

func errRespTesting(errString string, t *testing.T) {
	if handledErrorMessage(string(errString)) {
		t.Log("[Handled fail]")
		return
	}
	t.Errorf("Unhandled error:%v", errString)
}

func testUpdateMedicineAPI(t *testing.T, ent models.Entity) {
	req := update.NewPostMedicineUpdateParams()
	req.Info = &ent
	responder := UpdateMedicineAPI(req, nil)
	switch resp := responder.(type) {
	case *update.PostMedicineUpdateDefault:
		errRespTesting(string(resp.Payload), t)
	case *update.PostMedicineUpdateOK:
		return
	default:
		t.Errorf("Unhandled resp[%T]:%v", resp, resp)
	}
}

/** Register and updates (if possible) medicines **/
func TestEditMedicineAPIs(t *testing.T) {
	TestSetupDbConnection(t)
	for _, ent := range tests {
		medicine := register.NewPutMedicineRegisterParams()
		medicine.Info.Data = &ent
		responder := RegisterMedicineAPI(medicine, nil)
		switch resp := responder.(type) {
		case *register.PutMedicineRegisterDefault:
			errRespTesting(string(resp.Payload), t)
		case *register.PutMedicineRegisterOK:
			ent.ID = resp.Payload.ID
			ent.Name += "[MODIFIED]"
			testUpdateMedicineAPI(t, ent)
			continue
		default:
			t.Errorf("Unhandled error for ent| response type:%T| resp:%v", resp, resp)
		}
	}
}

func testUpdateMedTestAPI(t *testing.T, ent models.Entity) {
	req := update.NewPostMedicalTestUpdateParams()
	req.Info = &ent
	responder := UpdateMedTestAPI(req, nil)
	switch resp := responder.(type) {
	case *update.PostMedicineUpdateDefault:
		errRespTesting(string(resp.Payload), t)
	case *update.PostMedicineUpdateOK:
		return
	default:
		t.Errorf("Unhandled resp[%T]:%v", resp, resp)
	}
}

/** Register and updates (if possible) tests **/
func TestEditMedTestAPI(t *testing.T) {
	for _, ent := range tests {
		test := register.NewPutMedicalTestRegisterParams()
		test.Info.Data = &ent
		responder := RegisterTestAPI(test, nil)
		switch resp := responder.(type) {
		case *register.PutMedicalTestRegisterDefault:
			errRespTesting(string(resp.Payload), t)
		case *register.PutMedicalTestRegisterOK:
			ent.ID = resp.Payload.ID
			ent.Name += "[MODIFIED]"
			testUpdateMedTestAPI(t, ent)
			continue
		default:
			t.Errorf("Unhandled error for ent| response type:%T| resp:%v", resp, resp)
		}

	}
}

func testUpdateAdviceAPI(t *testing.T, ent models.Entity) {
	req := update.NewPostMedicalAdviceUpdateParams()
	req.Info = &ent
	responder := UpdateAdviceAPI(req, nil)
	switch resp := responder.(type) {
	case *update.PostMedicalAdviceUpdateDefault:
		errRespTesting(string(resp.Payload), t)
	case *update.PostMedicalAdviceUpdateOK:
		return
	default:
		t.Errorf("Unhandled resp[%T]:%v", resp, resp)
	}
}

func TestEditAdviceAPI(t *testing.T) {
	for _, ent := range tests {
		advice := register.NewPutMedicalAdviceRegisterParams()
		advice.Info.Data = &ent
		responder := RegisterAdviceAPI(advice, nil)
		switch resp := responder.(type) {
		case *register.PutMedicalAdviceRegisterDefault:
			errRespTesting(string(resp.Payload), t)
		case *register.PutMedicalAdviceRegisterOK:
			ent.ID = resp.Payload.ID
			ent.Name += "[MODIFIED]"
			testUpdateAdviceAPI(t, ent)
			continue
		default:
			t.Errorf("Unhandled error for ent| response type:%T| resp:%v", resp, resp)
		}
	}
}

func TestFindAndRemoveEntities(t *testing.T) {
	TestSetupDbConnection(t)
	for _, ent := range tests {
		findMedReq := read.NewGetMedicineFindParams()
		findMedReq.NameContaining = &ent.Name
		foundMedicineResponder := FindMedicineAPI(findMedReq, nil)
		var foundMeds [](*models.Entity)
		switch resp := foundMedicineResponder.(type) {
		case *read.GetMedicineFindDefault:
			errRespTesting(string(resp.Payload), t)
		case *read.GetMedicineFindOK:
			foundMeds = resp.Payload.Data
		default:
			t.Errorf("Unhandled error for | response type:%T| resp:%v", resp, resp)
		}
		for _, med := range foundMeds {
			deleteReq := remove.NewDeleteMedicineRemoveParams()
			deleteReq.ID = med.ID
			deleteResponder := RemoveMedicineAPI(deleteReq, nil)
			switch resp := deleteResponder.(type) {
			case *remove.DeleteMedicineRemoveDefault:
				errRespTesting(string(resp.Payload), t)
			case *remove.DeleteMedicineRemoveOK:
			default:
				t.Errorf("Unhandled error for | response type:%T| resp:%v", resp, resp)
			}
		}

		//test
		findMedTestReq := read.NewGetMedicalTestFindParams()
		findMedTestReq.NameContaining = &ent.Name
		foundMedTestResponder := FindMedTestAPI(findMedTestReq, nil)
		var foundTests [](*models.Entity)
		switch resp := foundMedTestResponder.(type) {
		case *read.GetMedicalTestFindDefault:
			errRespTesting(string(resp.Payload), t)
		case *read.GetMedicalTestFindOK:
			foundTests = resp.Payload.Data
		default:
			t.Errorf("Unhandled error for | response type:%T| resp:%v", resp, resp)
		}
		for _, med := range foundTests {
			deleteReq := remove.NewDeleteMedicineRemoveParams()
			deleteReq.ID = med.ID
			deleteResponder := RemoveMedicineAPI(deleteReq, nil)
			switch resp := deleteResponder.(type) {
			case *remove.DeleteMedicalTestRemoveDefault:
				errRespTesting(string(resp.Payload), t)
			case *remove.DeleteMedicalAdviceRemoveOK:
			default:
				t.Errorf("Unhandled error for | response type:%T| resp:%v", resp, resp)
			}
		}
		//advice
		findAdviceReq := read.NewGetMedicalAdviceFindParams()
		findAdviceReq.NameContaining = &ent.Name
		foundAdviceResponder := FindAdviceAPI(findAdviceReq, nil)
		var foundAdvices [](*models.Entity)
		switch resp := foundAdviceResponder.(type) {
		case *read.GetMedicalAdviceFindDefault:
			errRespTesting(string(resp.Payload), t)
		case *read.GetMedicalAdviceFindOK:
			foundAdvices = resp.Payload.Data
		default:
			t.Errorf("Unhandled error for | response type:%T| resp:%v", resp, resp)
		}
		for _, med := range foundAdvices {
			deleteReq := remove.NewDeleteMedicalAdviceRemoveParams()
			deleteReq.ID = med.ID
			deleteResponder := RemoveAdviceAPI(deleteReq, nil)
			switch resp := deleteResponder.(type) {
			case *remove.DeleteMedicalAdviceRemoveDefault:
				errRespTesting(string(resp.Payload), t)
			case *remove.DeleteMedicalAdviceRemoveOK:
			default:
				t.Errorf("Unhandled error for | response type:%T| resp:%v", resp, resp)
			}
		}
	}
}
