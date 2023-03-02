package openai

import "context"

func (c *Client) Edit(ctx context.Context, r EditReq) (*EditRes, error) {
	var res EditRes
	err := c.c.R().Post("edits").JSON(r).Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type EditReq struct {
	Model       string   `json:"model,omitempty"`
	Input       *string  `json:"input,omitempty"`
	Instruction string   `json:"instruction,omitempty"`
	N           *int     `json:"n,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	TopP        *float64 `json:"top_p,omitempty"`
}

type EditRes struct {
	Object  string     `json:"object"`
	Created int        `json:"created"`
	Choices []EditData `json:"choices"`
	Usage   Usage      `json:"usage"`
}

type EditData struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}
