package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// v1
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/transactions/summary", app.transactionsSummaryHandler)

	// v2
	router.HandlerFunc(http.MethodGet, "/v2/transactions", app.requireActivatedUser(app.listTransactionsHandler))
	router.HandlerFunc(http.MethodPost, "/v2/transactions", app.requireActivatedUser(app.saveTransactionHandler))

	router.HandlerFunc(http.MethodPost, "/v2/users/register", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v2/users/activate", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v2/auth/authenticate", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v2/auth/activation", app.createActivationTokenHandler)

	return app.recoverPanic(app.authenticate(router))
}
