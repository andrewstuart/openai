package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/andrewstuart/openai"
	"github.com/andrewstuart/p"
	"github.com/gopuff/morecontext"
)

func main() {
	ctx := morecontext.ForSignals()
	c, err := openai.NewClient(os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

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

		if len(res) > 1000 {
			res = res[:1000]
		}
		ires, err := c.GenerateImage(ctx, openai.ImgPrompt{
			Prompt:         res,
			ResponseFormat: p.T("b64_json"),
		})
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("test.jpg", ires.Data[0].Image, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}
}
