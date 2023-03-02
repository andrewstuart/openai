package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/andrewstuart/openai"
	"github.com/andrewstuart/p"
	"github.com/spf13/cobra"
)

// explainCodeCmd represents the explainCode command
var explainCodeCmd = &cobra.Command{
	Use:   "explain-code",
	Short: "Explain the code for a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		bs, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		pr := string(bs) + "---\n\nHere's what the above code is doing:"

		res, err := c.Complete(ctx, openai.CompleteReq{
			Model:     "code-davinci-002",
			Prompt:    pr,
			MaxTokens: p.T(256),
			Stop:      p.T("---"),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res.Choices[0].Text)
	},
}

func init() {
	rootCmd.AddCommand(explainCodeCmd)
}
