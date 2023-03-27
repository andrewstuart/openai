package openai

import (
	"net/http"

	"git.stuart.fun/andrew/rester/v2"
)

const defaultBase = "https://api.openai.com/v1"

// Client holds the base rester.Client and has methods for communicating with
// OpenAI.
type Client struct {
	c *rester.Client
}

type optStruct struct {
	Org        string
	BaseURL    string
	APIVersion string
}

type Opt func(*optStruct)

// Withorg sets the identifier for a specific org, for the case where a user may be part of more than one organization.
func WithOrg(o string) Opt {
	return func(opt *optStruct) {
		opt.Org = o
	}
}

// WithBaseURL sets the base URL of the API service to use. This can be used
// for alternative hosted versions of the OpenAI API, like MS Azure.
func WithBaseURL(base string) Opt {
	return func(opt *optStruct) {
		opt.BaseURL = base
	}
}

// WithAPIVersion sets the API Version header of the API service to use. This
// can be used for alternative hosted versions of the OpenAI API, like MS
// Azure.
func WithAPIVersion(base string) Opt {
	return func(opt *optStruct) {
		opt.APIVersion = base
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
	if os.APIVersion != "" {
		dh["OpenAI-Version"] = []string{os.APIVersion}
	}

	base := os.BaseURL
	if base == "" {
		base = defaultBase
	}

	c := rester.Must(rester.New(base))
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
