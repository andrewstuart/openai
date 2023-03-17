/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/andrewstuart/openai"
	"github.com/andrewstuart/p"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an input with a prompt",
	Run: func(cmd *cobra.Command, args []string) {
		f, _ := cmd.Flags().GetString("file")
		var comp string
		if f != "" && f != "-" {
			bs, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatal(err)
			}
			comp = string(bs)
		} else if f == "-" || (len(args) == 1 && args[0] == "-") || len(args) == 0 {
			bs, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
			comp = string(bs)
		} else {
			comp = strings.Join(args, " ")
		}

		i, _ := cmd.Flags().GetString("instruction")
		m, _ := cmd.Flags().GetString("model")
		n, _ := cmd.Flags().GetInt("number")
		t, _ := cmd.Flags().GetFloat64("temp")
		res, err := c.Edit(ctx, openai.EditReq{
			Model:       m,
			Input:       p.T(comp),
			Instruction: i,
			N:           &n,
			Temperature: &t,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res.Choices[0].Text)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().StringP("file", "f", "", "Optional file to load")
	editCmd.Flags().StringP("instruction", "i", "", "Instructions to the editor")
	editCmd.Flags().StringP("model", "m", openai.EditModelDavinciCode1, "The model to use")
	editCmd.Flags().IntP("number", "n", 1, "The model to use")
	editCmd.Flags().Float64P("temp", "t", 0, "The model to use")
	editCmd.RegisterFlagCompletionFunc("model", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{openai.EditModelDavinci1, openai.EditModelDavinciCode1}, 0
	})

}
