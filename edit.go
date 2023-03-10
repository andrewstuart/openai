package openai

import "context"

const (
	EditModelDavinci1     = "text-davinci-edit-001"
	EditModelDavinciCode1 = "code-davinci-edit-001"
)

// Edit calls the openai edits endpoint with the given input.
// https://platform.openai.com/docs/api-reference/edits
func (c *Client) Edit(ctx context.Context, r EditReq) (*EditRes, error) {
	var res EditRes
	err := c.c.R().Post("edits").JSON(r).Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// EditReq represents an edit on the OpenAI APIs.
type EditReq struct {
	Model       string   `json:"model,omitempty"`
	Input       *string  `json:"input,omitempty"`
	Instruction string   `json:"instruction,omitempty"`
	N           *int     `json:"n,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	TopP        *float64 `json:"top_p,omitempty"`
}

// EditRes contains the results from the edit invocation.
type EditRes struct {
	Object  string     `json:"object"`
	Created int        `json:"created"`
	Choices []EditData `json:"choices"`
	Usage   Usage      `json:"usage"`
}

// EditData holds the specific Choices returned by OpenAI.
type EditData struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}
