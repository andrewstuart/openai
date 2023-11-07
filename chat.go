package openai

import (
	"context"
)

// Well-known Chat constants
const (
	ChatRoleSystem    = "system"
	ChatRoleAssistant = "assistant"
	ChatRoleUser      = "user"

	ChatModelGPT35Turbo       = "gpt-3.5-turbo"
	ChatModelGPT35Turbo0301   = "gpt-3.5-turbo-0301"
	ChatModelGPT35Turbo16K    = "gpt-3.5-turbo-16k"
	ChatModelGPT35Turbo0613   = "gpt-3.5-turbo-0613"
	ChatModelGPT35Turbo1106   = "gpt-3.5-turbo-1106"
	ChatModelGPT4             = "gpt-4"
	ChatModelGPT40314         = "gpt-4-0314"
	ChatModelGPT40613         = "gpt-4-0613"
	ChatModelGPT432K          = "gpt-4-32k"
	ChatModelGPT432K0314      = "gpt-4-32k-0314"
	ChatModelGPT4TurboPreview = "gpt-4-1106-preview"
)

// ChatComplete is the raw chat/completions endpoint exposed for callers.
func (c Client) ChatComplete(ctx context.Context, r ChatReq) (*ChatRes, error) {
	r.Stream = nil
	var res ChatRes
	err := c.c.R().Post("chat/completions").JSON(r).Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// ChatSession is an abstraction that provides simple tracking of ChatMessages
// and sends the entire chat "session" to the OpenAI APIs for proper contextual
// chat.
type ChatSession struct {
	c        *Client
	Messages []ChatMessage
	Model    string
	Tpl      ChatReq
}

// NewChatSession returns a ChatSession object with the given prompt as a
// starting point (with the `system` role).
func (c *Client) NewChatSession(prompt string) ChatSession {
	return ChatSession{
		c:     c,
		Model: ChatModelGPT35Turbo,
		Messages: []ChatMessage{{
			Role:    ChatRoleSystem,
			Content: prompt,
		}},
	}
}

// Complete takes a message from the `user` and sends that, along with the
// Session context, to the OpenAI endpoints for completion.
func (s *ChatSession) Complete(ctx context.Context, msg string) (string, error) {
	s.Messages = append(s.Messages, ChatMessage{
		Role:    ChatRoleUser,
		Content: msg,
	})

	req := s.Tpl
	req.Messages = s.Messages
	req.Model = s.Model

	res, err := s.c.ChatComplete(ctx, req)
	if err != nil {
		return "", err
	}

	s.Messages = append(s.Messages, res.Choices[0].Message)
	return res.Choices[0].Message.Content, nil
}

// ChatReq is a Request to the chat/completions endpoints.
type ChatReq struct {
	Model            string             `json:"model"`
	Messages         []ChatMessage      `json:"messages"`
	Temperature      *float64           `json:"temperature,omitempty"`
	TopP             *int               `json:"top_p,omitempty"`
	N                *int               `json:"n,omitempty"`
	Stream           *bool              `json:"stream,omitempty"`
	MaxTokens        *int               `json:"max_tokens,omitempty"`
	PresencePenalty  *float64           `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64           `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]float64 `json:"logit_bias,omitempty"`
	User             *string            `json:"user,omitempty"`
}

// ChatMessage is the structure of the messages array in the request body.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRes represents the return response from OpenAI chat/completions
type ChatRes struct {
	Choices []ChatChoice `json:"choices"`
	Created int          `json:"created"`
	ID      string       `json:"id"`
	Model   string       `json:"model"`
	Object  string       `json:"object"`
	Usage   Usage        `json:"usage"`
}

// ChatChoice is one of the outputs from chatgpt (and sometimes others)
type ChatChoice struct {
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
}
