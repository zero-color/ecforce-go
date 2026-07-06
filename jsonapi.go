package ecforce

import (
	"encoding/json"
	"fmt"
)

// Resource is a JSON:API resource object: an ID/type pair with typed
// attributes and untyped relationship linkage.
type Resource[T any] struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    T             `json:"attributes"`
	Relationships Relationships `json:"relationships,omitempty"`
}

// Relationships maps relationship names (e.g. "billing_address", "orders") to
// the identifiers of related resources. The referenced resources themselves
// are side-loaded into Response.Included when requested via the include
// parameter.
type Relationships map[string]*Relationship

// Relationship holds the resource identifiers of one relationship. Its data
// may be a single object or an array in the wire format; both decode into the
// Data slice.
type Relationship struct {
	Data ResourceIdentifiers `json:"data"`
}

// ResourceIdentifier identifies a resource by ID and type.
type ResourceIdentifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// ResourceIdentifiers decodes from either a single JSON:API resource
// identifier object, an array of them, or null.
type ResourceIdentifiers []ResourceIdentifier

// UnmarshalJSON implements json.Unmarshaler.
func (r *ResourceIdentifiers) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*r = nil
		return nil
	}
	if data[0] == '[' {
		var ids []ResourceIdentifier
		if err := json.Unmarshal(data, &ids); err != nil {
			return err
		}
		*r = ids
		return nil
	}
	var id ResourceIdentifier
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	*r = ResourceIdentifiers{id}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (r ResourceIdentifiers) MarshalJSON() ([]byte, error) {
	if len(r) == 1 {
		return json.Marshal(r[0])
	}
	return json.Marshal([]ResourceIdentifier(r))
}

// IncludedResource is a side-loaded resource from a document's "included"
// array. Its attributes are kept raw; decode them into a typed struct with
// DecodeAttributes once you have dispatched on Type.
type IncludedResource struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Attributes    json.RawMessage `json:"attributes"`
	Relationships Relationships   `json:"relationships,omitempty"`
}

// DecodeAttributes decodes the raw attributes into v.
func (i *IncludedResource) DecodeAttributes(v any) error {
	if len(i.Attributes) == 0 {
		return fmt.Errorf("ecforce: included resource %s/%s has no attributes", i.Type, i.ID)
	}
	return json.Unmarshal(i.Attributes, v)
}

// FindIncluded returns the resource with the given type and ID from included,
// or nil if absent.
func FindIncluded(included []*IncludedResource, typ, id string) *IncludedResource {
	for _, inc := range included {
		if inc.Type == typ && inc.ID == id {
			return inc
		}
	}
	return nil
}

// Meta holds the pagination metadata of list responses.
type Meta struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	Per        int `json:"per"`
	Count      int `json:"count"`
	TotalPages int `json:"total_pages"`
}

// Links holds the pagination links of list responses.
type Links struct {
	Self  string `json:"self"`
	Prev  string `json:"prev"`
	First string `json:"first"`
	Next  string `json:"next"`
	Last  string `json:"last"`
}
