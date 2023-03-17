package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/andrewstuart/openai"
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
		m, _ := cmd.Flags().GetString("model")
		sess.Model = m

		go func() {
			<-ctx.Done()
			os.Stdin.SetDeadline(time.Now())
		}()

		for {
			select {
			case <-ctx.Done():
			default:
			}

			var in string
			err := survey.AskOne(&survey.Input{
				Message: "You: ",
			}, &in, survey.WithIcons(func(is *survey.IconSet) {
				is.Question.Text = ""
			}))
			if err != nil {
				log.Fatal(err)
			}
			res, err := sess.Stream(ctx, in)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(fn + ": ")
			for st := range res {
				fmt.Print(st)
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringP("prompt", "p", "", "A prompt to override the default")
	chatCmd.Flags().String("personality", "", "Shorthand for a personality to use as the speaking style for the prompt.")
	chatCmd.Flags().String("model", openai.ChatModelGPT35Turbo0301, "The model to use for chat completion")
	chatCmd.RegisterFlagCompletionFunc("model", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{openai.ChatModelGPT35Turbo, openai.ChatModelGPT35Turbo0301, openai.ChatModelGPT4}, 0
	})
}
