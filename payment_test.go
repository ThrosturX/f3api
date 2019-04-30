package f3api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

const DEFAULT_PAYMENT = `{"type":"Payment","id":"deadbeef-cab5-dad5-bad5-1337cafebabe","version":0,"organisation_id":"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb","attributes":{"amount":"100.21","beneficiary_party":{"account_name":"W Owens","account_number":"31926819","account_number_code":"BBAN","account_type":0,"address":"1 The Beneficiary Localtown SE2","bank_id":"403000","bank_id_code":"GBDSC","name":"Wilfred Jeremiah Owens"},"charges_information":{"bearer_code":"SHAR","sender_charges":[{"amount":"5.00","currency":"GBP"},{"amount":"10.00","currency":"USD"}],"receiver_charges_amount":"1.00","receiver_charges_currency":"USD"},"currency":"GBP","debtor_party":{"account_name":"EJ Brown Black","account_number":"GB29XABC10161234567801","account_number_code":"IBAN","address":"10 Debtor Crescent Sourcetown NE1","bank_id":"203301","bank_id_code":"GBDSC","name":"Emelia Jane Brown"},"end_to_end_reference":"Wil piano Jan","fx":{"contract_reference":"FX123","exchange_rate":"2.00000","original_amount":"200.42","original_currency":"USD"},"numeric_reference":"1002001","payment_id":"123456789012345678","payment_purpose":"Paying for goods/services","payment_scheme":"FPS","payment_type":"Credit","processing_date":"2017-01-18","reference":"Payment for Em's piano lessons","scheme_payment_sub_type":"InternetBanking","scheme_payment_type":"ImmediatePayment","sponsor_party":{"account_number":"56781234","bank_id":"123123","bank_id_code":"GBDSC"}}}`

// Utility for other tests
func defaultPayment() Payment {
	var p Payment

	err := json.Unmarshal([]byte(DEFAULT_PAYMENT), &p)
	if err != nil {
		panic(err)
	}

	return p
}

// Utility for checking if two json strings are "equal", borrowed from:
// https://bl.ocks.org/turtlemonvh/e4f7404e28387fadb8ad275a99596f67
func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

// Test idempotency between input example json and output json from payment resource
// NOTE: This test is quite strict, and will not allow any semantic difference
//       This test case helped reveal small details such as debtor_party not having an account_type field whereas beneficiary_party does have that field
func TestPaymentResourceJsonIdempotency(t *testing.T) {
	original := DEFAULT_PAYMENT

	// We'll start by unmarshaling the string into our payment resource struct
	var p Payment

	err := json.Unmarshal([]byte(original), &p)
	if err != nil {
		t.Fatal(err)
	}

	// We continue by marshaling the struct back into a string for comparison
	var bytes []byte
	var result string

	bytes, err = json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	result = string(bytes)

	// Compare the input and output strings with a comparison utility:
	ok, err := AreEqualJSON(result, original)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		// Something went wrong while testing, let's log the input and output to file for comparison with a tool such as `jq`
		err = ioutil.WriteFile("test_original.json", []byte(original), 0644)
		if err != nil {
			log.Println(err)
		}
		err = ioutil.WriteFile("test_result.json", []byte(result), 0644)
		if err != nil {
			log.Println(err)
		}
		t.Fatal("Idempotency test produced unequal output from input")
	}
}
