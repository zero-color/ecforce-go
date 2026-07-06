package admin

import (
	"context"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// PaymentMethod is a payment method (支払い方法).
type PaymentMethod struct {
	ID                         int64  `json:"id"`
	Name                       string `json:"name"`
	Description                string `json:"description"`
	DescriptionMobile          string `json:"description_mobile"`
	Active                     bool   `json:"active"`
	AmountUpperLimit           int    `json:"amount_upper_limit"`
	AmountUnderLimit           int    `json:"amount_under_limit"`
	DisplayOnMypage            bool   `json:"display_on_mypage"`
	DisplayOnAdminOiPage       bool   `json:"display_on_admin_oi_page"`
	UseCreditCardValidityCheck bool   `json:"use_credit_card_validity_check"`
	CvOfferRealtimePayment     bool   `json:"cv_offer_realtime_payment"`
	UseOPlux                   bool   `json:"use_o_plux"`
	UseSpiderAf                bool   `json:"use_spider_af"`
}

// List searches payment methods.
//
// Supported q attributes: id, name, active, display_on_mypage,
// display_on_admin_oi_page, with_deleted.
// Supported sort attributes: id, created_at, updated_at.
//
// ecforce API docs: GET /api/v2/admin/payment_methods
func (s *PaymentMethodsService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[PaymentMethod], *ecforce.Response, error) {
	return ecforce.DoResourceList[PaymentMethod](ctx, s.client, http.MethodGet, "admin/payment_methods.json", opts.Values(), nil)
}
