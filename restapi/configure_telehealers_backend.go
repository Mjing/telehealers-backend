package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/rs/cors"
	"telehealers.in/router/models"
	"telehealers.in/router/restapi/operations"
	"telehealers.in/router/restapi/operations/appointment"
	"telehealers.in/router/restapi/operations/conferrencing"
	"telehealers.in/router/restapi/operations/doctor"
	"telehealers.in/router/restapi/operations/patient"
	"telehealers.in/router/restapi/operations/patient_health_info"
	"telehealers.in/router/restapi/operations/read"
	"telehealers.in/router/restapi/operations/register"
	"telehealers.in/router/restapi/operations/remove"
	"telehealers.in/router/restapi/operations/update"
	customHandlers "telehealers.in/router/src/swagger_service_handler"
	vcCustomHandlers "telehealers.in/router/src/swagger_service_handler/conferrencing"
	dbApis "telehealers.in/router/src/swagger_service_handler/db_apis"
)

//go:generate swagger generate server --target ../../telehealers --name TelehealersBackend --spec ../swagger/swagger.yml --principal interface{}

type customServerOptions struct {
	Log    string `short:"l" long:"log" description:"Log file"`
	DBName string `long:"dbname" description:"data-base name, leave empty to toggle reading vars from env variables"`
	DBUser string `long:"dbuser" description:"db username for login" default:"root"`
	DBPass string `long:"dbpass" description:"db password for login" default:""`
	DBAddr string `long:"dbaddr" description:"db address in for ip:port" default:"127.0.0.1:3306"`
}

var currentOpts = new(customServerOptions)

func configureFlags(api *operations.TelehealersBackendAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{ShortDescription: "Server configuration flags", Options: currentOpts}}
}

