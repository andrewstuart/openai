package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/andrewstuart/openai"
	"github.com/andrewstuart/p"
	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Generate an image with prompts",
	Run: func(cmd *cobra.Command, args []string) {
		res := strings.Join(args, " ")
		if len(args) == 0 || args[0] == "-" {
			bs, err := io.ReadAll(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
			res = string(bs)
		}
		if len(res) > 1000 {
			res = res[:1000]
		}
		n, _ := cmd.Flags().GetInt("number")
		ires, err := c.GenerateImage(ctx, openai.ImgReq{
			Prompt:         res,
			ResponseFormat: p.T("b64_json"),
			N:              p.T(n),
		})
		if err != nil {
			log.Fatal(err)
		}

		file, _ := cmd.Flags().GetString("file")
		if n == 1 {
			err = ioutil.WriteFile(file, ires.Data[0].Image, 0600)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		ext := path.Ext(file)
		base := strings.TrimSuffix(file, ext)
		for i, d := range ires.Data {
			err = ioutil.WriteFile(fmt.Sprintf("%s-%03d%s", base, i+1, ext), d.Image, 0600)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.Flags().StringP("file", "f", "out.jpg", "Filename to write the output to")
	imageCmd.Flags().IntP("number", "n", 1, "The number of images to generate")
}
