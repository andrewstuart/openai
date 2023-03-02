package openai

import (
	"net/http"

	"git.stuart.fun/andrew/rester/v2"
)

type Client struct {
	c *rester.Client
}

func NewClient(tok string) (*Client, error) {
	c := rester.Must(rester.New("https://api.openai.com/v1"))
	c.Transport = rester.All{
		rester.MergeHeaders{
			"Authorization": {"Bearer " + tok},
			"Content-Type":  {"application/json"},
		},
		rester.ResponseFunc(parseOpenAIError),
	}.Wrap(http.DefaultTransport)
	return &Client{c: c}, nil
}
