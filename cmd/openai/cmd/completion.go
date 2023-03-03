package cmd

import (
	"bufio"
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
		// Read a complete line from either stdin or an entire file from `--file/-f`.
		var comp string
		if f != "" {
			bs, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatal(err)
			}
			comp = string(bs)
		} else if len(args) == 0 {
			inLine, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			comp = inLine
		} else {
			comp = strings.Join(args, " ")
		}

		res, err := c.Complete(ctx, openai.CompleteReq{
			Model:     model,
			Prompt:    comp,
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
	completeCmd.Flags().StringP("file", "f", "", "Opteal file to load")
}
