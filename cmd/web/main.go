package main

import (
	"github.com/CFF4HA/Dashboard/internal/backend/database"
	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/DAlba-sudo/verb"
	"github.com/DAlba-sudo/verbs"
)

func main() {
	cfg := parseConfig()

	core.BackendAddress = cfg.BackendAddr

	if err := database.Initialize(cfg.DatabaseConn); err != nil {
		core.Logger.Error("failed to initialize database", "error", err)
	}

	dbconn, err := database.Database()
	if err != nil {
		core.Logger.Error("failed to get database connection", "error", err)
	}

	v := verb.New(cfg.Address, cfg.Port, verb.Settings{
		Templates:  cfg.TemplateDir,
		Static:     cfg.StaticDir,
		LiveReload: cfg.LiveReload,
		Bridges: []verb.Bridge{
			verbs.Request{},
		},
	})

	registerFuncs(v)
	registerRoutes(v, cfg, dbconn)

	if err := v.Serve(); err != nil {
		panic(err)
	}
}
