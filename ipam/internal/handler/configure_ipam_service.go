// This file is safe to edit. Once it exists it will not be overwritten

package handler

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/alexvancasper/TunnelBroker/ipam/internal/handler/operations"
	serviceHanlders "github.com/alexvancasper/TunnelBroker/ipam/internal/servicesHandler"
)

//go:generate swagger generate server --target ../../../ipam --name IpamService --spec ../../api/api.yml --server-package internal/handler --principal interface{} --exclude-main

func configureFlags(api *operations.IpamServiceAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.IpamServiceAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf
	api.Logger = serviceHanlders.SHandler.Log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// if api.DeleteReleaseHandler == nil {
	// 	api.DeleteReleaseHandler = operations.DeleteReleaseHandlerFunc(func(params operations.DeleteReleaseParams) middleware.Responder {
	// 		return middleware.NotImplemented("operation operations.DeleteRelease has not yet been implemented")
	// 	})
	// }
	// if api.GetAcquireHandler == nil {
	// 	api.GetAcquireHandler = operations.GetAcquireHandlerFunc(func(params operations.GetAcquireParams) middleware.Responder {
	// 		return middleware.NotImplemented("operation operations.GetAcquire has not yet been implemented")
	// 	})
	// }

	api.DeleteReleaseHandler = operations.DeleteReleaseHandlerFunc(serviceHanlders.DeleteReleasePrefixHandler)
	api.GetAcquireHandler = operations.GetAcquireHandlerFunc(serviceHanlders.GetAcquirePrefixlenHandler)

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
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
