package admin

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/zero-color/ecforce-go"
)

func TestStoresService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/stores.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[name_cont]"); got != "本店" {
			t.Errorf("q[name_cont] = %q, want 本店", got)
		}
		if got := r.URL.Query().Get("sort"); got != "id" {
			t.Errorf("sort = %q, want id", got)
		}
		fmt.Fprint(w, `{"data": [{"id": "1", "type": "store", "attributes": {
			"id": 1, "name": "本店", "orders_count": 0, "lp_default": false,
			"cart_default": false, "cs_default": true, "order_input_default": true,
			"labels": ["ラベル1", "ラベル2", "ラベル3"],
			"expiry_date_start_at": null, "expiry_date_end_at": null, "cost": null,
			"created_at": "2025/07/17 10:33:54", "updated_at": "2025/07/23 07:14:19",
			"deleted_at": null}}]}`)
	})

	stores, _, err := client.Stores.List(context.Background(), &ecforce.ListOptions{
		Sort: []string{"id"},
		Q:    ecforce.Query{"name_cont": "本店"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(stores) != 1 {
		t.Fatalf("got %d stores, want 1", len(stores))
	}
	attrs := stores[0].Attributes
	if attrs.ID != 1 || attrs.Name != "本店" || !attrs.CSDefault || attrs.LPDefault {
		t.Errorf("unexpected store: %+v", attrs)
	}
	if len(attrs.Labels) != 3 || attrs.Labels[0] != "ラベル1" {
		t.Errorf("Labels = %v", attrs.Labels)
	}
}

func TestURLsService_Get(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/urls/489.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("include"); got != "url_group" {
			t.Errorf("include = %q, want url_group", got)
		}
		fmt.Fprint(w, `{"data": {"id": "489", "type": "url", "attributes": {
			"id": 489, "base_url": "foo3", "url_group_id": 1, "url_group_name": "test",
			"advertiser_id": 1, "advertiser_name": "test", "description": "test",
			"purchase_count": 2, "upsell": false, "cost": 1000, "default": false,
			"created_at": "2019/02/01 05:15:10", "updated_at": "2019/02/01 05:15:10",
			"deleted_at": null},
			"relationships": {"url_group": {"data": {"id": "36", "type": "url_group"}}}}}`)
	})

	url, _, err := client.URLs.Get(context.Background(), 489, &ecforce.GetOptions{
		Include: []string{"url_group"},
	})
	if err != nil {
		t.Fatal(err)
	}
	attrs := url.Attributes
	if attrs.ID != 489 || attrs.BaseURL != "foo3" || attrs.URLGroupName != "test" {
		t.Errorf("unexpected url: %+v", attrs)
	}
	if attrs.PurchaseCount != 2 || attrs.Cost != 1000 || attrs.Upsell {
		t.Errorf("purchaseCount=%d cost=%d upsell=%v", attrs.PurchaseCount, attrs.Cost, attrs.Upsell)
	}
	if rel := url.Relationships["url_group"]; rel == nil || len(rel.Data) != 1 || rel.Data[0].ID != "36" {
		t.Errorf("unexpected url_group relationship: %+v", url.Relationships)
	}
}

