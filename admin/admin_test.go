package admin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zero-color/ecforce-go"
)

// setup returns a test admin client wired to a test server, and the mux to
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
	mux.HandleFunc("POST /api/v2/admins/sign_in.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": 1, "email": "hoge@example.jp", "authentication_token": "issued-token"}`)
	})

	session, _, err := client.Sessions.SignIn(context.Background(), "hoge@example.jp", "pw")
	if err != nil {
		t.Fatal(err)
	}
	if session.AuthenticationToken != "issued-token" || session.ID != 1 {
		t.Errorf("unexpected session: %+v", session)
	}
	if got := client.Token(); got != "issued-token" {
		t.Errorf("client token = %q, want issued-token", got)
	}
}

func TestSessionsService_SignOut(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admins/sign_out.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	if _, err := client.Sessions.SignOut(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := client.Token(); got != "" {
		t.Errorf("client token = %q, want empty after sign out", got)
	}
}

func TestSexesService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/sexes.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[name_eq]"); got != "男性" {
			t.Errorf("q[name_eq] = %q", got)
		}
		fmt.Fprint(w, `{"data": [{"id": "1", "type": "sex", "attributes": {
			"id": 1, "name": "男性", "position": 0, "state": "active",
			"created_at": "2022-12-20T13:52:36.000+09:00", "deleted_at": null}}]}`)
	})

	sexes, _, err := client.Sexes.List(context.Background(), &ecforce.ListOptions{
		Q: ecforce.Query{"name_eq": "男性"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(sexes) != 1 || sexes[0].Attributes.Name != "男性" || sexes[0].Attributes.ID != 1 {
		t.Errorf("unexpected sexes: %+v", sexes)
	}
	if sexes[0].Attributes.DeletedAt != nil && !sexes[0].Attributes.DeletedAt.IsZero() {
		t.Errorf("DeletedAt = %v, want zero", sexes[0].Attributes.DeletedAt)
	}
}

func TestShippingCarriersService_Get(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/shipping_carriers/99.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data": {"id": "99", "type": "shipping_carrier", "attributes": {
			"id": 99, "name": "外部運輸", "state": "inactive",
			"pickup_locations": [{"id": 1, "pickup_location_name": "宅配ボックス"}],
			"pickup_location_primary_only": 0, "enable_doorbell_option": true}}}`)
	})

	carrier, _, err := client.ShippingCarriers.Get(context.Background(), 99, nil)
	if err != nil {
		t.Fatal(err)
	}
	attrs := carrier.Attributes
	if attrs.Name != "外部運輸" || len(attrs.PickupLocations) != 1 {
		t.Errorf("unexpected carrier: %+v", attrs)
	}
	if attrs.PickupLocationPrimaryOnly || !attrs.EnableDoorbellOption {
		t.Errorf("bool fields: primaryOnly=%v doorbell=%v", attrs.PickupLocationPrimaryOnly, attrs.EnableDoorbellOption)
	}
}
