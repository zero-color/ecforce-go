package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/zero-color/ecforce-go"
)

func TestOrdersService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/orders.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("page = %q, want 2", got)
		}
		if got := r.URL.Query().Get("q[state_eq]"); got != "complete" {
			t.Errorf("q[state_eq] = %q, want complete", got)
		}
		fmt.Fprint(w, `{
			"data": [{"id": "96704", "type": "order", "attributes": {
				"id": 96704, "number": "a905125d2a", "customer_id": 30962,
				"state": "complete", "human_state": "注文確定",
				"payment_state": "completed", "email": "ok@example.com",
				"subtotal": 500, "total": 799, "payment_total": 799,
				"shipping_carrier_id": 1, "shipping_slip": "abcde33333",
				"created_at": "2019/06/08 05:47:17"}}],
			"meta": {"total_count": 41708, "page": 2, "per": 1, "count": 1, "total_pages": 41708},
			"links": {"self": "http://localhost/api/v2/admin/orders?page=2&per=1"}}`)
	})

	orders, resp, err := client.Orders.List(context.Background(), &ecforce.ListOptions{
		Page: 2,
		Q:    ecforce.Query{"state_eq": "complete"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(orders) != 1 {
		t.Fatalf("len(orders) = %d, want 1", len(orders))
	}
	attrs := orders[0].Attributes
	if orders[0].ID != "96704" || attrs.ID != 96704 || attrs.Number != "a905125d2a" {
		t.Errorf("unexpected order identity: %+v", attrs)
	}
	if attrs.State != "complete" || attrs.Email != "ok@example.com" || attrs.Total != 799 {
		t.Errorf("unexpected order attributes: %+v", attrs)
	}
	if resp.Meta == nil || resp.Meta.TotalCount != 41708 || resp.Meta.Page != 2 {
		t.Errorf("unexpected meta: %+v", resp.Meta)
	}
	if !resp.HasNextPage() || resp.NextPage() != 3 {
		t.Errorf("HasNextPage/NextPage = %v/%d, want true/3", resp.HasNextPage(), resp.NextPage())
	}
}

func TestOrdersService_Get(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/orders/96704.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("include"); got != "customer" {
			t.Errorf("include = %q, want customer", got)
		}
		io.WriteString(w, `{"data": {"id": "96704", "type": "order", "attributes": {
			"id": 96704, "number": "a905125d2a", "customer_id": 30962,
			"state": "complete", "payment_state": "completed",
			"email": "ok@example.com", "total": 799, "kind": "定期受注", "nth": 1,
			"picked_list": true, "shipping_carrier_name": "ヤマト運輸",
			"order_free_columns": [{"id": 1, "values": ["回答"], "name": "質問"}],
			"created_at": "2019/06/08 05:47:17",
			"campaigns": [{"id": 1, "name": "spring10%off",
				"display_name": "春の10％offキャンペーン",
				"description": "商品価格から10%offするキャンペーンです"}]}}}`)
	})

	order, _, err := client.Orders.Get(context.Background(), 96704, &ecforce.GetOptions{
		Include: []string{"customer"},
	})
	if err != nil {
		t.Fatal(err)
	}
	attrs := order.Attributes
	if attrs.ID != 96704 || attrs.CustomerID != 30962 || attrs.Kind != "定期受注" {
		t.Errorf("unexpected order: %+v", attrs)
	}
	if !attrs.PickedList || attrs.ShippingCarrierName != "ヤマト運輸" {
		t.Errorf("unexpected order fields: %+v", attrs)
	}
	if len(attrs.OrderFreeColumns) != 1 || attrs.OrderFreeColumns[0].Name != "質問" ||
		len(attrs.OrderFreeColumns[0].Values) != 1 || attrs.OrderFreeColumns[0].Values[0] != "回答" {
		t.Errorf("unexpected order free columns: %+v", attrs.OrderFreeColumns)
	}
	if len(attrs.Campaigns) != 1 || attrs.Campaigns[0].Name != "spring10%off" {
		t.Errorf("unexpected campaigns: %+v", attrs.Campaigns)
	}
}

