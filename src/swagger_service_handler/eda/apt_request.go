package eda

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/vmware/transport-go/bus"
	"github.com/vmware/transport-go/model"
	"github.com/vmware/transport-go/plank/utils"
	"github.com/vmware/transport-go/service"
	"github.com/vmware/transport-go/stompserver"
	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations/doctor"
	"telehealers.in/router/restapi/operations/patient"
)

const (
	//Commands on AppointmentRequestChannel
	/**wait-for-apt: This is a data stream with Payload as Appointment-request
			F-Spec: Docs listen to this stream for apts
	**/
	WaitForAptCmd = "wait-for-apt"
	/**req-apt: Patients request for doctors in this message
	**/
	RequestAptCmd = "req-apt"
)

/** Model **/
type DocIDStruct struct {
	ID        int    `mapstructure:"id"`
	SessionID string `mapstructure:"session_id"`
}

/** Patient's appointment request**/
type PatAptReq struct {
	/** Doctor-ID if a particular doctor is to be requested, else leave empty **/
	DoctorID int `mapstructure:"doc_id,omitempty"`
	/** id x session_id is used for validation of session **/
	ID        int `mapstructure:"id"`
	SessionID int `mapstructure:"session_id"`
	/** appointment_types are used to filter applicable doctors **/
	AptTypes []string `mapstructure:"appointment_type,omitempty"`
	/** reason for appointment **/
	Reason string `mapstructure:"reason,omitempty"`
}

/** Confirm appointment request from doctor **/
type DocConfirmApt struct {
	PatientID int      `mapstructure:"patient_id"`
	AptTypes  []string `mapstructure:"appointment_type,omitempty"`
	Reason    string   `mapstructure:"reason,omitempty"`
}

type PatConfirmApt struct {
	DocID int `mapstructure:"doctor_id"`
}

const patAptReqPayloadSchema = "{id:int, [doc_id:int], session_id:string, [reason:string], [appointment_type:List of specialities(in conjunction)]}"

/** End of Model structs **/

type AppointmentService struct {
	lock          sync.RWMutex
	docIDToConnID map[int64]string
	/** inv: Request.Payload is decodable to return type of newDocInfoAndSessionID**/
	connIDToReq map[string]model.Request
	// core           service.FabricServiceCore
}

func NewAppointmentService() *AppointmentService {
	return &AppointmentService{docIDToConnID: make(map[int64]string), connIDToReq: make(map[string]model.Request)}
}

/** Returns error if multiple insert for same doc happens **/
func (ds *AppointmentService) addReq(docID int64, req *model.Request) (err error) {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	if _, connPresent := ds.connIDToReq[req.BrokerDestination.ConnectionId]; connPresent {
		utils.Log.Error("Multi apt request from doc ", docID)
		err = errors.New("Only single request allowed from single user")
	} else {
		utils.Log.Info("Doctor ", docID, " with connection ", req.BrokerDestination.ConnectionId, " inserted into the map")
		ds.docIDToConnID[docID] = req.BrokerDestination.ConnectionId
		ds.connIDToReq[req.BrokerDestination.ConnectionId] = *req
	}
	return
}

func (ds *AppointmentService) removeConnID(connID string) {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	if _, present := ds.connIDToReq[connID]; present {
		docData := newDocInfoAndSessionID()
		mapstructure.Decode(ds.connIDToReq[connID].Payload, &docData)
		docID := docData.Doctor.ID
		delete(ds.docIDToConnID, docID)
		delete(ds.connIDToReq, connID)
	} else {
		utils.Log.Error("Connection not found ", connID)
	}
}

func (s *AppointmentService) Init(core service.FabricServiceCore) error {
	return nil
}

func newDocInfoAndSessionID() *doctor.GetDoctorLoginOKBody {
	return &doctor.GetDoctorLoginOKBody{Doctor: &models.DoctorInfo{}}
}

func newPatInfoAndSessionID() *patient.GetPatientLoginOKBody {
	return &patient.GetPatientLoginOKBody{Patient: &models.PatientInfo{}}
}

