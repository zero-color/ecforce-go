package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/zero-color/ecforce-go"
)

func TestSessionsService_SignOut(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/customers/sign_out.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != `Token token="test-token"` {
			t.Errorf("Authorization = %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.Sessions.SignOut(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want 204", resp.StatusCode)
	}
	if got := client.Token(); got != "" {
		t.Errorf("client token = %q, want cleared", got)
	}
}

func TestRegistrationService_SignUp(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/customers.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		params, ok := body["customer"]
		if !ok {
			t.Fatalf(`body missing "customer" wrapper: %v`, body)
		}
		if params["email"] != "hoge@example.jp" || params["state"] != "member" {
			t.Errorf("customer params = %v", params)
		}
		addr, _ := params["billing_address_attributes"].(map[string]any)
		if addr["name01"] != "てすと" {
			t.Errorf("billing_address_attributes = %v", addr)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"data": {"id": "83897", "type": "customer", "attributes": {
			"id": 83897, "authentication_token": "CK8Xs-MNDg6PqDGjHAkeWgdAd5mHDN3T",
			"number": "dfe332987f", "state": "member", "human_state_name": "会員",
			"email": "hoge@example.jp", "birth": "2014/12/24", "optin": false,
			"link_number": "111", "created_at": "2020/06/06 18:16:46"}}}`)
	})

	created, resp, err := client.Registration.SignUp(context.Background(), &SignUpParams{
		Email:    ecforce.String("hoge@example.jp"),
		State:    ecforce.String("member"),
		Password: ecforce.String("secret"),
		SexID:    ecforce.Int64(1),
		Birth:    ecforce.String("2014/12/24"),
		BillingAddressAttributes: &AddressParams{
			Name01: ecforce.String("てすと"),
			Name02: ecforce.String("てすと"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want 201", resp.StatusCode)
	}
	if created.ID != "83897" || created.Attributes.Email != "hoge@example.jp" {
		t.Errorf("unexpected customer: %+v", created)
	}
	if created.Attributes.AuthenticationToken != "CK8Xs-MNDg6PqDGjHAkeWgdAd5mHDN3T" {
		t.Errorf("authentication_token = %q", created.Attributes.AuthenticationToken)
	}
}

func TestCustomerService_Update(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/customer.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		params, ok := body["customer"]
		if !ok {
			t.Fatalf(`body missing "customer" wrapper: %v`, body)
		}
		if params["email"] != "test@example.com" {
			t.Errorf("email = %v", params["email"])
		}
		if params["optin"] != float64(1) {
			t.Errorf("optin = %v, want 1", params["optin"])
		}
		fmt.Fprint(w, `{"data": {"id": "3172", "type": "customer", "attributes": {
			"id": 3172, "number": "6a54e6ebc5", "state": "member",
			"email": "test@example.com", "optin": true, "point": 92,
			"link_number": "111"}}}`)
	})

	updated, _, err := client.Customer.Update(context.Background(), &CustomerParams{
		Email: ecforce.String("test@example.com"),
		Optin: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Attributes.Email != "test@example.com" || !updated.Attributes.Optin {
		t.Errorf("unexpected customer: %+v", updated.Attributes)
	}
}

func TestInviteCodesService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/customer/invite_codes.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data": [{"id": "1", "type": "invite_code", "attributes": {
			"id": 1, "number": "xxxxx", "used_times": 5,
			"available_invited_times": 0, "max_limit": 5}}],
			"meta": {"total_count": 1, "page": 1, "per": 100, "count": 1, "total_pages": 1}}`)
	})

	codes, resp, err := client.InviteCodes.List(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(codes) != 1 || codes[0].Attributes.Number != "xxxxx" || codes[0].Attributes.UsedTimes != 5 {
		t.Errorf("unexpected invite codes: %+v", codes)
	}
	if resp.Meta == nil || resp.Meta.TotalCount != 1 || resp.Meta.TotalPages != 1 {
		t.Errorf("unexpected meta: %+v", resp.Meta)
	}
	if resp.HasNextPage() {
		t.Error("HasNextPage() = true, want false")
	}
}

