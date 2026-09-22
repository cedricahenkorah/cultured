package main

import (
	"cultured/internal/models"
	"encoding/json"
	"net/http"
	"strings"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var input loginRequest

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.clientError(w, http.StatusBadRequest, "")
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	if input.Email == "" || input.Password == "" {
		app.clientError(w, http.StatusBadRequest, models.ErrInvalidCredentials.Error())
		return
	}

	id, err := app.users.Authenticate(r.Context(), input.Email, input.Password)

	if err == models.ErrInvalidCredentials {
		app.clientError(w, http.StatusUnauthorized, models.ErrInvalidCredentials.Error())
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "userID", id)

	json.NewEncoder(w).Encode(id)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.Destroy(r.Context())

	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
