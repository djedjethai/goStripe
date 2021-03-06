package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	// add the middleware(the function SessionLoad)
	mux.Use(SessionLoad)

	mux.Get("/", app.Home)

	// middleware
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth)
		mux.Get("/virtual-terminal", app.VirtualTerminal)
		mux.Get("/all-sales", app.AllSales)
		mux.Get("/all-subscriptions", app.AllSubscriptions)
		mux.Get("/sales/{id}", app.ShowSale)
		mux.Get("/subscriptions/{id}", app.ShowSubscription)

	})

	mux.Get("/widget/{id}", app.ChargeOnce)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/receipt", app.Receipt)

	// creating a stub page and a stub handler
	// for the plan options/page
	mux.Get("/plans/bronze", app.BronzePlan)
	// route to redirect to receipt page after BronzePlan is validated
	mux.Get("/receipt/bronze", app.BronzePlanReceipt)

	// authentification routes
	mux.Get("/login", app.Login)
	mux.Post("/login", app.PostLogin)
	mux.Get("/logout", app.Logout)
	mux.Get("/forgot-password", app.ForgotPassword)
	mux.Get("/reset-password", app.ShowResetPassword)

	// static content could be embeded the same way we did with the template
	// but thats a little awkward, so we won't
	// let serve them from an external directory
	fileServer := http.FileServer(http.Dir("./static"))
	// StripPrefix() stripe off /static
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
