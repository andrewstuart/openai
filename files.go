package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"time"
)

const (
	PurposeFineTune = "fine-tune"
)

type FileUploadReq struct {
	Filename string
	File     io.Reader
	Purpose  string
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(bs []byte) error {
	var i int64
	err := json.Unmarshal(bs, &i)
	if err != nil {
		return err
	}
	t.Time = time.Unix(i, 0)
	return nil
}

type FileUploadRes struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt Time   `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// Upload uploads a file to be used within the OpenAI plaform. Most often for
// fine-tuning models.
func (c Client) Upload(ctx context.Context, v FileUploadReq) (*FileUploadRes, error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	image, err := w.CreateFormFile("file", v.Filename)
	if err != nil {
		return nil, fmt.Errorf("error creating audio multipart writer: %w", err)
	}
	io.Copy(image, v.File)

	n, err := w.CreateFormField("purpose")
	if err != nil {
		return nil, fmt.Errorf("error creating audio multipart writer model: %w", err)
	}
	fmt.Fprint(n, v.Purpose)

	w.Close()

	var res FileUploadRes
	err = c.c.R().
		Post("files").
		SetHeader("Content-Type", "multipart/form-data; boundary="+w.Boundary()).
		WithBody(body).
		Do(ctx).
		JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
