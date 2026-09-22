package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (app *application) readIDParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil || id < 1 {
		return 0, err
	}

	return id, nil
}
