package openai

import "context"

// Well-known Chat constants
const (
	ChatRoleSystem          = "system"
	ChatRoleUser            = "user"
	ChatModelGPT35Turbo     = "gpt-3.5-turbo"
	ChatModelGPT35Turbo0301 = "gpt-3.5-turbo-0301"
)

// ChatComplete is the raw chat/completions endpoint exposed for callers.
func (c Client) ChatComplete(ctx context.Context, r ChatReq) (*ChatRes, error) {
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

	res, err := s.c.ChatComplete(ctx, ChatReq{
		Model:    s.Model,
		Messages: s.Messages,
	})
	if err != nil {
		return "", err
	}

	s.Messages = append(s.Messages, res.Choices[0].Message)
	return res.Choices[0].Message.Content, nil
}

// ChatReq is a Request to the chat/completions endpoints.
type ChatReq struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
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
