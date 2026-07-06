package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// SubsOrder is a subscription order (定期受注). It merges the attributes of the
// normal and lighter response variants; fields absent from the variant in use
// are left at their zero values.
type SubsOrder struct {
	ID                                         int64                                `json:"id"`
	Number                                     string                               `json:"number"`
	CustomerNumber                             string                               `json:"customer_number"`
	Wrapping                                   string                               `json:"wrapping"`
	Times                                      int                                  `json:"times"`
	OrdersCount                                int                                  `json:"orders_count"`
	Store                                      string                               `json:"store"`
	State                                      string                               `json:"state"`
	HumanState                                 string                               `json:"human_state"`
	Subtotal                                   int                                  `json:"subtotal"`
	TBC                                        bool                                 `json:"tbc"`
	TBCSubsOrderReasons                        string                               `json:"tbc_subs_order_reasons"`
	RemainingNumberOfOrders                    int                                  `json:"remaining_number_of_orders"`
	AvailablePaymentSchedules                  []*SubsOrderAvailablePaymentSchedule `json:"available_payment_schedules"`
	PaymentSchedule                            string                               `json:"payment_schedule"`
	PaymentScheduleLocked                      ecforce.BoolInt                      `json:"payment_schedule_locked"`
	ScheduledToBeDeliveredEveryXMonth          int                                  `json:"scheduled_to_be_delivered_every_x_month"`
	ScheduledToBeDeliveredOnXthDay             int                                  `json:"scheduled_to_be_delivered_on_xth_day"`
	ScheduledToBeDeliveredEveryXDay            int                                  `json:"scheduled_to_be_delivered_every_x_day"`
	ScheduledToBeDeliveredOnXthDayOfWeek       int                                  `json:"scheduled_to_be_delivered_on_xth_day_of_week"`
	HumanScheduledToBeDeliveredEveryXDayOfWeek string                               `json:"human_scheduled_to_be_delivered_every_x_day_of_week"`
	ScheduledToBeShippedAt                     *ecforce.Time                        `json:"scheduled_to_be_shipped_at"`
	ScheduledToBeDeliveredAt                   *ecforce.Time                        `json:"scheduled_to_be_delivered_at"`
	ScheduledDeliveryTime                      string                               `json:"scheduled_delivery_time"`
	ScheduledDeliveryTimeID                    int64                                `json:"scheduled_delivery_time_id"`
	PreviousScheduledToBeShippedAt             *ecforce.Time                        `json:"previous_scheduled_to_be_shipped_at"`
	PreviousScheduledToBeDeliveredAt           *ecforce.Time                        `json:"previous_scheduled_to_be_delivered_at"`
	SuspendReasons                             string                               `json:"suspend_reasons"`
	CancelReasons                              string                               `json:"cancel_reasons"`
	SuspendReasonIDs                           []int64                              `json:"suspend_reason_ids"`
	CancelReasonIDs                            []int64                              `json:"cancel_reason_ids"`
	ElapsedDaysFromSuspended                   int                                  `json:"elapsed_days_from_suspended"`
	ElapsedDaysFromCanceled                    int                                  `json:"elapsed_days_from_canceled"`
	PaymentMethodName                          string                               `json:"payment_method_name"`
	PaymentTimes                               int                                  `json:"payment_times"`
	RecurringBlockTimes                        int                                  `json:"recurring_block_times"`
	RemainingBlockTimes                        int                                  `json:"remaining_block_times"`
	OrdersPaidCount                            int                                  `json:"orders_paid_count"`
	Labels                                     string                               `json:"labels"`
	CustomerLabels                             string                               `json:"customer_labels"`
	ProductLabels                              string                               `json:"product_labels"`
	CustomerContacts                           string                               `json:"customer_contacts"`
	Memo01                                     string                               `json:"memo01"`
	Memo02                                     string                               `json:"memo02"`
	Adjustment                                 int                                  `json:"adjustment"`
	OnetimeAdjustment                          bool                                 `json:"onetime_adjustment"`
	DefaultShippingCarrierID                   int64                                `json:"default_shipping_carrier_id"`
	DefaultShippingCarrier                     string                               `json:"default_shipping_carrier"`
	DefaultShippingCarrierLocked               bool                                 `json:"default_shipping_carrier_locked"`
	PayWithEpos                                bool                                 `json:"pay_with_epos"`
	TenantID                                   int64                                `json:"tenant_id"`
	AssignedAdminID                            int64                                `json:"assigned_admin_id"`
	RegistrantAdminID                          int64                                `json:"registrant_admin_id"`
	LinkedMemo                                 string                               `json:"linked_memo"`
	SkipNextScheduledDelivery                  bool                                 `json:"skip_next_scheduled_delivery"`
	LinkNumber                                 string                               `json:"link_number"`
	PickupLocation1                            string                               `json:"pickup_location1"`
	PickupLocation2                            string                               `json:"pickup_location2"`
	Doorbell                                   bool                                 `json:"doorbell"`
	OrderFreeColumns                           []*SubsOrderFreeColumn               `json:"order_free_columns"`
	CreatedAt                                  *ecforce.Time                        `json:"created_at"`
	UpdatedAt                                  *ecforce.Time                        `json:"updated_at"`
	SuspendedAt                                *ecforce.Time                        `json:"suspended_at"`
	CanceledAt                                 *ecforce.Time                        `json:"canceled_at"`
	DeletedAt                                  *ecforce.Time                        `json:"deleted_at"`
}

