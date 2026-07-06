package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/zero-color/ecforce-go"
)

func TestLabelsService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/labels.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[type_eq]"); got != "Customer" {
			t.Errorf("q[type_eq] = %q", got)
		}
		if got := r.URL.Query().Get("per"); got != "20" {
			t.Errorf("per = %q", got)
		}
		fmt.Fprint(w, `{
			"data": [
				{"id": 1, "name": "customer_label", "color": "#FFFFFF", "position": 0,
					"search_form_visibility": null, "type": "Customer"},
				{"id": 2, "name": "customer_label_2", "color": "", "position": 1,
					"search_form_visibility": null, "type": "Customer"}],
			"meta": {"total_count": 2, "page": 1, "per": 20, "count": 2, "total_pages": 1}}`)
	})

	labels, resp, err := client.Labels.List(context.Background(), &ecforce.ListOptions{
		Per: 20,
		Q:   ecforce.Query{"type_eq": "Customer"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 2 {
		t.Fatalf("labels = %+v", labels)
	}
	if labels[0].ID != 1 || labels[0].Name != "customer_label" || labels[0].Color != "#FFFFFF" {
		t.Errorf("unexpected label: %+v", labels[0])
	}
	if labels[1].Position != 1 || labels[1].Type != "Customer" {
		t.Errorf("unexpected label: %+v", labels[1])
	}
	if resp.Meta == nil || resp.Meta.TotalCount != 2 {
		t.Errorf("unexpected meta: %+v", resp.Meta)
	}
	if resp.HasNextPage() {
		t.Error("HasNextPage = true, want false")
	}
}

func TestLabelsService_Create(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("POST /api/v2/admin/labels.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := body["type"]; got != "Customer" {
			t.Errorf("type = %v", got)
		}
		if got := body["name"]; got != "customer_label" {
			t.Errorf("name = %v", got)
		}
		if got := body["color"]; got != "#FFFFFF" {
			t.Errorf("color = %v", got)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"data": {"id": 1, "name": "customer_label", "color": "#FFFFFF",
			"position": 0, "search_form_visibility": false, "type": "Customer"}}`)
	})

	label, _, err := client.Labels.Create(context.Background(), &LabelCreateRequest{
		Type:  ecforce.String("Customer"),
		Name:  ecforce.String("customer_label"),
		Color: ecforce.String("#FFFFFF"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if label.ID != 1 || label.Name != "customer_label" || label.Type != "Customer" {
		t.Errorf("unexpected label: %+v", label)
	}
	if label.SearchFormVisibility {
		t.Error("SearchFormVisibility = true, want false")
	}
}

func TestLabelsService_Apply(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("PUT /api/v2/admin/labeling.json", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if got := body["type"]; got != "Order" {
			t.Errorf("type = %v", got)
		}
		if got := body["method"]; got != "overwrite" {
			t.Errorf("method = %v", got)
		}
		if got := fmt.Sprint(body["target_ids"]); got != "[1 2 3 4 5]" {
			t.Errorf("target_ids = %v", got)
		}
		if got := fmt.Sprint(body["label_ids"]); got != "[1 2 3]" {
			t.Errorf("label_ids = %v", got)
		}
		fmt.Fprint(w, `{
			"success": [1, 2, 3],
			"failure": [4, 5],
			"errors": [
				{"id": 4, "code": "ALA0000", "message": "エラーが発生しました。",
					"errors": [{"code": "ALA3999", "message": "商品が重複しています"}]},
				{"id": 5, "code": "ALA0000", "message": "エラーが発生しました。",
					"errors": [{"code": "ALA9999", "message": "予期せぬエラーが発生しました。"}]}]}`)
	})

	result, _, err := client.Labels.Apply(context.Background(), &LabelingRequest{
		Type:      ecforce.String("Order"),
		Method:    ecforce.String("overwrite"),
		TargetIDs: []int64{1, 2, 3, 4, 5},
		LabelIDs:  []int64{1, 2, 3},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 3 || result.Success[0] != 1 {
		t.Errorf("success = %v", result.Success)
	}
	if len(result.Failure) != 2 || result.Failure[0] != 4 || result.Failure[1] != 5 {
		t.Errorf("failure = %v", result.Failure)
	}
	if len(result.Errors) != 2 || result.Errors[0].ID != 4 || result.Errors[0].Errors[0].Code != "ALA3999" {
		t.Errorf("errors = %+v", result.Errors)
	}
	if result.Errors[1].Code != "ALA0000" {
		t.Errorf("errors[1].code = %q", result.Errors[1].Code)
	}
}
