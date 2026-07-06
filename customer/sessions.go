package customer

import (
	"context"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// Session is a customer API session issued by SignIn.
type Session struct {
	ID                  int64  `json:"id"`
	Email               string `json:"email"`
	AuthenticationToken string `json:"authentication_token"`
}

// SignIn issues an API authentication token for the customer and stores it on
// the client for subsequent requests.
//
// The API rate-limits this endpoint (5 requests per hour); reuse tokens rather
// than signing in per request.
//
// ecforce API docs: POST /api/v2/customers/sign_in
func (s *SessionsService) SignIn(ctx context.Context, email, password string) (*Session, *ecforce.Response, error) {
	body := map[string]any{
		"customer": map[string]string{"email": email, "password": password},
	}
	session := new(Session)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPost, "customers/sign_in.json", nil, body, session)
	if err != nil {
		return nil, resp, err
	}
	// SMS-verification flows answer 202 Accepted without a token; don't
	// clobber an existing token in that case.
	if session.AuthenticationToken != "" {
		s.client.SetToken(session.AuthenticationToken)
	}
	return session, resp, nil
}

// SignOut revokes the client's API authentication token and clears it from
// the client.
//
// ecforce API docs: GET /api/v2/customers/sign_out
func (s *SessionsService) SignOut(ctx context.Context) (*ecforce.Response, error) {
	resp, err := ecforce.Do(ctx, s.client, http.MethodGet, "customers/sign_out.json", nil, nil, nil)
	if err != nil {
		return resp, err
	}
	s.client.SetToken("")
	return resp, nil
}
