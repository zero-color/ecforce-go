package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zero-color/ecforce-go"
)

// setup returns a test customer client wired to a test server, and the mux to
// register handlers on.
func setup(t *testing.T) (*Client, *http.ServeMux) {
	t.Helper()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL, ecforce.WithToken("test-token"))
	if err != nil {
		t.Fatal(err)
	}
	return client, mux
}

func TestSessionsService_SignIn(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/customers/sign_in.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["customer"]["email"] != "me@example.com" {
			t.Errorf("body = %v", body)
		}
		fmt.Fprint(w, `{"id": 1, "email": "me@example.com", "authentication_token": "issued-token"}`)
	})

	session, _, err := client.Sessions.SignIn(context.Background(), "me@example.com", "pw")
	if err != nil {
		t.Fatal(err)
	}
	if session.AuthenticationToken != "issued-token" {
		t.Errorf("unexpected session: %+v", session)
	}
	if got := client.Token(); got != "issued-token" {
		t.Errorf("client token = %q, want issued-token", got)
	}
}

func TestCustomerService_Get(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/customer.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("lighter"); got != "1" {
			t.Errorf("lighter = %q, want 1", got)
		}
		fmt.Fprint(w, `{"data": {"id": "30962", "type": "customer", "attributes": {
			"id": 30962, "number": "02ae6cdbef", "state": "member", "email": "ok@example.com",
			"point": 5, "optin": false}}}`)
	})

	me, _, err := client.Customer.Get(context.Background(), &ecforce.GetOptions{Lighter: ecforce.Bool(true)})
	if err != nil {
		t.Fatal(err)
	}
	if me.Attributes.Email != "ok@example.com" || me.Attributes.ID != 30962 {
		t.Errorf("unexpected customer: %+v", me.Attributes)
	}
}

func TestMultiPassService_Get(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/customers/multi_pass/tok123.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "email": "hoge@example.jp", "authentication_token": "t",
			"login_url": "https://x.com/customers/sign_in/multipass/ONE_TIME_TOKEN"}`)
	})

	mp, _, err := client.MultiPass.Get(context.Background(), "tok123")
	if err != nil {
		t.Fatal(err)
	}
	if mp.LoginURL == "" || mp.Email != "hoge@example.jp" {
		t.Errorf("unexpected multipass: %+v", mp)
	}
}

func TestSubsOrdersService_Update_dataArray(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/customer/subs_orders/34123.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data": [{"id": "34123", "type": "sub_order", "attributes": {
			"id": 34123, "number": "xxddg32dg", "state": "active"}}]}`)
	})

	updated, _, err := client.SubsOrders.Update(context.Background(), 34123, &SubsOrderParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(updated) != 1 || updated[0].Attributes.Number != "xxddg32dg" {
		t.Errorf("unexpected result: %+v", updated)
	}
}

func TestSubsOrdersService_DestroyOrderItem(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("DELETE /api/v2/customer/subs_orders/1/order_items/2.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.SubsOrders.DestroyOrderItem(context.Background(), 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want 204", resp.StatusCode)
	}
}
