package ai

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/ollama/ollama/api"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/types"
)

var (
	ProdudctBlurbSystemPrompt = `
	You are given a list of ingredients in a product, alongside a list of known hazards, effects, and symptomps. Your job is 
	to provide factual and actionable summaries based on the information made available to you. You will be given product 
	information expressed using JSON. 

	If you are ever uncertain or are making information up, state it and stop generating a response. 
	
	Do not provide misleading information.

	Your response should always be a valid JSON string that can be unmarshalled, and it should comply with the following schema:
	{
		"blurb": "This product..."
	}

	Keep your response to a single sentence, and make sure it is concise and informative including the types of people that should be 
	concerned.
	`
)

func ProductBlurb(ctx context.Context, product types.Product) error {
	if client == nil {
		return errors.New("ollama client is not initialized")
	} else if client.Heartbeat(ctx) != nil {
		return errors.New("ollama client is not healthy")
	}

	stream := false
	thinking := false

	var BlurbInput struct {
		Labels map[string][]types.Label
	}
	BlurbInput.Labels = make(map[string][]types.Label)

	for _, ingredient := range product.Ingredients {
		var label_list []types.Label
		for _, label := range ingredient.Labels {
			label_list = append(label_list, label)
		}

		BlurbInput.Labels[ingredient.PrimaryName] = label_list
	}
	core.Logger.Info("generating product blurb", "product_id", product.Id, "product_name", product.Name, "blurb_input", BlurbInput)
	data, err := json.Marshal(BlurbInput)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)

	err = client.Generate(ctx, &api.GenerateRequest{
		Model:   "gemma4:e4b",
		Options: map[string]any{"temperature": 0.7, "num_ctx": 8192},
		System:  strings.TrimSpace(ProdudctBlurbSystemPrompt),
		Prompt:  string(data),
		Stream:  &stream,
		Think: &api.ThinkValue{
			Value: &thinking,
		},
		KeepAlive: &api.Duration{time.Duration(10) * time.Second},
		Format:    json.RawMessage(`"json"`),
	}, func(gr api.GenerateResponse) error {
		var response struct {
			Blurb string `json:"blurb"`
		}
		defer wg.Done()

		if err := json.Unmarshal([]byte(gr.Response), &response); err != nil {
			return err
		}

		core.Logger.Info("generated product blurb", "blurb", response.Blurb)
		return core.DB.Model(&product.Metadata).Update("AiBlurb", response.Blurb).Error
	})
	if err != nil {
		core.Logger.Error("error generating product blurb", "error", err)
		return err
	}

	wg.Wait()
	return nil
}
