package admin

import (
	"context"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// OrderFreeColumn is an order free column master record (受注自由項目), the
// definition of a custom question attached to orders. Answered values appear
// on Order.OrderFreeColumns and in OrdersService.ListFreeColumnValues.
type OrderFreeColumn struct {
	ID               int64         `json:"id"`
	Label            string        `json:"label"`
	Required         bool          `json:"required"`
	Placeholder      string        `json:"placeholder"`
	Minlength        int           `json:"minlength"`
	Maxlength        int           `json:"maxlength"`
	Position         int           `json:"position"`
	ExportToCSV      bool          `json:"export_to_csv"`
	Enabled          bool          `json:"enabled"`
	MypageDisplayFlg bool          `json:"mypage_display_flg"`
	CreatedAt        *ecforce.Time `json:"created_at"`
	UpdatedAt        *ecforce.Time `json:"updated_at"`
}

// List searches order free column master records.
//
// Supported q attributes: id, label, required, placeholder, minlength,
// maxlength, position, export_to_csv, enabled, mypage_display_flg,
// created_at, updated_at.
// Supported sort attributes: id, created_at, updated_at.
//
// ecforce API docs: GET /api/v2/admin/order_free_columns
func (s *OrderFreeColumnsService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[OrderFreeColumn], *ecforce.Response, error) {
	return ecforce.DoResourceList[OrderFreeColumn](ctx, s.client, http.MethodGet, "admin/order_free_columns.json", opts.Values(), nil)
}
