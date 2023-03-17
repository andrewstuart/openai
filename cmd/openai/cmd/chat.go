package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/andrewstuart/openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with somebody",
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := "You are a helpful AI assistant."
		fn := "Assistant"
		p := viper.GetString("history.path")
		var out io.Writer
		if p != "" {
			fp := path.Join(p, time.Now().Format(time.RFC3339))
			var err error
			outF, err := os.OpenFile(fp, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			defer outF.Close()
			bw := bufio.NewWriter(outF)
			defer bw.Flush()
			out = bw
		}

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
				return err
			}
			if out != nil {
				fmt.Fprintf(out, "%s (%s): %s\n", "You", time.Now().Format(time.RFC3339), in)
			}

			res, err := sess.Stream(ctx, in)
			if err != nil {
				return err
			}

			outStr := ""
			fmt.Print(fn + ": ")
			for st := range res {
				outStr += st
				fmt.Print(st)
			}
			fmt.Println()
			fmt.Fprintf(out, "%s (%s): %s\n", fn, time.Now().Format(time.RFC3339), outStr)
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
