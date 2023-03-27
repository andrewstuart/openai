package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/andrewstuart/openai"
	"github.com/cenkalti/backoff"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with somebody",
	RunE: func(cmd *cobra.Command, args []string) error {

		prompt := "You are a helpful AI assistant."
		fn := "Assistant"
		p := viper.GetString("history.path")
		var out *os.File

		if p != "" {
			fp := path.Join(p, time.Now().Format(time.RFC3339)) + ".json"
			var err error
			out, err = os.OpenFile(fp, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0600)
			if err != nil {
				return err
			}
			defer out.Close()
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

		prompts := viper.GetStringMapString("prompts")
		if prompts != nil && prompts[prompt] != "" {
			prompt = prompts[prompt]
		}

		sess := c.NewChatSession(prompt)
		m, _ := cmd.Flags().GetString("model")
		models, err := c.Models(ctx)
		if err != nil {
			return err
		}
		if models.Has(m) {
			sess.Model = m
		} else {
			sess.Model = openai.ChatModelGPT35Turbo0301
		}

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
			}, &in, survey.WithIcons(func(is *survey.IconSet) { is.Question.Text = "" }))
			if err != nil {
				return err
			}

			if in == "" {
				return nil
			}

			bo := backoff.NewExponentialBackOff()
			bo.InitialInterval = time.Second
			var res <-chan string
			err = backoff.Retry(func() error {
				var err error
				res, err = sess.Stream(ctx, in)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				return err
			}, backoff.WithMaxRetries(backoff.WithContext(bo, ctx), 5))
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
			if out != nil {
				out.Seek(0, 0)
				json.NewEncoder(out).Encode(sess)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
	chatCmd.Flags().StringP("prompt", "p", "", "A prompt to override the default")
	chatCmd.Flags().String("personality", "", "Shorthand for a personality to use as the speaking style for the prompt.")
	chatCmd.Flags().String("model", openai.ChatModelGPT4, "The model to use for chat completion")
	// chatCmd.Flags().String("resume", "", "Resume a chat from file") TODO: add resume
	chatCmd.RegisterFlagCompletionFunc("model", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ourModels := []string{openai.ChatModelGPT35Turbo, openai.ChatModelGPT35Turbo0301}
		ms, err := c.Models(ctx)
		if err != nil {
			return ourModels, 0
		}
		if ms.Has(openai.ChatModelGPT4) {
			ourModels = append(ourModels, openai.ChatModelGPT4, openai.ChatModelGPT40314)
		}
		if ms.Has(openai.ChatModelGPT432K) {
			ourModels = append(ourModels, openai.ChatModelGPT432K, openai.ChatModelGPT432K0314)
		}
		return ourModels, 0
	})
	chatCmd.RegisterFlagCompletionFunc("prompt", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		m := viper.GetStringMapString("prompts")
		if m != nil {
			return maps.Keys(m), 0
		}

		return []string{}, cobra.ShellCompDirectiveDefault
	})

}
