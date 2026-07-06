package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zero-color/ecforce-go"
)

// StockItem is a stock record (在庫) of one variant at one stock location.
type StockItem struct {
	ID                int64         `json:"id"`
	StockLocationID   int64         `json:"stock_location_id"`
	StockLocationCode string        `json:"stock_location_code"`
	StockLocationName string        `json:"stock_location_name"`
	ProductID         int64         `json:"product_id"`
	ProductNumber     string        `json:"product_number"`
	VariantID         int64         `json:"variant_id"`
	VariantSKU        string        `json:"variant_sku"`
	Stock             int           `json:"stock"`
	StockUnlimited    bool          `json:"stock_unlimited"`
	StockAlert        bool          `json:"stock_alert"`
	StockBorderline   int           `json:"stock_borderline"`
	CreatedAt         *ecforce.Time `json:"created_at"`
	UpdatedAt         *ecforce.Time `json:"updated_at"`
	DeletedAt         *ecforce.Time `json:"deleted_at"`
}

// StockItemParams describes one stock record of a bulk update request. All
// three fields are required by the API.
type StockItemParams struct {
	VariantSKU        *string `json:"variant_sku,omitempty"`
	StockLocationCode *string `json:"stock_location_code,omitempty"`
	Stock             *int    `json:"stock,omitempty"`
}

// StockItemBulkResult is the outcome of StockItemsService.BulkUpdate. Unlike
// ecforce.BulkResult, its success and failure entries are the submitted stock
// records rather than record IDs.
type StockItemBulkResult struct {
	Success []*StockItemBulkEntry `json:"success"`
	Failure []*StockItemBulkEntry `json:"failure"`
	Errors  []*StockItemBulkError `json:"errors"`
}

// StockItemBulkEntry identifies one processed record of a stock item bulk
// update.
type StockItemBulkEntry struct {
	VariantSKU        string `json:"variant_sku"`
	StockLocationCode string `json:"stock_location_code"`
	Stock             int    `json:"stock"`
}

// UnmarshalJSON implements json.Unmarshaler. The API emits the stock count of
// an entry as either a JSON number or a quoted string; both decode into
// Stock.
func (e *StockItemBulkEntry) UnmarshalJSON(data []byte) error {
	var raw struct {
		VariantSKU        string          `json:"variant_sku"`
		StockLocationCode string          `json:"stock_location_code"`
		Stock             json.RawMessage `json:"stock"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.VariantSKU = raw.VariantSKU
	e.StockLocationCode = raw.StockLocationCode
	e.Stock = 0
	s := string(bytes.Trim(bytes.TrimSpace(raw.Stock), `"`))
	if s == "" || s == "null" {
		return nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("ecforce: cannot unmarshal %q into StockItemBulkEntry.Stock", raw.Stock)
	}
	e.Stock = n
	return nil
}

// StockItemBulkError describes why one record of a stock item bulk update
// failed. ID is the variant SKU of the failed record.
type StockItemBulkError struct {
	ID      string          `json:"id"`
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Errors  []ecforce.Error `json:"errors"`
}

// List searches stock records.
//
// Supported q attributes: variant_sku, stock_location_code, stock.
// Supported sort attributes: id, created_at, updated_at.
//
// ecforce API docs: GET /api/v2/admin/stock_items
func (s *StockItemsService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[StockItem], *ecforce.Response, error) {
	return ecforce.DoResourceList[StockItem](ctx, s.client, http.MethodGet, "admin/stock_items.json", opts.Values(), nil)
}

// BulkUpdate updates stock counts in bulk synchronously. Records are keyed by
// variant SKU and stock location code.
//
// ecforce API docs: PUT /api/v2/admin/stock_items
func (s *StockItemsService) BulkUpdate(ctx context.Context, stockItems []*StockItemParams) (*StockItemBulkResult, *ecforce.Response, error) {
	body := map[string][]*StockItemParams{"stock_items": stockItems}
	result := new(StockItemBulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, "admin/stock_items.json", nil, body, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}