// SubsOrderFreeColumn is an order free column (受注自由項目) answer attached to
// a subscription order.
type SubsOrderFreeColumn struct {
	ID     int64    `json:"id"`
	Values []string `json:"values"`
	Name   string   `json:"name"`
}

// SubsOrderAvailablePaymentSchedule describes one delivery cycle type
// selectable for a subscription order, with the choices allowed for each of
// its schedule attributes.
type SubsOrderAvailablePaymentSchedule struct {
	Value                                 string                     `json:"value"`
	Ja                                    string                     `json:"ja"`
	ScheduledToBeDeliveredEveryXMonth     []*SubsOrderScheduleChoice `json:"scheduled_to_be_delivered_every_x_month"`
	ScheduledToBeDeliveredOnXthDay        []*SubsOrderScheduleChoice `json:"scheduled_to_be_delivered_on_xth_day"`
	ScheduledToBeDeliveredEveryXDay       []*SubsOrderScheduleChoice `json:"scheduled_to_be_delivered_every_x_day"`
	ScheduledToBeDeliveredOnXthDayOfWeek  []*SubsOrderScheduleChoice `json:"scheduled_to_be_delivered_on_xth_day_of_week"`
	ScheduledToBeDeliveredEveryXDayOfWeek []*SubsOrderScheduleChoice `json:"scheduled_to_be_delivered_every_x_day_of_week"`
}

// SubsOrderScheduleChoice is one selectable value of a delivery schedule
// attribute, with its Japanese display label.
type SubsOrderScheduleChoice struct {
	Value int    `json:"value"`
	Ja    string `json:"ja"`
}

