package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	// add the middleware(the function SessionLoad)
	mux.Use(SessionLoad)

	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Get("/", app.Home)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/widget/{id}", app.ChargeOnce)

	// static content could be embeded the same way we did with the template
	// but thats a little awkward, so we won't
	// let serve them from an external directory
	fileServer := http.FileServer(http.Dir("./static"))
	// StripPrefix() stripe off /static
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
