package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Update the routes() method to return a http.Handler instead of a *httprouter.Router.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Add the route for the GET /v1/movies endpoint.
	// Use the requireActivatedUser() middleware on our five /v1/movies** endpoints.
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
	// Add the route for the PUT /v1/movies/:id endpoint.
	// Require a PATCH request, rather than PUT.
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	// Add the route for the DELETE /v1/movies/:id endpoint.
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))
	// Add the route for the POST /v1/users endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	// Add the PUT /v1/users/password endpoint.
	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)
	// Add the route for the POST /v1/tokens/authentication endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	// Add the POST /v1/tokens/activation endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
	// Add the POST /v1/tokens/password-reset endpoint.
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	// Register a new GET /debug/vars endpoint pointing to the expvar handler.
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// Wrap the router with the panic recovery middleware.
	// Wrap the router with the rateLimit() middleware.
	// Use the authenticate() middleware on all requests.
	// Add the enableCORS() middleware.

	// Use the new metrics() middleware at the start of the chain.
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
