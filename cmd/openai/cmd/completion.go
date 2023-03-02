/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

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

		res, err := c.Complete(ctx, openai.CompleteReq{
			Model: model,
			N:     &n,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res.Choices[0].Text)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().StringP("model", "m", openai.CompletionModelDavinci3, "Which completion model to use")
	completionCmd.Flags().IntP("number", "n", 1, "Which completion model to use")
}
