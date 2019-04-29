package f3api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type RestApi interface {
	// Fetch a payment resource
	GetPayment(rest.ResponseWriter, *rest.Request)

	// Create a new payment resource
	PostPayment(rest.ResponseWriter, *rest.Request)

	// Create or update a payment resource
	PutPayment(rest.ResponseWriter, *rest.Request)

	// Delete a payment resource
	DeletePayment(rest.ResponseWriter, *rest.Request)

	// TODO: List "a collection", whatever that means
	// For now, just list them all:
	// List all payment resources
	GetAllPayments(rest.ResponseWriter, *rest.Request)
}

// Example implementation of the API
// Does not implement stable storage, as it is literally an IN-MEMORY storage and therefore only suitable for testing
type InMemApi struct {
	// Map the ID to the payment object to prevent duplicates
	payments map[string]Payment
}

func NewInMemApi() *InMemApi {
	imapi := InMemApi{}
	imapi.payments = make(map[string]Payment)

	return &imapi
}

func (imapi *InMemApi) addPayment(p Payment) error {
	if _, ok := imapi.payments[p.ID]; ok {
		return errors.New(fmt.Sprintf("Payment with ID <%v> already exists!", p.ID))
	}
	imapi.payments[p.ID] = p
	return nil
}

func (imapi *InMemApi) GetPayment(w rest.ResponseWriter, r *rest.Request) {

}

func (imapi *InMemApi) PostPayment(w rest.ResponseWriter, r *rest.Request) {
	payment := Payment{}
	if err := r.DecodeJsonPayload(&payment); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := imapi.addPayment(payment); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&payment)
}

func (imapi *InMemApi) PutPayment(w rest.ResponseWriter, r *rest.Request) {

}

func (imapi *InMemApi) DeletePayment(w rest.ResponseWriter, r *rest.Request) {

}

func (imapi *InMemApi) GetAllPayments(w rest.ResponseWriter, r *rest.Request) {
	var payments []Payment
	for _, v := range imapi.payments {
		payments = append(payments, v)
	}
	w.WriteJson(&payments)
}
