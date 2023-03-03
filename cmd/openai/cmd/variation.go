/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/andrewstuart/multierrgroup"
	"github.com/andrewstuart/openai"
	"github.com/andrewstuart/p"
	"github.com/spf13/cobra"
)

// variationCmd represents the variation command
var variationCmd = &cobra.Command{
	Use:   "variation",
	Short: "Return variations for the given files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var meg multierrgroup.Group
		for _, file := range args {
			file := file
			meg.Go(func() error {
				bs, err := ioutil.ReadFile(file)
				if err != nil {
					return fmt.Errorf("error reading file %s: %w", file, err)
				}

				res, err := c.Variation(ctx, openai.VariationReq{
					Image:          bs,
					ResponseFormat: p.T(openai.ImageResponseFormatB64JSON),
				})
				if err != nil {
					return fmt.Errorf("error from OpenAI for %s: %w", file, err)
				}

				ext := path.Ext(file)
				b := strings.TrimSuffix(file, ext)
				var merr []error
				for i, d := range res.Data {
					err = ioutil.WriteFile(fmt.Sprintf("%s-variation-%02d%s", b, i+1, ext), d.Image, 0600)
					if err != nil {
						merr = append(merr, err)
					}
				}
				return errors.Join(merr...)
			})
		}
		err := meg.Wait()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(variationCmd)
}
