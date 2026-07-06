package ecforce

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Query is a ransack-style search query. Keys are an attribute joined with a
// predicate (e.g. "id_eq", "email_cont", "created_at_gteq"), values are the
// search targets. Slice values encode as repeated q[key][] parameters, which
// predicates such as "_in" expect.
//
// See https://activerecord-hackery.github.io/ransack/ for the predicate list.
type Query map[string]any

func (q Query) encode(v url.Values) {
	for key, val := range q {
		name := fmt.Sprintf("q[%s]", key)
		switch vv := val.(type) {
		case []string:
			for _, s := range vv {
				v.Add(name+"[]", s)
			}
		case []int:
			for _, n := range vv {
				v.Add(name+"[]", strconv.Itoa(n))
			}
		case []int64:
			for _, n := range vv {
				v.Add(name+"[]", strconv.FormatInt(n, 10))
			}
		case []any:
			for _, e := range vv {
				v.Add(name+"[]", queryValue(e))
			}
		default:
			v.Add(name, queryValue(val))
		}
	}
}

func queryValue(val any) string {
	switch vv := val.(type) {
	case bool:
		if vv {
			return "1"
		}
		return "0"
	case time.Time:
		return vv.Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprint(val)
	}
}

// ListOptions specifies the optional parameters supported by search (list)
// endpoints.
type ListOptions struct {
	// Lighter toggles the lightweight response variant (fewer attributes,
	// faster). Leave nil to use the shop's default behavior.
	Lighter *bool
	// Include names related data to side-load into Response.Included,
	// e.g. "billing_address".
	Include []string
	// Page is the 1-based page number.
	Page int
	// Per is the number of records per page (max 100).
	Per int
	// Sort lists attributes to order by; prefix with '-' for descending,
	// e.g. []string{"-created_at", "id"}.
	Sort []string
	// Q filters results with ransack-style predicates.
	Q Query
}

// Values encodes the options as URL query parameters. A nil receiver encodes
// to nil.
func (o *ListOptions) Values() url.Values {
	if o == nil {
		return nil
	}
	v := url.Values{}
	if o.Lighter != nil {
		v.Set("lighter", boolFlag(*o.Lighter))
	}
	if len(o.Include) > 0 {
		v.Set("include", strings.Join(o.Include, ","))
	}
	if o.Page > 0 {
		v.Set("page", strconv.Itoa(o.Page))
	}
	if o.Per > 0 {
		v.Set("per", strconv.Itoa(o.Per))
	}
	if len(o.Sort) > 0 {
		v.Set("sort", strings.Join(o.Sort, ","))
	}
	o.Q.encode(v)
	return v
}

// GetOptions specifies the optional parameters supported by detail (get)
// endpoints.
type GetOptions struct {
	// Lighter toggles the lightweight response variant. Leave nil to use the
	// shop's default behavior.
	Lighter *bool
	// Include names related data to side-load into Response.Included.
	Include []string
}

// Values encodes the options as URL query parameters. A nil receiver encodes
// to nil.
func (o *GetOptions) Values() url.Values {
	if o == nil {
		return nil
	}
	v := url.Values{}
	if o.Lighter != nil {
		v.Set("lighter", boolFlag(*o.Lighter))
	}
	if len(o.Include) > 0 {
		v.Set("include", strings.Join(o.Include, ","))
	}
	return v
}

func boolFlag(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
