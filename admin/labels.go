package admin

import (
	"context"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// Label is a label (ラベル) that can be attached to orders, subscription
// orders, customers, products, advertising URLs and LP templates. The label
// API is not JSON:API; labels are returned as plain objects.
type Label struct {
	ID                   int64  `json:"id"`
	Name                 string `json:"name"`
	Color                string `json:"color"`
	Position             int    `json:"position"`
	SearchFormVisibility bool   `json:"search_form_visibility"`
	Type                 string `json:"type"`
}

// LabelCreateRequest is the body of a LabelsService.Create request. Type is
// the label type (the API currently allows "Customer"), Name the label name,
// and Color an optional color code such as "#FFFFFF".
type LabelCreateRequest struct {
	Type  *string `json:"type,omitempty"`
	Name  *string `json:"name,omitempty"`
	Color *string `json:"color,omitempty"`
}

// LabelingRequest is the body of a LabelsService.Apply request.
type LabelingRequest struct {
	// Type is the record type to label: Order, SubsOrder, Customer, Product,
	// Url or Template.
	Type *string `json:"type,omitempty"`
	// Method is the operation: overwrite, add or remove.
	Method *string `json:"method,omitempty"`
	// TargetIDs lists the records to operate on; nonexistent IDs are skipped.
	TargetIDs []int64 `json:"target_ids"`
	// LabelIDs lists the labels to operate with. With Method "overwrite", an
	// empty (non-nil) slice detaches all labels from the targets.
	LabelIDs []int64 `json:"label_ids"`
}

// Apply attaches, detaches or overwrites labels on the given records in bulk
// (synchronous).
//
// ecforce API docs: PUT /api/v2/admin/labeling
func (s *LabelsService) Apply(ctx context.Context, req *LabelingRequest) (*ecforce.BulkResult, *ecforce.Response, error) {
	result := new(ecforce.BulkResult)
	resp, err := ecforce.Do(ctx, s.client, http.MethodPut, "admin/labeling.json", nil, req, result)
	if err != nil {
		return nil, resp, err
	}
	return result, resp, nil
}

// Create creates a label.
//
// ecforce API docs: POST /api/v2/admin/labels
func (s *LabelsService) Create(ctx context.Context, req *LabelCreateRequest) (*Label, *ecforce.Response, error) {
	var doc struct {
		Data *Label `json:"data"`
	}
	resp, err := ecforce.Do(ctx, s.client, http.MethodPost, "admin/labels.json", nil, req, &doc)
	if err != nil {
		return nil, resp, err
	}
	return doc.Data, resp, nil
}

// List searches labels.
//
// Supported q attributes: id, name, type.
//
// ecforce API docs: GET /api/v2/admin/labels
func (s *LabelsService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*Label, *ecforce.Response, error) {
	var doc struct {
		Data  []*Label       `json:"data"`
		Meta  *ecforce.Meta  `json:"meta"`
		Links *ecforce.Links `json:"links"`
	}
	resp, err := ecforce.Do(ctx, s.client, http.MethodGet, "admin/labels.json", opts.Values(), nil, &doc)
	if err != nil {
		return nil, resp, err
	}
	resp.Meta, resp.Links = doc.Meta, doc.Links
	return doc.Data, resp, nil
}
