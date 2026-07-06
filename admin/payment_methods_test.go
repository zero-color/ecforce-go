package admin

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/zero-color/ecforce-go"
)

func TestPaymentMethodsService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/payment_methods.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[active_eq]"); got != "1" {
			t.Errorf("q[active_eq] = %q, want 1", got)
		}
		fmt.Fprint(w, `{"data": [{"id": "1", "type": "payment_method", "attributes": {
			"id": 1, "name": "クレジットカード一括", "description": null,
			"description_mobile": null, "active": true, "amount_upper_limit": null,
			"amount_under_limit": 0, "display_on_mypage": true,
			"display_on_admin_oi_page": true, "use_credit_card_validity_check": false,
			"cv_offer_realtime_payment": false, "use_o_plux": false,
			"use_spider_af": true}}]}`)
	})

	methods, _, err := client.PaymentMethods.List(context.Background(), &ecforce.ListOptions{
		Q: ecforce.Query{"active_eq": true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(methods) != 1 {
		t.Fatalf("got %d payment methods, want 1", len(methods))
	}
	attrs := methods[0].Attributes
	if attrs.ID != 1 || attrs.Name != "クレジットカード一括" || !attrs.Active {
		t.Errorf("unexpected payment method: %+v", attrs)
	}
	if !attrs.DisplayOnMypage || attrs.UseCreditCardValidityCheck || !attrs.UseSpiderAf {
		t.Errorf("bool fields: mypage=%v validity=%v spiderAf=%v",
			attrs.DisplayOnMypage, attrs.UseCreditCardValidityCheck, attrs.UseSpiderAf)
	}
}
