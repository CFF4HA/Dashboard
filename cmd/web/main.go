package main

import (
	"flag"

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

	a.RegisterTemplate("/home", "dashboard/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/settings", "settings/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/admin", "admin/page.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: true,
	})

	a.RegisterTemplate("/overviewtable", "components/ingredientcard/overviewtable.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/keyinsights", "components/ingredientcard/keyinsights.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/risksummary", "components/ingredientcard/risksummary.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/personalproducts", "components/myproductscard/personalproducts.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/communityproducts", "components/myproductscard/communityproducts.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/saveproduct", "components/myproductscard/saveproduct.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/product1", "components/productcomparison/product1.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/product2", "components/productcomparison/product2.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/totalingredients", "components/productcomparison/totalingredients.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/lowriskcount", "components/productcomparison/lowriskcount.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false,
	})

	a.RegisterTemplate("/moderateriskcount", "components/productcomparison/moderateriskcount.html", pff.TemplateRegistrationOpts{
		IncludeBaseTemplate: false, 
	})

	a.RegisterTemplate("/highriskcount", "components/productcomparison/highriskcount.html", pff.TemplateRegistrationOpts{
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