func TestOrdersService_Update(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/orders/96704.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		order, ok := body["order"].(map[string]any)
		if !ok {
			t.Fatalf(`missing "order" wrapper key: %v`, body)
		}
		if got := order["memo01"]; got != "メモ1" {
			t.Errorf("order.memo01 = %v, want メモ1", got)
		}
		if got := order["adjustment"]; got != float64(-100) {
			t.Errorf("order.adjustment = %v, want -100", got)
		}
		if got := body["recalculate"]; got != float64(1) {
			t.Errorf("recalculate = %v, want 1", got)
		}
		fmt.Fprint(w, `{"data": {"id": "96704", "type": "order", "attributes": {
			"id": 96704, "number": "a905125d2a", "state": "complete",
			"memo01": "メモ1", "adjustment": -100, "total": 699}}}`)
	})

	order, _, err := client.Orders.Update(context.Background(), 96704, &OrderUpdateRequest{
		Order: &OrderParams{
			Memo01:     ecforce.String("メモ1"),
			Adjustment: ecforce.Int(-100),
		},
		Recalculate: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if order.Attributes.Memo01 != "メモ1" || order.Attributes.Adjustment != -100 || order.Attributes.Total != 699 {
		t.Errorf("unexpected order: %+v", order.Attributes)
	}
}

func TestOrdersService_BulkUpdate(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/orders/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		orders, ok := body["orders"].([]any)
		if !ok || len(orders) != 2 {
			t.Fatalf(`unexpected "orders" key: %v`, body)
		}
		first := orders[0].(map[string]any)
		if got := first["id"]; got != float64(88111) {
			t.Errorf("orders[0].id = %v, want 88111", got)
		}
		if got := first["link_number"]; got != "111" {
			t.Errorf("orders[0].link_number = %v, want 111", got)
		}
		if got := body["check_duplicate_link_numbers"]; got != float64(1) {
			t.Errorf("check_duplicate_link_numbers = %v, want 1", got)
		}
		fmt.Fprint(w, `{
			"success": [88111],
			"failure": [99999],
			"errors": [{"id": 99999, "code": "AOR0000", "message": "エラーが発生しました。",
				"errors": [
					{"code": "AOR2001", "message": "orderが見つかりません。"},
					{"code": "AOR3002", "message": "stateパラメータが正しくありません。"}]}]}`)
	})

	result, _, err := client.Orders.BulkUpdate(context.Background(), &OrderBulkUpdateRequest{
		Orders: []*OrderParams{
			{ID: ecforce.Int64(88111), LinkNumber: ecforce.String("111")},
			{ID: ecforce.Int64(99999), State: ecforce.String("bad")},
		},
		CheckDuplicateLinkNumbers: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 1 || result.Success[0] != 88111 {
		t.Errorf("Success = %v, want [88111]", result.Success)
	}
	if len(result.Failure) != 1 || result.Failure[0] != 99999 {
		t.Errorf("Failure = %v, want [99999]", result.Failure)
	}
	if len(result.Errors) != 1 || result.Errors[0].ID != 99999 || result.Errors[0].Code != "AOR0000" {
		t.Fatalf("unexpected errors: %+v", result.Errors)
	}
	if len(result.Errors[0].Errors) != 2 || result.Errors[0].Errors[0].Code != "AOR2001" {
		t.Errorf("unexpected nested errors: %+v", result.Errors[0].Errors)
	}
}

func TestOrdersService_BulkDestroy(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("DELETE /api/v2/admin/orders/bulk_destroy.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		ids, ok := body["order_ids"].([]any)
		if !ok || len(ids) != 2 || ids[0] != float64(88111) || ids[1] != float64(88112) {
			t.Errorf(`order_ids = %v, want [88111 88112]`, body["order_ids"])
		}
		fmt.Fprint(w, `{"success": [88111, 88112]}`)
	})

	result, _, err := client.Orders.BulkDestroy(context.Background(), []int64{88111, 88112})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 2 || result.Success[0] != 88111 || result.Success[1] != 88112 {
		t.Errorf("Success = %v, want [88111 88112]", result.Success)
	}
	if len(result.Failure) != 0 || len(result.Errors) != 0 {
		t.Errorf("Failure/Errors = %v/%v, want empty", result.Failure, result.Errors)
	}
}

