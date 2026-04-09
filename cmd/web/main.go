package main

import (
	"flag"
	"net/http"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/DAlba-sudo/pff"
	"github.com/DAlba-sudo/verb"
	"github.com/DAlba-sudo/verb/htmx"
)

// This is what the frontend will modify.
func Frontend(a *pff.App) {

	//PAGES//########//
	a.RegisterTemplate("", "index/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/admin", "admin/admin.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/user", "admin/user.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/login", "login/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/home", "dashboard/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/settings", "settings/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})
	//-------------------//

	//COMPONENTS//########//
	a.RegisterTemplate("/navbar", "components/navbar/navbar.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/products", "components/products/product_card.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	}).RegisterBridge("product", &bridges.ProductBridge{})

	a.RegisterTemplate("/ingredients", "components/ingredients/ingredients_card.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	}).RegisterBridge("ingredient", &bridges.IngredientBridge{})

	comparator := a.RegisterTemplate("/product-comparator", "components/products/product_comparator.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})
	comparator.RegisterBridge("product_comparator", &bridges.ProductComparator{})
	comparator.RegisterBridge("Count", &bridges.ProductCount{})

	a.RegisterTemplate("/searchbar", "components/searchbar/searchbar.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	search_ingredients := a.RegisterTemplate("/component/search/ingredient", "components/searchbar/ingredients.html", pff.TemplateRegistrationOpts{})
	search_ingredients.RegisterBridge("Ingredients", bridges.SearchIngredients{})

	a.RegisterTemplate("/product", "product/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	}).RegisterBridge("Products", bridges.ProductList{})

	product_draft := a.RegisterTemplate("/component/product/draft", "components/product/product.build.html", pff.TemplateRegistrationOpts{})
	product_draft.RegisterBridge("Product", bridges.DraftProduct{})

	a.RegisterTemplate("/component/product/draft/ingredient", "components/product/product.draft.ingredient.html", pff.TemplateRegistrationOpts{}).
		RegisterBridge("Ingredient", bridges.DraftProductIngredient{})

	a.RegisterTemplate("/component/product/recent", "components/product/product.html", pff.TemplateRegistrationOpts{}).
		RegisterBridge("Products", bridges.ProductList{})

	a.RegisterTemplate("/component/search/product", "components/product/search.html", pff.TemplateRegistrationOpts{})

	a.RegisterTemplate("/component/product/search", "components/product/search-res.html", pff.TemplateRegistrationOpts{}).
		RegisterBridge("Products", bridges.ProductSearch{})

	//-------------------//
}

// This is all basic boiler plate, as the frontend you will not have to touch this.
func main() {
	templateDir := flag.String("template", "templates", "the directory for the templates")
	staticDir := flag.String("static", "static", "the directory for the static files")
	address := flag.String("address", "127.0.0.1", "the address to bind the server to")
	port := flag.Int("port", 8080, "the port to bind the server to")
	live := flag.Bool("reload", true, "whether to do live reloads on the template files")
	backend := flag.String("backend", "http://localhost:8081", "the address of the backend server")
	flag.Parse()

	core.BackendAddress = strings.TrimRight(*backend, "/")

	app := pff.CreateApp(pff.Configuration{
		TemplateDirectoryPath: *templateDir,
		FileSystemPath:        *staticDir,
		Address:               *address,
		Port:                  *port,
		Live:                  *live,
	})
	Frontend(app)

	v := verb.New(*address, *port, verb.Settings{
		Templates:  *templateDir,
		Static:     *staticDir,
		LiveReload: *live,
	})

	v.Page("", "v2/pages/index.html").
		Bridge(verb.DataBridge{
			Key: "TotalIngredients",
			Provide: func(r *http.Request) (any, error) {
				return 3, nil
			},
		}).
		Bridge(verb.DataBridge{
			Key: "TotalProducts",
			Provide: func(r *http.Request) (any, error) {
				return 25, nil
			},
		})
	v.Component("v2/components/input-searchbar_generic.html", htmx.Create("div"))
	v.Component("v2/components/ingredient/ingredient-search_results.html", htmx.Create("div")).
		Bridge(verb.DataBridge{
			Key:     "Ingredients",
			Provide: bridges.IngredientSearchProvider,
		}).
		Bridge(verb.DataBridge{
			Key: "Showable",
			Provide: func(r *http.Request) (any, error) {
				var filter struct {
					Ingredient bool
					Product    bool
				}

				core.Logger.Info("Showable Bridge Payload", "Ingredient", r.FormValue("switchIngredient"), "Product", r.FormValue("switchProduct"))

				filter.Ingredient = r.FormValue("switchIngredient") == "on"
				filter.Product = r.FormValue("switchProduct") == "on"

				return filter, nil
			},
		})

	if err := v.Serve(); err != nil {
		panic(err)
	}
}
