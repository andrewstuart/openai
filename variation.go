package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
)

// Variation returns a Variation of an image. Convenience methods exist on
// images already returned from Client calls to easily vary those images.
func (c Client) Variation(ctx context.Context, v VariationReq) (*ImageRes, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	image, err := w.CreateFormFile("image", "image.png")
	if err != nil {
		return nil, fmt.Errorf("error creating image multipart writer: %w", err)
	}
	io.Copy(image, bytes.NewReader(v.Image))

	if v.N != nil {
		n, err := w.CreateFormField("n")
		if err != nil {
			return nil, fmt.Errorf("error creating image multipart writer n: %w", err)
		}
		fmt.Fprint(n, *v.N)
	}
	if v.Size != nil {
		n, err := w.CreateFormField("size")
		if err != nil {
			return nil, fmt.Errorf("error creating image multipart writer size: %w", err)
		}
		fmt.Fprint(n, *v.Size)
	}
	if v.ResponseFormat != nil {
		n, err := w.CreateFormField("response_format")
		if err != nil {
			return nil, fmt.Errorf("error creating image multipart writer ResponseFormat: %w", err)
		}
		fmt.Fprint(n, *v.ResponseFormat)
	}
	if v.User != nil {
		n, err := w.CreateFormField("user")
		if err != nil {
			return nil, fmt.Errorf("error creating image multipart writer User: %w", err)
		}
		fmt.Fprint(n, *v.User)
	}
	w.Close()

	var res ImageRes
	err = c.c.R().
		Post("images/variations").
		SetHeader("Content-Type", "multipart/form-data; boundary="+w.Boundary()).
		WithBody(body).
		Do(ctx).
		JSON(&res)
	if err != nil {
		return nil, err
	}
	for i := range res.Data {
		res.Data[i].c = &c
	}
	return &res, nil
}

// Create a variation on an image already returned.
func (i Images) Variation(ctx context.Context, v VariationReq) (*ImageRes, error) {
	v.Image = i.Image
	return i.c.Variation(ctx, v)
}

// VariationReq hold the data needed for image variation.
type VariationReq struct {
	Image          []byte  `json:"image,omitempty"`
	N              *int    `json:"n,omitempty"`
	Size           *string `json:"size,omitempty"`
	ResponseFormat *string `json:"response_format,omitempty"`
	User           *string `json:"user,omitempty"`
}