func TestOrdersService_SendMail(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/orders/96704/emails.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		if got := body["email_template_id"]; got != float64(5) {
			t.Errorf("email_template_id = %v, want 5", got)
		}
		fmt.Fprint(w, `{"data": {"id": "96704", "type": "order", "attributes": {
			"id": 96704, "number": "a905125d2a", "email": "ok@example.com",
			"state": "complete", "total": 799}}}`)
	})

	order, _, err := client.Orders.SendMail(context.Background(), 96704, 5)
	if err != nil {
		t.Fatal(err)
	}
	if order.Attributes.ID != 96704 || order.Attributes.Email != "ok@example.com" {
		t.Errorf("unexpected order: %+v", order.Attributes)
	}
}

func TestOrdersService_UpdateShipping(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/orders/shipping.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		orders, ok := body["orders"].([]any)
		if !ok || len(orders) != 2 {
			t.Fatalf(`unexpected "orders" key: %v`, body)
		}
		first := orders[0].(map[string]any)
		if got := first["id"]; got != float64(88111) {
			t.Errorf("orders[0].id = %v, want 88111", got)
		}
		if got := first["shipping_slip"]; got != "abcde12345" {
			t.Errorf("orders[0].shipping_slip = %v, want abcde12345", got)
		}
		if got := first["shipping_address_id"]; got != float64(120) {
			t.Errorf("orders[0].shipping_address_id = %v, want 120", got)
		}
		if got := body["with_sale"]; got != float64(1) {
			t.Errorf("with_sale = %v, want 1", got)
		}
		// Mixed numeric and string ("<order id>-<shipping address id>") IDs.
		fmt.Fprint(w, `{
			"success": ["88111-120", 88112],
			"failure": ["88111-125"],
			"errors": [{"id": "88111-125", "code": "AOR0000", "message": "エラーが発生しました。",
				"errors": [{"code": "AOR2001", "message": "orderが見つかりません。"}]}]}`)
	})

	shippedAt := ecforce.NewTime(time.Date(2017, 12, 13, 13, 22, 29, 0, time.UTC))
	result, _, err := client.Orders.UpdateShipping(context.Background(), &OrderShippingUpdateRequest{
		Orders: []*OrderShippingParams{
			{
				ID:                ecforce.Int64(88111),
				ShippingSlip:      ecforce.String("abcde12345"),
				ShippingCarrierID: ecforce.Int64(1),
				ShippedAt:         shippedAt,
				ShippingAddressID: ecforce.Int64(120),
			},
			{
				ID:           ecforce.Int64(88112),
				ShippingSlip: ecforce.String("abcde67890"),
				ShippedAt:    shippedAt,
			},
		},
		WithSale: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	// String-form and number-form IDs both decode into OrderShippingID.
	if len(result.Success) != 2 || result.Success[0] != "88111-120" || result.Success[1] != "88112" {
		t.Errorf(`Success = %v, want ["88111-120" "88112"]`, result.Success)
	}
	if len(result.Failure) != 1 || result.Failure[0] != "88111-125" {
		t.Errorf(`Failure = %v, want ["88111-125"]`, result.Failure)
	}
	if len(result.Errors) != 1 || result.Errors[0].ID != "88111-125" ||
		len(result.Errors[0].Errors) != 1 || result.Errors[0].Errors[0].Code != "AOR2001" {
		t.Errorf("unexpected errors: %+v", result.Errors)
	}
	if result.JobID != "" {
		t.Errorf("JobID = %q, want empty for synchronous response", result.JobID)
	}
}

func TestOrdersService_UpdatePaymentStatus(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/orders/96704/payment_status.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		if got := body["method"]; got != "void" {
			t.Errorf("method = %v, want void", got)
		}
		if got := body["need_to_send_email"]; got != float64(1) {
			t.Errorf("need_to_send_email = %v, want 1", got)
		}
		fmt.Fprint(w, `{"data": {"id": "96704", "type": "order", "attributes": {
			"id": 96704, "number": "a905125d2a", "payment_state": "voided",
			"payment_human_state": "取消済み", "payment_voided_at": "2019/06/09 10:00:00"}}}`)
	})

	order, _, err := client.Orders.UpdatePaymentStatus(context.Background(), 96704, &OrderPaymentStatusUpdateRequest{
		Method:          ecforce.String("void"),
		NeedToSendEmail: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if order.Attributes.PaymentState != "voided" || order.Attributes.PaymentVoidedAt == nil {
		t.Errorf("unexpected order: %+v", order.Attributes)
	}
}

func TestOrdersService_BulkUpdatePaymentStatus(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/orders/payment_status/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		if got := body["method"]; got != "force_void" {
			t.Errorf("method = %v, want force_void", got)
		}
		ids, ok := body["order_ids"].([]any)
		if !ok || len(ids) != 3 || ids[0] != float64(1) {
			t.Errorf("order_ids = %v, want [1 2 3]", body["order_ids"])
		}
		if got := body["decrement_subs_order_times"]; got != float64(0) {
			t.Errorf("decrement_subs_order_times = %v, want 0", got)
		}
		fmt.Fprint(w, `{"id": 12, "job_id": "1234"}`)
	})

	job, _, err := client.Orders.BulkUpdatePaymentStatus(context.Background(), &OrderPaymentStatusBulkUpdateRequest{
		Method:                  ecforce.String("force_void"),
		OrderIDs:                []int64{1, 2, 3},
		DecrementSubsOrderTimes: ecforce.Bool01(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	if job.ID != 12 || job.JobID != "1234" {
		t.Errorf("job = %+v, want ID=12 JobID=1234", job)
	}
}

func TestOrdersService_BulkCreateOrderItems(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/orders/96704/order_items/bulk_create.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		items, ok := body["order_items"].([]any)
		if !ok || len(items) != 2 {
			t.Fatalf(`unexpected "order_items" key: %v`, body)
		}
		first := items[0].(map[string]any)
		if got := first["variant_id"]; got != float64(1) {
			t.Errorf("order_items[0].variant_id = %v, want 1", got)
		}
		if got := first["quantity"]; got != float64(2) {
			t.Errorf("order_items[0].quantity = %v, want 2", got)
		}
		if got := body["refer_settings"]; got != float64(1) {
			t.Errorf("refer_settings = %v, want 1", got)
		}
		fmt.Fprint(w, `{"success": [1, 2]}`)
	})

	result, _, err := client.Orders.BulkCreateOrderItems(context.Background(), 96704, &OrderItemsBulkRequest{
		OrderItems: []*OrderItemParams{
			{VariantID: ecforce.Int64(1), Quantity: ecforce.Int(2)},
			{VariantID: ecforce.Int64(3), Quantity: ecforce.Int(4)},
		},
		ReferSettings: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 2 || result.Success[0] != 1 || result.Success[1] != 2 {
		t.Errorf("Success = %v, want [1 2]", result.Success)
	}
}

func TestOrdersService_BulkUpdateOrderItems(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/orders/96704/order_items/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		items, ok := body["order_items"].([]any)
		if !ok || len(items) != 2 {
			t.Fatalf(`unexpected "order_items" key: %v`, body)
		}
		first := items[0].(map[string]any)
		if got := first["id"]; got != float64(10) {
			t.Errorf("order_items[0].id = %v, want 10", got)
		}
		second := items[1].(map[string]any)
		if got := second["delete"]; got != float64(1) {
			t.Errorf("order_items[1].delete = %v, want 1", got)
		}
		fmt.Fprint(w, `{"success": [10, 11]}`)
	})

	result, _, err := client.Orders.BulkUpdateOrderItems(context.Background(), 96704, &OrderItemsBulkRequest{
		OrderItems: []*OrderItemParams{
			{ID: ecforce.Int64(10), Quantity: ecforce.Int(3)},
			{ID: ecforce.Int64(11), Delete: ecforce.Bool01(true)},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 2 || result.Success[0] != 10 || result.Success[1] != 11 {
		t.Errorf("Success = %v, want [10 11]", result.Success)
	}
}

func TestOrdersService_ListFreeColumnValues(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/orders/96704/free_column_values.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{"values": [{"free_column_category_seq": 1, "values": [
				{"free_column_id": 1, "free_column_label": "name", "free_column_value_values": ["タロウ"]}]}]},
			{"free_column_category_id": 1, "values": [{"free_column_category_seq": 1, "values": [
				{"free_column_id": 1, "free_column_label": "name", "free_column_value_values": ["タマ"]},
				{"free_column_id": 2, "free_column_label": "favorites", "free_column_value_values": ["ドライフード", "サバ缶"]}]}]}]`)
	})

	groups, _, err := client.Orders.ListFreeColumnValues(context.Background(), 96704)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("len(groups) = %d, want 2", len(groups))
	}
	if groups[0].FreeColumnCategoryID != 0 || groups[1].FreeColumnCategoryID != 1 {
		t.Errorf("category IDs = %d/%d, want 0/1", groups[0].FreeColumnCategoryID, groups[1].FreeColumnCategoryID)
	}
	entry := groups[1].Values[0]
	if entry.FreeColumnCategorySeq != 1 || len(entry.Values) != 2 {
		t.Fatalf("unexpected entry: %+v", entry)
	}
	v := entry.Values[1]
	if v.FreeColumnID != 2 || v.FreeColumnLabel != "favorites" ||
		len(v.FreeColumnValueValues) != 2 || v.FreeColumnValueValues[1] != "サバ缶" {
		t.Errorf("unexpected value: %+v", v)
	}
}

func TestOrdersService_BulkCreateFreeColumnValues(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/orders/96704/free_column_values/bulk_create.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		cols, ok := body["free_columns"].([]any)
		if !ok || len(cols) != 2 {
			t.Fatalf(`unexpected "free_columns" key: %v`, body)
		}
		first := cols[0].(map[string]any)
		if got := first["free_column_id"]; got != float64(1) {
			t.Errorf("free_columns[0].free_column_id = %v, want 1", got)
		}
		if got := first["free_column_value"]; got != "2000-01-01" {
			t.Errorf("free_columns[0].free_column_value = %v, want 2000-01-01", got)
		}
		second := cols[1].(map[string]any)
		if got := second["free_column_category_id"]; got != float64(1) {
			t.Errorf("free_columns[1].free_column_category_id = %v, want 1", got)
		}
		fmt.Fprint(w, `{
			"success": [1],
			"failure": [2],
			"errors": [{"id": 2, "code": "AOFV0000", "message": "エラーが発生しました。",
				"errors": [{"code": "AOFV3001", "message": "free_column_valueパラメータを入力してください。"}]}]}`)
	})

	result, _, err := client.Orders.BulkCreateFreeColumnValues(context.Background(), 96704, []*FreeColumnParams{
		{FreeColumnID: ecforce.Int64(1), FreeColumnValue: ecforce.String("2000-01-01")},
		{FreeColumnCategoryID: ecforce.Int64(1), Values: []*FreeColumnValueParams{
			{FreeColumnID: ecforce.Int64(3), FreeColumnOptionIDs: []int64{3}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 1 || result.Success[0] != 1 {
		t.Errorf("Success = %v, want [1]", result.Success)
	}
	if len(result.Failure) != 1 || result.Failure[0] != 2 {
		t.Errorf("Failure = %v, want [2]", result.Failure)
	}
	if len(result.Errors) != 1 || result.Errors[0].ID != 2 ||
		len(result.Errors[0].Errors) != 1 || result.Errors[0].Errors[0].Code != "AOFV3001" {
		t.Errorf("unexpected errors: %+v", result.Errors)
	}
}

func TestOrdersService_BulkUpdateFreeColumnValues(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/orders/96704/free_column_values/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request body: %v", err)
		}
		cols, ok := body["free_columns"].([]any)
		if !ok || len(cols) != 1 {
			t.Fatalf(`unexpected "free_columns" key: %v`, body)
		}
		first := cols[0].(map[string]any)
		if got := first["free_column_category_id"]; got != float64(1) {
			t.Errorf("free_columns[0].free_column_category_id = %v, want 1", got)
		}
		if got := first["free_column_category_seq"]; got != float64(2) {
			t.Errorf("free_columns[0].free_column_category_seq = %v, want 2", got)
		}
		if got := first["delete"]; got != float64(1) {
			t.Errorf("free_columns[0].delete = %v, want 1", got)
		}
		fmt.Fprint(w, `{"success": [1]}`)
	})

	result, _, err := client.Orders.BulkUpdateFreeColumnValues(context.Background(), 96704, []*FreeColumnParams{
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
		t.Errorf("Success = %v, want [1]", result.Success)
	}
	if len(result.Failure) != 0 || len(result.Errors) != 0 {
		t.Errorf("Failure/Errors = %v/%v, want empty", result.Failure, result.Errors)
	}
}
