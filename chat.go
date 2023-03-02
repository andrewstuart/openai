package openai

import "context"

func (c Client) ChatComplete(ctx context.Context, r ChatReq) (*ChatRes, error) {
	var res ChatRes
	err := c.c.R().Post("chat/completions").JSON(r).Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type ChatSession struct {
	c        *Client
	Messages []ChatMessage
}

func (c *Client) NewChatSession(prompt string) ChatSession {
	return ChatSession{
		c: c,
		Messages: []ChatMessage{{
			Role:    "system",
			Content: prompt,
		}},
	}
}

func (s *ChatSession) Complete(ctx context.Context, msg string) (string, error) {
	s.Messages = append(s.Messages, ChatMessage{
		Role:    "user",
		Content: msg,
	})

	res, err := s.c.ChatComplete(ctx, ChatReq{
		Model:    "gpt-3.5-turbo",
		Messages: s.Messages,
	})
	if err != nil {
		return "", err
	}

	s.Messages = append(s.Messages, res.Choices[0].Message)
	return res.Choices[0].Message.Content, nil
}

type ChatReq struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRes struct {
	Choices []Choice `json:"choices"`
	Created int      `json:"created"`
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
