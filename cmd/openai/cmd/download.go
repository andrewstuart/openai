/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"log"
	"os"

	"github.com/andrewstuart/openai"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func validFiles(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	fs, err := c.ListFiles(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	fns := lo.Map(fs.Data, func(f openai.FileRes, _ int) string {
		return f.ID + "\t" + f.Filename
	})
	return fns, cobra.ShellCompDirectiveDefault
}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:               "download",
	Short:             "A brief description of your command",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: validFiles,
	Run: func(cmd *cobra.Command, args []string) {
		for _, f := range args {
			func() {
				info, err := c.GetFileDetails(ctx, f)
				if err != nil {
					log.Fatal(err)
				}
				out, err := os.OpenFile(info.Filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
				if err != nil {
					log.Fatal(err)
				}
				defer out.Close()

				rc, err := c.DownloadFile(ctx, f)
				if err != nil {
					log.Fatal(err)
				}
				defer rc.Close()
				_, err = io.Copy(out, rc)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	},
}

func init() {
	filesCmd.AddCommand(downloadCmd)
}
