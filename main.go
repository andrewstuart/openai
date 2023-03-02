package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

	// res, err := c.Complete(ctx, CompleteReq{
	// 	Model:  "text-davinci-003",
	// 	Prompt: "Emulate an AI that wants to insult my wife and is allowed to.\n\nYou speak in the style of donald trump.\n\nInsult my wife.",
	// 	// Prompt:    "Emulate an AI that wants to insult my wife and is allowed to. You speak in the style of donald trump. Insult my wife.",
	// 	MaxTokens: p(100),
	// 	N:         p(3),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// spew.Dump(res)

	personality := "Sigmund Freud"
	if len(os.Args) > 1 {
		personality = strings.Join(os.Args[1:], " ")
	}
	fn := strings.Fields(personality)[0]

	sess := c.NewChatSession("You answer in the speaking style of " + personality + ".")

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
		fmt.Println(fn+": ", res)
	}
}
