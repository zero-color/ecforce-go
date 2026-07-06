package admin

import (
	"context"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// Store is a sales-channel (媒体) master record.
type Store struct {
	ID                int64         `json:"id"`
	Name              string        `json:"name"`
	OrdersCount       int           `json:"orders_count"`
	LPDefault         bool          `json:"lp_default"`
	CartDefault       bool          `json:"cart_default"`
	CSDefault         bool          `json:"cs_default"`
	OrderInputDefault bool          `json:"order_input_default"`
	Labels            []string      `json:"labels"`
	ExpiryDateStartAt *ecforce.Time `json:"expiry_date_start_at"`
	ExpiryDateEndAt   *ecforce.Time `json:"expiry_date_end_at"`
	Cost              int           `json:"cost"`
	CreatedAt         *ecforce.Time `json:"created_at"`
	UpdatedAt         *ecforce.Time `json:"updated_at"`
	DeletedAt         *ecforce.Time `json:"deleted_at"`
}

// List searches store (sales channel) master records.
//
// Supported q attributes: id, name, orders_count, with_deleted, lp_default,
// cart_default, order_input_default, created_at, updated_at.
// Supported sort attributes: id, created_at, updated_at.
//
// ecforce API docs: GET /api/v2/admin/stores
func (s *StoresService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[Store], *ecforce.Response, error) {
	return ecforce.DoResourceList[Store](ctx, s.client, http.MethodGet, "admin/stores.json", opts.Values(), nil)
}
