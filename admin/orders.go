package admin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/zero-color/ecforce-go"
)

// Order is an order (受注). It carries the union of the attributes of the
// normal and lighter response variants; the lighter variant (toggled via
// options) populates only a subset of the fields.
type Order struct {
	ID                               int64                   `json:"id"`
	Number                           string                  `json:"number"`
	Labels                           string                  `json:"labels"`
	CustomerID                       int64                   `json:"customer_id"`
	CustomerNumber                   string                  `json:"customer_number"`
	CustomerLabels                   string                  `json:"customer_labels"`
	SubsOrderID                      int64                   `json:"subs_order_id"`
	SubsOrderNumber                  string                  `json:"subs_order_number"`
	SubsOrderLabels                  string                  `json:"subs_order_labels"`
	ProductLabels                    string                  `json:"product_labels"`
	BundledItems                     string                  `json:"bundled_items"`
	ProductNameWithTax               string                  `json:"product_name_with_tax"`
	IP                               string                  `json:"ip"`
	DeviceVariant                    string                  `json:"device_variant"`
	ShippingCarrierCode              string                  `json:"shipping_carrier_code"`
	HumanOfferName                   string                  `json:"human_offer_name"`
	URLID                            int64                   `json:"url_id"`
	URL                              string                  `json:"url"`
	URLGroup                         string                  `json:"url_group"`
	Advertiser                       string                  `json:"advertiser"`
	PaymentAccessID                  string                  `json:"payment_access_id"`
	PaymentAccessPass                string                  `json:"payment_access_pass"`
	PaymentNumber                    string                  `json:"payment_number"`
	Wrapping                         string                  `json:"wrapping"`
	Remark                           string                  `json:"remark"`
	Store                            string                  `json:"store"`
	State                            string                  `json:"state"`
	HumanState                       string                  `json:"human_state"`
	PaymentState                     string                  `json:"payment_state"`
	PaymentHumanState                string                  `json:"payment_human_state"`
	Email                            string                  `json:"email"`
	Subtotal                         int                     `json:"subtotal"`
	Subtotal8                        int                     `json:"subtotal8"`
	Subtotal10                       int                     `json:"subtotal10"`
	Discount                         int                     `json:"discount"`
	Discount8                        int                     `json:"discount8"`
	Discount10                       int                     `json:"discount10"`
	Point                            int                     `json:"point"`
	Point8                           int                     `json:"point8"`
	Point10                          int                     `json:"point10"`
	DiscountWithPoint                int                     `json:"discount_with_point"`
	DiscountWithPoint8               int                     `json:"discount_with_point8"`
	DiscountWithPoint10              int                     `json:"discount_with_point10"`
	MiscFee                          int                     `json:"misc_fee"`
	DelivFee                         int                     `json:"deliv_fee"`
	Charge                           int                     `json:"charge"`
	Adjustment                       int                     `json:"adjustment"`
	Tax                              int                     `json:"tax"`
	Tax8                             int                     `json:"tax8"`
	Tax10                            int                     `json:"tax10"`
	Total                            int                     `json:"total"`
	Total8                           int                     `json:"total8"`
	Total10                          int                     `json:"total10"`
	PaymentTotal                     int                     `json:"payment_total"`
	RedeemPoint                      int                     `json:"redeem_point"`
	RewardPoint                      int                     `json:"reward_point"`
	GrantPlanPoint                   int                     `json:"grant_plan_point"`
	TotalRecurringSalesPriceDiscount int                     `json:"total_recurring_sales_price_discount"`
	ShippingCarrierID                int64                   `json:"shipping_carrier_id"`
	ShippingCarrierName              string                  `json:"shipping_carrier_name"`
	StockLocationID                  int64                   `json:"stock_location_id"`
	StockLocationName                string                  `json:"stock_location_name"`
	TrackingURL                      string                  `json:"tracking_url"`
	ShippingSlip                     string                  `json:"shipping_slip"`
	Times                            int                     `json:"times"`
	ShippingHistoriesCount           int                     `json:"shipping_histories_count"`
	PickedList                       bool                    `json:"picked_list"`
	PickedAt                         *ecforce.Time           `json:"picked_at"`
	PaymentLastErrorMessage          string                  `json:"payment_last_error_message"`
	CustomerContacts                 string                  `json:"customer_contacts"`
	Coupons                          string                  `json:"coupons"`
	PaymentMethodName                string                  `json:"payment_method_name"`
	PaymentTimes                     int                     `json:"payment_times"`
	Kind                             string                  `json:"kind"`
	Nth                              int                     `json:"nth"`
	ScheduledToBeShippedAt           *ecforce.Time           `json:"scheduled_to_be_shipped_at"`
	ScheduledToBeDeliveredAt         *ecforce.Time           `json:"scheduled_to_be_delivered_at"`
	ScheduledDeliveryTime            string                  `json:"scheduled_delivery_time"`
	ScheduledDeliveryTimeID          int64                   `json:"scheduled_delivery_time_id"`
	ScheduledDeliveryTimeCode        string                  `json:"scheduled_delivery_time_code"`
	PreviousScheduledToBeShippedAt   *ecforce.Time           `json:"previous_scheduled_to_be_shipped_at"`
	PreviousScheduledToBeDeliveredAt *ecforce.Time           `json:"previous_scheduled_to_be_delivered_at"`
	FirstRemindedAt                  *ecforce.Time           `json:"first_reminded_at"`
	ReminderedAt                     *ecforce.Time           `json:"remindered_at"`
	ReminderCount                    int                     `json:"reminder_count"`
	DueAt                            *ecforce.Time           `json:"due_at"`
	OrderFreeColumns                 []*OrderFreeColumnEntry `json:"order_free_columns"`
	CreatedAt                        *ecforce.Time           `json:"created_at"`
	UpdatedAt                        *ecforce.Time           `json:"updated_at"`
	CompletedAt                      *ecforce.Time           `json:"completed_at"`
	DeliveredAt                      *ecforce.Time           `json:"delivered_at"`
	ShippedAt                        *ecforce.Time           `json:"shipped_at"`
	PaymentAuthedAt                  *ecforce.Time           `json:"payment_authed_at"`
	PaymentCompletedAt               *ecforce.Time           `json:"payment_completed_at"`
	PaymentPaidAt                    *ecforce.Time           `json:"payment_paid_at"`
	PaymentVoidedAt                  *ecforce.Time           `json:"payment_voided_at"`
	LastOrder                        bool                    `json:"last_order"`
	Memo01                           string                  `json:"memo01"`
	Memo02                           string                  `json:"memo02"`
	SubsOrderMemo01                  string                  `json:"subs_order_memo01"`
	SubsOrderMemo02                  string                  `json:"subs_order_memo02"`
	InviteCode                       string                  `json:"invite_code"`
	PayWithEpos                      bool                    `json:"pay_with_epos"`
	MultipleShipping                 bool                    `json:"multiple_shipping"`
	TBC                              bool                    `json:"tbc"`
	TBCOrderReasons                  string                  `json:"tbc_order_reasons"`
	AtScoreResult                    string                  `json:"at_score_result"`
	OPluxResult                      string                  `json:"o_plux_result"`
	OPluxDescription                 string                  `json:"o_plux_description"`
	SpiderAFResult                   string                  `json:"spider_af_result"`
	SpiderAFDescription              string                  `json:"spider_af_description"`
	TenantID                         int64                   `json:"tenant_id"`
	AssignedAdminID                  int64                   `json:"assigned_admin_id"`
	RegistrantAdminID                int64                   `json:"registrant_admin_id"`
	LinkedMemo                       string                  `json:"linked_memo"`
	LinkNumber                       string                  `json:"link_number"`
	PickupLocation1                  string                  `json:"pickup_location1"`
	PickupLocation2                  string                  `json:"pickup_location2"`
	Doorbell                         bool                    `json:"doorbell"`
	CVRoute                          string                  `json:"cv_route"`
	Campaigns                        []*OrderCampaign        `json:"campaigns"`
}