func configureAPI(api *operations.TelehealersBackendAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	swaggerLogger := *dbApis.SetupLogFile(currentOpts.Log)
	swaggerLogger.SetPrefix("|GENCODE|")
	api.Logger = swaggerLogger.Printf
	if currentOpts.DBName != "" {
		api.Logger("[INFO]using cmd-line opts for db-conn values")
		dbApis.SetConnectionVars(currentOpts.DBName, currentOpts.DBUser, currentOpts.DBPass, currentOpts.DBAddr)
	}
	dbApis.InitConnection()

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.BinProducer = runtime.ByteStreamProducer()
	api.JSONProducer = runtime.JSONProducer()

	fmt.Println("Installing profile_picture endpoint.")
	/** Endpoint Handlers **/
	api.GetProfilePicturesNameHandler = operations.GetProfilePicturesNameHandlerFunc(
		customHandlers.GetProfilePictures)
	api.ConferrencingGetRoomAccessTokenHandler = conferrencing.GetRoomAccessTokenHandlerFunc(
		vcCustomHandlers.GetAccessToken)
	/** Doctor CRUD APIs **/
	api.DoctorPutDoctorRegisterHandler = doctor.PutDoctorRegisterHandlerFunc(
		func(pdrp doctor.PutDoctorRegisterParams, p *models.Principal) middleware.Responder {
			return dbApis.RegisterDoctor(pdrp)
		})
	api.DoctorPostDoctorUpdateHandler = doctor.PostDoctorUpdateHandlerFunc(
		func(pdup doctor.PostDoctorUpdateParams, p *models.Principal) middleware.Responder {
			return dbApis.UpdateDoctor(pdup)
		})
	api.DoctorDeleteDoctorRemoveHandler = doctor.DeleteDoctorRemoveHandlerFunc(
		func(ddrp doctor.DeleteDoctorRemoveParams, p *models.Principal) middleware.Responder {
			return dbApis.RemoveDoctor(ddrp)
		})
	api.DoctorGetDoctorFindHandler = doctor.GetDoctorFindHandlerFunc(
		func(gdfp doctor.GetDoctorFindParams, p *models.Principal) middleware.Responder {
			return dbApis.FindDoctor(gdfp)
		})
	api.DoctorPostDoctorRegisterApplyHandler = doctor.PostDoctorRegisterApplyHandlerFunc(
		func(pdrap doctor.PostDoctorRegisterApplyParams, p *models.Principal) middleware.Responder {
			return dbApis.DoctorRegistrationApplicationAPI(pdrap)
		})
	api.DoctorPostDoctorRegisterReviewHandler = doctor.PostDoctorRegisterReviewHandlerFunc(
		func(pdrrp doctor.PostDoctorRegisterReviewParams, p *models.Principal) middleware.Responder {
			return dbApis.DoctorRegistrationApplicationReviewAPI(pdrrp)
		})
	api.DoctorGetDoctorLoginHandler = doctor.GetDoctorLoginHandlerFunc(
		func(gdlp doctor.GetDoctorLoginParams, p *models.Principal) middleware.Responder {
			return dbApis.DoctorLoginAPI(gdlp)
		})
	api.DoctorGetDoctorRegisterPendingApplicationsHandler = doctor.GetDoctorRegisterPendingApplicationsHandlerFunc(
		func(gdrpap doctor.GetDoctorRegisterPendingApplicationsParams, p *models.Principal) middleware.Responder {
			return dbApis.GetDoctorRegistrationApplicationAPI(gdrpap)
		})
	api.DoctorGetDoctorPatientsHandler = doctor.GetDoctorPatientsHandlerFunc(
		dbApis.DoctorRelatedPatientsAPI)
	api.DoctorPostDoctorOnlineHandler = doctor.PostDoctorOnlineHandlerFunc(dbApis.DoctorOnlineAPI)
	/** Patient CRUD APIs **/
	api.PatientPutPatientRegisterHandler = patient.PutPatientRegisterHandlerFunc(
		func(pprp patient.PutPatientRegisterParams, p *models.Principal) middleware.Responder {
			return dbApis.RegisterPatient(pprp)
		})
	api.PatientPostPatientUpdateHandler = patient.PostPatientUpdateHandlerFunc(
		func(ppup patient.PostPatientUpdateParams, p *models.Principal) middleware.Responder {
			return dbApis.UpdatePatient(ppup)
		})
	api.PatientGetPatientFindHandler = patient.GetPatientFindHandlerFunc(
		func(gpfp patient.GetPatientFindParams, p *models.Principal) middleware.Responder {
			return dbApis.FindPatient(gpfp)
		})
	api.PatientDeletePatientRemoveHandler = patient.DeletePatientRemoveHandlerFunc(
		func(dprp patient.DeletePatientRemoveParams, p *models.Principal) middleware.Responder {
			return dbApis.RemovePatient(dprp)
		})
	api.PatientGetPatientLoginHandler = patient.GetPatientLoginHandlerFunc(
		func(gplp patient.GetPatientLoginParams, p *models.Principal) middleware.Responder {
			return dbApis.PatientLoginAPI(gplp)
		})
	/** Patient Health INfo APIs **/
	api.PatientHealthInfoPutPatientHealthInfoAddHandler = patient_health_info.PutPatientHealthInfoAddHandlerFunc(
		dbApis.AddHealthInfoAPI)
	api.PatientHealthInfoGetPatientHealthInfoFindHandler = patient_health_info.GetPatientHealthInfoFindHandlerFunc(
		dbApis.FindHealthInfoAPI)
	/** Appointment CRUD APIs **/
	api.AppointmentPutAppointmentRegisterHandler = appointment.PutAppointmentRegisterHandlerFunc(
		dbApis.RegisterppointmentAPI)
	api.AppointmentPostAppointmentUpdateHandler = appointment.PostAppointmentUpdateHandlerFunc(
		dbApis.UpdateAppointmentAPI)
	api.AppointmentDeleteAppointmentRemoveHandler = appointment.DeleteAppointmentRemoveHandlerFunc(
		dbApis.RemoveAppointmentAPI)
	api.AppointmentGetAppointmentFindHandler = appointment.GetAppointmentFindHandlerFunc(
		dbApis.FindAppointmentAPI)
	api.AppointmentGetAppointmentCountHandler = appointment.GetAppointmentCountHandlerFunc(
		dbApis.CountAppointmentAPI)
	/** Entities CRUDs: Medicine, Tests, Advices **/
	//Register
	api.RegisterPutMedicineRegisterHandler = register.PutMedicineRegisterHandlerFunc(dbApis.RegisterMedicineAPI)
	api.RegisterPutMedicalTestRegisterHandler = register.PutMedicalTestRegisterHandlerFunc(dbApis.RegisterTestAPI)
	api.RegisterPutMedicalAdviceRegisterHandler = register.PutMedicalAdviceRegisterHandlerFunc(dbApis.RegisterAdviceAPI)
	//Update
	api.UpdatePostMedicineUpdateHandler = update.PostMedicineUpdateHandlerFunc(dbApis.UpdateMedicineAPI)
	api.UpdatePostMedicalTestUpdateHandler = update.PostMedicalTestUpdateHandlerFunc(dbApis.UpdateMedTestAPI)
	api.UpdatePostMedicalAdviceUpdateHandler = update.PostMedicalAdviceUpdateHandlerFunc(dbApis.UpdateAdviceAPI)
	//remove
	api.RemoveDeleteMedicineRemoveHandler = remove.DeleteMedicineRemoveHandlerFunc(dbApis.RemoveMedicineAPI)
	api.RemoveDeleteMedicalTestRemoveHandler = remove.DeleteMedicalTestRemoveHandlerFunc(dbApis.RemoveMedTestAPI)
	api.RemoveDeleteMedicalAdviceRemoveHandler = remove.DeleteMedicalAdviceRemoveHandlerFunc(dbApis.RemoveAdviceAPI)
	//find
	api.ReadGetMedicineFindHandler = read.GetMedicineFindHandlerFunc(dbApis.FindMedicineAPI)
	api.ReadGetMedicalTestFindHandler = read.GetMedicalTestFindHandlerFunc(dbApis.FindMedTestAPI)
	api.ReadGetMedicalAdviceFindHandler = read.GetMedicalAdviceFindHandlerFunc(dbApis.FindAdviceAPI)
	/*****End OF Entity Registration **********/

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	api.KeyAuth = func(token string) (*models.Principal, error) {
		if token == "letmein" {
			prin := models.Principal(token)
			return &prin, nil
		}
		api.Logger("Access attempt with incorrect api key auth: %s", token)
		return nil, errors.New(401, "incorrect api key auth")
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	handleCORS := cors.Default().Handler
	return handleCORS(handler)
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
