/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/andrewstuart/openai"
	"github.com/andrewstuart/p"
	"github.com/spf13/cobra"
)

// translateCmd represents the translate command
var translateCmd = &cobra.Command{
	Use:   "translate <file>",
	Short: "Translate audio to English",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prompt, _ := cmd.Flags().GetString("prompt")
		force, _ := cmd.Flags().GetBool("force")
		if !force && prompt == "" {
			log.Println("Please use a prompt that explains the translation from/to because it seems to be necessary for openai. Force with -f to ignore")
		}
		bs, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal(err)
		}
		t := openai.TranscriptionReq{
			File:  bs,
			Model: openai.TranscriptionModelWhisper1,
		}
		if prompt != "" {
			t.Prompt = p.T(prompt)
		}

		res, err := c.Translation(ctx, t)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res.Text)
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().StringP("prompt", "p", "", "The prompt to send")
	translateCmd.Flags().BoolP("force", "f", false, "Ignore the prompt warning")
}
