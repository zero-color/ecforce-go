// Package admin provides a client for the ecforce v2 admin API, which
// authenticates as an ecforce administrator and offers the same capabilities
// as the admin console.
//
// Construct a client with NewClient, then either sign in with an
// administrator's credentials or supply a previously issued token:
//
//	client, err := admin.NewClient("https://example.ec-force.com")
//	if err != nil { ... }
//	session, _, err := client.Sessions.SignIn(ctx, "admin@example.com", "password")
//
//	// or, with a token in hand:
//	client, err := admin.NewClient("https://example.ec-force.com", ecforce.WithToken(token))
//
// API methods are grouped into services mirroring the API documentation:
//
//	customers, resp, err := client.Customers.List(ctx, &ecforce.ListOptions{
//		Q:    ecforce.Query{"email_cont": "@example.com"},
//		Sort: []string{"-created_at"},
//		Per:  100,
//	})
package admin

import (
	"github.com/zero-color/ecforce-go"
)

// Client is a client for the ecforce v2 admin API.
type Client struct {
	*ecforce.Client

	Advertisements    *AdvertisementsService
	Customers         *CustomersService
	Labels            *LabelsService
	OrderFreeColumns  *OrderFreeColumnsService
	Orders            *OrdersService
	PaymentMethods    *PaymentMethodsService
	ProductCategories *ProductCategoriesService
	Products          *ProductsService
	Sessions          *SessionsService
	Sexes             *SexesService
	ShippingCarriers  *ShippingCarriersService
	StockItems        *StockItemsService
	Stores            *StoresService
	SubsOrders        *SubsOrdersService
	URLGroups         *URLGroupsService
	URLs              *URLsService
}

// service is the shared backing store of all services.
type service struct {
	client *ecforce.Client
}

// One service type per API group. All share the same underlying value so the
// Client stays cheap to construct.
type (
	// AdvertisementsService handles the admin advertisement API.
	AdvertisementsService service
	// CustomersService handles the admin customer API.
	CustomersService service
	// LabelsService handles the admin label API.
	LabelsService service
	// OrderFreeColumnsService handles the admin order free column API.
	OrderFreeColumnsService service
	// OrdersService handles the admin order API.
	OrdersService service
	// PaymentMethodsService handles the admin payment method API.
	PaymentMethodsService service
	// ProductCategoriesService handles the admin product category API.
	ProductCategoriesService service
	// ProductsService handles the admin product API.
	ProductsService service
	// SessionsService handles admin authentication.
	SessionsService service
	// SexesService handles the admin sex master API.
	SexesService service
	// ShippingCarriersService handles the admin shipping carrier API.
	ShippingCarriersService service
	// StockItemsService handles the admin stock item API.
	StockItemsService service
	// StoresService handles the admin store master API.
	StoresService service
	// SubsOrdersService handles the admin subscription order API.
	SubsOrdersService service
	// URLGroupsService handles the admin advertising URL group API.
	URLGroupsService service
	// URLsService handles the admin advertising URL API.
	URLsService service
)

// NewClient returns a new ecforce admin API client for the shop at baseURL,
// e.g. "https://example.ec-force.com". Pass ecforce.WithToken to authenticate
// with an existing token, or call Sessions.SignIn to issue one.
func NewClient(baseURL string, opts ...ecforce.ClientOption) (*Client, error) {
	core, err := ecforce.NewClient(baseURL, opts...)
	if err != nil {
		return nil, err
	}
	c := &Client{Client: core}
	s := service{client: core}
	c.Advertisements = (*AdvertisementsService)(&s)
	c.Customers = (*CustomersService)(&s)
	c.Labels = (*LabelsService)(&s)
	c.OrderFreeColumns = (*OrderFreeColumnsService)(&s)
	c.Orders = (*OrdersService)(&s)
	c.PaymentMethods = (*PaymentMethodsService)(&s)
	c.ProductCategories = (*ProductCategoriesService)(&s)
	c.Products = (*ProductsService)(&s)
	c.Sessions = (*SessionsService)(&s)
	c.Sexes = (*SexesService)(&s)
	c.ShippingCarriers = (*ShippingCarriersService)(&s)
	c.StockItems = (*StockItemsService)(&s)
	c.Stores = (*StoresService)(&s)
	c.SubsOrders = (*SubsOrdersService)(&s)
	c.URLGroups = (*URLGroupsService)(&s)
	c.URLs = (*URLsService)(&s)
	return c, nil
}
