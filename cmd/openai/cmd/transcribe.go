/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/andrewstuart/openai"
	"github.com/spf13/cobra"
)

// transcribeCmd represents the transcribe command
var transcribeCmd = &cobra.Command{
	Use:   "transcribe",
	Short: "Transcribe an audio file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bs, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal(err)
		}

		res, err := c.Transcription(ctx, openai.TranscriptionParams{
			File:  bs,
			Model: openai.TranscriptionModelWhisper1,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res.Text)

	},
}

func init() {
	rootCmd.AddCommand(transcribeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transcribeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transcribeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
