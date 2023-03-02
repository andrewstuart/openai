package openai

import (
	"net/http"

	"git.stuart.fun/andrew/rester/v2"
)

// Client holds the base rester.Client and has methods for communicating with
// OpenAI.
type Client struct {
	c *rester.Client
}

type optStruct struct {
	Org string
}

type Opt func(*optStruct)

func WithOrg(o string) Opt {
	return func(opt *optStruct) {
		opt.Org = o
	}
}

// NewClient returns an OpenAI base client with the token.
func NewClient(tok string, opts ...Opt) (*Client, error) {
	var os optStruct
	for _, o := range opts {
		o(&os)
	}

	dh := rester.DefaultHeaders{
		"Authorization": {"Bearer " + tok},
		"Content-Type":  {"application/json"},
	}

	if os.Org != "" {
		dh["OpenAI-Organization"] = []string{os.Org}
	}

	c := rester.Must(rester.New("https://api.openai.com/v1"))
	c.Transport = rester.All{
		dh,
		rester.ResponseFunc(parseOpenAIError),
	}.Wrap(http.DefaultTransport)
	return &Client{c: c}, nil
}

// Usage is a record type returned from many different openai endpoints letting
// the user know how many tokens were used processing their request.
type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
