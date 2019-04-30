package f3api

import (
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

	// List all payment resources
	// -- List "a collection", whatever that means
	// Without further specification (e.g. about pagination), just list them all:
	GetAllPayments(rest.ResponseWriter, *rest.Request)
}

// Generic implementation of the API
type GenericApi struct {
	store ApiStore
}

func NewGenericApi(store ApiStore) *GenericApi {
	ga := GenericApi{
		store: store,
	}
	return &ga
}

func (api *GenericApi) handleError(w rest.ResponseWriter, r *rest.Request, err error) {
	switch err.(type) {
	case NotFoundError:
		rest.NotFound(w, r)
	default:
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (api *GenericApi) GetPayment(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	payment, err := api.store.GetPayment(id)
	if err != nil {
		api.handleError(w, r, err)
		return
	}

	w.WriteJson(payment)
}

func (api *GenericApi) PostPayment(w rest.ResponseWriter, r *rest.Request) {
	payment := Payment{}
	if err := r.DecodeJsonPayload(&payment); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// We probably shouldn't have an ID when POST-ing, but since I don't know how IDs are being generated,
	// we'll only accept payments with pre-generated IDs
	if payment.ID == "" {
		rest.Error(w, "Payment without ID is invalid", http.StatusInternalServerError)
		return
	}

	err := api.store.AddPayment(payment)
	if err != nil {
		api.handleError(w, r, err)
		return
	}

	w.WriteJson(&payment)
}

func (api *GenericApi) PutPayment(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")

	if id == "" {
		rest.Error(w, "PUT Request must include ID", http.StatusBadRequest)
	}

	payment := Payment{}
	if err := r.DecodeJsonPayload(&payment); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payment.ID = id

	api.store.StorePayment(payment)

	w.WriteJson(&payment)
}

func (api *GenericApi) DeletePayment(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")

	if id == "" {
		rest.Error(w, "DELETE Request must include ID", http.StatusBadRequest)
	}
	err := api.store.DeletePayment(id)
	if err != nil {
		api.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *GenericApi) GetAllPayments(w rest.ResponseWriter, r *rest.Request) {
	payments, err := api.store.GetAllPayments()

	if err != nil {
		api.handleError(w, r, err)
		return
	}

	w.WriteJson(&payments)
}
