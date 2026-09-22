package main

import (
	"net/http"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.apiResponse(w, http.StatusOK, data, "", nil)

	if err != nil {
		app.errorLog.Println(err)
		app.serverError(w, err)
	}
}
