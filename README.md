# ecforce-go

Go client library for the [ecforce](https://ec-force.com/) v2 API.

The library is split into three packages:

| Package | Import path | Purpose |
|---|---|---|
| `ecforce` | `github.com/zero-color/ecforce-go` | Shared plumbing: HTTP client, JSON:API types, search options, errors |
| `admin` | `github.com/zero-color/ecforce-go/admin` | v2 admin API — authenticates as an administrator, mirrors the admin console |
| `customer` | `github.com/zero-color/ecforce-go/customer` | v2 customer API — authenticates as a shop customer, mirrors "my page" |

## Installation

```sh
go get github.com/zero-color/ecforce-go
```

## Usage (admin API)

```go
import (
	"context"

	"github.com/zero-color/ecforce-go"
	"github.com/zero-color/ecforce-go/admin"
)

func main() {
	ctx := context.Background()

	client, err := admin.NewClient("https://example.ec-force.com")
	if err != nil {
		// ...
	}

	// Issue a token (rate limited — reuse tokens where possible). The token
	// is stored on the client automatically.
	session, _, err := client.Sessions.SignIn(ctx, "admin@example.com", "password")
	_ = session

	// Or construct the client with an existing token:
	client, err = admin.NewClient("https://example.ec-force.com",
		ecforce.WithToken("LcPTYtBqpUFQag3GwqAVFyv_EEcsFvEC"))
	if err != nil {
		// ...
	}

	// Search customers with ransack-style predicates.
	customers, resp, err := client.Customers.List(ctx, &ecforce.ListOptions{
		Q:    ecforce.Query{"email_cont": "@example.com", "state_eq": "member"},
		Sort: []string{"-created_at"},
		Per:  100,
	})
	if err != nil {
		// ...
	}
	for _, c := range customers {
		_ = c.Attributes.Email
	}
	_ = resp.Meta.TotalCount
}
```

### Pagination

List responses carry pagination metadata on `*ecforce.Response`:

```go
opts := &ecforce.ListOptions{Per: 100}
for {
	orders, resp, err := client.Orders.List(ctx, opts)
	if err != nil {
		return err
	}
	// ... use orders ...
	if !resp.HasNextPage() {
		break
	}
	opts.Page = resp.NextPage()
}
```

### Creating and updating records

Request structs use pointer fields; only non-nil fields are sent. Use
`ecforce.Ptr` / `ecforce.String` / `ecforce.Int` / `ecforce.Bool01` /
`ecforce.NewDate` to populate them:

```go
created, _, err := client.Customers.Create(ctx, &admin.CustomerCreateRequest{
	Customer: &admin.CustomerParams{
		Email:    ecforce.String("test@example.com"),
		Password: ecforce.String("********"),
		State:    ecforce.String("member"),
		Optin:    ecforce.Bool01(true),
		Birth:    ecforce.NewDate(1990, time.January, 2),
		BillingAddressAttributes: &admin.AddressParams{
			Name01:       ecforce.String("テスト"),
			Name02:       ecforce.String("太郎"),
			Kana01:       ecforce.String("テスト"),
			Kana02:       ecforce.String("タロウ"),
			Zip01:        ecforce.String("111"),
			Zip02:        ecforce.String("1111"),
			PrefectureID: ecforce.Int64(13),
			Addr01:       ecforce.String("目黒区下目黒"),
			Addr02:       ecforce.String("2-23-18"),
			Tel01:        ecforce.String("03"),
			Tel02:        ecforce.String("5759"),
			Tel03:        ecforce.String("6380"),
		},
	},
})
```

### Bulk operations

Synchronous bulk endpoints return `*ecforce.BulkResult` with per-record
success/failure detail; asynchronous ones return `*ecforce.JobResult`
identifying the background job:

```go
result, _, err := client.Orders.BulkUpdate(ctx, req)
if err != nil {
	return err
}
for _, e := range result.Errors {
	log.Printf("order %d failed: %v", e.ID, e.Errors)
}
```

### Side-loading related data (include)

Related resources requested via `Include` are returned raw on
`Response.Included`; decode them once you dispatch on `Type`:

```go
cust, resp, err := client.Customers.Get(ctx, 30962, &ecforce.GetOptions{
	Include: []string{"billing_address"},
})
if err != nil {
	return err
}
if rel := cust.Relationships["billing_address"]; rel != nil && len(rel.Data) > 0 {
	inc := ecforce.FindIncluded(resp.Included, rel.Data[0].Type, rel.Data[0].ID)
	var addr admin.Address
	if inc != nil {
		_ = inc.DecodeAttributes(&addr)
	}
}
```

### Error handling

API errors are returned as `*ecforce.ErrorResponse` carrying the documented
error codes:

```go
_, _, err := client.Customers.Get(ctx, 1, nil)
var apiErr *ecforce.ErrorResponse
if errors.As(err, &apiErr) {
	if apiErr.HasCode("ACU2001") { // customer not found
		// ...
	}
}
```

## Usage (customer API)

```go
import "github.com/zero-color/ecforce-go/customer"

client, err := customer.NewClient("https://example.ec-force.com")
if err != nil {
	// ...
}
if _, _, err := client.Sessions.SignIn(ctx, "me@example.com", "password"); err != nil {
	// ...
}

me, _, err := client.Customer.Get(ctx, nil)
subsOrders, _, err := client.SubsOrders.List(ctx, nil)
```

## Uncovered endpoints

Every documented v2 endpoint has a typed method, but if you need to call
something bespoke, the low-level helpers accept any path:

```go
var out map[string]any
resp, err := ecforce.Do(ctx, client.Client, http.MethodGet, "admin/whatever.json", nil, nil, &out)
```

## License

Apache 2.0 — see [LICENSE](LICENSE).
