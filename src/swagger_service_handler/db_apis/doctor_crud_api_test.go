package apis

import (
	"os"
	"strings"
	"testing"

	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/doctor"
)

var (
	//Test Doctors
	docOc = &models.DoctorInfo{Name: "Dr. Oto Octavius", Email: "8@sinister6",
		Phone: "8888666666"}
)

func TestSetupDbConnection(t *testing.T) {
	connErr := InitConnection(os.Getenv("DB_NAME"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"), os.Getenv("DB_ADDR"))
	if connErr != nil {
		t.Errorf("Check test-env vars.")
	}
}

func TestRegisterDoctor(t *testing.T) {
	TestSetupDbConnection(t)
	req := doctor.PutDoctorRegisterParams{Info: docOc}
	switch resp := RegisterDoctor(req).(type) {
	case *doctor.PutDoctorRegisterDefault:
		if strings.Contains(string(resp.Payload), "already") {
			t.Log("[Handled fail] Duplicate entry")
			return
		}
		t.Errorf("Unhandled error:%v", resp.Payload)
	}
}
