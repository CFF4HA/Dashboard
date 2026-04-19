package ai

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

var (
	client *api.Client
)

func Initialize(uri string) error {
	endpoint, err := url.Parse(uri)
	if err != nil {
		return err
	}

	client = api.NewClient(endpoint, http.DefaultClient)
	return client.Heartbeat(context.Background())

}
