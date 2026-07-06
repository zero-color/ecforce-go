package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/zero-color/ecforce-go"
)

func TestCustomersService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/customers.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[email_eq]"); got != "ok@example.com" {
			t.Errorf("q[email_eq] = %q", got)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("page = %q", got)
		}
		fmt.Fprint(w, `{
			"data": [{"id": "30962", "type": "customer", "attributes": {
				"id": 30962, "number": "02ae6cdbef", "state": "member",
				"human_state_name": "会員", "email": "ok@example.com",
				"buy_times": 1, "buy_total": 799, "point": 5,
				"optin": false, "blacklist": true, "blacklist_reasons": "重複注文",
				"created_at": "2019/06/08 05:47:17", "deleted_at": null}}],
			"meta": {"total_count": 29923, "page": 1, "per": 1, "count": 1, "total_pages": 29923}}`)
	})

	customers, resp, err := client.Customers.List(context.Background(), &ecforce.ListOptions{
		Page: 1,
		Q:    ecforce.Query{"email_eq": "ok@example.com"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(customers) != 1 || customers[0].ID != "30962" {
		t.Fatalf("unexpected customers: %+v", customers)
	}
	attrs := customers[0].Attributes
	if attrs.ID != 30962 || attrs.Email != "ok@example.com" || attrs.Number != "02ae6cdbef" {
		t.Errorf("unexpected attributes: %+v", attrs)
	}
	if !attrs.Blacklist || attrs.BuyTimes != 1 {
		t.Errorf("blacklist=%v buy_times=%d", attrs.Blacklist, attrs.BuyTimes)
	}
	if resp.Meta == nil || resp.Meta.TotalCount != 29923 {
		t.Errorf("unexpected meta: %+v", resp.Meta)
	}
	if !resp.HasNextPage() || resp.NextPage() != 2 {
		t.Errorf("HasNextPage=%v NextPage=%d", resp.HasNextPage(), resp.NextPage())
	}
}

func TestCustomersService_Get(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/customers/3172.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("include"); got != "billing_address" {
			t.Errorf("include = %q", got)
		}
		fmt.Fprint(w, `{"data": {"id": "3172", "type": "customer", "attributes": {
			"id": 3172, "number": "6a54e6ebc5", "state": "member",
			"human_state_name": "会員", "email": "test@example.com",
			"point": 92, "optin": true, "link_number": "111",
			"created_at": "2017/07/19 00:48:39", "deleted_at": null}}}`)
	})

	customer, _, err := client.Customers.Get(context.Background(), 3172, &ecforce.GetOptions{
		Include: []string{"billing_address"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if customer.ID != "3172" || customer.Type != "customer" {
		t.Errorf("unexpected resource: id=%q type=%q", customer.ID, customer.Type)
	}
	attrs := customer.Attributes
	if attrs.ID != 3172 || attrs.Email != "test@example.com" || attrs.Point != 92 || !attrs.Optin {
		t.Errorf("unexpected attributes: %+v", attrs)
	}
}

func TestCustomersService_Create(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/customers.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := body["check_duplicate_link_numbers"]; got != float64(1) {
			t.Errorf("check_duplicate_link_numbers = %v", got)
		}
		customer, _ := body["customer"].(map[string]any)
		if customer == nil {
			t.Fatalf("customer wrapper missing: %v", body)
		}
		if got := customer["email"]; got != "test@example.com" {
			t.Errorf("customer.email = %v", got)
		}
		if got := customer["state"]; got != "member" {
			t.Errorf("customer.state = %v", got)
		}
		if got := customer["optin"]; got != float64(1) {
			t.Errorf("customer.optin = %v", got)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"data": {"id": "77", "type": "customer", "attributes": {
			"id": 77, "number": "7cf6313e55", "state": "member",
			"human_state_name": "会員", "email": "test@example.com",
			"optin": true, "link_number": "111",
			"created_at": "2023/02/14 16:22:07", "deleted_at": null}}}`)
	})

	customer, _, err := client.Customers.Create(context.Background(), &CustomerCreateRequest{
		Customer: &CustomerParams{
			Email: ecforce.String("test@example.com"),
			State: ecforce.String("member"),
			Optin: ecforce.Bool01(true),
		},
		CheckDuplicateLinkNumbers: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if customer.ID != "77" || customer.Attributes.Email != "test@example.com" {
		t.Errorf("unexpected customer: %+v", customer)
	}
	if customer.Attributes.LinkNumber != "111" {
		t.Errorf("link_number = %q", customer.Attributes.LinkNumber)
	}
}

func TestCustomersService_Update(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/customers/77.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		customer, _ := body["customer"].(map[string]any)
		if customer == nil {
			t.Fatalf("customer wrapper missing: %v", body)
		}
		if got := customer["email"]; got != "updated@example.com" {
			t.Errorf("customer.email = %v", got)
		}
		if got := body["update_member_groups"]; got != float64(1) {
			t.Errorf("update_member_groups = %v", got)
		}
		if got := fmt.Sprint(body["member_group_ids"]); got != "[1 2]" {
			t.Errorf("member_group_ids = %v", got)
		}
		fmt.Fprint(w, `{"data": {"id": "77", "type": "customer", "attributes": {
			"id": 77, "state": "member", "email": "updated@example.com",
			"updated_at": "2023/02/15 10:00:00"}}}`)
	})

	customer, _, err := client.Customers.Update(context.Background(), 77, &CustomerUpdateRequest{
		Customer: &CustomerParams{
			Email: ecforce.String("updated@example.com"),
		},
		UpdateMemberGroups: ecforce.Bool01(true),
		MemberGroupIDs:     []int64{1, 2},
	})
	if err != nil {
		t.Fatal(err)
	}
	if customer.Attributes.ID != 77 || customer.Attributes.Email != "updated@example.com" {
		t.Errorf("unexpected customer: %+v", customer.Attributes)
	}
}

func TestCustomersService_BulkCreate(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/customers/bulk_create.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		customers, _ := body["customers"].([]any)
		if len(customers) != 2 {
			t.Fatalf("customers = %v", body["customers"])
		}
		first, _ := customers[0].(map[string]any)
		if got := first["email"]; got != "hoge@xxxx.yyyy" {
			t.Errorf("customers[0].email = %v", got)
		}
		fmt.Fprint(w, `{"id": 188, "job_id": "RsHBP5Z2DLoNJA"}`)
	})

	job, _, err := client.Customers.BulkCreate(context.Background(), &CustomerBulkRequest{
		Customers: []*CustomerParams{
			{Email: ecforce.String("hoge@xxxx.yyyy"), State: ecforce.String("member")},
			{Email: ecforce.String("hogehoge@xxxx.yyyy"), State: ecforce.String("member")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if job.ID != 188 || job.JobID != "RsHBP5Z2DLoNJA" {
		t.Errorf("unexpected job: %+v", job)
	}
}

func TestCustomersService_BulkUpdate(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/customers/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		customers, _ := body["customers"].([]any)
		if len(customers) != 2 {
			t.Fatalf("customers = %v", body["customers"])
		}
		first, _ := customers[0].(map[string]any)
		if got := first["id"]; got != float64(30) {
			t.Errorf("customers[0].id = %v", got)
		}
		if got := first["mail_delivery_stop"]; got != float64(1) {
			t.Errorf("customers[0].mail_delivery_stop = %v", got)
		}
		fmt.Fprint(w, `{
			"success": [30],
			"failure": [31],
			"errors": [{"id": 31, "code": "ACU0000", "message": "エラーが発生しました。",
				"errors": [{"code": "ACU3999", "message": "Eメールはすでに存在します"}]}]}`)
	})

	result, _, err := client.Customers.BulkUpdate(context.Background(), &CustomerBulkRequest{
		Customers: []*CustomerParams{
			{ID: ecforce.Int64(30), MailDeliveryStop: ecforce.Bool01(true)},
			{ID: ecforce.Int64(31), Email: ecforce.String("dup@example.com")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 1 || result.Success[0] != 30 {
		t.Errorf("success = %v", result.Success)
	}
	if len(result.Failure) != 1 || result.Failure[0] != 31 {
		t.Errorf("failure = %v", result.Failure)
	}
	if len(result.Errors) != 1 || result.Errors[0].ID != 31 || result.Errors[0].Errors[0].Code != "ACU3999" {
		t.Errorf("errors = %+v", result.Errors)
	}
}

func TestCustomersService_BulkDestroy(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("DELETE /api/v2/admin/customers/bulk_destroy.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprint(body["customer_ids"]); got != "[1 2 3]" {
			t.Errorf("customer_ids = %v", got)
		}
		fmt.Fprint(w, `{
			"success": [3],
			"failure": [1, 2],
			"errors": [
				{"id": 1, "code": "ACU0000", "message": "エラーが発生しました。",
					"errors": [{"code": "ACU3999", "message": "仮売上状態の受注が存在します"}]},
				{"id": 2, "code": "ACU0000", "message": "エラーが発生しました。",
					"errors": [{"code": "ACU3999", "message": "アクティブな定期受注が存在します"}]}]}`)
	})

	result, _, err := client.Customers.BulkDestroy(context.Background(), []int64{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 1 || result.Success[0] != 3 {
		t.Errorf("success = %v", result.Success)
	}
	if len(result.Failure) != 2 || result.Failure[0] != 1 || result.Failure[1] != 2 {
		t.Errorf("failure = %v", result.Failure)
	}
	if len(result.Errors) != 2 || result.Errors[1].ID != 2 || result.Errors[0].Code != "ACU0000" {
		t.Errorf("errors = %+v", result.Errors)
	}
}

func TestCustomersService_CreatePoint(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/customers/1/points.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := body["point"]; got != float64(1000) {
			t.Errorf("point = %v", got)
		}
		if got := body["name"]; got != "test" {
			t.Errorf("name = %v", got)
		}
		if got := body["addition_and_subtraction"]; got != float64(1) {
			t.Errorf("addition_and_subtraction = %v", got)
		}
		if got := body["with_create"]; got != float64(1) {
			t.Errorf("with_create = %v", got)
		}
		fmt.Fprint(w, `{"data": {"id": "1", "type": "point", "attributes": {
			"id": 1, "customer_id": 1, "point_event_id": 1000001,
			"point_event_name": "test", "point_event_description": "description",
			"point": 1000, "point_total": 1500,
			"created_at": "2021/11/01 23:59:59", "deleted_at": null}}}`)
	})

	point, _, err := client.Customers.CreatePoint(context.Background(), 1, &PointParams{
		Name:                   ecforce.String("test"),
		Description:            ecforce.String("description"),
		Point:                  ecforce.Int(1000),
		AdditionAndSubtraction: ecforce.Bool01(true),
		WithCreate:             ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	attrs := point.Attributes
	if attrs.ID != 1 || attrs.CustomerID != 1 || attrs.Point != 1000 || attrs.PointTotal != 1500 {
		t.Errorf("unexpected point: %+v", attrs)
	}
	if attrs.PointEventName != "test" {
		t.Errorf("point_event_name = %q", attrs.PointEventName)
	}
}

func TestCustomersService_BulkCreatePoints(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/customers/points/bulk_create.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		points, _ := body["customer_points"].([]any)
		if len(points) != 2 {
			t.Fatalf("customer_points = %v", body["customer_points"])
		}
		first, _ := points[0].(map[string]any)
		if got := first["customer_id"]; got != float64(1) {
			t.Errorf("customer_points[0].customer_id = %v", got)
		}
		if got := first["point"]; got != float64(100) {
			t.Errorf("customer_points[0].point = %v", got)
		}
		fmt.Fprint(w, `{"success": [1, 2]}`)
	})

	result, _, err := client.Customers.BulkCreatePoints(context.Background(), []*CustomerPointParams{
		{CustomerID: ecforce.Int64(1), PointEventID: ecforce.Int64(1), Point: ecforce.Int(100), AdditionAndSubtraction: ecforce.Bool01(true)},
		{CustomerID: ecforce.Int64(2), PointEventID: ecforce.Int64(1), Point: ecforce.Int(200), AdditionAndSubtraction: ecforce.Bool01(true)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 2 || result.Success[0] != 1 || result.Success[1] != 2 {
		t.Errorf("success = %v", result.Success)
	}
	if len(result.Failure) != 0 || len(result.Errors) != 0 {
		t.Errorf("failure = %v, errors = %+v", result.Failure, result.Errors)
	}
}

func TestCustomersService_SendMail(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/customers/3172/emails.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := body["subject"]; got != "メール件名" {
			t.Errorf("subject = %v", got)
		}
		if got := body["body"]; got != "メール本文" {
			t.Errorf("body = %v", got)
		}
		if got := body["use_html"]; got != float64(0) {
			t.Errorf("use_html = %v", got)
		}
		fmt.Fprint(w, `{"data": {"id": "3172", "type": "customer", "attributes": {
			"id": 3172, "number": "6a54e6ebc5", "state": "member",
			"email": "test@example.com", "point": 92}}}`)
	})

	customer, _, err := client.Customers.SendMail(context.Background(), 3172, &CustomerMailRequest{
		Subject: ecforce.String("メール件名"),
		Body:    ecforce.String("メール本文"),
		UseHTML: ecforce.Bool01(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	if customer.ID != "3172" || customer.Attributes.Email != "test@example.com" {
		t.Errorf("unexpected customer: %+v", customer)
	}
}

func TestCustomersService_SendSMS(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/customers/3172/sms.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := body["body"]; got != "SMS本文" {
			t.Errorf("body = %v", got)
		}
		fmt.Fprint(w, `{"data": {"id": "3172", "type": "customer", "attributes": {
			"id": 3172, "number": "6a54e6ebc5", "state": "member",
			"email": "test@example.com"}}}`)
	})

	customer, _, err := client.Customers.SendSMS(context.Background(), 3172, "SMS本文")
	if err != nil {
		t.Fatal(err)
	}
	if customer.Attributes.ID != 3172 || customer.Attributes.Number != "6a54e6ebc5" {
		t.Errorf("unexpected customer: %+v", customer.Attributes)
	}
}

func TestCustomersService_BulkCreateCreditCards(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/customers/1/credit_cards/bulk_create.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		cards, _ := body["credit_cards"].([]any)
		if len(cards) != 2 {
			t.Fatalf("credit_cards = %v", body["credit_cards"])
		}
		second, _ := cards[1].(map[string]any)
		if got := second["last_digits"]; got != "2222" {
			t.Errorf("credit_cards[1].last_digits = %v", got)
		}
		if got := second["default"]; got != float64(1) {
			t.Errorf("credit_cards[1].default = %v", got)
		}
		fmt.Fprint(w, `{
			"success": [1],
			"failure": [2],
			"errors": [{"id": 2, "code": "ACC0000", "message": "エラーが発生しました。",
				"errors": [{"code": "ACC3002", "message": "gateway_customer_profile_idパラメータが正しくありません。"}]}]}`)
	})

	result, _, err := client.Customers.BulkCreateCreditCards(context.Background(), 1, &CreditCardBulkRequest{
		CreditCards: []*CreditCardParams{
			{Month: ecforce.Int(5), Year: ecforce.Int(22), LastDigits: ecforce.String("1111"), Default: ecforce.Bool01(false)},
			{Month: ecforce.Int(6), Year: ecforce.Int(23), LastDigits: ecforce.String("2222"), Default: ecforce.Bool01(true)},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 1 || result.Success[0] != 1 {
		t.Errorf("success = %v", result.Success)
	}
	if len(result.Failure) != 1 || result.Failure[0] != 2 {
		t.Errorf("failure = %v", result.Failure)
	}
	if len(result.Errors) != 1 || result.Errors[0].Errors[0].Code != "ACC3002" {
		t.Errorf("errors = %+v", result.Errors)
	}
}

func TestCustomersService_BulkUpdateCreditCards(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/customers/1/credit_cards/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		cards, _ := body["credit_cards"].([]any)
		if len(cards) != 1 {
			t.Fatalf("credit_cards = %v", body["credit_cards"])
		}
		first, _ := cards[0].(map[string]any)
		if got := first["id"]; got != float64(10) {
			t.Errorf("credit_cards[0].id = %v", got)
		}
		if got := first["month"]; got != float64(12) {
			t.Errorf("credit_cards[0].month = %v", got)
		}
		if got := body["request_payment_service"]; got != float64(1) {
			t.Errorf("request_payment_service = %v", got)
		}
		fmt.Fprint(w, `{"success": [10]}`)
	})

	result, _, err := client.Customers.BulkUpdateCreditCards(context.Background(), 1, &CreditCardBulkRequest{
		CreditCards: []*CreditCardParams{
			{ID: ecforce.Int64(10), Month: ecforce.Int(12), Year: ecforce.Int(30)},
		},
		RequestPaymentService: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 1 || result.Success[0] != 10 {
		t.Errorf("success = %v", result.Success)
	}
}

func TestCustomersService_BulkDestroyCreditCards(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("DELETE /api/v2/admin/customers/1/credit_cards/bulk_destroy.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprint(body["credit_card_ids"]); got != "[10 11]" {
			t.Errorf("credit_card_ids = %v", got)
		}
		if got := body["request_payment_service"]; got != float64(0) {
			t.Errorf("request_payment_service = %v", got)
		}
		fmt.Fprint(w, `{"success": [10, 11]}`)
	})

	result, _, err := client.Customers.BulkDestroyCreditCards(context.Background(), 1, &CreditCardBulkDestroyRequest{
		CreditCardIDs:         []int64{10, 11},
		RequestPaymentService: ecforce.Bool01(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 2 || result.Success[0] != 10 || result.Success[1] != 11 {
		t.Errorf("success = %v", result.Success)
	}
}

func TestCustomersService_BulkDestroyShippingAddresses(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("DELETE /api/v2/admin/customers/1/shipping_addresses/bulk_destroy.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprint(body["shipping_address_ids"]); got != "[1 2 3]" {
			t.Errorf("shipping_address_ids = %v", got)
		}
		fmt.Fprint(w, `{"success": [1, 2, 3]}`)
	})

	result, _, err := client.Customers.BulkDestroyShippingAddresses(context.Background(), 1, []int64{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 3 || result.Success[2] != 3 {
		t.Errorf("success = %v", result.Success)
	}
}

func TestCustomersService_ListFreeColumnValues(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/customers/1/free_column_values.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{"values": [{"free_column_category_seq": 1, "values": [
				{"free_column_id": 1, "free_column_label": "name", "free_column_value_values": ["タロウ"]}]}]},
			{"free_column_category_id": 1, "values": [
				{"free_column_category_seq": 1, "values": [
					{"free_column_id": 1, "free_column_label": "name", "free_column_value_values": ["タマ"]},
					{"free_column_id": 2, "free_column_label": "favorites", "free_column_value_values": ["ドライフード", "サバ缶"]}]},
				{"free_column_category_seq": 2, "values": [
					{"free_column_id": 2, "free_column_label": "favorites", "free_column_value_values": ["サバ缶"]}]}]}]`)
	})

	groups, _, err := client.Customers.ListFreeColumnValues(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("groups = %+v", groups)
	}
	if groups[0].FreeColumnCategoryID != 0 || groups[1].FreeColumnCategoryID != 1 {
		t.Errorf("category ids = %d, %d", groups[0].FreeColumnCategoryID, groups[1].FreeColumnCategoryID)
	}
	if len(groups[1].Values) != 2 || groups[1].Values[1].FreeColumnCategorySeq != 2 {
		t.Fatalf("unexpected entries: %+v", groups[1].Values)
	}
	value := groups[1].Values[0].Values[1]
	if value.FreeColumnLabel != "favorites" || len(value.FreeColumnValueValues) != 2 {
		t.Errorf("unexpected value: %+v", value)
	}
}

