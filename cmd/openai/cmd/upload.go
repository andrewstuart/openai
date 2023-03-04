/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/andrewstuart/openai"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to openai",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fn := args[0]
		f, err := os.OpenFile(fn, os.O_RDONLY, 0400)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		res, err := c.Upload(ctx, openai.FileUploadReq{
			Filename: fn,
			File:     f,
			Purpose:  openai.PurposeFineTune,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(res.ID)
	},
}

func init() {
	filesCmd.AddCommand(uploadCmd)
}
