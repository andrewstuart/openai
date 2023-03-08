package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/andrewstuart/openai"
	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Get completion text from OpenAI",
	Run: func(cmd *cobra.Command, args []string) {
		model, _ := cmd.Flags().GetString("model")
		n, _ := cmd.Flags().GetInt("number")
		t, _ := cmd.Flags().GetInt("tokens")
		f, _ := cmd.Flags().GetString("file")
		p, _ := cmd.Flags().GetString("prefix")
		s, _ := cmd.Flags().GetString("suffix")
		// Read a complete line from either stdin or an entire file from `--file/-f`.
		var comp string
		if f != "" && f != "-" {
			bs, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatal(err)
			}
			comp = string(bs)
		} else if f == "-" || (len(args) == 1 && args[0] == "-") || len(args) == 0 {
			bs, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
			comp = string(bs)
		} else {
			comp = strings.Join(args, " ")
		}

		res, err := c.Complete(ctx, openai.CompleteReq{
			Model:     model,
			Prompt:    p + comp + s,
			N:         &n,
			MaxTokens: &t,
		})
		if err != nil {
			log.Fatal(err)
		}

		if n == 1 {
			fmt.Println(res.Choices[0].Text)
			return
		}
		for i, c := range res.Choices {
			fmt.Printf("Response %d:\n\n", i+1)
			fmt.Println(c.Text)
		}
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
	completeCmd.Flags().StringP("model", "m", openai.CompletionModelDavinci3, "Which completion model to use")
	completeCmd.Flags().IntP("number", "n", 1, "How many completions to generate")
	completeCmd.Flags().IntP("tokens", "t", 64, "How many completions to generate")
	completeCmd.Flags().StringP("file", "f", "", "Optional file to load")
	completeCmd.Flags().StringP("prefix", "p", "", "Prefix to prepend (for files/stdin, mostly)")
	completeCmd.Flags().StringP("suffix", "s", "", "Suffix to append (for files/stdin, mostly)")
}
