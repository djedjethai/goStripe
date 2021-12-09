package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/charge-once", app.ChargeOnce)

	// static content could be embeded the same way we did with the template
	// but thats a little awkward, so we won't
	// let serve them from an external directory
	fileServer := http.FileServer(http.Dir("./static"))
	// StripPrefix() stripe off /static
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
