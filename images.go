package openai

import "context"

func (c *Client) GenerateImage(ctx context.Context, p ImgPrompt) (*ImageRes, error) {
	var res ImageRes
	err := c.c.R().Post("images/generations").JSON(p).Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type ImgPrompt struct {
	Prompt         string  `json:"prompt,omitempty"`
	N              *int    `json:"n,omitempty"`
	Size           *string `json:"size,omitempty"`
	ResponseFormat *string `json:"response_format,omitempty"`
	User           *string `json:"user,omitempty"`
}

type ImageRes struct {
	Created int             `json:"created"`
	Dataa   []ImageResDatum `json:"data"`
}

type ImageResDatum struct {
	URL     string `json:"url"`
	B64JSON string `json:"b64_json"`
}
