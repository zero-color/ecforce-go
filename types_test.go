package ecforce

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		in   string
		want time.Time
	}{
		{`"2023/02/14 16:22:07"`, time.Date(2023, 2, 14, 16, 22, 7, 0, jst)},
		{`"2022-03-10 00:00:00"`, time.Date(2022, 3, 10, 0, 0, 0, 0, jst)},
		{`"2022-12-20T13:52:36.000+09:00"`, time.Date(2022, 12, 20, 13, 52, 36, 0, jst)},
		{`"2019/06/08"`, time.Date(2019, 6, 8, 0, 0, 0, 0, jst)},
		{`"2019-06-08"`, time.Date(2019, 6, 8, 0, 0, 0, 0, jst)},
		{`null`, time.Time{}},
	}
	for _, tt := range tests {
		var got Time
		if err := json.Unmarshal([]byte(tt.in), &got); err != nil {
			t.Errorf("Unmarshal(%s) error: %v", tt.in, err)
			continue
		}
		if !got.Time.Equal(tt.want) {
			t.Errorf("Unmarshal(%s) = %v, want %v", tt.in, got.Time, tt.want)
		}
	}

	var bad Time
	if err := json.Unmarshal([]byte(`"not a time"`), &bad); err == nil {
		t.Error("Unmarshal of invalid time succeeded, want error")
	}
}

func TestTimeMarshalJSON(t *testing.T) {
	b, err := json.Marshal(Time{Time: time.Date(2023, 2, 14, 16, 22, 7, 0, jst)})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(b), `"2023/02/14 16:22:07"`; got != want {
		t.Errorf("Marshal = %s, want %s", got, want)
	}

	b, err = json.Marshal(Time{})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(b), "null"; got != want {
		t.Errorf("Marshal zero = %s, want %s", got, want)
	}
}

func TestDateMarshalJSON(t *testing.T) {
	b, err := json.Marshal(NewDate(2022, time.March, 10))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(b), `"2022-03-10"`; got != want {
		t.Errorf("Marshal = %s, want %s", got, want)
	}
}

func TestBoolIntRoundTrip(t *testing.T) {
	for _, tt := range []struct {
		in   string
		want BoolInt
	}{
		{`1`, true}, {`0`, false}, {`true`, true}, {`false`, false},
		{`"1"`, true}, {`"0"`, false}, {`null`, false},
	} {
		var got BoolInt
		if err := json.Unmarshal([]byte(tt.in), &got); err != nil {
			t.Errorf("Unmarshal(%s) error: %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Unmarshal(%s) = %v, want %v", tt.in, got, tt.want)
		}
	}

	b, _ := json.Marshal(struct {
		On  BoolInt `json:"on"`
		Off BoolInt `json:"off"`
	}{On: true, Off: false})
	if got, want := string(b), `{"on":1,"off":0}`; got != want {
		t.Errorf("Marshal = %s, want %s", got, want)
	}

	var bad BoolInt
	if err := json.Unmarshal([]byte(`"maybe"`), &bad); err == nil {
		t.Error("Unmarshal of invalid BoolInt succeeded, want error")
	}
}
