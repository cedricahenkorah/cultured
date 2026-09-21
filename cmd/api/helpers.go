package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type contextKey string

var contextKeyUser = contextKey("user")

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) authenticatedUser(r *http.Request) int64 {
	return app.sessionManager.GetInt64(r.Context(), "userID")
}

func (app *application) getUserIDFromContext(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(contextKeyUser).(int64)

	return userID, ok && userID > 0
}
