package f3api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
)

func createRestRequest(method string, urlStr string, body io.Reader, params map[string]string, t *testing.T) *rest.Request {
	origReq, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		t.Fatal(err)
	}
	return &rest.Request{
		origReq,
		params,
		map[string]interface{}{},
	}
}

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

// Checks whether or not posting a new payment resource increases the number of items in the collection
// NOTE: This method ALWAYS posts the same payment resource (same ID). Consider adding it as a parameter or implementing an ID generator if needed.
func checkPostedResourceIncreasesCollectionSize(api RestApi, t *testing.T) error {
	var (
		err              error
		oldList, newList []Payment
		oldLen, newLen   int
	)
	// we'll send the request for all payments twice, so lets create it now:
	getPayments := createRestRequest("GET", "/payments", strings.NewReader(""), nil, t)

	params := make(map[string]string)
	// we will also be sending a new json payment to be created
	postPayment := createRestRequest("POST", "/payments", strings.NewReader(DEFAULT_PAYMENT), params, t)

	responseWriter := &testResponseWriter{}

	// get the original length, call the API, check the response
	api.GetAllPayments(responseWriter, getPayments)

	buf := responseWriter.Read()

	err = json.Unmarshal(buf, &oldList)
	if err != nil {
		t.Fatal(err)
	}

	// post the new payment
	api.PostPayment(responseWriter, postPayment)

	// get the new length as before
	api.GetAllPayments(responseWriter, getPayments)

	buf = responseWriter.Read()

	err = json.Unmarshal(buf, &newList)
	if err != nil {
		t.Fatal(err)
	}

	oldLen = len(oldList)
	newLen = len(newList)
	if newLen != oldLen+1 {
		t.Fatalf("Lengths did not match up! Could not satisfy %d == %d + 1", oldLen, newLen)
	}

	// Pass
	return nil
}

func TestInMemApi(t *testing.T) {
	// This line will not compile nuless GormApi implements the RestApi interface
	var restApi RestApi = NewInMemApi()

	// some basic tests
	err := checkPostedResourceIncreasesCollectionSize(restApi, t)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: More tests

	// Pass
}
