package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// SubsOrder is a subscription order (定期受注). It is the attribute union of
// the normal and lighter response variants; the lighter variant omits the
// resolved-name and reason attributes.
type SubsOrder struct {
	ID                      int64  `json:"id"`
	Number                  string `json:"number"`
	CustomerNumber          string `json:"customer_number"`
	Wrapping                string `json:"wrapping"`
	Times                   int    `json:"times"`
	OrdersCount             int    `json:"orders_count"`
	Store                   string `json:"store"`
	State                   string `json:"state"`
	HumanState              string `json:"human_state"`
	Subtotal                int    `json:"subtotal"`
	TBC                     bool   `json:"tbc"`
	TBCSubsOrderReasons     string `json:"tbc_subs_order_reasons"`
	RemainingNumberOfOrders int    `json:"remaining_number_of_orders"`
	// AvailablePaymentSchedules lists the delivery cycles selectable on
	// update; detail responses populate it.
	AvailablePaymentSchedules                  []*AvailablePaymentSchedule `json:"available_payment_schedules"`
	PaymentSchedule                            string                      `json:"payment_schedule"`
	PaymentScheduleLocked                      ecforce.BoolInt             `json:"payment_schedule_locked"`
	ScheduledToBeDeliveredEveryXMonth          int                         `json:"scheduled_to_be_delivered_every_x_month"`
	ScheduledToBeDeliveredOnXthDay             int                         `json:"scheduled_to_be_delivered_on_xth_day"`
	ScheduledToBeDeliveredEveryXDay            int                         `json:"scheduled_to_be_delivered_every_x_day"`
	ScheduledToBeDeliveredOnXthDayOfWeek       int                         `json:"scheduled_to_be_delivered_on_xth_day_of_week"`
	HumanScheduledToBeDeliveredEveryXDayOfWeek string                      `json:"human_scheduled_to_be_delivered_every_x_day_of_week"`
	ScheduledToBeShippedAt                     *ecforce.Time               `json:"scheduled_to_be_shipped_at"`
	ScheduledToBeDeliveredAt                   *ecforce.Time               `json:"scheduled_to_be_delivered_at"`
	ScheduledDeliveryTime                      string                      `json:"scheduled_delivery_time"`
	PreviousScheduledToBeShippedAt             *ecforce.Time               `json:"previous_scheduled_to_be_shipped_at"`
	PreviousScheduledToBeDeliveredAt           *ecforce.Time               `json:"previous_scheduled_to_be_delivered_at"`
	SuspendReasons                             string                      `json:"suspend_reasons"`
	CancelReasons                              string                      `json:"cancel_reasons"`
	ElapsedDaysFromSuspended                   int                         `json:"elapsed_days_from_suspended"`
	ElapsedDaysFromCanceled                    int                         `json:"elapsed_days_from_canceled"`
	PaymentMethodName                          string                      `json:"payment_method_name"`
	PaymentTimes                               int                         `json:"payment_times"`
	RecurringBlockTimes                        int                         `json:"recurring_block_times"`
	OrdersPaidCount                            int                         `json:"orders_paid_count"`
	RemainingBlockTimes                        int                         `json:"remaining_block_times"`
	SuspendReasonIDs                           []int64                     `json:"suspend_reason_ids"`
	CancelReasonIDs                            []int64                     `json:"cancel_reason_ids"`
	Labels                                     string                      `json:"labels"`
	CustomerLabels                             string                      `json:"customer_labels"`
	ProductLabels                              string                      `json:"product_labels"`
	CustomerContacts                           string                      `json:"customer_contacts"`
	Memo01                                     string                      `json:"memo01"`
	Memo02                                     string                      `json:"memo02"`
	DefaultShippingCarrierID                   int64                       `json:"default_shipping_carrier_id"`
	DefaultShippingCarrier                     string                      `json:"default_shipping_carrier"`
	DefaultShippingCarrierLocked               bool                        `json:"default_shipping_carrier_locked"`
	TenantID                                   int64                       `json:"tenant_id"`
	PayWithEpos                                bool                        `json:"pay_with_epos"`
	SkipNextScheduledDelivery                  bool                        `json:"skip_next_scheduled_delivery"`
	CreatedAt                                  *ecforce.Time               `json:"created_at"`
	UpdatedAt                                  *ecforce.Time               `json:"updated_at"`
	SuspendedAt                                *ecforce.Time               `json:"suspended_at"`
	CanceledAt                                 *ecforce.Time               `json:"canceled_at"`
	DeletedAt                                  *ecforce.Time               `json:"deleted_at"`
}

