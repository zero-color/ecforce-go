package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/zero-color/ecforce-go"
)

func TestStockItemsService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/stock_items.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[variant_sku_eq]"); got != "test_product_01" {
			t.Errorf("q[variant_sku_eq] = %q, want test_product_01", got)
		}
		if got := r.URL.Query().Get("sort"); got != "-created_at,id" {
			t.Errorf("sort = %q, want -created_at,id", got)
		}
		fmt.Fprint(w, `{"data": [{"id": "50113", "type": "stock_item", "attributes": {
			"id": 50113, "stock_location_id": 12, "stock_location_code": "test",
			"stock_location_name": "倉庫C", "product_id": 2312,
			"product_number": "test_product", "variant_id": 3228,
			"variant_sku": "test_product_01", "stock": 0, "stock_unlimited": false,
			"stock_alert": false, "stock_borderline": 0,
			"created_at": "2020/08/07 10:25:05", "updated_at": "2020/08/07 10:25:05",
			"deleted_at": null}}],
			"meta": {"total_count": 1, "page": 1, "per": 1, "count": 1, "total_pages": 1}}`)
	})

	items, resp, err := client.StockItems.List(context.Background(), &ecforce.ListOptions{
		Sort: []string{"-created_at", "id"},
		Q:    ecforce.Query{"variant_sku_eq": "test_product_01"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("got %d stock items, want 1", len(items))
	}
	attrs := items[0].Attributes
	if attrs.ID != 50113 || attrs.VariantSKU != "test_product_01" || attrs.StockLocationCode != "test" {
		t.Errorf("unexpected stock item: %+v", attrs)
	}
	if attrs.Stock != 0 || attrs.StockUnlimited || attrs.StockLocationName != "倉庫C" {
		t.Errorf("stock=%d unlimited=%v location=%q", attrs.Stock, attrs.StockUnlimited, attrs.StockLocationName)
	}
	if resp.Meta == nil || resp.Meta.TotalCount != 1 {
		t.Errorf("meta = %+v, want total_count 1", resp.Meta)
	}
	if resp.HasNextPage() {
		t.Error("HasNextPage() = true, want false")
	}
}

func TestStockItemsService_BulkUpdate(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/stock_items.json", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			StockItems []struct {
				VariantSKU        string `json:"variant_sku"`
				StockLocationCode string `json:"stock_location_code"`
				Stock             int    `json:"stock"`
			} `json:"stock_items"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.StockItems) != 2 {
			t.Fatalf("got %d stock_items in body, want 2", len(body.StockItems))
		}
		if it := body.StockItems[0]; it.VariantSKU != "S-001-001" || it.StockLocationCode != "default" || it.Stock != 100 {
			t.Errorf("unexpected stock item body: %+v", it)
		}
		// Partial failure: success stock is a string, failure stock is a number.
		fmt.Fprint(w, `{
			"success": [
				{"variant_sku": "S-001-001", "stock_location_code": "default", "stock": "100"}
			],
			"failure": [
				{"variant_sku": "SKU_02", "stock_location_code": "logistics_a", "stock": 15}
			],
			"errors": [
				{"id": "SKU_02", "code": "ASI0000", "message": "エラーが発生しました。",
				 "errors": [
					{"code": "ASI2002", "message": "stock_locationに在庫が見つかりません。"},
					{"code": "ASI3002", "message": "stockパラメータが正しくありません。"}
				]}
			]}`)
	})

	result, _, err := client.StockItems.BulkUpdate(context.Background(), []*StockItemParams{
		{
			VariantSKU:        ecforce.String("S-001-001"),
			StockLocationCode: ecforce.String("default"),
			Stock:             ecforce.Ptr(100),
		},
		{
			VariantSKU:        ecforce.String("SKU_02"),
			StockLocationCode: ecforce.String("logistics_a"),
			Stock:             ecforce.Ptr(15),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 1 {
		t.Fatalf("got %d success entries, want 1", len(result.Success))
	}
	// Stock emitted as a quoted string must decode into the int field.
	if s := result.Success[0]; s.VariantSKU != "S-001-001" || s.StockLocationCode != "default" || s.Stock != 100 {
		t.Errorf("unexpected success entry: %+v", s)
	}
	if len(result.Failure) != 1 {
		t.Fatalf("got %d failure entries, want 1", len(result.Failure))
	}
	// Stock emitted as a JSON number must decode too.
	if f := result.Failure[0]; f.VariantSKU != "SKU_02" || f.Stock != 15 {
		t.Errorf("unexpected failure entry: %+v", f)
	}
	if len(result.Errors) != 1 {
		t.Fatalf("got %d errors, want 1", len(result.Errors))
	}
	e := result.Errors[0]
	if e.ID != "SKU_02" || e.Code != "ASI0000" || len(e.Errors) != 2 || e.Errors[0].Code != "ASI2002" {
		t.Errorf("unexpected bulk error: %+v", e)
	}
}
