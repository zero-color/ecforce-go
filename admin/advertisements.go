package admin

import (
	"context"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// AdvertisementAccessCount aggregates access counts for one advertising URL.
type AdvertisementAccessCount struct {
	ID      int64  `json:"id"`
	BaseURL string `json:"base_url"`
	Total   int    `json:"total"`
	PC      int    `json:"pc"`
	SP      int    `json:"sp"`
	Other   int    `json:"other"`
}

// AdvertisementSearchResult is the aggregation returned by
// AdvertisementsService.List, along with the period it covers.
type AdvertisementSearchResult struct {
	AccessCounts  []*AdvertisementAccessCount
	CreatedAtGteq *ecforce.Time
	CreatedAtLt   *ecforce.Time
}

// List aggregates access counts of advertising URLs.
//
// Supported q attributes: url_id_in, created_at_gteq, created_at_lt. When the
// period is unspecified, the API covers the 8 days ending on the request date.
// The aggregation is expensive; restrict url_id_in to small batches.
//
// ecforce API docs: GET /api/v2/admin/advertisements
func (s *AdvertisementsService) List(ctx context.Context, opts *ecforce.ListOptions) (*AdvertisementSearchResult, *ecforce.Response, error) {
	var raw struct {
		Data struct {
			AccessCount []*AdvertisementAccessCount `json:"access_count"`
		} `json:"data"`
		CreatedAtGteq *ecforce.Time `json:"created_at_gteq"`
		CreatedAtLt   *ecforce.Time `json:"created_at_lt"`
	}
	resp, err := ecforce.Do(ctx, s.client, http.MethodGet, "admin/advertisements.json", opts.Values(), nil, &raw)
	if err != nil {
		return nil, resp, err
	}
	return &AdvertisementSearchResult{
		AccessCounts:  raw.Data.AccessCount,
		CreatedAtGteq: raw.CreatedAtGteq,
		CreatedAtLt:   raw.CreatedAtLt,
	}, resp, nil
}