// AvailablePaymentSchedule is one selectable delivery cycle (date, term or
// day_of_week) and the values selectable for each of its components.
type AvailablePaymentSchedule struct {
	Value                                 string                  `json:"value"`
	Ja                                    string                  `json:"ja"`
	ScheduledToBeDeliveredEveryXMonth     []*PaymentScheduleValue `json:"scheduled_to_be_delivered_every_x_month,omitempty"`
	ScheduledToBeDeliveredOnXthDay        []*PaymentScheduleValue `json:"scheduled_to_be_delivered_on_xth_day,omitempty"`
	ScheduledToBeDeliveredEveryXDay       []*PaymentScheduleValue `json:"scheduled_to_be_delivered_every_x_day,omitempty"`
	ScheduledToBeDeliveredOnXthDayOfWeek  []*PaymentScheduleValue `json:"scheduled_to_be_delivered_on_xth_day_of_week,omitempty"`
	ScheduledToBeDeliveredEveryXDayOfWeek []*PaymentScheduleValue `json:"scheduled_to_be_delivered_every_x_day_of_week,omitempty"`
}

// PaymentScheduleValue is one selectable value of a delivery cycle component,
// with its Japanese display label.
type PaymentScheduleValue struct {
	Value int    `json:"value"`
	Ja    string `json:"ja"`
}

