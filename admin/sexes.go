package admin

import (
	"context"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// Sex is a sex master record.
type Sex struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Position  int           `json:"position"`
	State     string        `json:"state"`
	CreatedAt *ecforce.Time `json:"created_at"`
	UpdatedAt *ecforce.Time `json:"updated_at"`
	DeletedAt *ecforce.Time `json:"deleted_at"`
}

// List searches sex master records.
//
// Supported q attributes: id, name, state, with_deleted.
// Supported sort attributes: id, created_at, updated_at, position.
//
// ecforce API docs: GET /api/v2/admin/sexes
func (s *SexesService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[Sex], *ecforce.Response, error) {
	return ecforce.DoResourceList[Sex](ctx, s.client, http.MethodGet, "admin/sexes.json", opts.Values(), nil)
}
