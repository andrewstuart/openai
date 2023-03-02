/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with somebody",
	Run: func(cmd *cobra.Command, args []string) {
		personality, _ := cmd.Flags().GetString("personality")
		fn := strings.Fields(personality)[0]

		prompt := "You answer in the speaking style of " + personality + "."
		if pr, err := cmd.Flags().GetString("prompt"); err == nil {
			prompt = pr
			fn = "Response"
		}

		sess := c.NewChatSession(prompt)

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
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().String("prompt", "", "A prompt to override the default")
	chatCmd.Flags().String("personality", "Sigmund Freud", "A prompt to override the default")
}
