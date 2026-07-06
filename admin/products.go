package admin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// Product is a product (商品). It carries the union of the attributes of the
// normal and lighter response variants; fields absent from the lighter
// variant (e.g. MakerName, MasterSKU, Labels) are zero when lighter is on.
type Product struct {
	ID                    int64         `json:"id"`
	Number                string        `json:"number"`
	State                 string        `json:"state"`
	HumanState            string        `json:"human_state"`
	Name                  string        `json:"name"`
	ForSale               bool          `json:"for_sale"`
	Position              int           `json:"position"`
	UpsellProductID       int64         `json:"upsell_product_id"`
	UpsellProductNumber   string        `json:"upsell_product_number"`
	UpsellProductName     string        `json:"upsell_product_name"`
	CvUpsellProductID     int64         `json:"cv_upsell_product_id"`
	CvUpsellProductNumber string        `json:"cv_upsell_product_number"`
	CvUpsellProductName   string        `json:"cv_upsell_product_name"`
	MakerID               int64         `json:"maker_id"`
	MakerName             string        `json:"maker_name"`
	Description           string        `json:"description"`
	DescriptionMobile     string        `json:"description_mobile"`
	SubDescription        string        `json:"sub_description"`
	SubDescriptionMobile  string        `json:"sub_description_mobile"`
	MetaDescription       string        `json:"meta_description"`
	MetaKeywords          string        `json:"meta_keywords"`
	IsRecurring           bool          `json:"is_recurring"`
	ProductCategoryNames  string        `json:"product_category_names"`
	MasterListPrice       int           `json:"master_list_price"`
	MasterSalesPrice      int           `json:"master_sales_price"`
	MasterSKU             string        `json:"master_sku"`
	TaxID                 int64         `json:"tax_id"`
	LinkNumber            string        `json:"link_number"`
	Option01              string        `json:"option01"`
	Option02              string        `json:"option02"`
	Option03              string        `json:"option03"`
	Option04              string        `json:"option04"`
	Option05              string        `json:"option05"`
	Option06              string        `json:"option06"`
	Option07              string        `json:"option07"`
	Option08              string        `json:"option08"`
	Option09              string        `json:"option09"`
	Option10              string        `json:"option10"`
	Labels                string        `json:"labels"`
	FirstPrice            int           `json:"first_price"`
	FirstPriceIncludeTax  int           `json:"first_price_include_tax"`
	CreatedAt             *ecforce.Time `json:"created_at"`
	UpdatedAt             *ecforce.Time `json:"updated_at"`
	DeletedAt             *ecforce.Time `json:"deleted_at"`
}

// Variant is a product variation (商品バリエーション/SKU) side-loaded via
// include. It carries the union of the normal and lighter attribute sets.
type Variant struct {
	ID                   int64            `json:"id"`
	SKU                  string           `json:"sku"`
	Name                 string           `json:"name"`
	State                string           `json:"state"`
	HumanStateName       string           `json:"human_state_name"`
	ForSale              bool             `json:"for_sale"`
	Position             int              `json:"position"`
	IsMaster             bool             `json:"is_master"`
	ProductID            int64            `json:"product_id"`
	Description          string           `json:"description"`
	DescriptionMobile    string           `json:"description_mobile"`
	Volume               int              `json:"volume"`
	ListPrice            int              `json:"list_price"`
	SalesPrice           int              `json:"sales_price"`
	Options              []*VariantOption `json:"options"`
	LinkNumber           string           `json:"link_number"`
	FirstPrice           int              `json:"first_price"`
	FirstPriceIncludeTax int              `json:"first_price_include_tax"`
	CreatedAt            *ecforce.Time    `json:"created_at"`
	UpdatedAt            *ecforce.Time    `json:"updated_at"`
	DeletedAt            *ecforce.Time    `json:"deleted_at"`
}

// VariantOption is one 規格 (option type) / 分類 (option value) pair of a
// variant.
type VariantOption struct {
	OptionType  string `json:"option_type"`
	OptionValue string `json:"option_value"`
}

