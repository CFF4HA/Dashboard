package main

import (
	"flag"

	"github.com/DAlba-sudo/pff"
)

// This is what the frontend will modify.
func Frontend(a *pff.App) {
	//pages
	a.RegisterTemplate("", "index/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/admin", "admin/admin.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/user", "admin/user.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	//components
	a.RegisterTemplate("/navbar", "components/navbar/navbar.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/products", "components/products/product_card.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/ingredients", "components/ingredients/ingredients_card.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

}

// This is all basic boiler plate, as the frontend you will not have to touch this.
func main() {
	templateDir := flag.String("template", "templates", "the directory for the templates")
	staticDir := flag.String("static", "static", "the directory for the static files")
	address := flag.String("address", "127.0.0.1", "the address to bind the server to")
	port := flag.Int("port", 8080, "the port to bind the server to")
	live := flag.Bool("reload", true, "whether to do live reloads on the template files")
	flag.Parse()

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
