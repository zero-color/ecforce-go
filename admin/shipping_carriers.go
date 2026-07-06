package admin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// ShippingCarrier is a shipping carrier (配送業者).
type ShippingCarrier struct {
	ID                        int64             `json:"id"`
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	DescriptionMobile         string            `json:"description_mobile"`
	Default                   bool              `json:"default"`
	State                     string            `json:"state"`
	HumanState                string            `json:"human_state"`
	MinVolume                 int               `json:"min_volume"`
	MaxVolume                 int               `json:"max_volume"`
	PickupLocations           []*PickupLocation `json:"pickup_locations"`
	PickupLocationPrimaryOnly ecforce.BoolInt   `json:"pickup_location_primary_only"`
	EnableDoorbellOption      ecforce.BoolInt   `json:"enable_doorbell_option"`
	CreatedAt                 *ecforce.Time     `json:"created_at"`
	UpdatedAt                 *ecforce.Time     `json:"updated_at"`
}

// PickupLocation is a pickup location (受取場所) offered by a shipping carrier.
type PickupLocation struct {
	ID                 int64  `json:"id"`
	PickupLocationName string `json:"pickup_location_name"`
}

// DeliveryTime is a delivery time slot (配送時間) side-loaded via include.
type DeliveryTime struct {
	ID   int64  `json:"id"`
	Time string `json:"time"`
}

// Get fetches a single shipping carrier.
//
// Supported include values: delivery_times.
//
// ecforce API docs: GET /api/v2/admin/shipping_carriers/:shipping_carrier_id
func (s *ShippingCarriersService) Get(ctx context.Context, shippingCarrierID int64, opts *ecforce.GetOptions) (*ecforce.Resource[ShippingCarrier], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/shipping_carriers/%d.json", shippingCarrierID)
	return ecforce.DoResource[ShippingCarrier](ctx, s.client, http.MethodGet, path, opts.Values(), nil)
}

// List searches shipping carriers.
//
// Supported q attributes: id, name, created_at, updated_at, with_deleted.
// Supported sort attributes: id, created_at, updated_at.
// Supported include values: delivery_times.
//
// ecforce API docs: GET /api/v2/admin/shipping_carriers
func (s *ShippingCarriersService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[ShippingCarrier], *ecforce.Response, error) {
	return ecforce.DoResourceList[ShippingCarrier](ctx, s.client, http.MethodGet, "admin/shipping_carriers.json", opts.Values(), nil)
}
