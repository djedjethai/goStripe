package main

// go get github.com/go-chi/cors

import (
	// "fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// use the chi cors package, it s passed as a middleware
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Post("/api/payment-intent", app.GetPaymentIntent)

	mux.Post("/api/create-customer-and-subscribe-to-plan", app.CreateCustomerAndSubscribeToPlan)

	// a db test
	mux.Get("/api/widget/{id}", app.GetWidgetByID)

	mux.Post("/api/authenticate", app.CreateAuthToken)

	mux.Post("/api/is-authenticated", app.CheckAuthentication)

	return mux
}
