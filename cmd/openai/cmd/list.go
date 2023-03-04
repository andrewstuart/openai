/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files in the openai backend",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := c.ListFiles(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if len(files.Data) == 0 {
			fmt.Println("no files")
		}
		tw := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)
		fmt.Fprintf(tw, "ID\tName\tSize\tPurpose\n")
		for _, f := range files.Data {
			fmt.Fprintf(tw, "%s\t%s\t%d\t%s\n", f.ID, f.Filename, f.Bytes, f.Purpose)
		}
		tw.Flush()
	},
}

func init() {
	filesCmd.AddCommand(listCmd)
}