// BundledItem is a bundled item (同梱物) of a product, side-loaded via
// include.
type BundledItem struct {
	ID        int64 `json:"id"`
	VariantID int64 `json:"variant_id"`
	Quantity  int   `json:"quantity"`
}

// Thumbnail is a product or variant image, side-loaded via include. The URL
// and path fields are absent from the lighter response variant.
type Thumbnail struct {
	ID          int64         `json:"id"`
	FileName    string        `json:"file_name"`
	ContentType string        `json:"content_type"`
	FileSize    int           `json:"file_size"`
	Position    int           `json:"position"`
	URL         string        `json:"url"`
	URLSmall    string        `json:"url_small"`
	URLMedium   string        `json:"url_medium"`
	URLLarge    string        `json:"url_large"`
	Path        string        `json:"path"`
	PathSmall   string        `json:"path_small"`
	PathMedium  string        `json:"path_medium"`
	PathLarge   string        `json:"path_large"`
	CreatedAt   *ecforce.Time `json:"created_at"`
	UpdatedAt   *ecforce.Time `json:"updated_at"`
}

// ProductBulkRequest is the request body of ProductsService.BulkCreate and
// ProductsService.BulkUpdate.
type ProductBulkRequest struct {
	// Products are the products to create or update.
	Products []*ProductParams `json:"products"`
	// CheckDuplicateLinkNumbers enables duplicate checking of link numbers
	// (0: skip the check, 1: check).
	CheckDuplicateLinkNumbers *ecforce.BoolInt `json:"check_duplicate_link_numbers,omitempty"`
}

// ProductParams describes one product of a bulk create or update request. It
// is the union of the fields of both operations: Name, Number and
// VariantAttributes apply to creates (Name, Number and VariantAttributes.SKU
// are required there); ID, State, Description, DescriptionMobile,
// UpsellProductID, OptionTypeIDs, ExclusionFilterIDs, ExportFeed and
// MasterAttributes apply to updates (ID is required there).
type ProductParams struct {
	ID                                    *int64                              `json:"id,omitempty"`
	Name                                  *string                             `json:"name,omitempty"`
	Number                                *string                             `json:"number,omitempty"`
	State                                 *string                             `json:"state,omitempty"`
	LabelIDs                              []int64                             `json:"label_ids,omitempty"`
	IsRecurring                           *ecforce.BoolInt                    `json:"is_recurring,omitempty"`
	TaxID                                 *int64                              `json:"tax_id,omitempty"`
	DeliveryFeeTemplateID                 *int64                              `json:"delivery_fee_template_id,omitempty"`
	PaymentMethodFeeTemplateID            *int64                              `json:"payment_method_fee_template_id,omitempty"`
	AgeCheckRequired                      *ecforce.BoolInt                    `json:"age_check_required,omitempty"`
	MinAgeLimit                           *int                                `json:"min_age_limit,omitempty"`
	CustomerRankThresholdName             *string                             `json:"customer_rank_threshold_name,omitempty"`
	Caution                               *string                             `json:"caution,omitempty"`
	Caution02                             *string                             `json:"caution02,omitempty"`
	Description                           *string                             `json:"description,omitempty"`
	DescriptionMobile                     *string                             `json:"description_mobile,omitempty"`
	SubDescription                        *string                             `json:"sub_description,omitempty"`
	SubDescriptionMobile                  *string                             `json:"sub_description_mobile,omitempty"`
	MetaDescription                       *string                             `json:"meta_description,omitempty"`
	MetaKeywords                          *string                             `json:"meta_keywords,omitempty"`
	MakerID                               *int64                              `json:"maker_id,omitempty"`
	UpsellProductID                       *int64                              `json:"upsell_product_id,omitempty"`
	ProductCategoryIDs                    []int64                             `json:"product_category_ids,omitempty"`
	OptionTypeIDs                         []int64                             `json:"option_type_ids,omitempty"`
	ExclusionFilterIDs                    []int64                             `json:"exclusion_filter_ids,omitempty"`
	RecurringBlockTimes                   *int                                `json:"recurring_block_times,omitempty"`
	ExportFeed                            *ecforce.BoolInt                    `json:"export_feed,omitempty"`
	PaymentSchedule                       *string                             `json:"payment_schedule,omitempty"`
	ScheduledToBeDeliveredEveryXMonth     *int                                `json:"scheduled_to_be_delivered_every_x_month,omitempty"`
	ScheduledToBeDeliveredOnXthDay        *int                                `json:"scheduled_to_be_delivered_on_xth_day,omitempty"`
	ScheduledToBeDeliveredEveryXDay       *int                                `json:"scheduled_to_be_delivered_every_x_day,omitempty"`
	ScheduledToBeDeliveredOnXthDayOfWeek  *int                                `json:"scheduled_to_be_delivered_on_xth_day_of_week,omitempty"`
	ScheduledToBeDeliveredEveryXDayOfWeek *int                                `json:"scheduled_to_be_delivered_every_x_day_of_week,omitempty"`
	DeliveryDateIDs                       []int64                             `json:"delivery_date_ids,omitempty"`
	DeliveryIntervalIDs                   []int64                             `json:"delivery_interval_ids,omitempty"`
	LinkNumber                            *string                             `json:"link_number,omitempty"`
	Option01                              *string                             `json:"option01,omitempty"`
	Option02                              *string                             `json:"option02,omitempty"`
	Option03                              *string                             `json:"option03,omitempty"`
	Option04                              *string                             `json:"option04,omitempty"`
	Option05                              *string                             `json:"option05,omitempty"`
	Option06                              *string                             `json:"option06,omitempty"`
	Option07                              *string                             `json:"option07,omitempty"`
	Option08                              *string                             `json:"option08,omitempty"`
	Option09                              *string                             `json:"option09,omitempty"`
	Option10                              *string                             `json:"option10,omitempty"`
	VariantAttributes                     *VariantParams                      `json:"variant_attributes,omitempty"`
	MasterAttributes                      *VariantParams                      `json:"master_attributes,omitempty"`
	BundledItemsAttributes                *BundledItemParams                  `json:"bundled_items_attributes,omitempty"`
	PaymentMethodDiscountsAttributes      *ProductPaymentMethodDiscountParams `json:"payment_method_discounts_attributes,omitempty"`
	ImageAttributes                       *ProductImageParams                 `json:"image_attributes,omitempty"`
}

