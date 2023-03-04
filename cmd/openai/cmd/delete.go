/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:               "delete",
	Short:             "Delete a file",
	ValidArgsFunction: validFiles,
	Run: func(cmd *cobra.Command, args []string) {
		var out bool
		err := survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to delete %d files?", len(args)),
		}, &out)
		if err != nil {
			log.Fatal(err)
		}
		if !out {
			return
		}
		for _, a := range args {
			fmt.Print("Deleting " + a)
			err := c.DeleteFile(ctx, a)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(" ✅")
		}
	},
}

func init() {
	filesCmd.AddCommand(deleteCmd)
}
