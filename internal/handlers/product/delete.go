package product

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.FormValue("id")

	// Explicitly remove join table rows before deleting the product.
	core.DB.Exec("DELETE FROM product_ingredients WHERE product_id = ?", id)
	core.DB.Exec("DELETE FROM product_tags WHERE product_id = ?", id)

	// ProductMetadata cascades via the FK constraint defined on the struct.
	return core.DB.Unscoped().Delete(&types.Product{}, "id = ?", id).Error
}