// OrderFreeColumnEntry is one answered order free column (受注自由項目) inside
// Order.OrderFreeColumns: the master record ID, the answers, and the label.
type OrderFreeColumnEntry struct {
	ID     int64    `json:"id"`
	Values []string `json:"values"`
	Name   string   `json:"name"`
}

// OrderCampaign is a campaign applied to an order. It is populated only for
// normal (non-lighter) responses.
type OrderCampaign struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

// OrderParams sets the attributes of an order on update requests. It covers
// both single updates (OrdersService.Update, where ID and LinkNumber are
// ignored) and bulk updates (OrdersService.BulkUpdate, where ID is required).
type OrderParams struct {
	ID                        *int64                        `json:"id,omitempty"`
	StoreID                   *int64                        `json:"store_id,omitempty"`
	State                     *string                       `json:"state,omitempty"`
	ShippingCarrierID         *int64                        `json:"shipping_carrier_id,omitempty"`
	Remark                    *string                       `json:"remark,omitempty"`
	ScheduledToBeShippedAt    *ecforce.Time                 `json:"scheduled_to_be_shipped_at,omitempty"`
	ScheduledToBeDeliveredAt  *ecforce.Time                 `json:"scheduled_to_be_delivered_at,omitempty"`
	ScheduledDeliveryTimeID   *int64                        `json:"scheduled_delivery_time_id,omitempty"`
	DeliveredAt               *ecforce.Time                 `json:"delivered_at,omitempty"`
	LinkNumber                *string                       `json:"link_number,omitempty"`
	Memo01                    *string                       `json:"memo01,omitempty"`
	Memo02                    *string                       `json:"memo02,omitempty"`
	AssignedAdminID           *int64                        `json:"assigned_admin_id,omitempty"`
	RegistrantAdminID         *int64                        `json:"registrant_admin_id,omitempty"`
	LinkedMemo                *string                       `json:"linked_memo,omitempty"`
	Adjustment                *int                          `json:"adjustment,omitempty"`
	BillingAddressAttributes  *AddressParams                `json:"billing_address_attributes,omitempty"`
	ShippingAddressAttributes *AddressParams                `json:"shipping_address_attributes,omitempty"`
	PickupLocationsAttributes *PickupLocationParams         `json:"pickup_locations_attributes,omitempty"`
	OrderItemsAttributes      []*OrderItemParams            `json:"order_items_attributes,omitempty"`
	PaymentAttributes         *OrderPaymentParams           `json:"payment_attributes,omitempty"`
	URL                       *string                       `json:"url,omitempty"`
	OrderFreeColumns          []*OrderFreeColumnEntryParams `json:"order_free_columns,omitempty"`
}

