package admin

import (
	"context"
	"net/http"

	"github.com/zero-color/ecforce-go"
)

// ProductCategory is a product category (商品カテゴリー). It is also
// side-loaded via include on the product endpoints, which additionally fill
// Position.
type ProductCategory struct {
	ID                        int64         `json:"id"`
	Name                      string        `json:"name"`
	ParentName                string        `json:"parent_name"`
	Path                      string        `json:"path"`
	Slug                      string        `json:"slug"`
	ProductCategoryVisibility bool          `json:"product_category_visibility"`
	Description               string        `json:"description"`
	ProductExtraText          string        `json:"product_extra_text"`
	Position                  int           `json:"position"`
	CreatedAt                 *ecforce.Time `json:"created_at"`
	UpdatedAt                 *ecforce.Time `json:"updated_at"`
}

// List searches product categories.
//
// Supported q attributes: id, name, slug, product_category_visibility,
// created_at, updated_at.
// Supported sort attributes: id, created_at, updated_at, position.
//
// ecforce API docs: GET /api/v2/admin/product_categories
func (s *ProductCategoriesService) List(ctx context.Context, opts *ecforce.ListOptions) ([]*ecforce.Resource[ProductCategory], *ecforce.Response, error) {
	return ecforce.DoResourceList[ProductCategory](ctx, s.client, http.MethodGet, "admin/product_categories.json", opts.Values(), nil)
}
