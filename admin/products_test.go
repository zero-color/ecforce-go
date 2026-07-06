package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/zero-color/ecforce-go"
)

func TestProductsService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/products.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[name_cont]"); got != "テスト" {
			t.Errorf("q[name_cont] = %q, want テスト", got)
		}
		if got := r.URL.Query().Get("per"); got != "10" {
			t.Errorf("per = %q, want 10", got)
		}
		fmt.Fprint(w, `{"data": [{"id": "1", "type": "product", "attributes": {
			"id": 1, "number": "S-001", "state": "active", "human_state": "表示",
			"name": "テスト商品（単品）", "for_sale": true, "position": 1,
			"is_recurring": false, "master_list_price": 500, "master_sales_price": 500,
			"master_sku": "S-001-001", "tax_id": 1, "link_number": "111",
			"labels": "ギフト商材A,ギフト商材B",
			"created_at": "2021/05/24 14:14:10", "updated_at": "2021/05/24 14:14:10",
			"deleted_at": null}}],
			"meta": {"total_count": 25, "page": 1, "per": 10, "count": 10, "total_pages": 3}}`)
	})

	products, resp, err := client.Products.List(context.Background(), &ecforce.ListOptions{
		Per: 10,
		Q:   ecforce.Query{"name_cont": "テスト"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 1 {
		t.Fatalf("got %d products, want 1", len(products))
	}
	attrs := products[0].Attributes
	if attrs.ID != 1 || attrs.Number != "S-001" || attrs.Name != "テスト商品（単品）" {
		t.Errorf("unexpected product: %+v", attrs)
	}
	if attrs.MasterSKU != "S-001-001" || !attrs.ForSale {
		t.Errorf("masterSKU=%q forSale=%v", attrs.MasterSKU, attrs.ForSale)
	}
	if resp.Meta == nil || resp.Meta.TotalCount != 25 {
		t.Errorf("meta = %+v, want total_count 25", resp.Meta)
	}
	if !resp.HasNextPage() || resp.NextPage() != 2 {
		t.Errorf("hasNextPage=%v nextPage=%d, want true/2", resp.HasNextPage(), resp.NextPage())
	}
}

func TestProductsService_Get(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/products/1.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("include"); got != "variants" {
			t.Errorf("include = %q, want variants", got)
		}
		fmt.Fprint(w, `{"data": {"id": "1", "type": "product", "attributes": {
			"id": 1, "number": "S-001", "state": "active", "name": "テスト商品（単品）",
			"for_sale": true, "master_list_price": 500, "master_sales_price": 500,
			"master_sku": "S-001-001", "tax_id": 1, "labels": "ギフト商材A,ギフト商材B",
			"created_at": "2021/05/24 14:14:10", "updated_at": "2021/05/24 14:14:10",
			"deleted_at": null},
			"relationships": {"variants": {"data": [{"id": "1", "type": "variant"}]}}}}`)
	})

	product, _, err := client.Products.Get(context.Background(), 1, &ecforce.GetOptions{
		Include: []string{"variants"},
	})
	if err != nil {
		t.Fatal(err)
	}
	attrs := product.Attributes
	if attrs.ID != 1 || attrs.Number != "S-001" || attrs.MasterSalesPrice != 500 {
		t.Errorf("unexpected product: %+v", attrs)
	}
	if attrs.Labels != "ギフト商材A,ギフト商材B" {
		t.Errorf("Labels = %q", attrs.Labels)
	}
	if rel := product.Relationships["variants"]; rel == nil || len(rel.Data) != 1 || rel.Data[0].ID != "1" {
		t.Errorf("unexpected variants relationship: %+v", product.Relationships)
	}
}