// OrderFreeColumnEntryParams sets the answers of one order free column
// (受注自由項目) on update requests. Columns not listed keep their current
// values; pass an empty non-nil Values slice ([]string{}) to clear a
// column's answers.
type OrderFreeColumnEntryParams struct {
	ID     *int64   `json:"id,omitempty"`
	Values []string `json:"values"`
}

// PickupLocationParams sets the pickup location preferences of an order
// (pickup_locations_attributes).
type PickupLocationParams struct {
	PickupLocationID1 *int64           `json:"pickup_location_id_1,omitempty"`
	PickupLocationID2 *int64           `json:"pickup_location_id_2,omitempty"`
	Doorbell          *ecforce.BoolInt `json:"doorbell,omitempty"`
}

// OrderItemParams sets one order item on update and bulk create/update
// requests. ID is required when updating or deleting an existing item;
// bulk creates take VariantID and Quantity only.
type OrderItemParams struct {
	ID                    *int64           `json:"id,omitempty"`
	VariantID             *int64           `json:"variant_id,omitempty"`
	Quantity              *int             `json:"quantity,omitempty"`
	PointExchangeQuantity *int             `json:"point_exchange_quantity,omitempty"`
	Delete                *ecforce.BoolInt `json:"delete,omitempty"`
}

// OrderPaymentParams sets the payment of an order (payment_attributes).
// SourceID and SourceType ("EcForce::CreditCard") are required when switching
// to a credit card payment method and must be omitted otherwise.
type OrderPaymentParams struct {
	PaymentMethodID *int64  `json:"payment_method_id,omitempty"`
	SourceID        *int64  `json:"source_id,omitempty"`
	SourceType      *string `json:"source_type,omitempty"`
	PaymentTimes    *int    `json:"payment_times,omitempty"`
}

// OrderUpdateRequest is the body of OrdersService.Update.
type OrderUpdateRequest struct {
	Order *OrderParams `json:"order"`
	// Preview, when true, returns the updated result without saving the order.
	Preview *ecforce.BoolInt `json:"preview,omitempty"`
	// UpdateSubsOrderMemo, when true, propagates the order memos to the
	// subscription order memos.
	UpdateSubsOrderMemo *ecforce.BoolInt `json:"update_subs_order_memo,omitempty"`
	// Recalculate, when true, recalculates the order amounts. Recalculation
	// always happens when Order.Adjustment is set, regardless of this flag.
	Recalculate *ecforce.BoolInt `json:"recalculate,omitempty"`
}

