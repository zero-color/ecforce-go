package customer

import (
	"context"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// InviteCode is a customer invite code (招待コード).
type InviteCode struct {
	ID        int64  `json:"id"`
	Number    string `json:"number"`
	UsedTimes int    `json:"used_times"`
	// AvailableInvitedTimes is the number of remaining uses. When MaxLimit is
	// unlimited (0), AvailableInvitedTimes is zero or negative.
	AvailableInvitedTimes int `json:"available_invited_times"`
	// MaxLimit is the maximum number of uses; 0 means unlimited.
	MaxLimit int `json:"max_limit"`
}

// List fetches the authenticated customer's invite codes.
//
// ecforce API docs: GET /api/v2/customer/invite_codes
func (s *InviteCodesService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[InviteCode], *ecforce.Response, error) {
	return ecforce.DoResourceList[InviteCode](ctx, s.client, http.MethodGet, "customer/invite_codes.json", opts.Values(), nil)
}

// Create issues an invite code for the authenticated customer. Issuing a
// second code fails with error code CIC3004.
//
// ecforce API docs: POST /api/v2/customer/invite_codes
func (s *InviteCodesService) Create(ctx context.Context) (*ecforce.Resource[InviteCode], *ecforce.Response, error) {
	return ecforce.DoResource[InviteCode](ctx, s.client, http.MethodPost, "customer/invite_codes.json", nil, nil)
}
