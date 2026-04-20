package main

import (
	"reflect"
	"strings"

	"github.com/CFF4HA/Dashboard/internal/bridges"
	"github.com/CFF4HA/Dashboard/internal/core"
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
			Bridges: []verb.Bridge{
				bridges.UserSessionRequired{
					Excludes: []string{"/login", "/user", "/user/login", "/htmx/form-user_signup"},
				},
			},
		},
	)

	v.Func("exists", func(m interface{}, key interface{}) bool {
		if m == nil {
			return false
		}
		rv := reflect.ValueOf(m)
		if rv.Kind() != reflect.Map {
			return false
		}
		return rv.MapIndex(reflect.ValueOf(key)).IsValid()
	})
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

	// this is the llm configuration portion
	if err := LlmConfig(); err != nil {
		core.Logger.Warn("ollama configuration error, llm features will be disabled", "error", err)
	} else {
		core.Logger.Info("ollama client configured successfully")
	}

	// The following are the general use component routes
	SearchBars(v)
	Ingredients(v)
	Products(v)
	User(v)

	// The following are the general use page creation routes
	Index(v)

	if err := v.Serve(); err != nil {
		panic(err)
	}
}
