package f3api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Time format of the "processing_Date" field of the "Payment" resource
const timeFmt string = "2006-01-02"

// Boxed Integer as a String
type StringedInt int64

func (si *StringedInt) UnmarshalJSON(buf []byte) error {
	ii, err := strconv.ParseInt(strings.Trim(string(buf), "\""), 10, 64)
	if err != nil {
		return err
	}
	*si = StringedInt(ii)
	return nil
}

func (si StringedInt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.d\"", si)), nil
}

// Boxed float64 for string as decimal amount marshaling/unmarshaling, 5 decimal places
type ExchangeRate float64

func (er *ExchangeRate) UnmarshalJSON(buf []byte) error {
	fl, err := strconv.ParseFloat(strings.Trim(string(buf), "\""), 64)
	if err != nil {
		return err
	}
	*er = ExchangeRate(fl)
	return nil
}

func (er ExchangeRate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.5f\"", er)), nil
}

// Boxed float64 for string as decimal amount marshaling/unmarshaling, 2 decimal places
type FractionalAmount float64

func (fa *FractionalAmount) UnmarshalJSON(buf []byte) error {
	fl, err := strconv.ParseFloat(strings.Trim(string(buf), "\""), 64)
	if err != nil {
		return err
	}
	*fa = FractionalAmount(fl)
	return nil
}

func (fa FractionalAmount) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%.2f\"", fa)), nil
}

// Boxed time.Time for marshaling/Unmarshaling timestamps in the timeFmt format
type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse(timeFmt, strings.Trim(string(buf), "\""))
	if err != nil {
		return err
	}
	d.Time = tt
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.Time.Format(timeFmt) + "\""), nil
}
