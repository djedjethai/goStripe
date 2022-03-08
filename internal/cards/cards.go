package cards

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/refund"
	"github.com/stripe/stripe-go/v72/sub"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

// to call this func to call the next one is not necessary
// but in case we use different payment provider like paypall,
// it can be usefull to create a layer of abstraction
func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	// payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	// at this point if i wanted to, i could add some metadata
	// params.AddMetadata("key", "value")

	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		// get the stripe err, from the stripe package
		// if ok is false
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}
	return pi, "", nil
}

// GetPaymentMethod gets the payment method by payment intend id
func (c *Card) GetPaymentMethod(s string) (*stripe.PaymentMethod, error) {
	stripe.Key = c.Secret

	pm, err := paymentmethod.Get(s, nil)
	if err != nil {
		return nil, err
	}

	return pm, nil
}

// RetreivePaymentIntent Get an existing payment Intent by id
func (c *Card) RetrievePaymentIntent(id string) (*stripe.PaymentIntent, error) {
	stripe.Key = c.Secret

	pi, err := paymentintent.Get(id, nil)
	if err != nil {
		return nil, err
	}

	return pi, nil
}

func (c *Card) SubscribeToPlan(cust *stripe.Customer, plan, email, last4, cardType string) (*stripe.Subscription, error) {
	stripeCustomerID := cust.ID
	// what we gonna subscribe this customer to
	items := []*stripe.SubscriptionItemsParams{
		{Plan: stripe.String(plan)}, // cast the plan to stripe.String() format
	}
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeCustomerID),
		Items:    items,
	}

	// add some metadata
	params.AddMetadata("last_four", last4)
	params.AddMetadata("card_type", cardType)

	// add the payment intent, as this people will be charge on a regular basis
	params.AddExpand("latest_invoice.payment_intent")

	subscription, err := sub.New(params)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (c *Card) CreateCustomer(pm, email string) (*stripe.Customer, string, error) {
	stripe.Key = c.Secret

	customerParams := &stripe.CustomerParams{
		PaymentMethod: stripe.String(pm),
		Email:         stripe.String(email),
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(pm),
		},
	}

	// where customer  come from ?????
	cust, err := customer.New(customerParams)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}

	return cust, "", nil
}

func (c *Card) Refund(pi string, amount int) error {
	stripe.Key = c.Secret
	amountToRefund := int64(amount)

	refundParams := &stripe.RefundParams{
		Amount:        &amountToRefund,
		PaymentIntent: &pi,
	}

	// the ignore param _ is the refund object return by stripe
	// see the doc, but we r not using it here
	_, err := refund.New(refundParams)
	if err != nil {
		return err
	}

	return nil
}

func (c *Card) CancelSubscription(subID string) error {

	stripe.Key = c.Secret
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err := sub.Update(subID, params)
	if err != nil {
		return err
	}

	return nil
}

func cardErrorMessage(code stripe.ErrorCode) string {
	var msg = ""
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Your card was declined"
	case stripe.ErrorCodeExpiredCard:
		msg = "Your card is expired"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "Incorrect CVC code"
	case stripe.ErrorCodeIncorrectZip:
		msg = "Incorrect zip/postal code"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "The amount is too large to charge to your card"
	case stripe.ErrorCodeAmountTooSmall:
		msg = "The amount is too small to charge to your card"
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "Insufficient balance"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "Your postal code is invalid"
	default:
		msg = "Your card was declined"
	}
	return msg
}
