package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/andrewstuart/openai"
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
	}
}
