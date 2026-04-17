package main

import "flag"

// Config holds all application configuration parsed from command-line flags.
type Config struct {
	TemplateDir    string
	StaticDir      string
	Address        string
	Port           int
	LiveReload     bool
	BackendAddr    string
	LLMAddr        string
	DatabaseConn   string
	CFClientID     string
	CFClientSecret string
}

// parseConfig reads command-line flags and returns a populated Config.
func parseConfig() Config {
	var cfg Config
	flag.StringVar(&cfg.TemplateDir, "template", "templates", "directory for template files")
	flag.StringVar(&cfg.StaticDir, "static", "static", "directory for static files")
	flag.StringVar(&cfg.Address, "address", "127.0.0.1", "address to bind the server to")
	flag.IntVar(&cfg.Port, "port", 8080, "port to bind the server to")
	flag.BoolVar(&cfg.LiveReload, "reload", true, "enable live template reloading")
	flag.StringVar(&cfg.BackendAddr, "backend", "http://localhost:8081", "backend service address")
	flag.StringVar(&cfg.LLMAddr, "llm", "http://localhost:8082", "LLM service address")
	flag.StringVar(&cfg.DatabaseConn, "db", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "PostgreSQL connection string")
	flag.StringVar(&cfg.CFClientID, "cf_client_id", "", "Cloudflare Access client ID")
	flag.StringVar(&cfg.CFClientSecret, "cf_client_secret", "", "Cloudflare Access client secret")
	flag.Parse()
	return cfg
}
