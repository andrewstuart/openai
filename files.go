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

type FileRes struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt Time   `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// Upload uploads a file to be used within the OpenAI plaform. Most often for
// fine-tuning models.
func (c Client) Upload(ctx context.Context, v FileUploadReq) (*FileRes, error) {
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

	var res FileRes
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

// FileListRes represents the list of files managed by OpenAI
type FileListRes struct {
	Data   []FileRes `json:"data"`
	Object string    `json:"object"`
}

// ListFiles returns the list of files known to OpenAI
func (c Client) ListFiles(ctx context.Context) (*FileListRes, error) {
	var res FileListRes
	err := c.c.R().
		Get("files").
		Do(ctx).
		JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetFileDetails returns the File details for an individual file managed by
// the OpenAI APIs.
func (c Client) GetFileDetails(ctx context.Context, id string) (*FileRes, error) {
	var res FileRes
	err := c.c.R().Get("/files/%s", id).Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// DeleteFile deletes the file specified.
func (c Client) DeleteFile(ctx context.Context, id string) error {
	return c.c.R().Delete("files/%s", id).Do(ctx).Err()
}

// DownloadFile returns a file as an io.ReadCloser. For now, any API errors will be
// returned in the io.Read method, though this will likely change, hence the
// error return.
func (c Client) DownloadFile(ctx context.Context, id string) (io.ReadCloser, error) {
	return c.c.R().Get("/files/%s/content", id).Do(ctx), nil
}
