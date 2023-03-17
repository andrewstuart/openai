package openai

import "context"

// ModelList is a simple convenience abstraction over a slice of ModelsRes.
type ModelList []ModelsRes

// Has returns whether or not the model ID is present in the list.
func (ml ModelList) Has(s string) bool {
	for _, m := range ml {
		if m.ID == s {
			return true
		}
	}
	return false
}

// Models calls the openai models endpoint and returns the results, as a slice
// but with some convenience methods on a slice type.
// https://platform.openai.com/docs/api-reference/models
func (c *Client) Models(ctx context.Context) (ModelList, error) {
	var res struct {
		Data []ModelsRes
	}
	err := c.c.R().Get("models").Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}

// ModelsRes returns the models currently available to the user, per the models
// API on OpenAI.
type ModelsRes struct {
	Created     int          `json:"created"`
	ID          string       `json:"id"`
	Object      string       `json:"object"`
	OwnedBy     string       `json:"owned_by"`
	Parent      *string      `json:"parent"`
	Permissions []Permission `json:"permission"`
	Root        string       `json:"root"`
}

// Permissions are the return response from OpenAI
type Permission struct {
	AllowCreateEngine  bool    `json:"allow_create_engine"`
	AllowFineTuning    bool    `json:"allow_fine_tuning"`
	AllowLogprobs      bool    `json:"allow_logprobs"`
	AllowSampling      bool    `json:"allow_sampling"`
	AllowSearchIndices bool    `json:"allow_search_indices"`
	AllowView          bool    `json:"allow_view"`
	Created            int     `json:"created"`
	Group              *string `json:"group"`
	ID                 string  `json:"id"`
	IsBlocking         bool    `json:"is_blocking"`
	Object             string  `json:"object"`
	Organization       string  `json:"organization"`
}
