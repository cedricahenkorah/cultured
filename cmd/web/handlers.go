package main

import (
	"context"
	"cultured/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.errorLog.Println("Invalid path:", r.URL.Path)
		app.notFound(w)
		return
	}

	reviews, err := app.reviews.GetAll(context.Background())

	if err != nil {
		app.errorLog.Println(err)
		app.serverError(w, err)
		return
	}

	json.NewEncoder(w).Encode(reviews)
}

func (app *application) getReview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	review, err := app.reviews.Get(context.Background(), int64(id))

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
