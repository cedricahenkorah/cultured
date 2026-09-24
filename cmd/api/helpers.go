package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type contextKey string

var contextKeyUser = contextKey("user")

type apiResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Output(2, trace)

	app.apiResponse(w, http.StatusInternalServerError, nil, http.StatusText(http.StatusInternalServerError), nil)
}

func (app *application) clientError(w http.ResponseWriter, status int, message string) {
	app.apiResponse(w, status, nil, message, nil)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.apiResponse(w, http.StatusNotFound, nil, http.StatusText(http.StatusNotFound), nil)
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
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) apiResponse(w http.ResponseWriter, code int, data any, message string, headers http.Header) error {
	status := "success"

	if code >= 400 {
		status = "error"
	}

	response := apiResponse{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    data,
	}

	js, err := json.Marshal(response)

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)

	return nil
}

func parseTimeDuration(value string, fallback time.Duration) time.Duration {
	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)

	if err != nil {
		log.Printf("invalid query timeout %q; using default %s", value, fallback)
		return fallback
	}

	return duration
}

func (app *application) background(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				app.errorLog.Println(err)
			}
		}()

		fn()
	}()
}
