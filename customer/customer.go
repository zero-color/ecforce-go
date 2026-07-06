// Package customer provides a client for the ecforce v2 customer API, which
// authenticates as a shop customer and offers the same capabilities as the
// customer "my page".
//
// Construct a client with NewClient, then either sign in with a customer's
// credentials or supply a previously issued token:
//
//	client, err := customer.NewClient("https://example.ec-force.com")
//	if err != nil { ... }
//	session, _, err := client.Sessions.SignIn(ctx, "customer@example.com", "password")
//
//	// or, with a token in hand:
//	client, err := customer.NewClient("https://example.ec-force.com", ecforce.WithToken(token))
//
// API methods are grouped into services mirroring the API documentation:
//
//	subsOrders, resp, err := client.SubsOrders.List(ctx, &ecforce.ListOptions{
//		Q:    ecforce.Query{"state_eq": "active"},
//		Sort: []string{"-created_at"},
//		Per:  100,
//	})
package customer

import (
	"github.com/zero-color/ecforce-go"
)

// Client is a client for the ecforce v2 customer API.
type Client struct {
	*ecforce.Client

	Customer     *CustomerService
	InviteCodes  *InviteCodesService
	MultiPass    *MultiPassService
	Registration *RegistrationService
	Sessions     *SessionsService
	SubsOrders   *SubsOrdersService
}

// service is the shared backing store of all services.
type service struct {
	client *ecforce.Client
}

// One service type per API group. All share the same underlying value so the
// Client stays cheap to construct.
type (
	// CustomerService handles the authenticated customer's own record.
	CustomerService service
	// InviteCodesService handles the customer invite code API.
	InviteCodesService service
	// MultiPassService handles MultiPass authentication.
	MultiPassService service
	// RegistrationService handles customer registration.
	RegistrationService service
	// SessionsService handles customer authentication.
	SessionsService service
	// SubsOrdersService handles the customer subscription order API.
	SubsOrdersService service
)

// NewClient returns a new ecforce customer API client for the shop at baseURL,
// e.g. "https://example.ec-force.com". Pass ecforce.WithToken to authenticate
// with an existing token, or call Sessions.SignIn to issue one.
func NewClient(baseURL string, opts ...ecforce.ClientOption) (*Client, error) {
	core, err := ecforce.NewClient(baseURL, opts...)
	if err != nil {
		return nil, err
	}
	c := &Client{Client: core}
	s := service{client: core}
	c.Customer = (*CustomerService)(&s)
	c.InviteCodes = (*InviteCodesService)(&s)
	c.MultiPass = (*MultiPassService)(&s)
	c.Registration = (*RegistrationService)(&s)
	c.Sessions = (*SessionsService)(&s)
	c.SubsOrders = (*SubsOrdersService)(&s)
	return c, nil
}
