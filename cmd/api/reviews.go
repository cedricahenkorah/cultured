package main

import (
	"cultured/internal/data"
	"encoding/json"
	"net/http"
	"strings"
)

type createReviewRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Rating  int    `json:"rating"`
}

type updateReviewRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Rating  *int    `json:"rating"`
}

func (app *application) getReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := app.models.Reviews.GetAll(r.Context())

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
		app.notFound(w, r)
		return
	}

	review, err := app.models.Reviews.Get(r.Context(), id)

	if err == data.ErrNoRecord {
		app.notFound(w, r)
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

	review, err := app.models.Reviews.Insert(r.Context(), input.Title, input.Content, input.Rating, userID)

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

func (app *application) updateReview(w http.ResponseWriter, r *http.Request) {
	var input updateReviewRequest

	id, err := app.readIDParam(r)

	if err != nil {
		app.notFound(w, r)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		app.clientError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		app.clientError(w, http.StatusBadRequest, "Title cannot be blank")
		return
	}

	if input.Content != nil && strings.TrimSpace(*input.Content) == "" {
		app.clientError(w, http.StatusBadRequest, "Content cannot be blank")
		return
	}

	if input.Rating != nil && (*input.Rating < 1 || *input.Rating > 5) {
		app.clientError(w, http.StatusBadRequest, "Rating must be between 1 and 5")
		return
	}

	if input.Title == nil && input.Content == nil && input.Rating == nil {
		app.clientError(w, http.StatusBadRequest, "Provide a field to update it's value")
		return
	}

	userID, ok := app.getUserIDFromContext(r)

	if !ok {
		app.clientError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	review, err := app.models.Reviews.Update(r.Context(), input.Title, input.Content, input.Rating, id, userID)

	if err == data.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.clientError(w, http.StatusFailedDependency, http.StatusText(http.StatusFailedDependency))
		return
	}

	err = app.apiResponse(w, http.StatusOK, review, "", nil)

	if err != nil {
		app.errorLog.Println(err)
		app.serverError(w, err)
	}
}

func (app *application) deleteReview(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFound(w, r)
		return
	}

	userID, ok := app.getUserIDFromContext(r)

	if !ok {
		app.clientError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	err = app.models.Reviews.Delete(r.Context(), id, userID)

	if err == data.ErrNoRecord {
		app.notFound(w, r)
		return
	} else if err != nil {
		app.clientError(w, http.StatusFailedDependency, http.StatusText(http.StatusFailedDependency))
		return
	}

	err = app.apiResponse(w, http.StatusOK, nil, "Review deleted", nil)

	if err != nil {
		app.errorLog.Println(err)
		app.serverError(w, err)
	}
}