type AptScheduledResponse struct {
	ParticipantIDs []int64 `json:"participant_ids"`
	AccessToken    string  `json:"access_token"`
	RoomName       string  `json:"room_name"`
	AppointmentID  int64   `json:"appointment_id"`
}

func newAptScheduleResponse(id, aptID int64, roomName, at string) *AptScheduledResponse {
	return &AptScheduledResponse{ParticipantIDs: []int64{id}, AppointmentID: aptID, RoomName: roomName, AccessToken: at}
}

/*
  - App Flow
    Doctors listen to goOnlineChannel -> Patient request on it -> Doctor receives call and patient cred ->
    patient receives doc and call cred

*
*/
/** TODO maintain doctor availability data in DB, using Telehealers REST APIs **/
func (ds *AppointmentService) HandleServiceRequest(request *model.Request, core service.FabricServiceCore) {
	switch request.Request {
	case WaitForAptCmd:
		_req := newDocInfoAndSessionID()
		if decodeErr := mapstructure.Decode(request.Payload, _req); decodeErr != nil {
			utils.Log.Error("In decoding input:", decodeErr, "|Payload:", request.Payload)
			core.SendErrorResponse(request, 400, "Payload error")
		} else {
			insertionError := ds.addReq(_req.Doctor.ID, request)
			if insertionError != nil {
				core.SendErrorResponse(request, 400, insertionError.Error())
			}
		}
	case RequestAptCmd:
		ds.handleRequestAptCMD(request, core)
	default:
		core.HandleUnknownRequest(request)
	}
}

/*
* This function create access token for doctor(one of the random) and patient and creates access token for a common room
to both *
*/
func (ds *AppointmentService) handleRequestAptCMD(patientRequest *model.Request, core service.FabricServiceCore) {
	_patReq := newPatInfoAndSessionID()
	if decodeErr := mapstructure.Decode(patientRequest.Payload, _patReq); decodeErr != nil {
		utils.Log.Error(RequestAptCmd, "| decode err:", decodeErr.Error(), " payload:%v", patientRequest.Payload.(map[string]interface{}))
		core.SendErrorResponse(patientRequest, 400, "Payload error.schema:"+patAptReqPayloadSchema)
		return
	}
	ds.lock.RLock()
	var docCachedRequest *model.Request = nil
	for _, _req := range ds.connIDToReq {
		docCachedRequest = &_req
	}
	ds.lock.RUnlock()

	if docCachedRequest == nil {
		utils.Log.Error("No doctor available. requested by patient ", _patReq.Patient.ID)
		core.SendErrorResponse(patientRequest, 404, "No doctor available, try later.")
		return
	}

	docInfo := newDocInfoAndSessionID()
	mapstructure.Decode(docCachedRequest.Payload, &docInfo)
	docID := docInfo.Doctor.ID
	roomName, patAT, docAT, atErr := getAccessTokensForDocAndPat(docID, _patReq.Patient.ID)
	if atErr != nil {
		utils.Log.Error("In fetching accesstoken:", atErr.Error())
		core.SendErrorResponse(docCachedRequest, 500, "Internal server error")
		core.SendErrorResponse(patientRequest, 500, "Internal server error")
	} else {
		if aptID, aptRegError := registerAppointment(_patReq.Patient.ID, docID); aptRegError != nil {
			utils.Log.Error("In registering appointment", aptRegError.Error())
			core.SendErrorResponse(docCachedRequest, 500, "Internal server error")
			core.SendErrorResponse(patientRequest, 500, "Internal server error")
		} else {
			core.SendResponse(patientRequest, newAptScheduleResponse(docID, aptID, roomName, patAT))
			core.SendResponse(docCachedRequest, newAptScheduleResponse(_patReq.Patient.ID, aptID, roomName, docAT))
			utils.Log.Info("Scheduled doc ", docID, " to patient ", _patReq.Patient.ID, " appointment id ", aptID)
		}
	}
}

