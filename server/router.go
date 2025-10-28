package server

import (
	"account-transactions/store"

	_ "account-transactions/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(db store.Store) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", HandleAccountPost(db))

		r.Route("/{accountId}", func(r chi.Router) {
			r.Get("/", HandleGetAccount(db))
		})
	})
	r.Route("/transactions", func(r chi.Router) {
		r.Post("/", HandleTransactionPost(db))
	})

	return r
}
