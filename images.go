package openai

import "context"

// Well-known responseformats
const (
	ImageResponseFormatURL     = "url"
	ImageResponseFormatB64JSON = "b64_json"

	ImageSize256  = "256x256"
	ImageSize512  = "512x512"
	ImageSize1024 = "1024x1024"
)

// GenerateImage calls the images/generations API and returns the API response
// or any error.
func (c *Client) GenerateImage(ctx context.Context, p ImgReq) (*ImageRes, error) {
	var res ImageRes
	err := c.c.R().Post("images/generations").JSON(p).Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	for i := range res.Data {
		res.Data[i].c = c
	}
	return &res, nil
}

// ImgReq is the input type for image generation.
type ImgReq struct {
	Prompt         string  `json:"prompt,omitempty"`
	N              *int    `json:"n,omitempty"`
	Size           *string `json:"size,omitempty"`
	ResponseFormat *string `json:"response_format,omitempty"`
	User           *string `json:"user,omitempty"`
}

// ImageRes is returned by the Image generation.
type ImageRes struct {
	Created int      `json:"created"`
	Data    []Images `json:"data"`
}

// Images are returned as apart of the Image generation calls, and will contain
// either a URL to the image generated, or the bytes as `.Image` if the input
// specified `b64_json` as the ResponseFormat.
type Images struct {
	c     *Client
	URL   string `json:"url"`
	Image []byte `json:"b64_json"`
}
