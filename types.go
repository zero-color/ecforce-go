package ecforce

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Ptr returns a pointer to v. It helps populate optional request fields:
//
//	req := &admin.CustomerParams{Email: ecforce.Ptr("a@example.com")}
func Ptr[T any](v T) *T { return &v }

// String returns a pointer to s.
func String(s string) *string { return &s }

// Int returns a pointer to i.
func Int(i int) *int { return &i }

// Int64 returns a pointer to i.
func Int64(i int64) *int64 { return &i }

// Bool returns a pointer to b.
func Bool(b bool) *bool { return &b }

// Bool01 returns a pointer to a BoolInt with value b. Use it for request
// flags the API documents as 0/1 booleans.
func Bool01(b bool) *BoolInt { bi := BoolInt(b); return &bi }

// BoolInt is a boolean that marshals to the 0/1 integers the ecforce API uses
// for request flags, and unmarshals from 0/1, true/false, or their string
// forms.
type BoolInt bool

// MarshalJSON implements json.Marshaler.
func (b BoolInt) MarshalJSON() ([]byte, error) {
	if b {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BoolInt) UnmarshalJSON(data []byte) error {
	s := string(bytes.Trim(bytes.TrimSpace(data), `"`))
	switch s {
	case "1", "true":
		*b = true
	case "0", "false", "null", "":
		*b = false
	default:
		return fmt.Errorf("ecforce: cannot unmarshal %q into BoolInt", data)
	}
	return nil
}

// timeLayouts are the timestamp formats the API is known to emit or accept,
// tried in order.
var timeLayouts = []string{
	"2006/01/02 15:04:05",
	"2006-01-02 15:04:05",
	"2006/01/02",
	"2006-01-02",
	time.RFC3339,
}

// Time wraps time.Time to handle the ecforce API's timestamp format
// ("2006/01/02 15:04:05"). JSON null unmarshals to the zero Time.
type Time struct {
	time.Time
}

// NewTime returns a pointer to a Time wrapping t.
func NewTime(t time.Time) *Time { return &Time{Time: t} }

// MarshalJSON implements json.Marshaler using the API's timestamp format.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	return strconv.AppendQuote(nil, t.Format("2006/01/02 15:04:05")), nil
}

// UnmarshalJSON implements json.Unmarshaler, accepting the formats the API
// emits across endpoints.
func (t *Time) UnmarshalJSON(data []byte) error {
	s := string(bytes.TrimSpace(data))
	if s == "null" || s == `""` {
		t.Time = time.Time{}
		return nil
	}
	unquoted, err := strconv.Unquote(s)
	if err != nil {
		return fmt.Errorf("ecforce: cannot unmarshal %s into Time", data)
	}
	for _, layout := range timeLayouts {
		if parsed, err := time.ParseInLocation(layout, unquoted, jst); err == nil {
			t.Time = parsed
			return nil
		}
	}
	return fmt.Errorf("ecforce: cannot parse time %q", unquoted)
}

// Equal reports whether t and u represent the same time instant.
func (t Time) Equal(u Time) bool { return t.Time.Equal(u.Time) }

// Date wraps time.Time for date-only fields, using the "2006-01-02" format
// the API documents for values such as birth dates.
type Date struct {
	time.Time
}

// NewDate returns a pointer to a Date for the given year, month and day.
func NewDate(year int, month time.Month, day int) *Date {
	return &Date{Time: time.Date(year, month, day, 0, 0, 0, 0, jst)}
}

// MarshalJSON implements json.Marshaler using the yyyy-mm-dd format.
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return strconv.AppendQuote(nil, d.Format("2006-01-02")), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Date) UnmarshalJSON(data []byte) error {
	var t Time
	if err := t.UnmarshalJSON(data); err != nil {
		return err
	}
	d.Time = t.Time
	return nil
}

// Equal reports whether d and u represent the same time instant.
func (d Date) Equal(u Date) bool { return d.Time.Equal(u.Time) }

// jst is the timezone ecforce shops operate in; timestamps in API responses
// carry no offset.
var jst = time.FixedZone("Asia/Tokyo", 9*60*60)
