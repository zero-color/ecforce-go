package admin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// URLGroup is an advertising URL group (広告URLグループ).
type URLGroup struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CreatedAt   *ecforce.Time `json:"created_at"`
	UpdatedAt   *ecforce.Time `json:"updated_at"`
}

// Get fetches a single advertising URL group.
//
// Supported include values: advertiser.
//
// ecforce API docs: GET /api/v2/admin/url_groups/:url_group_id
func (s *URLGroupsService) Get(ctx context.Context, urlGroupID int64, opts *ecforce.GetOptions) (*ecforce.Resource[URLGroup], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/url_groups/%d.json", urlGroupID)
	return ecforce.DoResource[URLGroup](ctx, s.client, http.MethodGet, path, opts.Values(), nil)
}

// List searches advertising URL groups.
//
// Supported q attributes: id, name, advertiser_id, created_at, updated_at,
// with_deleted.
// Supported sort attributes: id, created_at, updated_at.
// Supported include values: advertiser.
//
// ecforce API docs: GET /api/v2/admin/url_groups
func (s *URLGroupsService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[URLGroup], *ecforce.Response, error) {
	return ecforce.DoResourceList[URLGroup](ctx, s.client, http.MethodGet, "admin/url_groups.json", opts.Values(), nil)
}
