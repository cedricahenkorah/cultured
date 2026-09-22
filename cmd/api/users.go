package main

import (
	"cultured/internal/data"
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
)

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	var input signupRequest

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.clientError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	parsedEmail, err := mail.ParseAddress(input.Email)

	if err != nil || parsedEmail.Address != input.Email {
		app.clientError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if input.Name == "" || len(input.Email) > 254 || len(input.Password) < 8 {
		app.clientError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	id, err := app.models.Users.CreateUser(r.Context(), input.Name, input.Email, input.Password)

	if err == data.ErrDuplicateEmail {
		app.clientError(w, http.StatusConflict, data.ErrDuplicateEmail.Error())
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	json.NewEncoder(w).Encode(id)
}
