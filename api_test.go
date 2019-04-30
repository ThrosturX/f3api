package f3api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
)

func TestGenericApi(t *testing.T) {
	var restApi RestApi = NewGenericApi(NewInMemStore())

	// a basic test I wrote to get myself started -- throstur
	err := checkPostedResourceIncreasesCollectionSize(restApi, t)
	if err != nil {
		t.Fatal(err)
	}

	// NOTE: Could put more tests here
}

// Tests whether or not a resource can be created with a PUT request.
func TestPutCreates(t *testing.T) {
	var (
		api     RestApi
		payment Payment
	)
	responseWriter := &testResponseWriter{}

	api = NewGenericApi(NewInMemStore())

	err := json.Unmarshal([]byte(DEFAULT_PAYMENT), &payment)
	if err != nil {
		t.Fatal(err)
	}

	// verify that the resource doesn't exist
	if newList, err := getPayments(api); err != nil || len(newList) > 0 {
		if err != nil {
			t.Fatal(err)
		}
		t.Fatal("Test initialized with non-empty data set")
	}

	sendPutRequest(api, payment, responseWriter)

	// verify that the resource exists
	found, err := fetchPayment(api, payment.ID)
	if err != nil {
		t.Fatal(err)
	}

	// the ID should always be the same, regardless of race conditions, so check this first
	if payment.ID != found.ID {
		t.Fatalf("Mismatched IDs on payments: %s and %s", payment.ID, found.ID)
	}

	if !reflect.DeepEqual(payment, found) {
		t.Fatal("Mismatched payments: DeepEqual returned false")
	}
}

// Tests whether or not a resource can be updated with a PUT request.
// Also tests the GetPayment (get single) endpoint, since it makes sense to do it somewhere.
func TestPutCreatesAndUpdates(t *testing.T) {
	var (
		api     RestApi
		payment Payment
	)
	responseWriter := &testResponseWriter{}

	api = NewGenericApi(NewInMemStore())

	err := json.Unmarshal([]byte(DEFAULT_PAYMENT), &payment)
	if err != nil {
		t.Fatal(err)
	}

	// PUT the payment resource
	sendPutRequest(api, payment, responseWriter)

	// Modify the resource and PUT it again
	payment.Attributes.PaymentID = StringedInt(1337)

	// PUT the mutated payment resource
	sendPutRequest(api, payment, responseWriter)

	// Verify that the resource exists
	found, err := fetchPayment(api, payment.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that the found payment is the same as the newly mutated payment
	if !reflect.DeepEqual(payment, found) {
		t.Fatal("Mismatched payments: DeepEqual returned false")
	}

	// Sanity check: also verify that the PaymentID changed:
	if found.Attributes.PaymentID != 1337 {
		t.Fatalf("The PaymentID did not get transformed to 1337!")
	}
}

// Tests the DeletePayment method exclusively, leverages direct ApiStore access.
func TestDeletePayment(t *testing.T) {
	var (
		store          ApiStore
		api            RestApi
		responseWriter rest.ResponseWriter
		p1, p2         Payment
		err            error
	)

	store = NewInMemStore()
	api = NewGenericApi(store)
	responseWriter = &testResponseWriter{}

	p1 = defaultPayment()

	store.AddPayment(p1)

	_, err = store.GetPayment(p1.ID)
	if err != nil {
		t.Fatal(err)
	}

	params := make(map[string]string)
	params["id"] = p1.ID
	request := createRestRequest("DELETE", "/payments/:id", strings.NewReader(""), params)
	api.DeletePayment(responseWriter, request)

	p2, err = store.GetPayment(p1.ID)
	if err == nil || p2.ID == p1.ID {
		t.Fatalf("The resource wasn't supposed to exist!")
	}
}

// Utilities below

// A basic ResponseWriter that writes the result into a string
type testResponseWriter struct {
	http.ResponseWriter
	result string
}

func (trw *testResponseWriter) EncodeJson(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (trw *testResponseWriter) WriteJson(v interface{}) error {
	b, err := trw.EncodeJson(v)
	if err != nil {
		return err
	}
	_, err = trw.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (w *testResponseWriter) Write(b []byte) (int, error) {
	w.result = string(b)
	return len(w.result), nil
}

func (w *testResponseWriter) Read() []byte {
	return []byte(w.result)
}

func (w *testResponseWriter) WriteHeader(int) {
	// this method is currently not used for testing, and only exists to implement rest.ResponseWriter
	return
}

func createRestRequest(method string, urlStr string, body io.Reader, params map[string]string) *rest.Request {
	origReq, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		panic(err)
	}
	return &rest.Request{
		origReq,
		params,
		map[string]interface{}{},
	}
}

// Checks whether or not posting a new payment resource increases the number of items in the collection.
// Tests POST and GET(all) endpoints
// NOTE: This method ALWAYS posts the same payment resource (same ID). Consider adding it as a parameter or implementing an ID generator if needed.
func checkPostedResourceIncreasesCollectionSize(api RestApi, t *testing.T) error {
	var (
		err              error
		oldList, newList []Payment
		oldLen, newLen   int
	)

	params := make(map[string]string)
	// we will also be sending a new json payment to be created
	postPayment := createRestRequest("POST", "/payments", strings.NewReader(DEFAULT_PAYMENT), params)

	// get the original length, call the API, check the response
	if oldList, err = getPayments(api); err != nil {
		return err
	}

	// post the new payment
	responseWriter := &testResponseWriter{}
	api.PostPayment(responseWriter, postPayment)

	// get the new length as before
	if newList, err = getPayments(api); err != nil {
		return err
	}

	oldLen = len(oldList)
	newLen = len(newList)
	if newLen != oldLen+1 {
		return errors.New(fmt.Sprintf("Lengths did not match up! Could not satisfy %d == %d + 1", oldLen, newLen))
	}

	// Pass
	return nil
}

func getPayments(api RestApi) ([]Payment, error) {
	var payments []Payment

	responseWriter := &testResponseWriter{}
	getRequest := createRestRequest("GET", "/payments", strings.NewReader(""), nil)

	api.GetAllPayments(responseWriter, getRequest)

	buf := responseWriter.Read()

	err := json.Unmarshal(buf, &payments)
	if err != nil {
		return payments, err
	}

	return payments, nil
}

func sendPaymentsRequest(api RestApi, payload Payment, responseWriter rest.ResponseWriter, requestType string, params map[string]string) error {
	var bytes []byte
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	request := createRestRequest(requestType, "/payments", strings.NewReader(string(bytes)), params)

	switch requestType {
	case "PUT":
		api.PutPayment(responseWriter, request)
	case "POST":
		api.PostPayment(responseWriter, request)
	}

	return nil
}

func sendPostRequest(api RestApi, payload Payment, responseWriter rest.ResponseWriter) error {
	return sendPaymentsRequest(api, payload, responseWriter, "POST", nil)
}

func sendPutRequest(api RestApi, payload Payment, responseWriter rest.ResponseWriter) error {
	params := make(map[string]string)
	params["id"] = payload.ID
	return sendPaymentsRequest(api, payload, responseWriter, "PUT", params)
}

func fetchPayment(api RestApi, id string) (Payment, error) {
	var payment Payment
	responseWriter := &testResponseWriter{}

	params := make(map[string]string)
	params["id"] = id
	request := createRestRequest("GET", "/payments", strings.NewReader(""), params)
	api.GetPayment(responseWriter, request)

	buf := responseWriter.Read()
	err := json.Unmarshal(buf, &payment)

	return payment, err
}