// VariantParams describes the variant of a product in a bulk create request
// (ProductParams.VariantAttributes, where SKU is required) or the master
// variant in a bulk update request (ProductParams.MasterAttributes, which
// supports the price/quantity/flag subset only).
type VariantParams struct {
	SKU                       *string          `json:"sku,omitempty"`
	MasterSKU                 *string          `json:"master_sku,omitempty"`
	State                     *string          `json:"state,omitempty"`
	ForSale                   *ecforce.BoolInt `json:"for_sale,omitempty"`
	ListPrice                 *int             `json:"list_price,omitempty"`
	SalesPrice                *int             `json:"sales_price,omitempty"`
	Description               *string          `json:"description,omitempty"`
	DescriptionMobile         *string          `json:"description_mobile,omitempty"`
	OptionType01              *string          `json:"option_type01,omitempty"`
	OptionValue01             *string          `json:"option_value01,omitempty"`
	OptionType02              *string          `json:"option_type02,omitempty"`
	OptionValue02             *string          `json:"option_value02,omitempty"`
	OptionType03              *string          `json:"option_type03,omitempty"`
	OptionValue03             *string          `json:"option_value03,omitempty"`
	LimitQuantity             *ecforce.BoolInt `json:"limit_quantity,omitempty"`
	MinQuantity               *int             `json:"min_quantity,omitempty"`
	MaxQuantity               *int             `json:"max_quantity,omitempty"`
	CustomerLimitQuantity     *int             `json:"customer_limit_quantity,omitempty"`
	ChangeableVariantGroupIDs []int64          `json:"changeable_variant_group_ids,omitempty"`
	Cost                      *int             `json:"cost,omitempty"`
	Volume                    *int             `json:"volume,omitempty"`
	IsNew                     *ecforce.BoolInt `json:"is_new,omitempty"`
	IsSale                    *ecforce.BoolInt `json:"is_sale,omitempty"`
	VisibleInCart             *ecforce.BoolInt `json:"visible_in_cart,omitempty"`
	VisibleInEmail            *ecforce.BoolInt `json:"visible_in_email,omitempty"`
	VisibleInReport           *ecforce.BoolInt `json:"visible_in_report,omitempty"`
}