// OrderBulkUpdateRequest is the body of OrdersService.BulkUpdate. Each order
// must carry its ID.
type OrderBulkUpdateRequest struct {
	Orders []*OrderParams `json:"orders"`
	// CheckDuplicateLinkNumbers, when true, rejects duplicate link_number
	// values.
	CheckDuplicateLinkNumbers *ecforce.BoolInt `json:"check_duplicate_link_numbers,omitempty"`
	// UpdateSubsOrderMemo, when true, propagates the order memos to the
	// subscription order memos.
	UpdateSubsOrderMemo *ecforce.BoolInt `json:"update_subs_order_memo,omitempty"`
}

// OrderShippingParams reports the shipping result of one order for
// OrdersService.UpdateShipping. ShippingAddressID is required for orders
// with multiple shipping addresses.
type OrderShippingParams struct {
	ID                *int64        `json:"id,omitempty"`
	ShippingSlip      *string       `json:"shipping_slip,omitempty"`
	ShippingCarrierID *int64        `json:"shipping_carrier_id,omitempty"`
	ShippedAt         *ecforce.Time `json:"shipped_at,omitempty"`
	ShippingAddressID *int64        `json:"shipping_address_id,omitempty"`
}

// OrderShippingUpdateRequest is the body of OrdersService.UpdateShipping.
type OrderShippingUpdateRequest struct {
	Orders []*OrderShippingParams `json:"orders"`
	// WithAsync, when true, processes the request asynchronously (recommended)
	// and the result carries only a JobID.
	WithAsync *ecforce.BoolInt `json:"with_async,omitempty"`
	// WithSale, when true, also performs the sales (capture) processing.
	WithSale *ecforce.BoolInt `json:"with_sale,omitempty"`
}

// OrderShippingID identifies one processed record of a shipping result
// update: the order ID, or "<order id>-<shipping address id>" for orders
// with multiple shipping addresses. It decodes from either a JSON number or
// string.
type OrderShippingID string

// UnmarshalJSON implements json.Unmarshaler.
func (i *OrderShippingID) UnmarshalJSON(data []byte) error {
	s := string(data)
	if unquoted, err := strconv.Unquote(s); err == nil {
		s = unquoted
	}
	*i = OrderShippingID(strings.TrimSpace(s))
	return nil
}

// OrderShippingError describes why one record of a shipping result update
// failed.
type OrderShippingError struct {
	ID      OrderShippingID `json:"id"`
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Errors  []ecforce.Error `json:"errors"`
}

// OrderShippingResult is the outcome of OrdersService.UpdateShipping. For
// synchronous requests Success/Failure/Errors identify the processed records;
// for asynchronous requests (WithAsync) only JobID is set.
type OrderShippingResult struct {
	Success []OrderShippingID     `json:"success"`
	Failure []OrderShippingID     `json:"failure"`
	Errors  []*OrderShippingError `json:"errors"`
	JobID   string                `json:"job_id"`
}

// OrderPaymentStatusUpdateRequest is the body of
// OrdersService.UpdatePaymentStatus. Method is one of "void" (cancel) or
// "reauth" (re-authorize).
type OrderPaymentStatusUpdateRequest struct {
	Method *string `json:"method,omitempty"`
	// DecrementSubsOrderTimes, when true, decrements the subscription order's
	// recurring count by one. Effective only when Method is "void".
	DecrementSubsOrderTimes *ecforce.BoolInt `json:"decrement_subs_order_times,omitempty"`
	// RecalculateSubsOrder, when true, recalculates the subscription order.
	// Effective only with Method "void" and DecrementSubsOrderTimes set.
	RecalculateSubsOrder *ecforce.BoolInt `json:"recalculate_subs_order,omitempty"`
	// NeedToSendEmail, when true, sends the cancellation confirmation email.
	// Effective only when Method is "void".
	NeedToSendEmail *ecforce.BoolInt `json:"need_to_send_email,omitempty"`
}

