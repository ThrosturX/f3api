package f3api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Time format of the "processing_date" field of the "Payment" resource
// ISO 8601
const timeFmt string = "2006-01-02"

// Boxed Integer as a String
type StringedInt int64

// Unmarshal an integer string into an integer as StringedInt
func (si *StringedInt) UnmarshalJSON(buf []byte) error {
	ii, err := strconv.ParseInt(strings.Trim(string(buf), "\""), 10, 64)
	if err != nil {
		return err
	}
	*si = StringedInt(ii)
	return nil
}

// Marshal an integer into a string
func (si StringedInt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.d\"", si)), nil
}

// Boxed float64 for string as decimal amount marshaling/unmarshaling, 5 decimal places
type ExchangeRate float64

// Unarshal a float64 as ExchangeRate from a string
func (er *ExchangeRate) UnmarshalJSON(buf []byte) error {
	fl, err := strconv.ParseFloat(strings.Trim(string(buf), "\""), 64)
	if err != nil {
		return err
	}
	*er = ExchangeRate(fl)
	return nil
}

// Marshal a float64 into a string of five decimal places
func (er ExchangeRate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.5f\"", er)), nil
}

// Boxed float64 for string as decimal amount marshaling/unmarshaling, 2 decimal places
type FractionalAmount float64

// Unmarshal a float64 as FractionalAmount from a string
func (fa *FractionalAmount) UnmarshalJSON(buf []byte) error {
	fl, err := strconv.ParseFloat(strings.Trim(string(buf), "\""), 64)
	if err != nil {
		return err
	}
	*fa = FractionalAmount(fl)
	return nil
}

// Marshal a float64 into a string of two decimal places
func (fa FractionalAmount) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.2f\"", fa)), nil
}

// Boxed time.Time for marshaling/Unmarshaling timestamps in the timeFmt format
type Date struct {
	time.Time
}

// Unmarshal a time.Time object in the timeFmt format (ISO 8601) from a string
func (d *Date) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse(timeFmt, strings.Trim(string(buf), "\""))
	if err != nil {
		return err
	}
	d.Time = tt
	return nil
}

// Marshal a time.Time object in the timeFmt format (ISO 8601) into a string
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.Time.Format(timeFmt) + "\""), nil
}