func TestInviteCodesService_Create(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/customer/invite_codes.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data": {"id": "1", "type": "invite_code", "attributes": {
			"id": 1, "number": "xxxxx", "used_times": 0,
			"available_invited_times": 5, "max_limit": 5}}}`)
	})

	code, _, err := client.InviteCodes.Create(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if code.Attributes.Number != "xxxxx" || code.Attributes.MaxLimit != 5 {
		t.Errorf("unexpected invite code: %+v", code.Attributes)
	}
}

func TestSubsOrdersService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/customer/subs_orders.json", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("q[state_eq]"); got != "active" {
			t.Errorf("q[state_eq] = %q, want active", got)
		}
		if got := q.Get("per"); got != "1" {
			t.Errorf("per = %q, want 1", got)
		}
		fmt.Fprint(w, `{"data": [{"id": "4909", "type": "sub_order", "attributes": {
			"id": 4909, "number": "a2886ff4f2", "customer_number": "02ae6cdbef",
			"times": 1, "orders_count": 1, "state": "active", "human_state": "有効",
			"subtotal": 100, "payment_schedule_locked": 1,
			"scheduled_to_be_delivered_at": "2019/07/01 00:00:00"}}],
			"meta": {"total_count": 3505, "page": 1, "per": 1, "count": 1, "total_pages": 3505}}`)
	})

	orders, resp, err := client.SubsOrders.List(context.Background(), &ecforce.ListOptions{
		Per: 1,
		Q:   ecforce.Query{"state_eq": "active"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(orders) != 1 || orders[0].Attributes.Number != "a2886ff4f2" || orders[0].Attributes.State != "active" {
		t.Errorf("unexpected subs orders: %+v", orders)
	}
	if resp.Meta == nil || resp.Meta.TotalCount != 3505 || resp.Meta.Page != 1 {
		t.Errorf("unexpected meta: %+v", resp.Meta)
	}
	if !resp.HasNextPage() {
		t.Error("HasNextPage() = false, want true")
	}
	if got := resp.NextPage(); got != 2 {
		t.Errorf("NextPage() = %d, want 2", got)
	}
}

func TestSubsOrdersService_Get(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/customer/subs_orders/4909.json", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("lighter"); got != "1" {
			t.Errorf("lighter = %q, want 1", got)
		}
		if got := q.Get("include"); got != "order_items,payment" {
			t.Errorf("include = %q, want order_items,payment", got)
		}
		fmt.Fprint(w, `{"data": {"id": "4909", "type": "sub_order", "attributes": {
			"id": 4909, "number": "a2886ff4f2", "state": "active", "subtotal": 100,
			"available_payment_schedules": [{"value": "date", "ja": "日付で指定",
				"scheduled_to_be_delivered_every_x_month": [
					{"value": 1, "ja": "1ヶ月"}, {"value": 2, "ja": "2ヶ月"}]}],
			"payment_schedule": "1ヶ月ごとの1日に配送", "payment_schedule_locked": 1}}}`)
	})

	order, _, err := client.SubsOrders.Get(context.Background(), 4909, &ecforce.GetOptions{
		Lighter: ecforce.Bool(true),
		Include: []string{"order_items", "payment"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if order.Attributes.ID != 4909 || order.Attributes.Number != "a2886ff4f2" {
		t.Errorf("unexpected subs order: %+v", order.Attributes)
	}
	schedules := order.Attributes.AvailablePaymentSchedules
	if len(schedules) != 1 || schedules[0].Value != "date" || len(schedules[0].ScheduledToBeDeliveredEveryXMonth) != 2 {
		t.Errorf("unexpected available_payment_schedules: %+v", schedules)
	}
}

func TestSubsOrdersService_AddOrderItem(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/customer/subs_orders/4909/order_items.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		params, ok := body["order_item"]
		if !ok {
			t.Fatalf(`body missing "order_item" wrapper: %v`, body)
		}
		if params["variant_id"] != float64(10) || params["quantity"] != float64(1) {
			t.Errorf("order_item params = %v", params)
		}
		if params["onetime"] != float64(1) {
			t.Errorf("onetime = %v, want 1", params["onetime"])
		}
		fmt.Fprint(w, `{"data": [{"id": "34123", "type": "order_item", "attributes": {
			"id": 34123, "variant_id": 10, "product_number": "TA-001",
			"product_name": "定期テストA", "variant_sku": "TA-001-001",
			"price": 100, "quantity": 1, "tax_rate": 10}}]}`)
	})

	items, _, err := client.SubsOrders.AddOrderItem(context.Background(), 4909, &OrderItemParams{
		VariantID: ecforce.Int64(10),
		Quantity:  ecforce.Ptr(1),
		Onetime:   ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Attributes.VariantID != 10 || items[0].Attributes.VariantSKU != "TA-001-001" {
		t.Errorf("unexpected order items: %+v", items)
	}
}

func TestSubsOrdersService_UpdateOrderItem(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/customer/subs_orders/4909/order_items/34123.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("include"); got != "variant" {
			t.Errorf("include = %q, want variant", got)
		}
		var body map[string]map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		params, ok := body["order_item"]
		if !ok {
			t.Fatalf(`body missing "order_item" wrapper: %v`, body)
		}
		if params["variant_id"] != float64(10) || params["quantity"] != float64(2) {
			t.Errorf("order_item params = %v", params)
		}
		fmt.Fprint(w, `{"data": [{"id": "34123", "type": "order_item", "attributes": {
			"id": 34123, "variant_id": 10, "product_number": "TA-001",
			"variant_sku": "TA-001-001", "price": 100, "quantity": 2}}],
			"included": [{"id": "10", "type": "variant", "attributes": {
				"id": 10, "sku": "TA-001-001", "product_id": 1}}]}`)
	})

	items, resp, err := client.SubsOrders.UpdateOrderItem(context.Background(), 4909, 34123, &OrderItemParams{
		VariantID: ecforce.Int64(10),
		Quantity:  ecforce.Ptr(2),
	}, &ecforce.GetOptions{Include: []string{"variant"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Attributes.Quantity != 2 {
		t.Errorf("unexpected order items: %+v", items)
	}
	if ecforce.FindIncluded(resp.Included, "variant", "10") == nil {
		t.Errorf("included missing variant 10: %+v", resp.Included)
	}
}

func TestPaymentScheduleValue_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		data string
		want PaymentScheduleValue
	}{
		{name: "ja as string", data: `{"value": 1, "ja": "1ヶ月"}`, want: PaymentScheduleValue{Value: 1, Ja: "1ヶ月"}},
		{name: "ja as bare number", data: `{"value": 2, "ja": 2}`, want: PaymentScheduleValue{Value: 2, Ja: "2"}},
		{name: "ja null", data: `{"value": 3, "ja": null}`, want: PaymentScheduleValue{Value: 3, Ja: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got PaymentScheduleValue
			if err := json.Unmarshal([]byte(tt.data), &got); err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
