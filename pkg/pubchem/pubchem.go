package pubchem

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	pugBase     = "https://pubchem.ncbi.nlm.nih.gov/rest/pug"
	pugViewBase = "https://pubchem.ncbi.nlm.nih.gov/rest/pug_view"
)

// Client makes requests to the PubChem REST API.
type Client struct {
	http   *http.Client
	logger *slog.Logger
}

// New returns a Client using the given HTTP client and logger.
// Pass nil for either to use http.DefaultClient or no logging, respectively.
func New(httpClient *http.Client, logger *slog.Logger) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{http: httpClient, logger: logger}
}

// Default is a package-level Client with no logging and http.DefaultClient.
var Default = New(nil, nil)

type cidResponse struct {
	IdentifierList struct {
		CID []int `json:"CID"`
	} `json:"IdentifierList"`
}

type sidResponse struct {
	IdentifierList struct {
		SID []int `json:"SID"`
	} `json:"IdentifierList"`
}

func (c *Client) get(rawURL string) (*http.Response, error) {
	if c.logger != nil {
		c.logger.Debug("pubchem request", "url", rawURL)
	}
	resp, err := c.http.Get(rawURL)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("pubchem request failed", "url", rawURL, "err", err)
		}
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		if c.logger != nil {
			c.logger.Warn("pubchem non-2xx response", "url", rawURL, "status", resp.StatusCode)
		}
		return nil, fmt.Errorf("pubchem: status %d for %s", resp.StatusCode, rawURL)
	}
	return resp, nil
}

// GetCompoundId returns the first CID matching the given compound name.
func (c *Client) GetCompoundId(name string) (int, error) {
	u := fmt.Sprintf("%s/compound/name/%s/cids/JSON", pugBase, url.PathEscape(name))
	resp, err := c.get(u)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result cidResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	if len(result.IdentifierList.CID) == 0 {
		return 0, fmt.Errorf("pubchem: no compound found for %q", name)
	}
	return result.IdentifierList.CID[0], nil
}

// GetSubstanceId returns the first SID matching the given substance name.
func (c *Client) GetSubstanceId(name string) (int, error) {
	u := fmt.Sprintf("%s/substance/name/%s/sids/JSON", pugBase, url.PathEscape(name))
	resp, err := c.get(u)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result sidResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	if len(result.IdentifierList.SID) == 0 {
		return 0, fmt.Errorf("pubchem: no substance found for %q", name)
	}
	return result.IdentifierList.SID[0], nil
}

// GetCompoundAsJSON returns the full PubChem view data for a compound by CID.
func (c *Client) GetCompoundAsJSON(id int) (json.RawMessage, error) {
	u := fmt.Sprintf("%s/data/compound/%d/JSON", pugViewBase, id)
	resp, err := c.get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}

// GetSubstanceAsJSON returns the full PubChem view data for a substance by SID.
func (c *Client) GetSubstanceAsJSON(id int) (json.RawMessage, error) {
	u := fmt.Sprintf("%s/data/substance/%d/JSON", pugViewBase, id)
	resp, err := c.get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}

// Package-level convenience functions using Default.

func GetCompoundId(name string) (int, error)             { return Default.GetCompoundId(name) }
func GetSubstanceId(name string) (int, error)            { return Default.GetSubstanceId(name) }
func GetCompoundAsJSON(id int) (json.RawMessage, error)  { return Default.GetCompoundAsJSON(id) }
func GetSubstanceAsJSON(id int) (json.RawMessage, error) { return Default.GetSubstanceAsJSON(id) }
