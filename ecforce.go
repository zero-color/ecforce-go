// Package ecforce provides the shared HTTP plumbing for the ecforce v2 API
// client libraries. Most users should use the admin or customer subpackages,
// which expose typed services on top of this package.
//
// The ecforce v2 API lives under https://<shop domain>/api/v2 and serves
// JSON:API-style documents. Authentication uses a token issued by the
// sign-in endpoints, sent as `Authorization: Token token="..."`.
package ecforce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Version is the version of this library.
const Version = "0.1.0"

const (
	defaultAPIPath   = "/api/v2/"
	defaultUserAgent = "ecforce-go/" + Version
	mediaTypeJSON    = "application/json"
)

// Client manages communication with the ecforce v2 API. It is safe for
// concurrent use by multiple goroutines.
type Client struct {
	client *http.Client

	// BaseURL for API requests, e.g. https://example.ec-force.com/api/v2/.
	// BaseURL should always have a trailing slash.
	BaseURL *url.URL

	// UserAgent used when communicating with the ecforce API.
	UserAgent string

	mu    sync.RWMutex
	token string
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithHTTPClient sets the underlying *http.Client used to make requests.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) { c.client = hc }
}

// WithToken sets the API authentication token. Tokens are issued by the
// sign-in endpoints (admin.SessionsService.SignIn / customer.SessionsService.SignIn).
func WithToken(token string) ClientOption {
	return func(c *Client) { c.token = token }
}

// WithUserAgent sets the User-Agent header sent with every request.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) { c.UserAgent = ua }
}

// NewClient returns a new ecforce API client for the shop at baseURL.
//
// baseURL is the shop origin, e.g. "https://example.ec-force.com". If it has
// no path, the default API path "/api/v2/" is appended; a custom path is
// preserved (with a trailing slash added if missing).
func NewClient(baseURL string, opts ...ClientOption) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("ecforce: invalid base URL %q: %w", baseURL, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("ecforce: base URL %q must be absolute, e.g. https://example.ec-force.com", baseURL)
	}
	if u.Path == "" || u.Path == "/" {
		u.Path = defaultAPIPath
	} else if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	c := &Client{
		client:    http.DefaultClient,
		BaseURL:   u,
		UserAgent: defaultUserAgent,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// SetToken replaces the API authentication token used for subsequent requests.
func (c *Client) SetToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
}

// Token returns the API authentication token currently in use.
func (c *Client) Token() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

// NewRequest creates an API request. path is resolved relative to BaseURL and
// should not have a leading slash, e.g. "admin/customers.json". If body is
// non-nil it is JSON-encoded as the request body.
func (c *Client) NewRequest(ctx context.Context, method, path string, query url.Values, body any) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("ecforce: invalid request path %q: %w", path, err)
	}
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}

	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("ecforce: encoding request body: %w", err)
		}
		buf = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", mediaTypeJSON)
	}
	req.Header.Set("Accept", mediaTypeJSON)
	if ua := c.UserAgent; ua != "" {
		req.Header.Set("User-Agent", ua)
	}
	if token := c.Token(); token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Token token=%q", token))
	}
	return req, nil
}

// Response wraps an ecforce API response. It carries the pagination metadata
// and included resources of JSON:API documents when present.
type Response struct {
	*http.Response

	// Meta holds pagination metadata for list responses.
	Meta *Meta
	// Links holds pagination links for list responses.
	Links *Links
	// Included holds side-loaded resources requested via the include parameter.
	Included []*IncludedResource
}

// HasNextPage reports whether a further page of results exists.
func (r *Response) HasNextPage() bool {
	return r.Meta != nil && r.Meta.Page < r.Meta.TotalPages
}

// NextPage returns the next page number, or 0 if there is none.
func (r *Response) NextPage() int {
	if !r.HasNextPage() {
		return 0
	}
	return r.Meta.Page + 1
}

// Do sends an API request and decodes the JSON response body into v when v is
// non-nil and the response has content. The response body is closed. If the
// API returns an error status, the returned error is an *ErrorResponse.
func (c *Client) Do(req *http.Request, v any) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}

	if err := CheckResponse(resp); err != nil {
		return response, err
	}

	if v == nil || resp.StatusCode == http.StatusNoContent {
		// Drain so the connection can be reused.
		_, _ = io.Copy(io.Discard, resp.Body)
		return response, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("ecforce: reading response body: %w", err)
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return response, nil
	}
	if err := json.Unmarshal(body, v); err != nil {
		return response, fmt.Errorf("ecforce: decoding response body: %w", err)
	}
	return response, nil
}

// Do sends a request with an optional JSON body and decodes the response into
// v when v is non-nil. It is the low-level escape hatch for endpoints not
// covered by a typed service method.
func Do(ctx context.Context, c *Client, method, path string, query url.Values, body, v any) (*Response, error) {
	req, err := c.NewRequest(ctx, method, path, query, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

// DoResource performs a request whose response is a JSON:API document with a
// single resource in "data", and returns that resource.
func DoResource[T any](ctx context.Context, c *Client, method, path string, query url.Values, body any) (*Resource[T], *Response, error) {
	var doc struct {
		Data     *Resource[T]        `json:"data"`
		Included []*IncludedResource `json:"included"`
		Meta     *Meta               `json:"meta"`
		Links    *Links              `json:"links"`
	}
	resp, err := Do(ctx, c, method, path, query, body, &doc)
	if err != nil {
		return nil, resp, err
	}
	resp.Meta, resp.Links, resp.Included = doc.Meta, doc.Links, doc.Included
	return doc.Data, resp, nil
}

// DoResourceList performs a request whose response is a JSON:API document with
// an array of resources in "data", and returns those resources. Pagination
// metadata is attached to the returned Response.
func DoResourceList[T any](ctx context.Context, c *Client, method, path string, query url.Values, body any) ([]*Resource[T], *Response, error) {
	var doc struct {
		Data     []*Resource[T]      `json:"data"`
		Included []*IncludedResource `json:"included"`
		Meta     *Meta               `json:"meta"`
		Links    *Links              `json:"links"`
	}
	resp, err := Do(ctx, c, method, path, query, body, &doc)
	if err != nil {
		return nil, resp, err
	}
	resp.Meta, resp.Links, resp.Included = doc.Meta, doc.Links, doc.Included
	return doc.Data, resp, nil
}
