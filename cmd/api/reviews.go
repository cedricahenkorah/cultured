package main

import (
	"cultured/internal/models"
	"encoding/json"
	"net/http"
)

type createReviewRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Rating  int    `json:"rating"`
}

func (app *application) getReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := app.reviews.GetAll(r.Context())

	if err != nil {
		app.errorLog.Println(err)
		app.serverError(w, err)
		return
	}

	err = app.apiResponse(w, http.StatusOK, reviews, "", nil)

	if err != nil {
		app.errorLog.Println(err)
		app.serverError(w, err)
	}
}

func (app *application) getReview(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
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

	err = app.apiResponse(w, http.StatusOK, review, "", nil)

	if err != nil {
		app.errorLog.Println(err)
		app.serverError(w, err)
	}
}

func (app *application) createReview(w http.ResponseWriter, r *http.Request) {
	var input createReviewRequest

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.clientError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if input.Title == "" || input.Content == "" || input.Rating < 1 || input.Rating > 5 {
		app.clientError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	userID, ok := app.getUserIDFromContext(r)

	if !ok {
		app.clientError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	review, err := app.reviews.Insert(r.Context(), input.Title, input.Content, input.Rating, userID)

	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.apiResponse(w, http.StatusCreated, review, "", nil)

	if err != nil {
		app.errorLog.Println(err)
		app.serverError(w, err)
	}
}
