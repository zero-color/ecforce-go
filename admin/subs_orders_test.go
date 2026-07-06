package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/zero-color/ecforce-go"
)

func TestSubsOrdersService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/subs_orders.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("page = %q, want 2", got)
		}
		if got := r.URL.Query().Get("q[state_eq]"); got != "active" {
			t.Errorf("q[state_eq] = %q, want active", got)
		}
		fmt.Fprint(w, `{
			"data": [{
				"id": "4909", "type": "sub_order",
				"attributes": {
					"id": 4909, "number": "a2886ff4f2", "customer_number": "02ae6cdbef",
					"times": 1, "orders_count": 1, "state": "active", "human_state": "有効",
					"subtotal": 100, "tbc": true,
					"available_payment_schedules": [{
						"value": "date", "ja": "日付で指定",
						"scheduled_to_be_delivered_every_x_month": [{"value": 1, "ja": "1ヶ月"}],
						"scheduled_to_be_delivered_on_xth_day": [{"value": 99, "ja": "末日"}]
					}],
					"payment_schedule": "1ヶ月ごとの1日に配送", "payment_schedule_locked": 1,
					"scheduled_to_be_shipped_at": "2019/06/26 00:00:00",
					"link_number": "111",
					"order_free_columns": [{"id": 1, "values": ["回答"], "name": "質問"}],
					"created_at": "2019/06/08 05:47:20", "deleted_at": null
				}
			}],
			"meta": {"total_count": 3505, "page": 2, "per": 1, "count": 1, "total_pages": 3505}
		}`)
	})

	subsOrders, resp, err := client.SubsOrders.List(context.Background(), &ecforce.ListOptions{
		Page: 2,
		Q:    ecforce.Query{"state_eq": "active"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(subsOrders) != 1 {
		t.Fatalf("len(subsOrders) = %d, want 1", len(subsOrders))
	}
	attrs := subsOrders[0].Attributes
	if subsOrders[0].ID != "4909" || attrs.ID != 4909 || attrs.Number != "a2886ff4f2" {
		t.Errorf("unexpected subs order identity: %+v", subsOrders[0])
	}
	if attrs.State != "active" || attrs.Subtotal != 100 || !attrs.TBC || !bool(attrs.PaymentScheduleLocked) {
		t.Errorf("unexpected attributes: %+v", attrs)
	}
	if len(attrs.AvailablePaymentSchedules) != 1 || attrs.AvailablePaymentSchedules[0].Value != "date" {
		t.Errorf("unexpected available payment schedules: %+v", attrs.AvailablePaymentSchedules)
	}
	if got := attrs.ScheduledToBeShippedAt.Format("2006-01-02"); got != "2019-06-26" {
		t.Errorf("ScheduledToBeShippedAt = %s, want 2019-06-26", got)
	}
	if len(attrs.OrderFreeColumns) != 1 || attrs.OrderFreeColumns[0].Name != "質問" {
		t.Errorf("unexpected order free columns: %+v", attrs.OrderFreeColumns)
	}
	if !resp.HasNextPage() || resp.NextPage() != 3 || resp.Meta.TotalCount != 3505 {
		t.Errorf("unexpected pagination: meta=%+v", resp.Meta)
	}
}

func TestSubsOrdersService_Get(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/subs_orders/4909.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("include"); got != "payment" {
			t.Errorf("include = %q, want payment", got)
		}
		fmt.Fprint(w, `{
			"data": {
				"id": "4909", "type": "sub_order",
				"attributes": {
					"id": 4909, "number": "a2886ff4f2", "state": "active",
					"payment_method_name": "NP後払い wiz", "recurring_block_times": 0,
					"linked_memo": "メモ内容", "doorbell": true
				}
			},
			"included": [{"id": "98647", "type": "payment", "attributes": {"id": 98647}}]
		}`)
	})

	subsOrder, resp, err := client.SubsOrders.Get(context.Background(), 4909, &ecforce.GetOptions{
		Include: []string{"payment"},
	})
	if err != nil {
		t.Fatal(err)
	}
	attrs := subsOrder.Attributes
	if subsOrder.ID != "4909" || attrs.Number != "a2886ff4f2" || attrs.PaymentMethodName != "NP後払い wiz" {
		t.Errorf("unexpected subs order: %+v", attrs)
	}
	if attrs.LinkedMemo != "メモ内容" || !attrs.Doorbell {
		t.Errorf("unexpected attributes: %+v", attrs)
	}
	if len(resp.Included) != 1 || resp.Included[0].Type != "payment" {
		t.Errorf("unexpected included: %+v", resp.Included)
	}
}

