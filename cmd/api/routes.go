package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	router := chi.NewRouter()

	router.Use(app.sessionManager.LoadAndSave)

	router.Group(func(r chi.Router) {
		r.Use(app.requireAuth)

		r.Get("/", app.getReviews)
		r.Get("/review/{id}", app.getReview)
		r.Post("/review/create", app.createReview)
	})

	router.Post("/user/signup", app.signUp)
	router.Post("/user/login", app.login)
	router.Post("/user/logout", app.logout)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(router)
}
