package admin

import (
	"github.com/zero-color/ecforce-go"
)

// Free column (自由項目) values are attached to customers, orders and
// subscription orders through identically shaped endpoints; the types here
// are shared by the three services.

// FreeColumnParams is one entry of a free column value bulk create/update:
// either a single column value (FreeColumnID and a value), or a category of
// values (FreeColumnCategoryID and Values).
type FreeColumnParams struct {
	FreeColumnID        *int64  `json:"free_column_id,omitempty"`
	FreeColumnValue     *string `json:"free_column_value,omitempty"`
	FreeColumnOptionIDs []int64 `json:"free_column_option_ids,omitempty"`

	FreeColumnCategoryID *int64                   `json:"free_column_category_id,omitempty"`
	Values               []*FreeColumnValueParams `json:"values,omitempty"`

	// FreeColumnCategorySeq targets a specific answer sequence within a
	// category on bulk updates.
	FreeColumnCategorySeq *int64 `json:"free_column_category_seq,omitempty"`
	// Delete removes the targeted category answer on bulk updates.
	Delete *ecforce.BoolInt `json:"delete,omitempty"`
}

// FreeColumnValueParams is a single column value inside a category entry of
// FreeColumnParams.
type FreeColumnValueParams struct {
	FreeColumnID        *int64  `json:"free_column_id,omitempty"`
	FreeColumnValue     *string `json:"free_column_value,omitempty"`
	FreeColumnOptionIDs []int64 `json:"free_column_option_ids,omitempty"`
}

// FreeColumnValueGroup is one element of a free column value search response:
// values grouped by (optional) category.
type FreeColumnValueGroup struct {
	FreeColumnCategoryID int64                   `json:"free_column_category_id"`
	Values               []*FreeColumnValueEntry `json:"values"`
}

// FreeColumnValueEntry groups values by their sequence within a category.
type FreeColumnValueEntry struct {
	FreeColumnCategorySeq int                `json:"free_column_category_seq"`
	Values                []*FreeColumnValue `json:"values"`
}

// FreeColumnValue is a stored free column value.
type FreeColumnValue struct {
	FreeColumnID          int64    `json:"free_column_id"`
	FreeColumnLabel       string   `json:"free_column_label"`
	FreeColumnValueValues []string `json:"free_column_value_values"`
}
