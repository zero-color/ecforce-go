package customer

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zero-color/ecforce-go"
)

// MultiPass is the result of a MultiPass authentication: the identified
// customer, an API authentication token and a one-time login URL.
type MultiPass struct {
	ID                  int64  `json:"id"`
	Email               string `json:"email"`
	AuthenticationToken string `json:"authentication_token"`
	// LoginURL logs anyone in as the customer without credentials; handle it
	// with care. It expires after 5 minutes and can be used only once.
	LoginURL string `json:"login_url"`
}

// Get exchanges a MultiPass token for an API authentication token and a
// one-time login URL. A MultiPass token that has already issued an
// authentication token cannot be used again.
//
// ecforce API docs: GET /api/v2/customers/multi_pass/:token
func (s *MultiPassService) Get(ctx context.Context, token string) (*MultiPass, *ecforce.Response, error) {
	path := fmt.Sprintf("customers/multi_pass/%s.json", url.PathEscape(token))
	mp := new(MultiPass)
	resp, err := ecforce.Do(ctx, s.client, http.MethodGet, path, nil, nil, mp)
	if err != nil {
		return nil, resp, err
	}
	return mp, resp, nil
}