// UnmarshalJSON implements json.Unmarshaler. The API emits the label as
// either a string or a number depending on the attribute; both decode into Ja.
func (c *SubsOrderScheduleChoice) UnmarshalJSON(data []byte) error {
	var raw struct {
		Value int             `json:"value"`
		Ja    json.RawMessage `json:"ja"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	c.Value = raw.Value
	c.Ja = ""
	if len(raw.Ja) > 0 && string(raw.Ja) != "null" {
		if raw.Ja[0] == '"' {
			if err := json.Unmarshal(raw.Ja, &c.Ja); err != nil {
				return err
			}
		} else {
			c.Ja = string(raw.Ja)
		}
	}
	return nil
}

// CancelReason is a subscription order cancel reason (定期受注キャンセル理由).
type CancelReason struct {
	ID             int64         `json:"id"`
	Name           string        `json:"name"`
	ParentReasonID int64         `json:"parent_reason_id"`
	CreatedAt      *ecforce.Time `json:"created_at"`
	UpdatedAt      *ecforce.Time `json:"updated_at"`
}

// SuspendReason is a subscription order suspend reason (定期受注停止理由).
type SuspendReason struct {
	ID             int64         `json:"id"`
	Name           string        `json:"name"`
	ParentReasonID int64         `json:"parent_reason_id"`
	CreatedAt      *ecforce.Time `json:"created_at"`
	UpdatedAt      *ecforce.Time `json:"updated_at"`
}

// SubsOrderParams sets the attributes of a subscription order on update and
// bulk create/update requests. Fields that only one endpoint accepts (e.g. ID
// for bulk updates, CustomerID and OrderItemsAttributes for bulk creates) are
// ignored by the others.
//
// SuspendReasonIDs and CancelReasonIDs take reason IDs as decimal strings;
// pass []string{"null"} (or []string{""}) to remove all attached reasons.
type SubsOrderParams struct {
	// ID identifies the subscription order to update (bulk update only).
	ID *int64 `json:"id,omitempty"`
	// CustomerID is the owning customer (bulk create only).
	CustomerID                            *int64                          `json:"customer_id,omitempty"`
	LinkNumber                            *string                         `json:"link_number,omitempty"`
	State                                 *string                         `json:"state,omitempty"`
	Times                                 *int                            `json:"times,omitempty"`
	RecurringBlockTimes                   *int                            `json:"recurring_block_times,omitempty"`
	RemainingNumberOfOrders               *int                            `json:"remaining_number_of_orders,omitempty"`
	SuspendReasonIDs                      []string                        `json:"suspend_reason_ids,omitempty"`
	CancelReasonIDs                       []string                        `json:"cancel_reason_ids,omitempty"`
	PaymentSchedule                       *string                         `json:"payment_schedule,omitempty"`
	PaymentScheduleLocked                 *ecforce.BoolInt                `json:"payment_schedule_locked,omitempty"`
	ScheduledToBeDeliveredEveryXMonth     *int                            `json:"scheduled_to_be_delivered_every_x_month,omitempty"`
	ScheduledToBeDeliveredOnXthDay        *int                            `json:"scheduled_to_be_delivered_on_xth_day,omitempty"`
	ScheduledToBeDeliveredEveryXDay       *int                            `json:"scheduled_to_be_delivered_every_x_day,omitempty"`
	ScheduledToBeDeliveredOnXthDayOfWeek  *int                            `json:"scheduled_to_be_delivered_on_xth_day_of_week,omitempty"`
	ScheduledToBeDeliveredEveryXDayOfWeek *int                            `json:"scheduled_to_be_delivered_every_x_day_of_week,omitempty"`
	ScheduledToBeShippedAt                *ecforce.Time                   `json:"scheduled_to_be_shipped_at,omitempty"`
	SyncScheduledToBeDeliveredAt          *ecforce.BoolInt                `json:"sync_scheduled_to_be_delivered_at,omitempty"`
	ScheduledToBeDeliveredAt              *ecforce.Time                   `json:"scheduled_to_be_delivered_at,omitempty"`
	SyncScheduledToBeShippedAt            *ecforce.BoolInt                `json:"sync_scheduled_to_be_shipped_at,omitempty"`
	EnablePastScheduledToBeShippedAtError *ecforce.BoolInt                `json:"enable_past_scheduled_to_be_shipped_at_error,omitempty"`
	SkipNextScheduledDelivery             *ecforce.BoolInt                `json:"skip_next_scheduled_delivery,omitempty"`
	OnetimePointRedeem                    *ecforce.BoolInt                `json:"onetime_point_redeem,omitempty"`
	PointRedeem                           *int                            `json:"point_redeem,omitempty"`
	DefaultShippingCarrierID              *int64                          `json:"default_shipping_carrier_id,omitempty"`
	DefaultShippingCarrierLocked          *ecforce.BoolInt                `json:"default_shipping_carrier_locked,omitempty"`
	DeliveryTimeID                        *int64                          `json:"delivery_time_id,omitempty"`
	StoreID                               *int64                          `json:"store_id,omitempty"`
	Memo01                                *string                         `json:"memo01,omitempty"`
	Memo02                                *string                         `json:"memo02,omitempty"`
	AssignedAdminID                       *int64                          `json:"assigned_admin_id,omitempty"`
	RegistrantAdminID                     *int64                          `json:"registrant_admin_id,omitempty"`
	LinkedMemo                            *string                         `json:"linked_memo,omitempty"`
	Adjustment                            *int                            `json:"adjustment,omitempty"`
	OnetimeAdjustment                     *ecforce.BoolInt                `json:"onetime_adjustment,omitempty"`
	ShippingAddressAttributes             *AddressParams                  `json:"shipping_address_attributes,omitempty"`
	OrderItemsAttributes                  []*SubsOrderItemParams          `json:"order_items_attributes,omitempty"`
	PaymentAttributes                     *SubsOrderPaymentParams         `json:"payment_attributes,omitempty"`
	PickupLocationsAttributes             *SubsOrderPickupLocationsParams `json:"pickup_locations_attributes,omitempty"`
}

// SubsOrderPaymentParams sets the payment of a subscription order.
type SubsOrderPaymentParams struct {
	PaymentMethodID *int64  `json:"payment_method_id,omitempty"`
	SourceID        *int64  `json:"source_id,omitempty"`
	SourceType      *string `json:"source_type,omitempty"`
	PaymentTimes    *int    `json:"payment_times,omitempty"`
	// AmazonV2TransactionAttributes is required for Amazon Pay V2 payment
	// methods on bulk creates.
	AmazonV2TransactionAttributes *SubsOrderAmazonV2TransactionParams `json:"amazon_v2_transaction_attributes,omitempty"`
}

// SubsOrderAmazonV2TransactionParams sets the Amazon Pay V2 transaction of a
// subscription order payment.
type SubsOrderAmazonV2TransactionParams struct {
	ChargePermissionID *string `json:"charge_permission_id,omitempty"`
}

// SubsOrderPickupLocationsParams sets the preferred pickup locations of a
// subscription order.
type SubsOrderPickupLocationsParams struct {
	PickupLocationID1 *int64           `json:"pickup_location_id_1,omitempty"`
	PickupLocationID2 *int64           `json:"pickup_location_id_2,omitempty"`
	Doorbell          *ecforce.BoolInt `json:"doorbell,omitempty"`
}

// SubsOrderUpdateRequest is the body of SubsOrdersService.Update.
type SubsOrderUpdateRequest struct {
	SubsOrder *SubsOrderParams `json:"subs_order"`
	// AvoidHolidaysForScheduledToBeShippedAt pushes a next shipping date that
	// falls on a configured holiday to the next business day. Requires
	// SubsOrder.ScheduledToBeShippedAt.
	AvoidHolidaysForScheduledToBeShippedAt *ecforce.BoolInt `json:"avoid_holidays_for_scheduled_to_be_shipped_at,omitempty"`
}

// SubsOrderBulkCreateRequest is the body of SubsOrdersService.BulkCreate.
type SubsOrderBulkCreateRequest struct {
	SubsOrders                []*SubsOrderParams `json:"subs_orders"`
	CheckDuplicateLinkNumbers *ecforce.BoolInt   `json:"check_duplicate_link_numbers,omitempty"`
}

// SubsOrderBulkUpdateRequest is the body of SubsOrdersService.BulkUpdate.
type SubsOrderBulkUpdateRequest struct {
	SubsOrders                []*SubsOrderParams `json:"subs_orders"`
	CheckDuplicateLinkNumbers *ecforce.BoolInt   `json:"check_duplicate_link_numbers,omitempty"`
	// AvoidHolidaysForScheduledToBeShippedAt pushes a next shipping date that
	// falls on a configured holiday to the next business day. Requires
	// ScheduledToBeShippedAt on each subscription order.
	AvoidHolidaysForScheduledToBeShippedAt *ecforce.BoolInt `json:"avoid_holidays_for_scheduled_to_be_shipped_at,omitempty"`
}

// SubsOrderItemParams sets the attributes of a subscription order item on
// create/update requests. ID identifies the item on (bulk) updates; Skip and
// Delete apply to updates only.
type SubsOrderItemParams struct {
	ID        *int64           `json:"id,omitempty"`
	VariantID *int64           `json:"variant_id,omitempty"`
	Quantity  *int             `json:"quantity,omitempty"`
	Price     *int             `json:"price,omitempty"`
	Onetime   *ecforce.BoolInt `json:"onetime,omitempty"`
	Skip      *ecforce.BoolInt `json:"skip,omitempty"`
	Delete    *ecforce.BoolInt `json:"delete,omitempty"`
}

// SubsOrderItemCreateRequest is the body of SubsOrdersService.CreateOrderItem.
type SubsOrderItemCreateRequest struct {
	OrderItem *SubsOrderItemParams `json:"order_item"`
	// RecalculatePrice recalculates unit prices treating the next order as the
	// Times-th; it conflicts with OrderItem.Price.
	RecalculatePrice           *ecforce.BoolInt `json:"recalculate_price,omitempty"`
	RecalculatePaymentSchedule *ecforce.BoolInt `json:"recalculate_payment_schedule,omitempty"`
	// Times is required when RecalculatePrice is true.
	Times *int `json:"times,omitempty"`
	// Sync is required when RecalculatePrice is true; it syncs Times back to
	// the subscription order.
	Sync *ecforce.BoolInt `json:"sync,omitempty"`
	// Include names related data to side-load into Response.Included,
	// e.g. "variant" or "variant.changeable_variants".
	Include *string `json:"include,omitempty"`
}

// SubsOrderItemUpdateRequest is the body of SubsOrdersService.UpdateOrderItem.
type SubsOrderItemUpdateRequest struct {
	OrderItem *SubsOrderItemParams `json:"order_item"`
	// RecalculatePrice recalculates unit prices treating the next order as the
	// Times-th; it conflicts with OrderItem.Price.
	RecalculatePrice           *ecforce.BoolInt `json:"recalculate_price,omitempty"`
	RecalculatePaymentSchedule *ecforce.BoolInt `json:"recalculate_payment_schedule,omitempty"`
	// Times is required when RecalculatePrice is true.
	Times *int `json:"times,omitempty"`
	// Sync is required when RecalculatePrice is true; it syncs Times back to
	// the subscription order.
	Sync *ecforce.BoolInt `json:"sync,omitempty"`
	// Include names related data to side-load into Response.Included,
	// e.g. "variant" or "variant.changeable_variants".
	Include *string `json:"include,omitempty"`
}

// SubsOrderItemBulkCreateRequest is the body of
// SubsOrdersService.BulkCreateOrderItems.
type SubsOrderItemBulkCreateRequest struct {
	OrderItems []*SubsOrderItemParams `json:"order_items"`
	// RecalculatePrice recalculates unit prices treating the next order as the
	// Times-th; it conflicts with per-item Price.
	RecalculatePrice           *ecforce.BoolInt `json:"recalculate_price,omitempty"`
	RecalculatePaymentSchedule *ecforce.BoolInt `json:"recalculate_payment_schedule,omitempty"`
	// Times is required when RecalculatePrice is true.
	Times *int `json:"times,omitempty"`
	// Sync is required when RecalculatePrice is true; it syncs Times back to
	// the subscription order.
	Sync *ecforce.BoolInt `json:"sync,omitempty"`
}

// SubsOrderItemBulkUpdateRequest is the body of
// SubsOrdersService.BulkUpdateOrderItems.
type SubsOrderItemBulkUpdateRequest struct {
	OrderItems []*SubsOrderItemParams `json:"order_items"`
	// RecalculatePrice recalculates unit prices treating the next order as the
	// Times-th; it conflicts with per-item Price.
	RecalculatePrice           *ecforce.BoolInt `json:"recalculate_price,omitempty"`
	RecalculatePaymentSchedule *ecforce.BoolInt `json:"recalculate_payment_schedule,omitempty"`
	// Times is required when RecalculatePrice is true.
	Times *int `json:"times,omitempty"`
	// Sync is required when RecalculatePrice is true; it syncs Times back to
	// the subscription order.
	Sync *ecforce.BoolInt `json:"sync,omitempty"`
}

// SubsOrderSetParams is one set item (variant and quantity) of a subscription
// order set bulk update.
type SubsOrderSetParams struct {
	ID       *int64 `json:"id,omitempty"`
	Quantity *int   `json:"quantity,omitempty"`
}

// SubsOrderSplitDeliveryCycleRequest is the body of
// SubsOrdersService.SplitDeliveryCycle.
type SubsOrderSplitDeliveryCycleRequest struct {
	// OrderItemIDs are the order items to move to the new subscription order.
	OrderItemIDs []int64 `json:"order_item_ids"`
	// PaymentSchedule is one of "date", "term" or "day_of_week"; each choice
	// requires its matching schedule fields below.
	PaymentSchedule                                     *string          `json:"payment_schedule,omitempty"`
	ScheduledToBeDeliveredEveryXMonth                   *int             `json:"scheduled_to_be_delivered_every_x_month,omitempty"`
	ScheduledToBeDeliveredOnXthDay                      *int             `json:"scheduled_to_be_delivered_on_xth_day,omitempty"`
	ScheduledToBeDeliveredEveryXDay                     *int             `json:"scheduled_to_be_delivered_every_x_day,omitempty"`
	ScheduledToBeDeliveredOnXthDayOfWeek                *int             `json:"scheduled_to_be_delivered_on_xth_day_of_week,omitempty"`
	ScheduledToBeDeliveredEveryXDayOfWeek               *int             `json:"scheduled_to_be_delivered_every_x_day_of_week,omitempty"`
	RecalculateScheduledToBeDeliveredAtBasedOnLastOrder *ecforce.BoolInt `json:"recalculate_scheduled_to_be_delivered_at_based_on_last_order,omitempty"`
}

// List searches subscription orders. Set opts.Lighter to choose between the
// normal and lighter response variants; SubsOrder covers the attributes of
// both.
//
// Supported q attributes: id, number, times, orders_count, state, link_number,
// created_at, updated_at, suspended_at, canceled_at, with_deleted
// (orders_count: normal variant only).
// Supported sort attributes: id, created_at, updated_at.
// Supported include values: customer, customer.billing_address,
// customer.shipping_addresses, customer.notes, payment, payment.source,
// orders, order_items, shipping_address (lighter variant: customer, orders,
// order_items, shipping_address, payment).
//
// ecforce API docs: GET /api/v2/admin/subs_orders
func (s *SubsOrdersService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[SubsOrder], *ecforce.Response, error) {
	return ecforce.DoResourceList[SubsOrder](ctx, s.client, http.MethodGet, "admin/subs_orders.json", opts.Values(), nil)
}

// Get fetches a single subscription order. Set opts.Lighter to choose between
// the normal and lighter response variants; SubsOrder covers the attributes
// of both.
//
// Supported include values: customer, customer.billing_address,
// customer.shipping_addresses, customer.notes, payment, payment.source,
// orders, order_items, shipping_address (lighter variant: customer, orders,
// order_items, shipping_address, payment).
//
// ecforce API docs: GET /api/v2/admin/subs_orders/:subs_order_id
func (s *SubsOrdersService) Get(ctx context.Context, subsOrderID int64, opts *ecforce.GetOptions) (*ecforce.Resource[SubsOrder], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d.json", subsOrderID)
	return ecforce.DoResource[SubsOrder](ctx, s.client, http.MethodGet, path, opts.Values(), nil)
}

// Update updates a subscription order (synchronous).
//
// ecforce API docs: PUT /api/v2/admin/subs_orders/:subs_order_id
func (s *SubsOrdersService) Update(ctx context.Context, subsOrderID int64, req *SubsOrderUpdateRequest) (*ecforce.Resource[SubsOrder], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d.json", subsOrderID)
	return ecforce.DoResource[SubsOrder](ctx, s.client, http.MethodPut, path, nil, req)
}

// BulkCreate registers subscription orders in bulk (asynchronous). The
// returned JobResult identifies the accepted background job.
//
// ecforce API docs: POST /api/v2/admin/subs_orders/bulk_create
func (s *SubsOrdersService) BulkCreate(ctx context.Context, req *SubsOrderBulkCreateRequest) (*ecforce.JobResult, *ecforce.Response, error) {
	result := new(ecforce.JobResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPost, "admin/subs_orders/bulk_create.json", nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkUpdate updates subscription orders in bulk (synchronous).
//
// ecforce API docs: PUT /api/v2/admin/subs_orders/bulk_update
func (s *SubsOrdersService) BulkUpdate(ctx context.Context, req *SubsOrderBulkUpdateRequest) (*ecforce.BulkResult, *ecforce.Response, error) {
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, "admin/subs_orders/bulk_update.json", nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkDestroy deletes subscription orders in bulk (synchronous). Deleted
// subscription orders cannot be restored; nonexistent or already deleted IDs
// are skipped.
//
// ecforce API docs: DELETE /api/v2/admin/subs_orders/bulk_destroy
func (s *SubsOrdersService) BulkDestroy(ctx context.Context, subsOrderIDs []int64) (*ecforce.BulkResult, *ecforce.Response, error) {
	body := map[string][]int64{"subs_order_ids": subsOrderIDs}
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodDelete, "admin/subs_orders/bulk_destroy.json", nil, body, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkRecalculate recalculates subscription orders in bulk (asynchronous).
// Nonexistent IDs are skipped. The returned JobResult identifies the accepted
// background job.
//
// ecforce API docs: POST /api/v2/admin/subs_orders/bulk_recalculate
func (s *SubsOrdersService) BulkRecalculate(ctx context.Context, subsOrderIDs []int64) (*ecforce.JobResult, *ecforce.Response, error) {
	body := map[string][]int64{"subs_order_ids": subsOrderIDs}
	result := new(ecforce.JobResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPost, "admin/subs_orders/bulk_recalculate.json", nil, body, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// CreateOrder manually creates the next order (child order) of a subscription
// order.
//
// ecforce API docs: POST /api/v2/admin/subs_orders/:subs_order_id/orders
func (s *SubsOrdersService) CreateOrder(ctx context.Context, subsOrderID int64) (*ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d/orders.json", subsOrderID)
	return ecforce.Do(ctx, s.client, http.MethodPost, path, nil, nil, nil)
}

// CreateOrderItem adds an item to a subscription order (synchronous).
//
// ecforce API docs: POST /api/v2/admin/subs_orders/:subs_order_id/order_items
func (s *SubsOrdersService) CreateOrderItem(ctx context.Context, subsOrderID int64, req *SubsOrderItemCreateRequest) (*ecforce.Resource[OrderItem], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d/order_items.json", subsOrderID)
	return ecforce.DoResource[OrderItem](ctx, s.client, http.MethodPost, path, nil, req)
}

// UpdateOrderItem updates an item of a subscription order (synchronous). When
// the request deletes the item, the API answers 204 No Content and the
// returned resource is nil.
//
// ecforce API docs: PUT /api/v2/admin/subs_orders/:subs_order_id/order_items/:order_item_id
func (s *SubsOrdersService) UpdateOrderItem(ctx context.Context, subsOrderID, orderItemID int64, req *SubsOrderItemUpdateRequest) (*ecforce.Resource[OrderItem], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d/order_items/%d.json", subsOrderID, orderItemID)
	return ecforce.DoResource[OrderItem](ctx, s.client, http.MethodPut, path, nil, req)
}

// BulkCreateOrderItems adds items to a subscription order in bulk
// (synchronous). Success/failure entries refer to 1-based positions in the
// request array.
//
// ecforce API docs: POST /api/v2/admin/subs_orders/:subs_order_id/order_items/bulk_create
func (s *SubsOrdersService) BulkCreateOrderItems(ctx context.Context, subsOrderID int64, req *SubsOrderItemBulkCreateRequest) (*ecforce.BulkResult, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d/order_items/bulk_create.json", subsOrderID)
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPost, path, nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkUpdateOrderItems updates items of a subscription order in bulk
// (synchronous). Success/failure entries refer to order item IDs.
//
// ecforce API docs: PUT /api/v2/admin/subs_orders/:subs_order_id/order_items/bulk_update
func (s *SubsOrdersService) BulkUpdateOrderItems(ctx context.Context, subsOrderID int64, req *SubsOrderItemBulkUpdateRequest) (*ecforce.BulkResult, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d/order_items/bulk_update.json", subsOrderID)
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, path, nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkUpdateSets replaces the set items of a subscription order in bulk
// (synchronous). The API resets the existing set items first, so include the
// items you want to keep. Success/failure entries refer to 1-based positions
// in the request array.
//
// ecforce API docs: PUT /api/v2/admin/subs_orders/:subs_order_id/sets/bulk_update
func (s *SubsOrdersService) BulkUpdateSets(ctx context.Context, subsOrderID int64, variants []*SubsOrderSetParams) (*ecforce.BulkResult, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d/sets/bulk_update.json", subsOrderID)
	body := map[string][]*SubsOrderSetParams{"variants": variants}
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, path, nil, body, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// SplitDeliveryCycle moves some items of a subscription order onto a new
// delivery cycle, creating and returning a new subscription order
// (synchronous).
//
// ecforce API docs: POST /api/v2/admin/subs_orders/:subs_order_id/split_delivery_cycle
func (s *SubsOrdersService) SplitDeliveryCycle(ctx context.Context, subsOrderID int64, req *SubsOrderSplitDeliveryCycleRequest) (*ecforce.Resource[SubsOrder], *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d/split_delivery_cycle.json", subsOrderID)
	return ecforce.DoResource[SubsOrder](ctx, s.client, http.MethodPost, path, nil, req)
}

// ListFreeColumnValues searches the free column (custom field) values of a
// subscription order.
//
// ecforce API docs: GET /api/v2/admin/subs_orders/:subs_order_id/free_column_values
func (s *SubsOrdersService) ListFreeColumnValues(ctx context.Context, subsOrderID int64) ([]*FreeColumnValueGroup, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d/free_column_values.json", subsOrderID)
	var groups []*FreeColumnValueGroup
	resp, err := ecforce.Do(ctx, s.client, http.MethodGet, path, nil, nil, &groups)
	if err != nil {
		return nil, resp, err
	}
	return groups, resp, nil
}

// BulkCreateFreeColumnValues registers free column (custom field) values of a
// subscription order in bulk. Success/failure entries refer to 1-based
// positions in freeColumns.
//
// ecforce API docs: POST /api/v2/admin/subs_orders/:subs_order_id/free_column_values/bulk_create
func (s *SubsOrdersService) BulkCreateFreeColumnValues(ctx context.Context, subsOrderID int64, freeColumns []*FreeColumnParams) (*ecforce.BulkResult, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d/free_column_values/bulk_create.json", subsOrderID)
	body := map[string][]*FreeColumnParams{"free_columns": freeColumns}
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPost, path, nil, body, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// BulkUpdateFreeColumnValues updates free column (custom field) values of a
// subscription order in bulk. Success/failure entries refer to 1-based
// positions in freeColumns.
//
// ecforce API docs: PUT /api/v2/admin/subs_orders/:subs_order_id/free_column_values/bulk_update
func (s *SubsOrdersService) BulkUpdateFreeColumnValues(ctx context.Context, subsOrderID int64, freeColumns []*FreeColumnParams) (*ecforce.BulkResult, *ecforce.Response, error) {
	path := fmt.Sprintf("admin/subs_orders/%d/free_column_values/bulk_update.json", subsOrderID)
	body := map[string][]*FreeColumnParams{"free_columns": freeColumns}
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, path, nil, body, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// ListCancelReasons lists subscription order cancel reasons.
//
// Supported sort attributes: id, created_at, updated_at.
//
// ecforce API docs: GET /api/v2/admin/cancel_reasons
func (s *SubsOrdersService) ListCancelReasons(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[CancelReason], *ecforce.Response, error) {
	return ecforce.DoResourceList[CancelReason](ctx, s.client, http.MethodGet, "admin/cancel_reasons.json", opts.Values(), nil)
}

// ListSuspendReasons lists subscription order suspend reasons.
//
// Supported sort attributes: id, created_at, updated_at.
//
// ecforce API docs: GET /api/v2/admin/suspend_reasons
func (s *SubsOrdersService) ListSuspendReasons(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[SuspendReason], *ecforce.Response, error) {
	return ecforce.DoResourceList[SuspendReason](ctx, s.client, http.MethodGet, "admin/suspend_reasons.json", opts.Values(), nil)
}
