package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.metrics, app.recoverPanic, app.logRequest, secureHeaders)

	router := chi.NewRouter()

	router.NotFound(app.notFound)

	router.Use(app.sessionManager.LoadAndSave)

	router.Group(func(r chi.Router) {
		r.Use(app.requireAuth)

		r.Get("/v1/review", app.getReviews)
		r.Get("/v1/review/{id}", app.getReview)
		r.Patch("/v1/review/{id}", app.updateReview)
		r.Delete("/v1/review/{id}", app.deleteReview)
		r.Post("/v1/review/create", app.createReview)
	})

	router.Get("/debug/vars", expvar.Handler().ServeHTTP)
	router.Get("/v1/healthcheck", app.healthcheck)

	router.Post("/v1/user/signup", app.signUp)
	router.Post("/v1/user/login", app.login)
	router.Post("/v1/user/logout", app.logout)

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(router)
}