func TestCustomersService_BulkCreateFreeColumnValues(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/customers/1/free_column_values/bulk_create.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		columns, _ := body["free_columns"].([]any)
		if len(columns) != 2 {
			t.Fatalf("free_columns = %v", body["free_columns"])
		}
		first, _ := columns[0].(map[string]any)
		if got := first["free_column_id"]; got != float64(1) {
			t.Errorf("free_columns[0].free_column_id = %v", got)
		}
		if got := first["free_column_value"]; got != "タロウ" {
			t.Errorf("free_columns[0].free_column_value = %v", got)
		}
		second, _ := columns[1].(map[string]any)
		if got := second["free_column_category_id"]; got != float64(1) {
			t.Errorf("free_columns[1].free_column_category_id = %v", got)
		}
		fmt.Fprint(w, `{"success": [1, 2]}`)
	})

	result, _, err := client.Customers.BulkCreateFreeColumnValues(context.Background(), 1, []*FreeColumnParams{
		{FreeColumnID: ecforce.Int64(1), FreeColumnValue: ecforce.String("タロウ")},
		{FreeColumnCategoryID: ecforce.Int64(1), Values: []*FreeColumnValueParams{
			{FreeColumnID: ecforce.Int64(2), FreeColumnOptionIDs: []int64{3}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 2 || result.Success[0] != 1 {
		t.Errorf("success = %v", result.Success)
	}
}

func TestCustomersService_BulkUpdateFreeColumnValues(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/customers/1/free_column_values/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		columns, _ := body["free_columns"].([]any)
		if len(columns) != 1 {
			t.Fatalf("free_columns = %v", body["free_columns"])
		}
		first, _ := columns[0].(map[string]any)
		if got := first["free_column_category_seq"]; got != float64(2) {
			t.Errorf("free_columns[0].free_column_category_seq = %v", got)
		}
		if got := first["delete"]; got != float64(1) {
			t.Errorf("free_columns[0].delete = %v", got)
		}
		fmt.Fprint(w, `{"success": [1]}`)
	})

	result, _, err := client.Customers.BulkUpdateFreeColumnValues(context.Background(), 1, []*FreeColumnParams{
		{
			FreeColumnCategoryID:  ecforce.Int64(1),
			FreeColumnCategorySeq: ecforce.Int64(2),
			Delete:                ecforce.Bool01(true),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 1 || result.Success[0] != 1 {
		t.Errorf("success = %v", result.Success)
	}
}

func TestCustomersService_ListBlacklistReasons(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/blacklist_reasons.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("sort"); got != "-created_at" {
			t.Errorf("sort = %q", got)
		}
		fmt.Fprint(w, `{"data": [
			{"id": "2", "type": "blacklist_reason", "attributes": {
				"id": 2, "name": "重複購入", "reason": "同一人物による重複購入",
				"created_at": "2021/08/18 18:54:28", "updated_at": "2021/08/18 18:54:28"}},
			{"id": "3", "type": "blacklist_reason", "attributes": {
				"id": 3, "name": "未払い", "reason": "過去に代金未払い有り",
				"created_at": "2021/10/06 13:32:07", "updated_at": "2021/10/06 13:32:07"}}]}`)
	})

	reasons, _, err := client.Customers.ListBlacklistReasons(context.Background(), &ecforce.ListOptions{
		Sort: []string{"-created_at"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(reasons) != 2 {
		t.Fatalf("reasons = %+v", reasons)
	}
	if reasons[0].ID != "2" || reasons[0].Attributes.Name != "重複購入" {
		t.Errorf("unexpected reason: %+v", reasons[0])
	}
	if reasons[1].Attributes.ID != 3 || reasons[1].Attributes.Reason != "過去に代金未払い有り" {
		t.Errorf("unexpected reason: %+v", reasons[1].Attributes)
	}
}
