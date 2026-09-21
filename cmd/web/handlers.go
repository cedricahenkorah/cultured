package main

import (
	"cultured/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	reviews, err := app.reviews.GetAll(r.Context())

	if err != nil {
		app.errorLog.Println(err)
		app.serverError(w, err)
		return
	}

	json.NewEncoder(w).Encode(reviews)
}

func (app *application) getReview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	review, err := app.reviews.Get(r.Context(), id)

	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	json.NewEncoder(w).Encode(review)
}

func (app *application) createReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("create a new review"))
}
