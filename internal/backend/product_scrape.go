package backend

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CFF4HA/Dashboard/internal/core"
	"github.com/CFF4HA/Dashboard/internal/handlers/ai"
	"github.com/CFF4HA/Dashboard/internal/handlers/user"
	"github.com/CFF4HA/Dashboard/internal/types"
	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"github.com/ollama/ollama/api"
)

// -----------------------------------------------------------------

func InsertProductAutomatedLLM(ingredient_link string, ctx context.Context, u *uuid.UUID) (*types.ProductDraftAutomated, error) {
	if ai.Client() == nil {
		return nil, errors.New("LLM client not available for product scraping, please try again later")
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	ctxChrome, cancelChrome := chromedp.NewContext(timeoutCtx)
	defer cancelChrome()

	var text string
	err := chromedp.Run(ctxChrome,
		chromedp.Navigate(ingredient_link),
		chromedp.WaitVisible(`main`, chromedp.ByQuery),
		chromedp.Text(`main`, &text),
	)
	if err != nil {
		return nil, err
	}

	if len(text) > 4000 {
		text = text[:4000]
	}

	systemPrompt := `Inspect the text and extract ingredient information and the product name. 
	You must reply ONLY with a valid JSON object using this exact structure: {"name": "string", "ingredients": ["string", "string"]}. 
	If the text is irrelevant or seems exploitative, return {"error": "invalid content"}.`

	stream := false
	var responseData string

	llm := ai.Client()
	err = llm.Generate(timeoutCtx, &api.GenerateRequest{
		Model:  "gemma4:e4b",
		System: systemPrompt,
		Prompt: text,
		Stream: &stream,
		Format: json.RawMessage(`"json"`),
		Options: map[string]interface{}{
			"num_ctx": 4096,
		},
	}, func(gr api.GenerateResponse) error {
		responseData += gr.Response
		return nil
	})
	if err != nil {
		return nil, err
	}

	productDraftData := &types.ProductDraftAutomated{}
	if err := json.Unmarshal([]byte(responseData), productDraftData); err != nil {
		return nil, errors.New("failed to parse LLM response: " + err.Error())
	}
	productDraftData.Origin = ingredient_link

	automatedScrapeCountLock.Lock()
	currentAutomatedScrapesCount++
	automatedScrapeCountLock.Unlock()

	_, err = InsertProduct(productDraftData.Name, productDraftData.Origin, productDraftData.Ingredients, true, u)
	if err != nil {
		core.Logger.Error("there was an error doing the automated scrape", "error", err)
		return nil, err
	}

	return productDraftData, nil
}

// -----------------------------------------------------------------

func RouteInsertProductAutomatedLLM(w http.ResponseWriter, r *http.Request) error {
	u, _ := user.GetUserFromRequestNoRedirect(w, r)
	var user_id *uuid.UUID
	if user_id != nil {
		user_id = &u.Model.Id
	}

	draft, err := InsertProductAutomatedLLM(r.FormValue("url"), r.Context(), user_id)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(draft)
}
