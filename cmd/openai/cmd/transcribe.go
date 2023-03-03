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

		res, err := c.Transcription(ctx, openai.TranscriptionReq{
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
}
