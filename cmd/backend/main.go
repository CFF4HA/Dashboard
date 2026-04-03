package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/backend/http/routes"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/DAlba-sudo/pbf"
)

func main() {
	address := flag.String("address", "localhost", "the address to listen on")
	port := flag.Int("port", 8080, "the port to listen on")
	db := flag.String("db", "postgres://user:password@localhost:5432/database?sslmode=disable", "the database connection string")
	flag.Parse()

	r := pbf.CreateRouter()

	// Initialiation is done here
	r.Address = *address
	r.Port = *port
	database.Initialize(*db)

	// Adding the routes here.

	// This route will be used to get an ingredient. If it doesn't exist it will
	// tap into the python server with the relevant information to update
	// the database.
	r.Add(pbf.RouteOptions{
		Method:   http.MethodGet,
		Endpoint: "/ingredient",
		Handler:  routes.IngredientGET,
	})

	// Base Product Routes
	r.Add(pbf.RouteOptions{
		Method:   http.MethodGet,
		Endpoint: "/product",
		Handler:  routes.ProductGET,
	})

	r.Add(pbf.RouteOptions{
		Method:   http.MethodPut,
		Endpoint: "/product",
		Handler:  routes.ProductPUT,
	})

	r.Add(pbf.RouteOptions{
		Method:   http.MethodPatch,
		Endpoint: "/product",
		Handler:  routes.ProductPATCH,
	})

	r.Add(pbf.RouteOptions{
		Method:   http.MethodDelete,
		Endpoint: "/product",
		Handler:  routes.ProductDELETE,
	})

	// Starting!
	core.Logger.Info("the backend is starting", "address", r.Address, "port", r.Port, "link", fmt.Sprintf("http://%s:%d", r.Address, r.Port))
	if err := r.Start(); err != nil {
		panic(err)
	}
}