func TestSubsOrdersService_Update(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/subs_orders/4909.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		subsOrder, ok := body["subs_order"].(map[string]any)
		if !ok {
			t.Fatalf("subs_order wrapper missing: %v", body)
		}
		if subsOrder["state"] != "suspend" || subsOrder["times"] != 3.0 {
			t.Errorf("unexpected subs_order body: %v", subsOrder)
		}
		if body["avoid_holidays_for_scheduled_to_be_shipped_at"] != 1.0 {
			t.Errorf("avoid_holidays = %v, want 1", body["avoid_holidays_for_scheduled_to_be_shipped_at"])
		}
		fmt.Fprint(w, `{"data": {"id": "4909", "type": "sub_order", "attributes": {
			"id": 4909, "state": "suspend", "human_state": "停止", "times": 3,
			"suspend_reasons": "金銭的な問題", "suspended_at": "2022/01/05 10:00:00"
		}}}`)
	})

	subsOrder, _, err := client.SubsOrders.Update(context.Background(), 4909, &SubsOrderUpdateRequest{
		SubsOrder: &SubsOrderParams{
			State:            ecforce.String("suspend"),
			Times:            ecforce.Ptr(3),
			SuspendReasonIDs: []string{"1"},
		},
		AvoidHolidaysForScheduledToBeShippedAt: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	attrs := subsOrder.Attributes
	if attrs.State != "suspend" || attrs.Times != 3 || attrs.SuspendReasons != "金銭的な問題" {
		t.Errorf("unexpected subs order: %+v", attrs)
	}
	if attrs.SuspendedAt == nil || attrs.SuspendedAt.IsZero() {
		t.Errorf("SuspendedAt = %v, want set", attrs.SuspendedAt)
	}
}

func TestSubsOrdersService_BulkCreate(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/subs_orders/bulk_create.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		subsOrders, ok := body["subs_orders"].([]any)
		if !ok || len(subsOrders) != 1 {
			t.Fatalf("subs_orders wrapper missing or wrong length: %v", body)
		}
		first := subsOrders[0].(map[string]any)
		if first["customer_id"] != 1.0 || first["payment_schedule"] != "date" {
			t.Errorf("unexpected subs_orders[0]: %v", first)
		}
		items, ok := first["order_items_attributes"].([]any)
		if !ok || len(items) != 1 || items[0].(map[string]any)["variant_id"] != 1.0 {
			t.Errorf("unexpected order_items_attributes: %v", first["order_items_attributes"])
		}
		if body["check_duplicate_link_numbers"] != 1.0 {
			t.Errorf("check_duplicate_link_numbers = %v, want 1", body["check_duplicate_link_numbers"])
		}
		fmt.Fprint(w, `{"id": 555, "job_id": "1234"}`)
	})

	job, _, err := client.SubsOrders.BulkCreate(context.Background(), &SubsOrderBulkCreateRequest{
		SubsOrders: []*SubsOrderParams{{
			CustomerID:                        ecforce.Int64(1),
			PaymentSchedule:                   ecforce.String("date"),
			ScheduledToBeDeliveredEveryXMonth: ecforce.Ptr(1),
			ScheduledToBeDeliveredOnXthDay:    ecforce.Ptr(1),
			OrderItemsAttributes: []*SubsOrderItemParams{
				{VariantID: ecforce.Int64(1), Quantity: ecforce.Ptr(1)},
			},
			PaymentAttributes: &SubsOrderPaymentParams{PaymentMethodID: ecforce.Int64(81)},
		}},
		CheckDuplicateLinkNumbers: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if job.ID != 555 || job.JobID != "1234" {
		t.Errorf("unexpected job result: %+v", job)
	}
}

func TestSubsOrdersService_BulkUpdate(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/subs_orders/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		subsOrders, ok := body["subs_orders"].([]any)
		if !ok || len(subsOrders) != 2 {
			t.Fatalf("subs_orders wrapper missing or wrong length: %v", body)
		}
		first := subsOrders[0].(map[string]any)
		if first["id"] != 4909.0 || first["state"] != "canceled" {
			t.Errorf("unexpected subs_orders[0]: %v", first)
		}
		fmt.Fprint(w, `{
			"success": [4909],
			"failure": [9999],
			"errors": [{
				"id": 9999, "code": "ASO0000", "message": "エラーが発生しました。",
				"errors": [{"code": "ASO2001", "message": "subs_orderが見つかりません。"}]
			}]
		}`)
	})

	result, _, err := client.SubsOrders.BulkUpdate(context.Background(), &SubsOrderBulkUpdateRequest{
		SubsOrders: []*SubsOrderParams{
			{ID: ecforce.Int64(4909), State: ecforce.String("canceled"), CancelReasonIDs: []string{"1"}},
			{ID: ecforce.Int64(9999), State: ecforce.String("canceled")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 1 || result.Success[0] != 4909 {
		t.Errorf("Success = %v, want [4909]", result.Success)
	}
	if len(result.Failure) != 1 || result.Failure[0] != 9999 {
		t.Errorf("Failure = %v, want [9999]", result.Failure)
	}
	if len(result.Errors) != 1 || result.Errors[0].ID != 9999 || result.Errors[0].Errors[0].Code != "ASO2001" {
		t.Errorf("unexpected errors: %+v", result.Errors)
	}
}

func TestSubsOrdersService_BulkDestroy(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("DELETE /api/v2/admin/subs_orders/bulk_destroy.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		ids, ok := body["subs_order_ids"].([]any)
		if !ok || len(ids) != 3 || ids[0] != 123.0 {
			t.Errorf("unexpected subs_order_ids: %v", body["subs_order_ids"])
		}
		fmt.Fprint(w, `{
			"success": [123, 125],
			"failure": [124],
			"errors": [{
				"id": 124, "code": "ASO0000", "message": "エラーが発生しました。",
				"errors": [{"code": "ASO3999", "message": "ステータスが有効のため削除できません。"}]
			}]
		}`)
	})

	result, _, err := client.SubsOrders.BulkDestroy(context.Background(), []int64{123, 124, 125})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 2 || result.Success[0] != 123 || result.Success[1] != 125 {
		t.Errorf("Success = %v, want [123 125]", result.Success)
	}
	if len(result.Failure) != 1 || result.Failure[0] != 124 {
		t.Errorf("Failure = %v, want [124]", result.Failure)
	}
	if len(result.Errors) != 1 || result.Errors[0].Errors[0].Code != "ASO3999" {
		t.Errorf("unexpected errors: %+v", result.Errors)
	}
}

func TestSubsOrdersService_BulkRecalculate(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/subs_orders/bulk_recalculate.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		ids, ok := body["subs_order_ids"].([]any)
		if !ok || len(ids) != 2 || ids[0] != 1.0 || ids[1] != 2.0 {
			t.Errorf("unexpected subs_order_ids: %v", body["subs_order_ids"])
		}
		fmt.Fprint(w, `{"id": 777, "job_id": "5678"}`)
	})

	job, _, err := client.SubsOrders.BulkRecalculate(context.Background(), []int64{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	if job.ID != 777 || job.JobID != "5678" {
		t.Errorf("unexpected job result: %+v", job)
	}
}

func TestSubsOrdersService_CreateOrder(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/subs_orders/4909/orders.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.SubsOrders.CreateOrder(context.Background(), 4909)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want 204", resp.StatusCode)
	}
}

func TestSubsOrdersService_CreateOrderItem(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/subs_orders/4909/order_items.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		item, ok := body["order_item"].(map[string]any)
		if !ok {
			t.Fatalf("order_item wrapper missing: %v", body)
		}
		if item["variant_id"] != 1.0 || item["quantity"] != 1.0 {
			t.Errorf("unexpected order_item body: %v", item)
		}
		if body["recalculate_price"] != 1.0 || body["times"] != 2.0 || body["sync"] != 1.0 {
			t.Errorf("unexpected recalculate params: %v", body)
		}
		fmt.Fprint(w, `{"data": {
			"id": "1185", "type": "order_item",
			"attributes": {
				"id": 1185, "variant_id": 1, "product_number": "S-001",
				"product_name": "テスト商品（単品）", "variant_sku": "S-001-001",
				"list_price": 500, "sales_price": 500, "price": 250,
				"quantity": 1, "tax_rate": 10,
				"created_at": "2023/03/20 17:09:32"
			},
			"relationships": {"variant": {"data": {"id": "1", "type": "variant"}}}
		}}`)
	})

	item, _, err := client.SubsOrders.CreateOrderItem(context.Background(), 4909, &SubsOrderItemCreateRequest{
		OrderItem:        &SubsOrderItemParams{VariantID: ecforce.Int64(1), Quantity: ecforce.Ptr(1)},
		RecalculatePrice: ecforce.Bool01(true),
		Times:            ecforce.Ptr(2),
		Sync:             ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	attrs := item.Attributes
	if item.ID != "1185" || attrs.ID != 1185 || attrs.VariantID != 1 {
		t.Errorf("unexpected order item identity: %+v", item)
	}
	if attrs.ProductName != "テスト商品（単品）" || attrs.Price != 250 || attrs.Quantity != 1 || attrs.TaxRate != 10 {
		t.Errorf("unexpected attributes: %+v", attrs)
	}
}

func TestSubsOrdersService_UpdateOrderItem(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/subs_orders/4909/order_items/1185.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		item, ok := body["order_item"].(map[string]any)
		if !ok {
			t.Fatalf("order_item wrapper missing: %v", body)
		}
		if item["delete"] != 1.0 {
			t.Errorf("order_item.delete = %v, want 1", item["delete"])
		}
		w.WriteHeader(http.StatusNoContent)
	})

	item, resp, err := client.SubsOrders.UpdateOrderItem(context.Background(), 4909, 1185, &SubsOrderItemUpdateRequest{
		OrderItem: &SubsOrderItemParams{Delete: ecforce.Bool01(true)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if item != nil {
		t.Errorf("item = %+v, want nil on 204 delete response", item)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want 204", resp.StatusCode)
	}
}

func TestSubsOrdersService_BulkCreateOrderItems(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/subs_orders/4909/order_items/bulk_create.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		items, ok := body["order_items"].([]any)
		if !ok || len(items) != 2 {
			t.Fatalf("order_items wrapper missing or wrong length: %v", body)
		}
		first := items[0].(map[string]any)
		if first["variant_id"] != 1.0 || first["quantity"] != 1.0 {
			t.Errorf("unexpected order_items[0]: %v", first)
		}
		fmt.Fprint(w, `{
			"success": [1],
			"failure": [2],
			"errors": [{
				"id": 2, "code": "ASOI0000", "message": "エラーが発生しました。",
				"errors": [{"code": "ASOI2001", "message": "variantが見つかりません。"}]
			}]
		}`)
	})

	result, _, err := client.SubsOrders.BulkCreateOrderItems(context.Background(), 4909, &SubsOrderItemBulkCreateRequest{
		OrderItems: []*SubsOrderItemParams{
			{VariantID: ecforce.Int64(1), Quantity: ecforce.Ptr(1)},
			{VariantID: ecforce.Int64(999), Quantity: ecforce.Ptr(2)},
		},
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
	if len(result.Errors) != 1 || result.Errors[0].ID != 2 || result.Errors[0].Errors[0].Code != "ASOI2001" {
		t.Errorf("unexpected errors: %+v", result.Errors)
	}
}

func TestSubsOrdersService_BulkUpdateOrderItems(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/subs_orders/4909/order_items/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		items, ok := body["order_items"].([]any)
		if !ok || len(items) != 2 {
			t.Fatalf("order_items wrapper missing or wrong length: %v", body)
		}
		first := items[0].(map[string]any)
		second := items[1].(map[string]any)
		if first["id"] != 1234.0 || first["quantity"] != 3.0 || second["delete"] != 1.0 {
			t.Errorf("unexpected order_items: %v", items)
		}
		fmt.Fprint(w, `{"success": [1234, 1235]}`)
	})

	result, _, err := client.SubsOrders.BulkUpdateOrderItems(context.Background(), 4909, &SubsOrderItemBulkUpdateRequest{
		OrderItems: []*SubsOrderItemParams{
			{ID: ecforce.Int64(1234), Quantity: ecforce.Ptr(3)},
			{ID: ecforce.Int64(1235), Delete: ecforce.Bool01(true)},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 2 || result.Success[0] != 1234 || result.Success[1] != 1235 {
		t.Errorf("Success = %v, want [1234 1235]", result.Success)
	}
	if len(result.Failure) != 0 || len(result.Errors) != 0 {
		t.Errorf("unexpected failures: %+v", result)
	}
}

func TestSubsOrdersService_BulkUpdateSets(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/subs_orders/4909/sets/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		variants, ok := body["variants"].([]any)
		if !ok || len(variants) != 2 {
			t.Fatalf("variants wrapper missing or wrong length: %v", body)
		}
		first := variants[0].(map[string]any)
		if first["id"] != 1234.0 || first["quantity"] != 1.0 {
			t.Errorf("unexpected variants[0]: %v", first)
		}
		fmt.Fprint(w, `{"success": [1, 2]}`)
	})

	result, _, err := client.SubsOrders.BulkUpdateSets(context.Background(), 4909, []*SubsOrderSetParams{
		{ID: ecforce.Int64(1234), Quantity: ecforce.Ptr(1)},
		{ID: ecforce.Int64(5678), Quantity: ecforce.Ptr(2)},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 2 || result.Success[0] != 1 || result.Success[1] != 2 {
		t.Errorf("Success = %v, want [1 2]", result.Success)
	}
}

func TestSubsOrdersService_SplitDeliveryCycle(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/subs_orders/4909/split_delivery_cycle.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		ids, ok := body["order_item_ids"].([]any)
		if !ok || len(ids) != 1 || ids[0] != 157328.0 {
			t.Errorf("unexpected order_item_ids: %v", body["order_item_ids"])
		}
		if body["payment_schedule"] != "term" || body["scheduled_to_be_delivered_every_x_day"] != 7.0 {
			t.Errorf("unexpected schedule params: %v", body)
		}
		fmt.Fprint(w, `{"data": {"id": "5000", "type": "sub_order", "attributes": {
			"id": 5000, "number": "b3997ee5a1", "state": "active", "times": 1,
			"payment_schedule": "7日ごとに配送"
		}}}`)
	})

	subsOrder, _, err := client.SubsOrders.SplitDeliveryCycle(context.Background(), 4909, &SubsOrderSplitDeliveryCycleRequest{
		OrderItemIDs:                    []int64{157328},
		PaymentSchedule:                 ecforce.String("term"),
		ScheduledToBeDeliveredEveryXDay: ecforce.Ptr(7),
		RecalculateScheduledToBeDeliveredAtBasedOnLastOrder: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	attrs := subsOrder.Attributes
	if subsOrder.ID != "5000" || attrs.ID != 5000 || attrs.Number != "b3997ee5a1" {
		t.Errorf("unexpected new subs order: %+v", subsOrder)
	}
	if attrs.PaymentSchedule != "7日ごとに配送" {
		t.Errorf("PaymentSchedule = %q", attrs.PaymentSchedule)
	}
}

func TestSubsOrdersService_ListFreeColumnValues(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/subs_orders/4909/free_column_values.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
			{
				"values": [{
					"free_column_category_seq": 1,
					"values": [
						{"free_column_id": 1, "free_column_label": "name", "free_column_value_values": ["タロウ"]}
					]
				}]
			},
			{
				"free_column_category_id": 1,
				"values": [{
					"free_column_category_seq": 1,
					"values": [
						{"free_column_id": 1, "free_column_label": "name", "free_column_value_values": ["タマ"]},
						{"free_column_id": 2, "free_column_label": "favorites", "free_column_value_values": ["ドライフード", "サバ缶"]}
					]
				}]
			}
		]`)
	})

	groups, _, err := client.SubsOrders.ListFreeColumnValues(context.Background(), 4909)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("len(groups) = %d, want 2", len(groups))
	}
	if groups[0].FreeColumnCategoryID != 0 || groups[1].FreeColumnCategoryID != 1 {
		t.Errorf("unexpected category IDs: %d, %d", groups[0].FreeColumnCategoryID, groups[1].FreeColumnCategoryID)
	}
	entry := groups[1].Values[0]
	if entry.FreeColumnCategorySeq != 1 || len(entry.Values) != 2 {
		t.Fatalf("unexpected entry: %+v", entry)
	}
	if entry.Values[1].FreeColumnLabel != "favorites" || len(entry.Values[1].FreeColumnValueValues) != 2 {
		t.Errorf("unexpected value: %+v", entry.Values[1])
	}
}

func TestSubsOrdersService_BulkCreateFreeColumnValues(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/subs_orders/4909/free_column_values/bulk_create.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		freeColumns, ok := body["free_columns"].([]any)
		if !ok || len(freeColumns) != 2 {
			t.Fatalf("free_columns wrapper missing or wrong length: %v", body)
		}
		first := freeColumns[0].(map[string]any)
		if first["free_column_id"] != 1.0 || first["free_column_value"] != "タロウ" {
			t.Errorf("unexpected free_columns[0]: %v", first)
		}
		second := freeColumns[1].(map[string]any)
		if second["free_column_category_id"] != 1.0 {
			t.Errorf("unexpected free_columns[1]: %v", second)
		}
		fmt.Fprint(w, `{
			"success": [1],
			"failure": [2],
			"errors": [{
				"id": 2, "code": "ASFV0000", "message": "エラーが発生しました。",
				"errors": [{"code": "ASFV3001", "message": "free_column_valueパラメータを入力してください。"}]
			}]
		}`)
	})

	result, _, err := client.SubsOrders.BulkCreateFreeColumnValues(context.Background(), 4909, []*FreeColumnParams{
		{FreeColumnID: ecforce.Int64(1), FreeColumnValue: ecforce.String("タロウ")},
		{
			FreeColumnCategoryID: ecforce.Int64(1),
			Values:               []*FreeColumnValueParams{{FreeColumnID: ecforce.Int64(2)}},
		},
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
	if len(result.Errors) != 1 || result.Errors[0].Errors[0].Code != "ASFV3001" {
		t.Errorf("unexpected errors: %+v", result.Errors)
	}
}

func TestSubsOrdersService_BulkUpdateFreeColumnValues(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/subs_orders/4909/free_column_values/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		freeColumns, ok := body["free_columns"].([]any)
		if !ok || len(freeColumns) != 1 {
			t.Fatalf("free_columns wrapper missing or wrong length: %v", body)
		}
		first := freeColumns[0].(map[string]any)
		if first["free_column_category_id"] != 1.0 || first["free_column_category_seq"] != 2.0 || first["delete"] != 1.0 {
			t.Errorf("unexpected free_columns[0]: %v", first)
		}
		fmt.Fprint(w, `{"success": [1]}`)
	})

	result, _, err := client.SubsOrders.BulkUpdateFreeColumnValues(context.Background(), 4909, []*FreeColumnParams{
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
	if len(result.Failure) != 0 {
		t.Errorf("Failure = %v, want empty", result.Failure)
	}
}

func TestSubsOrdersService_ListCancelReasons(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/cancel_reasons.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("sort"); got != "id" {
			t.Errorf("sort = %q, want id", got)
		}
		fmt.Fprint(w, `{"data": [
			{"id": "1", "type": "cancel_reason", "attributes": {
				"id": 1, "name": "マイページからのキャンセル", "parent_reason_id": null,
				"created_at": "2021/08/18 18:54:28", "updated_at": "2021/08/18 18:54:28"}},
			{"id": "100000", "type": "cancel_reason", "attributes": {
				"id": 100000, "name": "連絡不通", "parent_reason_id": 1,
				"created_at": "2021/10/06 13:32:07", "updated_at": "2021/10/06 13:32:07"}}
		]}`)
	})

	reasons, _, err := client.SubsOrders.ListCancelReasons(context.Background(), &ecforce.ListOptions{
		Sort: []string{"id"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(reasons) != 2 {
		t.Fatalf("len(reasons) = %d, want 2", len(reasons))
	}
	if reasons[0].Attributes.Name != "マイページからのキャンセル" || reasons[0].Attributes.ParentReasonID != 0 {
		t.Errorf("unexpected reasons[0]: %+v", reasons[0].Attributes)
	}
	if reasons[1].Attributes.ID != 100000 || reasons[1].Attributes.ParentReasonID != 1 {
		t.Errorf("unexpected reasons[1]: %+v", reasons[1].Attributes)
	}
}

func TestSubsOrdersService_ListSuspendReasons(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/suspend_reasons.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("per"); got != "50" {
			t.Errorf("per = %q, want 50", got)
		}
		fmt.Fprint(w, `{"data": [
			{"id": "1", "type": "suspend_reason", "attributes": {
				"id": 1, "name": "マイページからの停止", "parent_reason_id": null,
				"created_at": "2021/08/18 18:56:06", "updated_at": "2021/08/18 18:56:06"}},
			{"id": "2", "type": "suspend_reason", "attributes": {
				"id": 2, "name": "金銭的な問題", "parent_reason_id": 1,
				"created_at": "2021/10/07 12:23:41", "updated_at": "2021/10/07 12:23:41"}}
		]}`)
	})

	reasons, _, err := client.SubsOrders.ListSuspendReasons(context.Background(), &ecforce.ListOptions{Per: 50})
	if err != nil {
		t.Fatal(err)
	}
	if len(reasons) != 2 {
		t.Fatalf("len(reasons) = %d, want 2", len(reasons))
	}
	if reasons[0].Attributes.Name != "マイページからの停止" {
		t.Errorf("unexpected reasons[0]: %+v", reasons[0].Attributes)
	}
	if reasons[1].Attributes.ID != 2 || reasons[1].Attributes.ParentReasonID != 1 {
		t.Errorf("unexpected reasons[1]: %+v", reasons[1].Attributes)
	}
}

func TestSubsOrderScheduleChoice_UnmarshalJSON(t *testing.T) {
	var choice SubsOrderScheduleChoice
	if err := json.Unmarshal([]byte(`{"value": 1, "ja": "1ヶ月"}`), &choice); err != nil {
		t.Fatal(err)
	}
	if choice.Value != 1 || choice.Ja != "1ヶ月" {
		t.Errorf("string ja: got %+v, want {Value:1 Ja:1ヶ月}", choice)
	}

	if err := json.Unmarshal([]byte(`{"value": 2, "ja": 2}`), &choice); err != nil {
		t.Fatal(err)
	}
	if choice.Value != 2 || choice.Ja != "2" {
		t.Errorf("numeric ja: got %+v, want {Value:2 Ja:2}", choice)
	}

	if err := json.Unmarshal([]byte(`{"value": 3, "ja": null}`), &choice); err != nil {
		t.Fatal(err)
	}
	if choice.Value != 3 || choice.Ja != "" {
		t.Errorf("null ja: got %+v, want {Value:3 Ja:}", choice)
	}
}
