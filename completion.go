package openai

import "context"

func (c Client) Complete(ctx context.Context, req CompleteReq) (*CompleteRes, error) {
	var res CompleteRes
	err := c.c.R().Post("completions").JSON(req).Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type CompleteReq struct {
	Model             string          `json:"model"`
	Prompt            string          `json:"prompt"`
	Suffix            *string         `json:"suffix,omitempty"`
	MaxTokens         *int            `json:"max_tokens,omitempty"`
	Temperature       *float64        `json:"temperature,omitempty"`
	TopP              *int            `json:"top_p,omitempty"`
	N                 *int            `json:"n,omitempty"`
	Logprobs          *int            `json:"logprobs,omitempty"`
	Echo              *bool           `json:"echo,omitempty"`
	Stop              *string         `json:"stop,omitempty"`
	PresencePentalty  *float64        `json:"presence_pentalty,omitempty"`
	FrequencyPentalty *float64        `json:"frequency_pentalty,omitempty"`
	BestOf            *int            `json:"best_of,omitempty"`
	LogitBias         *map[string]any `json:"logit_bias,omitempty"`
	User              *string         `json:"user,omitempty"`
}

func p[T any](t T) *T {
	return &t
}

type CompleteRes struct {
	Choices []CompleteChoice `json:"choices"`
	Created int              `json:"created"`
	ID      string           `json:"id"`
	Model   string           `json:"model"`
	Object  string           `json:"object"`
	Usage   Usage            `json:"usage"`
}

type CompleteChoice struct {
	Text         string
	Index        int
	LogProbs     any
	FinishReason string
}
