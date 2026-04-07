package main

import (
	"flag"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/DAlba-sudo/pff"
)

// This is what the frontend will modify.
func Frontend(a *pff.App) {
	a.RegisterTemplate("", "index/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/login", "login/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/navbar", "components/navbar/navbar.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/searchbar", "components/searchbar/searchbar.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	search_ingredients := a.RegisterTemplate("/component/search/ingredient", "components/searchbar/ingredients.html", pff.TemplateRegistrationOpts{})
	search_ingredients.RegisterBridge("Ingredients", bridges.SearchIngredients{})

	a.RegisterTemplate("/home", "dashboard/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/settings", "settings/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/admin", "admin/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/product", "product/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	product_draft := a.RegisterTemplate("/component/product/draft", "components/product/product.build.html", pff.TemplateRegistrationOpts{})
	product_draft.RegisterBridge("Product", bridges.DraftProduct{})

	a.RegisterTemplate("/component/product/draft/ingredient", "components/product/product.draft.ingredient.html", pff.TemplateRegistrationOpts{}).
		RegisterBridge("Ingredient", bridges.DraftProductIngredient{})

	a.RegisterTemplate("/component/search/product", "components/product/search.html", pff.TemplateRegistrationOpts{})
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

	if err := app.Start(); err != nil {
		panic(err)
	}
}
