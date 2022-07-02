package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/rs/cors"
	"telehealers.in/router/restapi/operations"
	"telehealers.in/router/restapi/operations/conferrencing"
	customHandlers "telehealers.in/router/src/swagger_service_handler"
	vcCustomHandlers "telehealers.in/router/src/swagger_service_handler/conferrencing"
)

//go:generate swagger generate server --target ../../telehealers --name TelehealersBackend --spec ../swagger/swagger.yml --principal interface{}

func configureFlags(api *operations.TelehealersBackendAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TelehealersBackendAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.BinProducer = runtime.ByteStreamProducer()
	api.JSONProducer = runtime.JSONProducer()

	fmt.Println("Installing profile_picture endpoint.")
	api.GetProfilePicturesNameHandler = operations.GetProfilePicturesNameHandlerFunc(
		customHandlers.GetProfilePictures)
	api.ConferrencingGetRoomAccessTokenHandler = conferrencing.GetRoomAccessTokenHandlerFunc(
		vcCustomHandlers.GetAccessToken)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

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
