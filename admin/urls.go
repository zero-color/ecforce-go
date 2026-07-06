package admin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// URL is an advertising URL (広告URL).
type URL struct {
	ID             int64         `json:"id"`
	BaseURL        string        `json:"base_url"`
	URLGroupID     int64         `json:"url_group_id"`
	URLGroupName   string        `json:"url_group_name"`
	AdvertiserID   int64         `json:"advertiser_id"`
	AdvertiserName string        `json:"advertiser_name"`
	Description    string        `json:"description"`
	PurchaseCount  int           `json:"purchase_count"`
	Upsell         bool          `json:"upsell"`
	Cost           int           `json:"cost"`
	Default        bool          `json:"default"`
	CreatedAt      *ecforce.Time `json:"created_at"`
	UpdatedAt      *ecforce.Time `json:"updated_at"`
	DeletedAt      *ecforce.Time `json:"deleted_at"`
}

// Advertiser is an advertiser (広告主) side-loaded via include.
type Advertiser struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Number    string        `json:"number"`
	Email     string        `json:"email"`
	CreatedAt *ecforce.Time `json:"created_at"`
	UpdatedAt *ecforce.Time `json:"updated_at"`
}

// Get fetches a single advertising URL.
//
// Supported include values: url_group, url_group.advertiser.
//
// ecforce API docs: GET /api/v2/admin/urls/:url_id
func (s *URLsService) Get(ctx context.Context, urlID int64, opts *ecforce.GetOptions) (*ecforce.Resource[URL], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/urls/%d.json", urlID)
	return ecforce.DoResource[URL](ctx, s.client, http.MethodGet, path, opts.Values(), nil)
}

// List searches advertising URLs.
//
// Supported q attributes: id, base_url, url_group_id, created_at, updated_at,
// with_deleted.
// Supported sort attributes: id, created_at, updated_at.
// Supported include values: url_group, url_group.advertiser.
//
// ecforce API docs: GET /api/v2/admin/urls
func (s *URLsService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[URL], *ecforce.Response, error) {
	return ecforce.DoResourceList[URL](ctx, s.client, http.MethodGet, "admin/urls.json", opts.Values(), nil)
}