// BundledItemParams sets the bundled items of a product in a bulk create
// request by their SKU codes.
type BundledItemParams struct {
	VariantSKU []string `json:"variant_sku,omitempty"`
}

// ProductPaymentMethodDiscountParams sets a payment method discount of a
// product in a bulk create request. Set either DiscountAmount or
// DiscountRatio, not both.
type ProductPaymentMethodDiscountParams struct {
	PaymentMethodID *int64           `json:"payment_method_id,omitempty"`
	DiscountAmount  *int             `json:"discount_amount,omitempty"`
	DiscountRatio   *int             `json:"discount_ratio,omitempty"`
	FirstTimeOnly   *ecforce.BoolInt `json:"first_time_only,omitempty"`
}

// ProductImageParams sets the product and variant images of a product in a
// bulk create request by URL.
type ProductImageParams struct {
	ProductURL   *string `json:"product_url,omitempty"`
	VariantURL01 *string `json:"variant_url01,omitempty"`
	VariantURL02 *string `json:"variant_url02,omitempty"`
	VariantURL03 *string `json:"variant_url03,omitempty"`
	VariantURL04 *string `json:"variant_url04,omitempty"`
	VariantURL05 *string `json:"variant_url05,omitempty"`
}

// List searches products. Set opts.Lighter to toggle the lightweight response
// variant, which omits some Product attributes.
//
// Supported q attributes: id, number, name, link_number, created_at,
// updated_at, with_deleted.
// Supported sort attributes: id, created_at, updated_at.
// Supported include values: bundled_items, variants, variants.thumbnails,
// thumbnail, product_categories, labels (lighter only).
//
// ecforce API docs: GET /api/v2/admin/products
func (s *ProductsService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[Product], *ecforce.Response, error) {
	return ecforce.DoResourceList[Product](ctx, s.client, http.MethodGet, "admin/products.json", opts.Values(), nil)
}

// Get fetches a single product. Set opts.Lighter to toggle the lightweight
// response variant, which omits some Product attributes.
//
// Supported include values: bundled_items, variants, variants.thumbnails,
// thumbnail, product_categories, labels (lighter only).
//
// ecforce API docs: GET /api/v2/admin/products/:product_id
func (s *ProductsService) Get(ctx context.Context, productID int64, opts *ecforce.GetOptions) (*ecforce.Resource[Product], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/products/%d.json", productID)
	return ecforce.DoResource[Product](ctx, s.client, http.MethodGet, path, opts.Values(), nil)
}

// BulkCreate registers products in bulk as an asynchronous job. Each entry
// requires Name, Number and VariantAttributes (with SKU); SKU-only
// registration is also possible.
//
// ecforce API docs: POST /api/v2/admin/products/bulk_create
func (s *ProductsService) BulkCreate(ctx context.Context, req *ProductBulkRequest) (*ecforce.JobResult, *ecforce.Response, error) {
	result := new(ecforce.JobResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPost, "admin/products/bulk_create.json", nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkUpdate updates products in bulk synchronously. Each entry requires ID;
// master variant fields go in MasterAttributes.
//
// ecforce API docs: PUT /api/v2/admin/products/bulk_update
func (s *ProductsService) BulkUpdate(ctx context.Context, req *ProductBulkRequest) (*ecforce.BulkResult, *ecforce.Response, error) {
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, "admin/products/bulk_update.json", nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkDestroy deletes products in bulk synchronously. Deleted products cannot
// be restored; nonexistent or already deleted IDs are skipped.
//
// ecforce API docs: DELETE /api/v2/admin/products/bulk_destroy
func (s *ProductsService) BulkDestroy(ctx context.Context, productIDs []int64) (*ecforce.BulkResult, *ecforce.Response, error) {
	body := map[string][]int64{"product_ids": productIDs}
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodDelete, "admin/products/bulk_destroy.json", nil, body, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}
