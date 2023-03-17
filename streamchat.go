package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/andrewstuart/p"
)

type ChatStreamRes struct {
	ID        string             `json:"id,omitempty"`
	Object    string             `json:"object,omitempty"`
	CreatedAt int64              `json:"created_at,omitempty"`
	Choices   []*ChatStreamChoce `json:"choices,omitempty"`
}

type ChatStreamChoce struct {
	Delta        ChatMessage `json:"delta,omitempty"`
	Index        int         `json:"index,omitempty"`
	LogProbs     int         `json:"logprobs,omitempty"`
	FinishReason string      `json:"finish_reason,omitempty"`
}

func getEvent(bs []byte, r *ChatStreamRes) (bool, error) {
	// fmt.Println(string(bs))
	bs = bytes.TrimPrefix(bytes.TrimSpace(bs), []byte("data:"))
	if len(bs) == 0 {
		return false, nil
	}
	if bytes.Equal(bs, []byte("[DONE]")) {
		return false, io.EOF
	}
	err := json.Unmarshal(bs, r)
	if err != nil {
		return false, fmt.Errorf("error getting json from %s: %w", bs, err)
	}
	return true, nil
}

// ChatStream takes a request and streams the response from OpenAI.
func (c Client) ChatStream(ctx context.Context, r ChatReq) (<-chan ChatStreamRes, error) {
	r.Stream = p.T(true)
	res := c.c.R().Post("chat/completions").JSON(r).Do(ctx)

	br := bufio.NewReader(res)
	bs, err := br.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	var ev ChatStreamRes
	ok, err := getEvent(bs, &ev)
	if err != nil {
		return nil, err
	}

	ch := make(chan ChatStreamRes)
	go func() {
		defer close(ch)
		defer res.HTTP.Body.Close()
		if ok {
			select {
			case ch <- ev:
			case <-ctx.Done():
				return
			}
		}

		for {
			bs, err := br.ReadBytes('\n')
			if err != nil {
				return
			}
			var ev ChatStreamRes
			ok, err := getEvent(bs, &ev)
			if err != nil {
				return
			}
			if ok && len(ev.Choices) > 0 {
				select {
				case ch <- ev:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return ch, nil
}

// Stream takes a message from the `user` and sends that, along with the
// Session context, to the OpenAI endpoints for completion. The response will be streamed.
func (s *ChatSession) Stream(ctx context.Context, msg string) (<-chan string, error) {
	s.Messages = append(s.Messages, ChatMessage{
		Role:    ChatRoleUser,
		Content: msg,
	})

	res, err := s.c.ChatStream(ctx, ChatReq{
		Model:    s.Model,
		Messages: s.Messages,
	})
	if err != nil {
		return nil, err
	}

	content := ""
	ch := make(chan string)
	go func() {
		defer close(ch)
		defer func() {
			s.Messages = append(s.Messages, ChatMessage{
				Role:    ChatRoleSystem,
				Content: content,
			})
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case r, ok := <-res:
				if !ok {
					return
				}
				content += r.Choices[0].Delta.Content
				select {
				case <-ctx.Done():
					return
				case ch <- r.Choices[0].Delta.Content:
				}
			}
		}
	}()
	return ch, nil
}
