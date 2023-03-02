package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"git.stuart.fun/andrew/rester/v2"
	"github.com/gopuff/morecontext"
)

func parseOpenAIError(r *http.Response) (*http.Response, error) {
	if r.StatusCode > 399 {
		defer r.Body.Close()
		var e OpenAIError
		err := json.NewDecoder(r.Body).Decode(&e)
		if err != nil {
			return r, err
		}
		return r, &e
	}
	return r, nil
}

type Client struct {
	c *rester.Client
}

func (c Client) Complete(ctx context.Context, r ChatReq) (*ChatRes, error) {
	var res ChatRes
	err := c.c.R().Post("chat/completions").JSON(r).Do(ctx).JSON(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
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

func main() {
	ctx := morecontext.ForSignals()
	c, err := NewClient(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	sess := c.NewChatSession("You answer in the speaking style of Donald Trump.")

	r := bufio.NewReader(os.Stdin)
	go func() {
		<-ctx.Done()
		os.Stdin.SetDeadline(time.Now())
	}()

	for {
		select {
		case <-ctx.Done():
		default:
		}
		fmt.Print("You: ")
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		res, err := sess.Complete(ctx, msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Trump: ", res)
	}
}
