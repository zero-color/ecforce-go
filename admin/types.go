package admin

import (
	"github.com/zero-color/ecforce-go"
)

// Address is a billing or shipping address. Shipping addresses additionally
// carry ShippingAddressNote.
type Address struct {
	ID                   int64         `json:"id"`
	Name01               string        `json:"name01"`
	Name02               string        `json:"name02"`
	Kana01               string        `json:"kana01"`
	Kana02               string        `json:"kana02"`
	CompanyName          string        `json:"company_name"`
	Zip01                string        `json:"zip01"`
	Zip02                string        `json:"zip02"`
	Addr01               string        `json:"addr01"`
	Addr02               string        `json:"addr02"`
	Addr03               string        `json:"addr03"`
	Tel01                string        `json:"tel01"`
	Tel02                string        `json:"tel02"`
	Tel03                string        `json:"tel03"`
	Tel01Received        string        `json:"tel01_received"`
	Tel02Received        string        `json:"tel02_received"`
	Tel03Received        string        `json:"tel03_received"`
	Fax01                string        `json:"fax01"`
	Fax02                string        `json:"fax02"`
	Fax03                string        `json:"fax03"`
	ShippingAddressNote  string        `json:"shipping_address_note"`
	PrefectureID         int64         `json:"prefecture_id"`
	PrefectureName       string        `json:"prefecture_name"`
	FullName             string        `json:"full_name"`
	FullKana             string        `json:"full_kana"`
	FullTel              string        `json:"full_tel"`
	FullFax              string        `json:"full_fax"`
	FullZip              string        `json:"full_zip"`
	FullAddress          string        `json:"full_address"`
	FullAddressWithSpace string        `json:"full_address_with_space"`
	CreatedAt            *ecforce.Time `json:"created_at"`
	UpdatedAt            *ecforce.Time `json:"updated_at"`
}

// AddressParams sets a billing or shipping address on create/update requests.
// ShippingAddressNote applies to shipping addresses only.
type AddressParams struct {
	ID                  *int64  `json:"id,omitempty"`
	Name01              *string `json:"name01,omitempty"`
	Name02              *string `json:"name02,omitempty"`
	Kana01              *string `json:"kana01,omitempty"`
	Kana02              *string `json:"kana02,omitempty"`
	CompanyName         *string `json:"company_name,omitempty"`
	Zip01               *string `json:"zip01,omitempty"`
	Zip02               *string `json:"zip02,omitempty"`
	PrefectureID        *int64  `json:"prefecture_id,omitempty"`
	Addr01              *string `json:"addr01,omitempty"`
	Addr02              *string `json:"addr02,omitempty"`
	Addr03              *string `json:"addr03,omitempty"`
	Tel01               *string `json:"tel01,omitempty"`
	Tel02               *string `json:"tel02,omitempty"`
	Tel03               *string `json:"tel03,omitempty"`
	Fax01               *string `json:"fax01,omitempty"`
	Fax02               *string `json:"fax02,omitempty"`
	Fax03               *string `json:"fax03,omitempty"`
	ShippingAddressNote *string `json:"shipping_address_note,omitempty"`
}

// OrderItem is a line item of an order or subscription order.
type OrderItem struct {
	ID                    int64         `json:"id"`
	VariantID             int64         `json:"variant_id"`
	ProductBundledItemID  int64         `json:"product_bundled_item_id"`
	ProductNumber         string        `json:"product_number"`
	ProductName           string        `json:"product_name"`
	ProductCategoryNames  string        `json:"product_category_names"`
	ProductMakerName      string        `json:"product_maker_name"`
	VariantSKU            string        `json:"variant_sku"`
	StockLocationID       int64         `json:"stock_location_id"`
	StockLocationName     string        `json:"stock_location_name"`
	ListPrice             int           `json:"list_price"`
	SalesPrice            int           `json:"sales_price"`
	Price                 int           `json:"price"`
	Quantity              int           `json:"quantity"`
	SetItemQuantity       int           `json:"set_item_quantity"`
	PointExchangeQuantity int           `json:"point_exchange_quantity"`
	TaxRate               int           `json:"tax_rate"`
	Onetime               bool          `json:"onetime"`
	Skip                  bool          `json:"skip"`
	GiftTargetList        []*GiftTarget `json:"gift_target_list"`
	CreatedAt             *ecforce.Time `json:"created_at"`
	UpdatedAt             *ecforce.Time `json:"updated_at"`
}

// GiftTarget is gift wrapping/noshi information attached to an order item.
type GiftTarget struct {
	Type                string `json:"type"`
	TargetVariantSKU    string `json:"target_variant_sku"`
	TargetProductNumber string `json:"target_product_number"`
	Quantity            int    `json:"quantity"`
	NoshiOmotegaki      string `json:"noshi_omotegaki"`
	NoshiNaire          string `json:"noshi_naire"`
}
