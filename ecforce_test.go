package ecforce

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// setup returns a test client wired to a test server, and the mux to register
// handlers on.
func setup(t *testing.T) (*Client, *http.ServeMux) {
	t.Helper()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL, WithToken("test-token"))
	if err != nil {
		t.Fatal(err)
	}
	return client, mux
}

func TestNewClient(t *testing.T) {
	c, err := NewClient("https://example.ec-force.com")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.BaseURL.String(), "https://example.ec-force.com/api/v2/"; got != want {
		t.Errorf("BaseURL = %q, want %q", got, want)
	}

	c, err = NewClient("https://example.ec-force.com/custom/path")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := c.BaseURL.String(), "https://example.ec-force.com/custom/path/"; got != want {
		t.Errorf("BaseURL = %q, want %q", got, want)
	}

	if _, err := NewClient("not a url\x7f"); err == nil {
		t.Error("NewClient with invalid URL succeeded, want error")
	}
	if _, err := NewClient("/relative"); err == nil {
		t.Error("NewClient with relative URL succeeded, want error")
	}
}

func TestNewRequest_headers(t *testing.T) {
	c, err := NewClient("https://example.ec-force.com", WithToken("secret"))
	if err != nil {
		t.Fatal(err)
	}
	req, err := c.NewRequest(context.Background(), http.MethodPost, "admin/customers.json", nil, map[string]string{"a": "b"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := req.URL.String(), "https://example.ec-force.com/api/v2/admin/customers.json"; got != want {
		t.Errorf("URL = %q, want %q", got, want)
	}
	if got, want := req.Header.Get("Authorization"), `Token token="secret"`; got != want {
		t.Errorf("Authorization = %q, want %q", got, want)
	}
	if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
	if got, want := req.Header.Get("User-Agent"), defaultUserAgent; got != want {
		t.Errorf("User-Agent = %q, want %q", got, want)
	}
}

func TestDoResourceList_pagination(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/things.json", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("page = %q, want 2", got)
		}
		fmt.Fprint(w, `{
			"data": [{"id": "7", "type": "thing", "attributes": {"name": "x"},
				"relationships": {"parent": {"data": {"id": "1", "type": "thing"}}}}],
			"included": [{"id": "1", "type": "thing", "attributes": {"name": "parent"}}],
			"meta": {"total_count": 30, "page": 2, "per": 10, "count": 10, "total_pages": 3},
			"links": {"next": "..."}
		}`)
	})

	type thing struct {
		Name string `json:"name"`
	}
	things, resp, err := DoResourceList[thing](context.Background(), client, http.MethodGet, "admin/things.json", url.Values{"page": {"2"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(things) != 1 || things[0].ID != "7" || things[0].Attributes.Name != "x" {
		t.Errorf("unexpected resources: %+v", things)
	}
	rel := things[0].Relationships["parent"]
	if rel == nil || len(rel.Data) != 1 || rel.Data[0].ID != "1" {
		t.Errorf("unexpected relationships: %+v", things[0].Relationships)
	}
	if resp.Meta == nil || resp.Meta.TotalPages != 3 {
		t.Errorf("unexpected meta: %+v", resp.Meta)
	}
	if !resp.HasNextPage() || resp.NextPage() != 3 {
		t.Errorf("HasNextPage/NextPage = %v/%d, want true/3", resp.HasNextPage(), resp.NextPage())
	}
	if len(resp.Included) != 1 || resp.Included[0].ID != "1" {
		t.Errorf("unexpected included: %+v", resp.Included)
	}
	var parent thing
	if err := resp.Included[0].DecodeAttributes(&parent); err != nil || parent.Name != "parent" {
		t.Errorf("DecodeAttributes = %+v, %v", parent, err)
	}
}

func TestDoResource_single(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/things/7.json", func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), `Token token="test-token"`; got != want {
			t.Errorf("Authorization = %q, want %q", got, want)
		}
		fmt.Fprint(w, `{"data": {"id": "7", "type": "thing", "attributes": {"name": "x"}}}`)
	})

	type thing struct {
		Name string `json:"name"`
	}
	got, _, err := DoResource[thing](context.Background(), client, http.MethodGet, "admin/things/7.json", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "7" || got.Attributes.Name != "x" {
		t.Errorf("unexpected resource: %+v", got)
	}
}

func TestDo_errorResponse(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admin/boom.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"errors": [{"code": "ASS1001", "message": "認証に失敗しました。"}]}`)
	})

	_, err := Do(context.Background(), client, http.MethodGet, "admin/boom.json", nil, nil, nil)
	var errResp *ErrorResponse
	if !errors.As(err, &errResp) {
		t.Fatalf("error = %T(%v), want *ErrorResponse", err, err)
	}
	if !errResp.HasCode("ASS1001") {
		t.Errorf("HasCode(ASS1001) = false, want true; errors: %+v", errResp.Errors)
	}
	if errResp.HasCode("XXX") {
		t.Error("HasCode(XXX) = true, want false")
	}
	if errResp.Response.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want 401", errResp.Response.StatusCode)
	}
}

func TestDo_noContent(t *testing.T) {
	client, mux := setup(t)
	mux.HandleFunc("GET /api/v2/admins/sign_out.json", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	resp, err := Do(context.Background(), client, http.MethodGet, "admins/sign_out.json", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d, want 204", resp.StatusCode)
	}
}

func TestResourceIdentifiers_objectAndArray(t *testing.T) {
	var rel Relationship
	if err := json.Unmarshal([]byte(`{"data": {"id": "1", "type": "address"}}`), &rel); err != nil {
		t.Fatal(err)
	}
	if len(rel.Data) != 1 || rel.Data[0].Type != "address" {
		t.Errorf("object form = %+v", rel.Data)
	}

	if err := json.Unmarshal([]byte(`{"data": [{"id": "1", "type": "note"}, {"id": "2", "type": "note"}]}`), &rel); err != nil {
		t.Fatal(err)
	}
	if len(rel.Data) != 2 || rel.Data[1].ID != "2" {
		t.Errorf("array form = %+v", rel.Data)
	}

	if err := json.Unmarshal([]byte(`{"data": null}`), &rel); err != nil {
		t.Fatal(err)
	}
	if rel.Data != nil {
		t.Errorf("null form = %+v, want nil", rel.Data)
	}
}

func TestBulkResult_partialFailure(t *testing.T) {
	var result BulkResult
	body := `{
		"success": [88111, 88112],
		"failure": [99999],
		"errors": [{"id": 99999, "code": "ACU0000", "message": "エラーが発生しました。",
			"errors": [{"code": "ACU2001", "message": "customerが見つかりません。"}]}]
	}`
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.Success) != 2 || len(result.Failure) != 1 {
		t.Errorf("success/failure = %v/%v", result.Success, result.Failure)
	}
	if len(result.Errors) != 1 || result.Errors[0].ID != 99999 || result.Errors[0].Errors[0].Code != "ACU2001" {
		t.Errorf("errors = %+v", result.Errors)
	}
}
