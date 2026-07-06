package ecforce

import (
	"testing"
	"time"
)

func TestListOptionsValues(t *testing.T) {
	opts := &ListOptions{
		Lighter: Bool(true),
		Include: []string{"billing_address", "orders"},
		Page:    2,
		Per:     50,
		Sort:    []string{"-created_at", "id"},
		Q: Query{
			"email_cont":      "@example.com",
			"id_in":           []int64{1, 2, 3},
			"with_deleted":    true,
			"created_at_gteq": time.Date(2024, 1, 2, 3, 4, 5, 0, jst),
		},
	}
	v := opts.Values()

	for key, want := range map[string]string{
		"lighter":            "1",
		"include":            "billing_address,orders",
		"page":               "2",
		"per":                "50",
		"sort":               "-created_at,id",
		"q[email_cont]":      "@example.com",
		"q[with_deleted]":    "1",
		"q[created_at_gteq]": "2024-01-02 03:04:05",
	} {
		if got := v.Get(key); got != want {
			t.Errorf("Values()[%q] = %q, want %q", key, got, want)
		}
	}
	if got := v["q[id_in][]"]; len(got) != 3 || got[0] != "1" || got[2] != "3" {
		t.Errorf("Values()[q[id_in][]] = %v, want [1 2 3]", got)
	}
}

func TestListOptionsValuesNil(t *testing.T) {
	var opts *ListOptions
	if v := opts.Values(); v != nil {
		t.Errorf("nil options Values() = %v, want nil", v)
	}
	var gopts *GetOptions
	if v := gopts.Values(); v != nil {
		t.Errorf("nil get options Values() = %v, want nil", v)
	}
}

func TestGetOptionsValues(t *testing.T) {
	v := (&GetOptions{Lighter: Bool(false), Include: []string{"notes"}}).Values()
	if got := v.Get("lighter"); got != "0" {
		t.Errorf("lighter = %q, want 0", got)
	}
	if got := v.Get("include"); got != "notes" {
		t.Errorf("include = %q, want notes", got)
	}
}