/** This function creates a room, and generates access token for patient and doctor both **/
func getAccessTokensForDocAndPat(docID, patID int64) (roomName, patAT, docAT string, err error) {
	roomName = fmt.Sprintf("doc-%v-pat-%v", docID, patID)
	docAT, err = getAccessToken(roomName, docID)
	if err != nil {
		utils.Log.Error("In creating doc-access-token for doc:", docID)
		return
	}
	patAT, err = getAccessToken(roomName, patID)
	if err != nil {
		utils.Log.Error("In creating patient-access-token for patient:", patID)
	}
	return
}

/** Make REST Call and fetch access token**/
func getAccessToken(roomName string, id int64) (accessToken string, err error) {
	url := fmt.Sprintf("%v/room_access_token?room=%v&id=%v", restBackendAddress, roomName, id)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add(restAuthTokenName, restAuthPass)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.Log.Error("In creating access token for ", roomName, " user:", id, " error: ", err.Error())
		return
	}
	accessTokenBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Log.Error("In reading access token response for ", roomName, " user:", id, ". error:", err.Error())
		return
	}
	accessToken = string(accessTokenBytes)
	accessTokenSize := len(accessToken)
	if accessTokenSize >= 3 {
		accessToken = accessToken[1 : accessTokenSize-2]
	}
	return
}

/** Create new appointment in DB**/
func registerAppointment(patientID, docID int64) (aptID int64, err error) {
	url := fmt.Sprintf("%v/appointment/register", restBackendAddress)
	reqBody, _ := json.Marshal(models.AppointmentInfo{
		PatientID:       patientID,
		DoctorID:        docID,
		InitializeToNow: []string{"date", "start_time_requested"}})
	req, reqErr := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
	if reqErr != nil {
		utils.Log.Error("In creating request", reqErr.Error())
		return
	}
	req.Header.Add(restAuthTokenName, restAuthPass)
	req.Header.Add("Content-Type", "application/json")
	resp, putAPIErr := http.DefaultClient.Do(req)
	if putAPIErr != nil {
		utils.Log.Error("Error in APT-ID. error:", putAPIErr.Error())
		return
	}
	aptIDByte, respReadErr := ioutil.ReadAll(resp.Body)
	if respReadErr != nil {
		utils.Log.Error("Error in reading response:", respReadErr.Error())
		return
	}
	aptIDString := string(aptIDByte)
	aptIDShort, err := strconv.Atoi(aptIDString[:len(aptIDString)-1])
	aptID = int64(aptIDShort)
	utils.Log.Info("Appointment for ", docID, patientID, " registered with id ", aptID)
	return
}

// func (ds *AppointmentService) GetDoctor()

func (s *AppointmentService) OnServiceReady() chan bool {
	sessionNotifyHandler, _ := bus.GetBus().ListenStream(bus.STOMP_SESSION_NOTIFY_CHANNEL)
	utils.Log.Info("Appointment service custom event handler registered")
	sessionNotifyHandler.Handle(func(message *model.Message) {
		stompSessionEvt := message.Payload.(*bus.StompSessionEvent)
		switch stompSessionEvt.EventType {
		case stompserver.ConnectionClosed, stompserver.UnsubscribeFromTopic:
			utils.Log.Warn("Conn: ", stompSessionEvt.Id, " went offline")
			s.removeConnID(stompSessionEvt.Id)
		case stompserver.ConnectionEstablished, stompserver.ConnectionStarting:
			utils.Log.Info("New connection established:", stompSessionEvt.Id)
		default:
			utils.Log.Error(" Unhandled connection on ", AppointmentRequestChannel, " event ", stompSessionEvt.EventType)
		}
	}, func(err error) {
		utils.Log.Error("Session Notification Handler Error:", err.Error())
	})
	readyChan := make(chan bool, 1)
	readyChan <- true
	return readyChan
}

// OnServerShutdown removes the running tickers
func (*AppointmentService) OnServerShutdown() {
	return
}

// GetRESTBridgeConfig returns a config for a REST endpoint that performs the same action as the STOMP variant
// except that there will be only one response instead of every 30 seconds.
func (*AppointmentService) GetRESTBridgeConfig() []*service.RESTBridgeConfig {
	return []*service.RESTBridgeConfig{}
}
