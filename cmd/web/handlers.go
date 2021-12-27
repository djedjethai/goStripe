package main

import (
	// "github.com/djedjethai/goStripe/internal/models"
	"github.com/djedjethai/goStripe/internal/cards"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {

	if err := app.renderTemplate(w, r, "terminal", &templateData{}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	// app.infoLog.Println("hit the handlers")
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	cardHolder := r.Form.Get("cardholder_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")

	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}

	// see func in cards
	pi, err := card.RetrievePaymentIntent(paymentIntent)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// see func in cards
	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear

	data := make(map[string]interface{})
	data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency
	data["last_four"] = lastFour
	data["expiry_month"] = expiryMonth
	data["expiry_year"] = expiryYear
	data["bank_return_code"] = pi.Charges.Data[0].ID

	if err := app.renderTemplate(w, r, "succeeded", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

// CargeOnce display the page to buy one widget
func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, _ := strconv.Atoi(id)

	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := make(map[string]interface{})
	data["widget"] = widget
	if err := app.renderTemplate(w, r, "buy-once", &templateData{
		Data: data,
	}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}
