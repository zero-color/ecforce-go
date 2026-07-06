package customer

import (
	"context"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// SignUpParams are the attributes of RegistrationService.SignUp.
type SignUpParams struct {
	// Email is the customer's email address (required).
	Email *string `json:"email,omitempty"`
	// State is the membership status (required): visitor, guest or member.
	State    *string `json:"state,omitempty"`
	Password *string `json:"password,omitempty"`
	SexID    *int64  `json:"sex_id,omitempty"`
	JobID    *int64  `json:"job_id,omitempty"`
	// Birth is the birth date in yyyy/mm/dd format.
	Birth                    *string        `json:"birth,omitempty"`
	LinkNumber               *string        `json:"link_number,omitempty"`
	TenantID                 *int64         `json:"tenant_id,omitempty"`
	BillingAddressAttributes *AddressParams `json:"billing_address_attributes,omitempty"`
}

// SignUp registers a new customer (synchronous). The returned customer record
// includes an API authentication token in its attributes.
//
// ecforce API docs: POST /api/v2/customers
func (s *RegistrationService) SignUp(ctx context.Context, params *SignUpParams) (*ecforce.Resource[Customer], *ecforce.Response, error) {
	body := map[string]any{"customer": params}
	return ecforce.DoResource[Customer](ctx, s.client, http.MethodPost, "customers.json", nil, body)
}
