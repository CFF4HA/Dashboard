package main

import (
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/handlers/product"
	"github.com/DAlba-sudo/verb"
)

func Index(v *verb.Verb) {

}

func Ingredients(v *verb.Verb) {

}

func Products(v *verb.Verb) {
	// this will instantiate the create product route
	v.Action(http.MethodPut, "/product", product.CreateProduct)
}
