package main

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/CFF4HA/Dashboard/internal/handlers/product"
	"github.com/DAlba-sudo/verb"
	"github.com/DAlba-sudo/verb/htmx"
)

func SearchBars(v *verb.Verb) {
	v.Component("v2/searchbars/searchbar-manual_search.html", htmx.Div().Classes("p-2")).
		Bridge(bridges.ManualSearchMultiplexer)
}

func Index(v *verb.Verb) {
	v.Page("", "v2/pages/index.html")
}

func Ingredients(v *verb.Verb) {

}

func Products(v *verb.Verb) {
	// this will instantiate the create product route
	v.Action(http.MethodPut, "/product", product.CreateProduct)

	// this is the components that are rendered as part of the
	// product page.
	v.Component("v2/forms/form-product_create.html", htmx.Div())

	// this will be the primary products page.
	v.Page("/products", "v2/pages/products.html").
		Bridge(bridges.ProductBridge(20, map[string]string{
			"name": " ILIKE ",
		}))
}
