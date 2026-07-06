package ecforce

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Error is a single API error with an ecforce error code (e.g. "ACU2001") and
// a human-readable Japanese message.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// ErrorResponse reports errors returned by an API request with a non-2xx
// status code. The body has the form {"errors": [{"code", "message"}, ...]}.
type ErrorResponse struct {
	Response *http.Response `json:"-"`
	Errors   []Error        `json:"errors"`
}

func (r *ErrorResponse) Error() string {
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	detail := strings.Join(msgs, "; ")
	if detail == "" {
		detail = "unexpected error"
	}
	if r.Response == nil {
		return fmt.Sprintf("ecforce: %s", detail)
	}
	return fmt.Sprintf("ecforce: %v %v: %d %s",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, detail)
}

// HasCode reports whether the response contains an error with the given code.
func (r *ErrorResponse) HasCode(code string) bool {
	for _, e := range r.Errors {
		if e.Code == code {
			return true
		}
	}
	return false
}

// CheckResponse checks the API response for an error status and, if found,
// returns it as an *ErrorResponse. A response is considered an error if its
// status code is outside the 2xx range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errResp := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		// A non-JSON body (e.g. an HTML error page) still yields a usable
		// ErrorResponse keyed on the status code.
		_ = json.Unmarshal(data, errResp)
	}
	return errResp
}

// BulkError describes why one record of a bulk operation failed. ID refers to
// the failed record: for bulk updates it is the record's ID, for bulk creates
// it is the 1-based position in the request array.
type BulkError struct {
	ID      int64   `json:"id"`
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Errors  []Error `json:"errors"`
}

// BulkResult is the outcome of a synchronous bulk operation. Success and
// Failure identify processed records: record IDs for updates/destroys, or
// 1-based positions in the request array for creates.
type BulkResult struct {
	Success []int64     `json:"success"`
	Failure []int64     `json:"failure"`
	Errors  []BulkError `json:"errors"`
}

// JobResult identifies an asynchronous bulk job accepted by the API.
type JobResult struct {
	ID    int64  `json:"id"`
	JobID string `json:"job_id"`
}
