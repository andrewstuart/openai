/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/andrewstuart/openai"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Get completion text from OpenAI",
	Run: func(cmd *cobra.Command, args []string) {
		model, _ := cmd.Flags().GetString("model")
		n, _ := cmd.Flags().GetInt("number")
		t, _ := cmd.Flags().GetInt("tokens")
		f, _ := cmd.Flags().GetString("file")

		// Read a completion line from either stdin or an entire file from `--file/-f`.
		var comp string
		if f != "" {
			bs, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatal(err)
			}
			comp = string(bs)
		} else {
			inLine, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			comp = inLine
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
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().StringP("model", "m", openai.CompletionModelDavinci3, "Which completion model to use")
	completionCmd.Flags().IntP("number", "n", 1, "How many completions to generate")
	completionCmd.Flags().IntP("tokens", "t", 64, "How many completions to generate")
	completionCmd.Flags().StringP("file", "f", "", "Optional file to load")
}
