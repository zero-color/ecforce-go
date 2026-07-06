package admin

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/zero-color/ecforce-go"
)

func TestOrderFreeColumnsService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/order_free_columns.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[id_eq]"); got != "1" {
			t.Errorf("q[id_eq] = %q, want 1", got)
		}
		fmt.Fprint(w, `{
			"data": [{"id": "1", "type": "order_free_column", "attributes": {
				"id": 1, "label": "ラベル名0", "required": false, "placeholder": null,
				"minlength": null, "maxlength": null, "position": 13,
				"export_to_csv": true, "enabled": true, "mypage_display_flg": true,
				"created_at": "2025/09/05 16:01:26", "updated_at": "2025/09/05 16:01:26"}}],
			"meta": {"total_count": 31, "page": 1, "per": 20, "count": 1, "total_pages": 2},
			"links": {"self": "http://localhost/api/v2/admin/order_free_columns?page=1&per=20"}}`)
	})

	columns, resp, err := client.OrderFreeColumns.List(context.Background(), &ecforce.ListOptions{
		Q: ecforce.Query{"id_eq": 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(columns) != 1 {
		t.Fatalf("len(columns) = %d, want 1", len(columns))
	}
	attrs := columns[0].Attributes
	if columns[0].ID != "1" || attrs.ID != 1 || attrs.Label != "ラベル名0" {
		t.Errorf("unexpected column: %+v", attrs)
	}
	if attrs.Required || !attrs.Enabled || !attrs.ExportToCSV || attrs.Position != 13 {
		t.Errorf("unexpected column flags: %+v", attrs)
	}
	if resp.Meta == nil || resp.Meta.TotalCount != 31 || resp.Meta.TotalPages != 2 {
		t.Errorf("unexpected meta: %+v", resp.Meta)
	}
	if !resp.HasNextPage() || resp.NextPage() != 2 {
		t.Errorf("HasNextPage/NextPage = %v/%d, want true/2", resp.HasNextPage(), resp.NextPage())
	}
}