func TestProductsService_BulkCreate(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/products/bulk_create.json", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Products []struct {
				Name              string `json:"name"`
				Number            string `json:"number"`
				IsRecurring       int    `json:"is_recurring"`
				VariantAttributes struct {
					SKU        string `json:"sku"`
					SalesPrice int    `json:"sales_price"`
					ForSale    int    `json:"for_sale"`
				} `json:"variant_attributes"`
			} `json:"products"`
			CheckDuplicateLinkNumbers int `json:"check_duplicate_link_numbers"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Products) != 1 {
			t.Fatalf("got %d products in body, want 1", len(body.Products))
		}
		p := body.Products[0]
		if p.Name != "サンプル単品" || p.Number != "sample-single001" || p.IsRecurring != 0 {
			t.Errorf("unexpected product body: %+v", p)
		}
		if p.VariantAttributes.SKU != "sample-single001" || p.VariantAttributes.SalesPrice != 1000 || p.VariantAttributes.ForSale != 1 {
			t.Errorf("unexpected variant_attributes: %+v", p.VariantAttributes)
		}
		if body.CheckDuplicateLinkNumbers != 1 {
			t.Errorf("check_duplicate_link_numbers = %d, want 1", body.CheckDuplicateLinkNumbers)
		}
		fmt.Fprint(w, `{"id": 123, "job_id": "6cb0d3e5-0e5a-4a4b-a4b5-111111111111"}`)
	})

	result, _, err := client.Products.BulkCreate(context.Background(), &ProductBulkRequest{
		Products: []*ProductParams{{
			Name:        ecforce.String("サンプル単品"),
			Number:      ecforce.String("sample-single001"),
			IsRecurring: ecforce.Bool01(false),
			VariantAttributes: &VariantParams{
				SKU:        ecforce.String("sample-single001"),
				SalesPrice: ecforce.Ptr(1000),
				ForSale:    ecforce.Bool01(true),
			},
		}},
		CheckDuplicateLinkNumbers: ecforce.Bool01(true),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != 123 || result.JobID != "6cb0d3e5-0e5a-4a4b-a4b5-111111111111" {
		t.Errorf("unexpected job result: %+v", result)
	}
}

func TestProductsService_BulkUpdate(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/products/bulk_update.json", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Products []struct {
				ID               int64  `json:"id"`
				State            string `json:"state"`
				MasterAttributes struct {
					SalesPrice int `json:"sales_price"`
				} `json:"master_attributes"`
			} `json:"products"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Products) != 3 {
			t.Fatalf("got %d products in body, want 3", len(body.Products))
		}
		if p := body.Products[0]; p.ID != 77777 || p.State != "active" || p.MasterAttributes.SalesPrice != 1200 {
			t.Errorf("unexpected product body: %+v", p)
		}
		// Partial failure response.
		fmt.Fprint(w, `{
			"success": [77777],
			"failure": [88888, 99999],
			"errors": [
				{"id": 88888, "code": "APR0000", "message": "エラーが発生しました。",
				 "errors": [{"code": "APR3007", "message": "単品商品のため配送サイクルを更新できません。"}]},
				{"id": 99999, "code": "APR0000", "message": "エラーが発生しました。",
				 "errors": [{"code": "APR3006", "message": "配送サイクルを正しく入力してください"}]}
			]}`)
	})

	result, _, err := client.Products.BulkUpdate(context.Background(), &ProductBulkRequest{
		Products: []*ProductParams{
			{
				ID:               ecforce.Int64(77777),
				State:            ecforce.String("active"),
				MasterAttributes: &VariantParams{SalesPrice: ecforce.Ptr(1200)},
			},
			{ID: ecforce.Int64(88888)},
			{ID: ecforce.Int64(99999)},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 1 || result.Success[0] != 77777 {
		t.Errorf("Success = %v, want [77777]", result.Success)
	}
	if len(result.Failure) != 2 || result.Failure[0] != 88888 || result.Failure[1] != 99999 {
		t.Errorf("Failure = %v, want [88888 99999]", result.Failure)
	}
	if len(result.Errors) != 2 {
		t.Fatalf("got %d errors, want 2", len(result.Errors))
	}
	e := result.Errors[0]
	if e.ID != 88888 || e.Code != "APR0000" || len(e.Errors) != 1 || e.Errors[0].Code != "APR3007" {
		t.Errorf("unexpected bulk error: %+v", e)
	}
}

func TestProductsService_BulkDestroy(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("DELETE /api/v2/admin/products/bulk_destroy.json", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ProductIDs []int64 `json:"product_ids"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.ProductIDs) != 3 || body.ProductIDs[0] != 1 || body.ProductIDs[2] != 3 {
			t.Errorf("product_ids = %v, want [1 2 3]", body.ProductIDs)
		}
		fmt.Fprint(w, `{"success": [1, 2, 3]}`)
	})

	result, _, err := client.Products.BulkDestroy(context.Background(), []int64{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 3 || result.Success[0] != 1 || result.Success[2] != 3 {
		t.Errorf("Success = %v, want [1 2 3]", result.Success)
	}
	if len(result.Failure) != 0 || len(result.Errors) != 0 {
		t.Errorf("Failure = %v, Errors = %v, want empty", result.Failure, result.Errors)
	}
}