// UnmarshalJSON implements json.Unmarshaler. The API emits the label as
// either a string ("1ヶ月") or a bare number (1); both decode into Ja.
func (v *PaymentScheduleValue) UnmarshalJSON(data []byte) error {
	var raw struct {
		Value int             `json:"value"`
		Ja    json.RawMessage `json:"ja"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	v.Value = raw.Value
	v.Ja = ""
	if s := string(raw.Ja); len(raw.Ja) > 0 && s != "null" {
		if raw.Ja[0] == '"' {
			if err := json.Unmarshal(raw.Ja, &v.Ja); err != nil {
				return err
			}
		} else {
			v.Ja = s
		}
	}
	return nil
}

// SubsOrderParams are the attributes of SubsOrdersService.Update. Only the
// combinations the API documents may be updated together (e.g. State with
// SuspendReasonIDs or CancelReasonIDs; PaymentSchedule with its component
// values; ScheduledToBeDeliveredAt and SkipNextScheduledDelivery alone).
type SubsOrderParams struct {
	// State is the subscription status: active, suspend or canceled.
	State            *string `json:"state,omitempty"`
	SuspendReasonIDs []int64 `json:"suspend_reason_ids,omitempty"`
	CancelReasonIDs  []int64 `json:"cancel_reason_ids,omitempty"`
	// PaymentSchedule is the delivery cycle kind: date, term or day_of_week.
	PaymentSchedule                       *string `json:"payment_schedule,omitempty"`
	ScheduledToBeDeliveredEveryXMonth     *int    `json:"scheduled_to_be_delivered_every_x_month,omitempty"`
	ScheduledToBeDeliveredOnXthDay        *int    `json:"scheduled_to_be_delivered_on_xth_day,omitempty"`
	ScheduledToBeDeliveredEveryXDay       *int    `json:"scheduled_to_be_delivered_every_x_day,omitempty"`
	ScheduledToBeDeliveredOnXthDayOfWeek  *int    `json:"scheduled_to_be_delivered_on_xth_day_of_week,omitempty"`
	ScheduledToBeDeliveredEveryXDayOfWeek *int    `json:"scheduled_to_be_delivered_every_x_day_of_week,omitempty"`
	// ScheduledToBeDeliveredAt sets the next scheduled delivery date.
	ScheduledToBeDeliveredAt  *ecforce.Time    `json:"scheduled_to_be_delivered_at,omitempty"`
	SkipNextScheduledDelivery *ecforce.BoolInt `json:"skip_next_scheduled_delivery,omitempty"`
	// RecalculateScheduledToBeDeliveredAtBasedOnLastOrder recalculates the
	// scheduled delivery/shipping dates from the last order.
	RecalculateScheduledToBeDeliveredAtBasedOnLastOrder *ecforce.BoolInt `json:"recalculate_scheduled_to_be_delivered_at_based_on_last_order,omitempty"`
	// EnableRecalcPastDateError controls whether recalculated dates in the
	// past raise an error; effective only with the recalculate flag set.
	EnableRecalcPastDateError *ecforce.BoolInt `json:"enable_recalc_past_date_error,omitempty"`
}

// OrderItemParams are the attributes of SubsOrdersService.AddOrderItem and
// SubsOrdersService.UpdateOrderItem. Onetime applies to AddOrderItem only.
type OrderItemParams struct {
	VariantID *int64 `json:"variant_id,omitempty"`
	Quantity  *int   `json:"quantity,omitempty"`
	// Onetime delivers the item on the next order only (1) or on every order
	// (0). AddOrderItem only.
	Onetime *ecforce.BoolInt `json:"onetime,omitempty"`
}

// List searches the authenticated customer's subscription orders.
//
// Supported q attributes: id, number, times, orders_count, state, created_at,
// updated_at, suspended_at, canceled_at, with_deleted.
// Supported sort attributes: id, created_at, updated_at.
// Supported include values: customer, customer.billing_address,
// customer.shipping_addresses, customer.notes, orders, order_items,
// shipping_address, payment.
//
// ecforce API docs: GET /api/v2/customer/subs_orders
func (s *SubsOrdersService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[SubsOrder], *ecforce.Response, error) {
	return ecforce.DoResourceList[SubsOrder](ctx, s.client, http.MethodGet, "customer/subs_orders.json", opts.Values(), nil)
}

// Get fetches a single subscription order of the authenticated customer.
//
// Supported include values: customer, customer.billing_address,
// customer.shipping_addresses, customer.notes, orders, order_items,
// shipping_address, payment.
//
// ecforce API docs: GET /api/v2/customer/subs_orders/:subs_order_id
func (s *SubsOrdersService) Get(ctx context.Context, subsOrderID int64, opts *ecforce.GetOptions) (*ecforce.Resource[SubsOrder], *ecforce.Response, error) {
	path := fmt.Sprintf("customer/subs_orders/%d.json", subsOrderID)
	return ecforce.DoResource[SubsOrder](ctx, s.client, http.MethodGet, path, opts.Values(), nil)
}

// Update updates a subscription order (synchronous). The API returns the
// updated subscription order as a one-element array.
//
// ecforce API docs: PUT /api/v2/customer/subs_orders/:subs_order_id
func (s *SubsOrdersService) Update(ctx context.Context, subsOrderID int64, params *SubsOrderParams) ([]*ecforce.Resource[SubsOrder], *ecforce.Response, error) {
	path := fmt.Sprintf("customer/subs_orders/%d.json", subsOrderID)
	body := map[string]any{"subs_order": params}
	return ecforce.DoResourceList[SubsOrder](ctx, s.client, http.MethodPut, path, nil, body)
}

// AddOrderItem adds an item to a subscription order (synchronous). The API
// returns the added order items as an array.
//
// ecforce API docs: POST /api/v2/customer/subs_orders/:subs_order_id/order_items
func (s *SubsOrdersService) AddOrderItem(ctx context.Context, subsOrderID int64, params *OrderItemParams) ([]*ecforce.Resource[OrderItem], *ecforce.Response, error) {
	path := fmt.Sprintf("customer/subs_orders/%d/order_items.json", subsOrderID)
	body := map[string]any{"order_item": params}
	return ecforce.DoResourceList[OrderItem](ctx, s.client, http.MethodPost, path, nil, body)
}

// UpdateOrderItem updates an item of a subscription order (synchronous). The
// API returns the updated order items as an array.
//
// Supported include values: variant, variant.changeable_variants.
//
// ecforce API docs: PUT /api/v2/customer/subs_orders/:subs_order_id/order_items/:order_item_id
func (s *SubsOrdersService) UpdateOrderItem(ctx context.Context, subsOrderID, orderItemID int64, params *OrderItemParams, opts *ecforce.GetOptions) ([]*ecforce.Resource[OrderItem], *ecforce.Response, error) {
	path := fmt.Sprintf("customer/subs_orders/%d/order_items/%d.json", subsOrderID, orderItemID)
	body := map[string]any{"order_item": params}
	return ecforce.DoResourceList[OrderItem](ctx, s.client, http.MethodPut, path, opts.Values(), body)
}

// DestroyOrderItem removes an item from a subscription order (synchronous).
//
// ecforce API docs: DELETE /api/v2/customer/subs_orders/:subs_order_id/order_items/:order_item_id
func (s *SubsOrdersService) DestroyOrderItem(ctx context.Context, subsOrderID, orderItemID int64) (*ecforce.Response, error) {
	path := fmt.Sprintf("customer/subs_orders/%d/order_items/%d.json", subsOrderID, orderItemID)
	return ecforce.Do(ctx, s.client, http.MethodDelete, path, nil, nil, nil)
}
