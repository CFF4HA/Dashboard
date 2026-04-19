package main

import (
	"errors"
	"flag"
	"log/slog"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/ai"
)

var (
	Cfg Config
)

type Config struct {
	// the address and port to bind to for the web server
	Address     string
	Port        int
	TemplateDir string
	StaticDir   string
	Reload      bool

	// the ollama information required for the different pieces like
	// routing, generation, etc.
	Ollama struct {
		HttpAddress string
	}

	// the database information rquired for the web app to function.
	Database struct {
		URL string
	}

	// the following is app specific misc configuartion like debug levels, etc.
	LogLevel int
}

func RegisterConfigFlags() {
	flag.StringVar(&Cfg.Address, "address", "127.0.0.1", "the itnerface address to bind to")
	flag.IntVar(&Cfg.Port, "port", 8080, "the port to bind to")
	flag.StringVar(&Cfg.TemplateDir, "template-dir", "./templates", "the directory where the html templates are located")
	flag.StringVar(&Cfg.StaticDir, "static-dir", "./static", "the directory where the static files are located")
	flag.BoolVar(&Cfg.Reload, "reload", false, "whether to enable live reloading of templates and static files")
	flag.StringVar(&Cfg.Ollama.HttpAddress, "ollama-http-address", "http://localhost:11434", "the address of the ollama http server")
	flag.StringVar(&Cfg.Database.URL, "database-url", "postgres://user:password@localhost:5432/dbname?sslmode=disable", "the database connection URL")
	flag.IntVar(&Cfg.LogLevel, "log-level", int(slog.LevelInfo), "the log level for the application (0=debug, 1=info, 2=warn, 3=error)")
	flag.Parse()
}

func DatabaseConfig() error {
	if Cfg.Database.URL == "" {
		return errors.New("no database url is provided")
	}
	return core.InitializeDatabase(Cfg.Database.URL)
}

func LlmConfig() error {
	if Cfg.Ollama.HttpAddress == "" {
		return nil
	}

	return ai.Initialize(Cfg.Ollama.HttpAddress)
}
