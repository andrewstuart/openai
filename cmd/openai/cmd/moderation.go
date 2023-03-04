/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/andrewstuart/openai"
	"github.com/spf13/cobra"
)

// moderationCmd represents the moderation command
var moderationCmd = &cobra.Command{
	Use:   "moderation",
	Short: "Use the openai moderation endpoint to check for policy violations",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := c.Moderation(ctx, openai.ModerationReq{
			Input: strings.Join(args, " "),
		})
		if err != nil {
			log.Fatal(err)
		}
		bs, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(os.Stdout, bytes.NewReader(bs))
	},
}

func init() {
	rootCmd.AddCommand(moderationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moderationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moderationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