// OrderPaymentStatusBulkUpdateRequest is the body of
// OrdersService.BulkUpdatePaymentStatus. Method is one of "void" (cancel),
// "sales" (capture), "reauth" (re-authorize) or "force_void" (force cancel
// without notifying the payment provider).
type OrderPaymentStatusBulkUpdateRequest struct {
	Method   *string `json:"method,omitempty"`
	OrderIDs []int64 `json:"order_ids,omitempty"`
	// DecrementSubsOrderTimes, when true, decrements the subscription order's
	// recurring count by one. Effective only when Method is "void" or
	// "force_void".
	DecrementSubsOrderTimes *ecforce.BoolInt `json:"decrement_subs_order_times,omitempty"`
	// RecalculateSubsOrder, when true, recalculates the subscription order.
	// Effective only with DecrementSubsOrderTimes set.
	RecalculateSubsOrder *ecforce.BoolInt `json:"recalculate_subs_order,omitempty"`
}

// OrderItemsBulkRequest is the body of OrdersService.BulkCreateOrderItems and
// OrdersService.BulkUpdateOrderItems.
type OrderItemsBulkRequest struct {
	OrderItems []*OrderItemParams `json:"order_items"`
	// ReferSettings, when true, resolves sales prices and bundled items from
	// the subscription detail settings (falling back to the SKU sales price
	// and the product's bundled items); when false, the SKU sales price is
	// used and bundled items are ignored.
	ReferSettings *ecforce.BoolInt `json:"refer_settings,omitempty"`
}

// Get fetches a single order.
//
// Set opts.Lighter to toggle the lightweight response variant, which
// populates only a subset of the Order fields.
// Supported include values: customer, customer.billing_address,
// customer.shipping_addresses, customer.notes, subs_order, billing_address,
// shipping_address, shipping_addresses, and more (the apiDoc truncates the
// list).
//
// ecforce API docs: GET /api/v2/admin/orders/:order_id
func (s *OrdersService) Get(ctx context.Context, orderID int64, opts *ecforce.GetOptions) (*ecforce.Resource[Order], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/orders/%d.json", orderID)
	return ecforce.DoResource[Order](ctx, s.client, http.MethodGet, path, opts.Values(), nil)
}

// List searches orders.
//
// Set opts.Lighter to toggle the lightweight response variant, which
// populates only a subset of the Order fields.
// Supported q attributes: id, number, customer_id, subs_order_id, url_id,
// state, payment_state, email, total, payment_total, shipping_carrier_id,
// stock_location_id, times, scheduled_to_be_shipped_at, and more (the apiDoc
// truncates the list); q[with_deleted]=1 includes deleted orders.
// Supported sort attributes: id, created_at, updated_at.
// Supported include values: customer, customer.billing_address,
// customer.shipping_addresses, customer.notes, subs_order, billing_address,
// shipping_address, shipping_addresses, and more (the apiDoc truncates the
// list).
//
// ecforce API docs: GET /api/v2/admin/orders
func (s *OrdersService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[Order], *ecforce.Response, error) {
	return ecforce.DoResourceList[Order](ctx, s.client, http.MethodGet, "admin/orders.json", opts.Values(), nil)
}

// Update updates an order synchronously and returns the updated order.
//
// ecforce API docs: PUT /api/v2/admin/orders/:order_id
func (s *OrdersService) Update(ctx context.Context, orderID int64, req *OrderUpdateRequest) (*ecforce.Resource[Order], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/orders/%d.json", orderID)
	return ecforce.DoResource[Order](ctx, s.client, http.MethodPut, path, nil, req)
}

