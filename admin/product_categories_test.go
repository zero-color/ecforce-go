package admin

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/zero-color/ecforce-go"
)

func TestProductCategoriesService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/product_categories.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[name_cont]"); got != "category" {
			t.Errorf("q[name_cont] = %q, want category", got)
		}
		if got := r.URL.Query().Get("sort"); got != "position" {
			t.Errorf("sort = %q, want position", got)
		}
		fmt.Fprint(w, `{"data": [
			{"id": "1", "type": "product_category", "attributes": {
				"id": 1, "name": "category_1", "parent_name": null, "path": "category_1",
				"slug": "slug", "product_category_visibility": true, "description": "",
				"product_extra_text": "",
				"created_at": "2025/05/26 10:26:54", "updated_at": "2025/05/27 10:03:20"}},
			{"id": "2", "type": "product_category", "attributes": {
				"id": 2, "name": "category_2", "parent_name": "category_1",
				"path": "category_1 > category_2", "slug": "slug2",
				"product_category_visibility": true, "description": "詳細テキスト",
				"product_extra_text": "商品テキスト",
				"created_at": "2025/05/26 17:49:41", "updated_at": "2025/05/26 18:05:11"}}]}`)
	})

	categories, _, err := client.ProductCategories.List(context.Background(), &ecforce.ListOptions{
		Sort: []string{"position"},
		Q:    ecforce.Query{"name_cont": "category"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(categories) != 2 {
		t.Fatalf("got %d categories, want 2", len(categories))
	}
	if attrs := categories[0].Attributes; attrs.ID != 1 || attrs.Name != "category_1" || attrs.ParentName != "" {
		t.Errorf("unexpected category: %+v", attrs)
	}
	attrs := categories[1].Attributes
	if attrs.ParentName != "category_1" || attrs.Path != "category_1 > category_2" || attrs.Slug != "slug2" {
		t.Errorf("unexpected category: %+v", attrs)
	}
	if !attrs.ProductCategoryVisibility || attrs.Description != "詳細テキスト" {
		t.Errorf("visibility=%v description=%q", attrs.ProductCategoryVisibility, attrs.Description)
	}
}
