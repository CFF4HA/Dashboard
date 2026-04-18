package main

import (
	"github.com/DAlba-sudo/verb"
)

func main() {
	RegisterConfigFlags()

	v := verb.New(
		Cfg.Address,
		Cfg.Port,
		verb.Settings{
			Templates:  Cfg.TemplateDir,
			Static:     Cfg.StaticDir,
			LiveReload: Cfg.Reload,
			Bridges:    []verb.Bridge{},
		},
	)

	// this is some further configuration that we have to do with regards to
	// database, etc.

	// this registers all components, actions, and pages related to rendering
	// the index page.
	Index(v)

	if err := v.Serve(); err != nil {
		panic(err)
	}
}