// BulkUpdate updates multiple orders synchronously.
//
// ecforce API docs: PUT /api/v2/admin/orders/bulk_update
func (s *OrdersService) BulkUpdate(ctx context.Context, req *OrderBulkUpdateRequest) (*ecforce.BulkResult, *ecforce.Response, error) {
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, "admin/orders/bulk_update.json", nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkDestroy deletes multiple orders synchronously. Deleted orders cannot be
// restored.
//
// ecforce API docs: DELETE /api/v2/admin/orders/bulk_destroy
func (s *OrdersService) BulkDestroy(ctx context.Context, orderIDs []int64) (*ecforce.BulkResult, *ecforce.Response, error) {
	body := struct {
		OrderIDs []int64 `json:"order_ids"`
	}{OrderIDs: orderIDs}
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodDelete, "admin/orders/bulk_destroy.json", nil, body, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// SendMail sends an order email using the given mail template, whose category
// must be order management (受注管理), and returns the order.
//
// ecforce API docs: POST /api/v2/admin/orders/:order_id/emails
func (s *OrdersService) SendMail(ctx context.Context, orderID, emailTemplateID int64) (*ecforce.Resource[Order], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/orders/%d/emails.json", orderID)
	body := struct {
		EmailTemplateID int64 `json:"email_template_id"`
	}{EmailTemplateID: emailTemplateID}
	return ecforce.DoResource[Order](ctx, s.client, http.MethodPost, path, nil, body)
}

// UpdateShipping reflects shipping results (発送結果戻し) onto orders.
//
// ecforce API docs: PUT /api/v2/admin/orders/shipping
func (s *OrdersService) UpdateShipping(ctx context.Context, req *OrderShippingUpdateRequest) (*OrderShippingResult, *ecforce.Response, error) {
	result := new(OrderShippingResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, "admin/orders/shipping.json", nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// UpdatePaymentStatus updates the payment status of an order synchronously
// and returns the order.
//
// ecforce API docs: POST /api/v2/admin/orders/:order_id/payment_status
func (s *OrdersService) UpdatePaymentStatus(ctx context.Context, orderID int64, req *OrderPaymentStatusUpdateRequest) (*ecforce.Resource[Order], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/orders/%d/payment_status.json", orderID)
	return ecforce.DoResource[Order](ctx, s.client, http.MethodPost, path, nil, req)
}

// BulkUpdatePaymentStatus updates the payment status of multiple orders
// asynchronously and returns the accepted job.
//
// ecforce API docs: POST /api/v2/admin/orders/payment_status/bulk_update
func (s *OrdersService) BulkUpdatePaymentStatus(ctx context.Context, req *OrderPaymentStatusBulkUpdateRequest) (*ecforce.JobResult, *ecforce.Response, error) {
	result := new(ecforce.JobResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPost, "admin/orders/payment_status/bulk_update.json", nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkCreateOrderItems adds order items to an order synchronously. Success
// and Failure hold 1-based positions in the request array.
//
// ecforce API docs: POST /api/v2/admin/orders/:order_id/order_items/bulk_create
func (s *OrdersService) BulkCreateOrderItems(ctx context.Context, orderID int64, req *OrderItemsBulkRequest) (*ecforce.BulkResult, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/orders/%d/order_items/bulk_create.json", orderID)
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPost, path, nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkUpdateOrderItems updates or deletes order items of an order
// synchronously. Each item must carry its ID.
//
// ecforce API docs: PUT /api/v2/admin/orders/:order_id/order_items/bulk_update
func (s *OrdersService) BulkUpdateOrderItems(ctx context.Context, orderID int64, req *OrderItemsBulkRequest) (*ecforce.BulkResult, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/orders/%d/order_items/bulk_update.json", orderID)
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, path, nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// ListFreeColumnValues fetches the custom free column values of an order,
// grouped by category.
//
// ecforce API docs: GET /api/v2/admin/orders/:order_id/free_column_values
func (s *OrdersService) ListFreeColumnValues(ctx context.Context, orderID int64) ([]*FreeColumnValueGroup, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/orders/%d/free_column_values.json", orderID)
	var groups []*FreeColumnValueGroup
	resp, err := ecforce.Do(ctx, s.client, http.MethodGet, path, nil, nil, &groups)
	if err != nil {
		return nil, resp, err
	}
	return groups, resp, nil
}

// BulkCreateFreeColumnValues registers custom free column values on an order.
// Success and Failure hold 1-based positions in the request array.
//
// ecforce API docs: POST /api/v2/admin/orders/:order_id/free_column_values/bulk_create
func (s *OrdersService) BulkCreateFreeColumnValues(ctx context.Context, orderID int64, freeColumns []*FreeColumnParams) (*ecforce.BulkResult, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/orders/%d/free_column_values/bulk_create.json", orderID)
	body := struct {
		FreeColumns []*FreeColumnParams `json:"free_columns"`
	}{FreeColumns: freeColumns}
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPost, path, nil, body, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkUpdateFreeColumnValues updates custom free column values of an order.
// Success and Failure hold 1-based positions in the request array.
//
// ecforce API docs: PUT /api/v2/admin/orders/:order_id/free_column_values/bulk_update
func (s *OrdersService) BulkUpdateFreeColumnValues(ctx context.Context, orderID int64, freeColumns []*FreeColumnParams) (*ecforce.BulkResult, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/orders/%d/free_column_values/bulk_update.json", orderID)
	body := struct {
		FreeColumns []*FreeColumnParams `json:"free_columns"`
	}{FreeColumns: freeColumns}
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, path, nil, body, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}