func TestURLsService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/urls.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("page = %q, want 1", got)
		}
		if got := r.URL.Query().Get("per"); got != "1" {
			t.Errorf("per = %q, want 1", got)
		}
		fmt.Fprint(w, `{"data": [{"id": "489", "type": "url", "attributes": {
			"id": 489, "base_url": "foo3", "url_group_id": 1, "url_group_name": "test",
			"advertiser_id": 1, "advertiser_name": "test", "description": "test",
			"purchase_count": 2, "upsell": false, "cost": 1000, "default": false,
			"created_at": "2019/02/01 05:15:10", "updated_at": "2019/02/01 05:15:10",
			"deleted_at": null}}],
			"meta": {"total_count": 299, "page": 1, "per": 1, "count": 1, "total_pages": 299}}`)
	})

	urls, resp, err := client.URLs.List(context.Background(), &ecforce.ListOptions{Page: 1, Per: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(urls) != 1 {
		t.Fatalf("got %d urls, want 1", len(urls))
	}
	if attrs := urls[0].Attributes; attrs.ID != 489 || attrs.BaseURL != "foo3" || attrs.URLGroupID != 1 {
		t.Errorf("unexpected url: %+v", attrs)
	}
	if resp.Meta == nil || resp.Meta.TotalCount != 299 {
		t.Errorf("meta = %+v, want total_count 299", resp.Meta)
	}
	if !resp.HasNextPage() || resp.NextPage() != 2 {
		t.Errorf("hasNextPage=%v nextPage=%d, want true/2", resp.HasNextPage(), resp.NextPage())
	}
}

func TestURLGroupsService_Get(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/url_groups/66.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("include"); got != "advertiser" {
			t.Errorf("include = %q, want advertiser", got)
		}
		fmt.Fprint(w, `{"data": {"id": "66", "type": "url_group", "attributes": {
			"id": 66, "name": "test173", "description": "",
			"created_at": "2018/10/05 16:53:41", "updated_at": "2018/10/05 16:53:41"},
			"relationships": {"advertiser": {"data": {"id": "65", "type": "advertiser"}}}}}`)
	})

	group, _, err := client.URLGroups.Get(context.Background(), 66, &ecforce.GetOptions{
		Include: []string{"advertiser"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if attrs := group.Attributes; attrs.ID != 66 || attrs.Name != "test173" || attrs.Description != "" {
		t.Errorf("unexpected url group: %+v", attrs)
	}
	if rel := group.Relationships["advertiser"]; rel == nil || len(rel.Data) != 1 || rel.Data[0].ID != "65" {
		t.Errorf("unexpected advertiser relationship: %+v", group.Relationships)
	}
}

func TestURLGroupsService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/url_groups.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[advertiser_id_eq]"); got != "65" {
			t.Errorf("q[advertiser_id_eq] = %q, want 65", got)
		}
		fmt.Fprint(w, `{"data": [{"id": "66", "type": "url_group", "attributes": {
			"id": 66, "name": "test173", "description": "",
			"created_at": "2018/10/05 16:53:41", "updated_at": "2018/10/05 16:53:41"}}],
			"meta": {"total_count": 36, "page": 1, "per": 1, "count": 1, "total_pages": 36}}`)
	})

	groups, resp, err := client.URLGroups.List(context.Background(), &ecforce.ListOptions{
		Q: ecforce.Query{"advertiser_id_eq": 65},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 1 {
		t.Fatalf("got %d url groups, want 1", len(groups))
	}
	if attrs := groups[0].Attributes; attrs.ID != 66 || attrs.Name != "test173" {
		t.Errorf("unexpected url group: %+v", attrs)
	}
	if resp.Meta == nil || resp.Meta.TotalPages != 36 || !resp.HasNextPage() {
		t.Errorf("meta = %+v hasNextPage=%v, want total_pages 36/true", resp.Meta, resp.HasNextPage())
	}
}

func TestShippingCarriersService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/shipping_carriers.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[name_cont]"); got != "運輸" {
			t.Errorf("q[name_cont] = %q, want 運輸", got)
		}
		fmt.Fprint(w, `{"data": [{"id": "99", "type": "shipping_carrier", "attributes": {
			"id": 99, "name": "外部運輸", "description": null, "description_mobile": null,
			"default": false, "state": "inactive", "human_state": "無効",
			"min_volume": 0, "max_volume": 0,
			"pickup_locations": [
				{"id": 1, "pickup_location_name": "宅配ボックス"},
				{"id": 2, "pickup_location_name": "玄関"}],
			"pickup_location_primary_only": 0, "enable_doorbell_option": true,
			"created_at": "2019/03/29 18:38:36", "updated_at": "2019/03/29 18:38:36"}}],
			"meta": {"total_count": 27, "page": 1, "per": 1, "count": 1, "total_pages": 27}}`)
	})

	carriers, resp, err := client.ShippingCarriers.List(context.Background(), &ecforce.ListOptions{
		Q: ecforce.Query{"name_cont": "運輸"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(carriers) != 1 {
		t.Fatalf("got %d carriers, want 1", len(carriers))
	}
	attrs := carriers[0].Attributes
	if attrs.ID != 99 || attrs.Name != "外部運輸" || attrs.State != "inactive" {
		t.Errorf("unexpected carrier: %+v", attrs)
	}
	if len(attrs.PickupLocations) != 2 || attrs.PickupLocations[1].PickupLocationName != "玄関" {
		t.Errorf("PickupLocations = %+v", attrs.PickupLocations)
	}
	if attrs.PickupLocationPrimaryOnly || !attrs.EnableDoorbellOption {
		t.Errorf("bool fields: primaryOnly=%v doorbell=%v", attrs.PickupLocationPrimaryOnly, attrs.EnableDoorbellOption)
	}
	if resp.Meta == nil || resp.Meta.TotalCount != 27 {
		t.Errorf("meta = %+v, want total_count 27", resp.Meta)
	}
}

func TestAdvertisementsService_List(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/advertisements.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q[created_at_gteq]"); got != "2020-10-01 00:00:00" {
			t.Errorf("q[created_at_gteq] = %q, want 2020-10-01 00:00:00", got)
		}
		fmt.Fprint(w, `{
			"data": {"access_count": [
				{"id": 1, "base_url": "index", "total": 5593, "pc": 1233, "sp": 4360, "other": 0},
				{"id": 2, "base_url": "sample", "total": 10510, "pc": 789, "sp": 9721, "other": 0}
			]},
			"created_at_gteq": "2020/10/12 00:00:00",
			"created_at_lt": "2020/10/19 23:59:59"}`)
	})

	result, _, err := client.Advertisements.List(context.Background(), &ecforce.ListOptions{
		Q: ecforce.Query{"created_at_gteq": "2020-10-01 00:00:00"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.AccessCounts) != 2 {
		t.Fatalf("got %d access counts, want 2", len(result.AccessCounts))
	}
	ac := result.AccessCounts[0]
	if ac.ID != 1 || ac.BaseURL != "index" || ac.Total != 5593 || ac.PC != 1233 || ac.SP != 4360 {
		t.Errorf("unexpected access count: %+v", ac)
	}
	if result.CreatedAtGteq == nil || result.CreatedAtGteq.Format("2006/01/02 15:04:05") != "2020/10/12 00:00:00" {
		t.Errorf("CreatedAtGteq = %v, want 2020/10/12 00:00:00", result.CreatedAtGteq)
	}
	if result.CreatedAtLt == nil || result.CreatedAtLt.Format("2006/01/02 15:04:05") != "2020/10/19 23:59:59" {
		t.Errorf("CreatedAtLt = %v, want 2020/10/19 23:59:59", result.CreatedAtLt)
	}
}
