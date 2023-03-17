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

		prompt := "You are a helpful AI assistant."
		fn := "Assistant"

		personality, _ := cmd.Flags().GetString("personality")
		if personality != "" {
			fn = strings.Fields(personality)[0]
			prompt = "You answer in the speaking style of " + personality + "."
		}

		if pr, _ := cmd.Flags().GetString("prompt"); pr != "" {
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

			fmt.Println(fn+":", strings.TrimSpace(res))

		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringP("prompt", "p", "", "A prompt to override the default")
	chatCmd.Flags().String("personality", "", "Shorthand for a personality to use as the speaking style for the prompt.")
}
