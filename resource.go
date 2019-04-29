package f3api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Time format of the "processing_date" field of the "Payment" resource
const timeFmt string = "2006-01-02"

// Boxed Integer as a String
type stringedInt int64

func (si *stringedInt) UnmarshalJSON(buf []byte) error {
	ii, err := strconv.ParseInt(strings.Trim(string(buf), "\""), 10, 64)
	if err != nil {
		return err
	}
	*si = stringedInt(ii)
	return nil
}

func (si stringedInt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.d\"", si)), nil
}

// Boxed float64 for string as decimal amount marshaling/unmarshaling, 5 decimal places
type exchangeRate float64

func (er *exchangeRate) UnmarshalJSON(buf []byte) error {
	fl, err := strconv.ParseFloat(strings.Trim(string(buf), "\""), 64)
	if err != nil {
		return err
	}
	*er = exchangeRate(fl)
	return nil
}

func (er exchangeRate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.5f\"", er)), nil
}

// Boxed float64 for string as decimal amount marshaling/unmarshaling, 2 decimal places
type fractionalAmount float64

func (fa *fractionalAmount) UnmarshalJSON(buf []byte) error {
	fl, err := strconv.ParseFloat(strings.Trim(string(buf), "\""), 64)
	if err != nil {
		return err
	}
	*fa = fractionalAmount(fl)
	return nil
}

func (fa fractionalAmount) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.2f\"", fa)), nil
}

// Boxed time.Time for marshaling/Unmarshaling timestamps in the timeFmt format
type date struct {
	time.Time
}

func (d *date) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse(timeFmt, strings.Trim(string(buf), "\""))
	if err != nil {
		return err
	}
	d.Time = tt
	return nil
}

func (d date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.Time.Format(timeFmt) + "\""), nil
}
