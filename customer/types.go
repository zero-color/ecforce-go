package customer

import (
	"encoding/json"

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
// ID selects an existing shipping address to update; ShippingAddressNote
// applies to shipping addresses only.
type AddressParams struct {
	ID                  *int64  `json:"id,omitempty"`
	Name01              *string `json:"name01,omitempty"`
	Name02              *string `json:"name02,omitempty"`
	Kana01              *string `json:"kana01,omitempty"`
	Kana02              *string `json:"kana02,omitempty"`
	Zip01               *string `json:"zip01,omitempty"`
	Zip02               *string `json:"zip02,omitempty"`
	PrefectureID        *int64  `json:"prefecture_id,omitempty"`
	Addr01              *string `json:"addr01,omitempty"`
	Addr02              *string `json:"addr02,omitempty"`
	Addr03              *string `json:"addr03,omitempty"`
	Tel01               *string `json:"tel01,omitempty"`
	Tel02               *string `json:"tel02,omitempty"`
	Tel03               *string `json:"tel03,omitempty"`
	ShippingAddressNote *string `json:"shipping_address_note,omitempty"`
}

// Note is a memo attached to a customer, side-loaded via include.
type Note struct {
	ID         int64         `json:"id"`
	Content    string        `json:"content"`
	CreatedAt  *ecforce.Time `json:"created_at"`
	UpdatedAt  *ecforce.Time `json:"updated_at"`
	OperatedAt *ecforce.Time `json:"operated_at"`
}

// NoteParams sets a memo on customer update requests. ID selects an existing
// memo to update.
type NoteParams struct {
	ID         *int64        `json:"id,omitempty"`
	Content    *string       `json:"content,omitempty"`
	OperatedAt *ecforce.Time `json:"operated_at,omitempty"`
	OperatedBy *int64        `json:"operated_by,omitempty"`
}

// OrderItem is a line item of a subscription order. It is the attribute union
// of the normal and lighter response variants; subscription order detail
// responses additionally populate the Product* delivery schedule fields.
type OrderItem struct {
	ID                                                int64         `json:"id"`
	VariantID                                         int64         `json:"variant_id"`
	ProductBundledItemID                              int64         `json:"product_bundled_item_id"`
	ProductNumber                                     string        `json:"product_number"`
	ProductName                                       string        `json:"product_name"`
	ProductCategoryNames                              string        `json:"product_category_names"`
	ProductMakerName                                  string        `json:"product_maker_name"`
	ProductPaymentSchedule                            string        `json:"product_payment_schedule"`
	ProductPaymentScheduleLocked                      bool          `json:"product_payment_schedule_locked"`
	ProductScheduledToBeDeliveredEveryXMonth          int           `json:"product_scheduled_to_be_delivered_every_x_month"`
	ProductScheduledToBeDeliveredOnXthDay             int           `json:"product_scheduled_to_be_delivered_on_xth_day"`
	ProductScheduledToBeDeliveredEveryXDay            int           `json:"product_scheduled_to_be_delivered_every_x_day"`
	ProductScheduledToBeDeliveredOnXthDayOfWeek       int           `json:"product_scheduled_to_be_delivered_on_xth_day_of_week"`
	ProductHumanScheduledToBeDeliveredEveryXDayOfWeek string        `json:"product_human_scheduled_to_be_delivered_every_x_day_of_week"`
	VariantSKU                                        string        `json:"variant_sku"`
	StockLocationID                                   int64         `json:"stock_location_id"`
	StockLocationName                                 string        `json:"stock_location_name"`
	ListPrice                                         int           `json:"list_price"`
	SalesPrice                                        int           `json:"sales_price"`
	Price                                             int           `json:"price"`
	Quantity                                          int           `json:"quantity"`
	SetItemQuantity                                   int           `json:"set_item_quantity"`
	TaxRate                                           int           `json:"tax_rate"`
	Onetime                                           bool          `json:"onetime"`
	GiftTargetList                                    []*GiftTarget `json:"gift_target_list"`
	CreatedAt                                         *ecforce.Time `json:"created_at"`
	UpdatedAt                                         *ecforce.Time `json:"updated_at"`
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

// Payment is a subscription order payment, side-loaded via include. It is the
// attribute union of the normal and lighter response variants.
type Payment struct {
	ID                  int64         `json:"id"`
	PaymentMethodID     int64         `json:"payment_method_id"`
	PaymentMethodName   string        `json:"payment_method_name"`
	State               string        `json:"state"`
	HumanState          string        `json:"human_state"`
	PaymentTimes        int           `json:"payment_times"`
	Amount              int           `json:"amount"`
	AccessID            string        `json:"access_id"`
	AccessPass          string        `json:"access_pass"`
	AuthedAt            *ecforce.Time `json:"authed_at"`
	ProviderAuthorizeID string        `json:"provider_authorize_id"`
	CompletedAt         *ecforce.Time `json:"completed_at"`
	PaidAt              *ecforce.Time `json:"paid_at"`
	VoidedAt            *ecforce.Time `json:"voided_at"`
	CreatedAt           *ecforce.Time `json:"created_at"`
	UpdatedAt           *ecforce.Time `json:"updated_at"`
	OPluxResult         string        `json:"o_plux_result"`
	OPluxDescription    string        `json:"o_plux_description"`
	SpiderAfResult      string        `json:"spider_af_result"`
	SpiderAfDescription string        `json:"spider_af_description"`
}

// Variant is a product variation, side-loaded via include.
type Variant struct {
	ID                int64           `json:"id"`
	SKU               string          `json:"sku"`
	Name              string          `json:"name"`
	State             string          `json:"state"`
	HumanStateName    string          `json:"human_state_name"`
	IsMaster          bool            `json:"is_master"`
	ProductID         int64           `json:"product_id"`
	Description       string          `json:"description"`
	DescriptionMobile string          `json:"description_mobile"`
	Volume            int             `json:"volume"`
	ListPrice         int             `json:"list_price"`
	SalesPrice        int             `json:"sales_price"`
	Options           json.RawMessage `json:"options"`
	CreatedAt         *ecforce.Time   `json:"created_at"`
	UpdatedAt         *ecforce.Time   `json:"updated_at"`
	DeletedAt         *ecforce.Time   `json:"deleted_at"`
}
