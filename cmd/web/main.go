package main

import (
	"github.com/DAlba-sudo/verb"
	"strings"
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

	v.Func("deref", func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	})
	v.Func("sub", func(a, b int) int { return a - b })
	v.Func("lower", func(a string) string { return strings.ToLower(a) })

	// this is some further configuration that we have to do with regards to
	// database, etc.
	if err := DatabaseConfig(); err != nil {
		// do potential recover here if we want to save to to a local
		// storage?
		panic(err)
	}

	// The following are the general use component routes
	SearchBars(v)
	Ingredients(v)
	Products(v)

	// The following are the general use page creation routes
	Index(v)

	if err := v.Serve(); err != nil {
		panic(err)
	}
}
