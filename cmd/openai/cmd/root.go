package cmd

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/andrewstuart/openai"
	"github.com/gopuff/morecontext"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var c *openai.Client
var ctx = morecontext.ForSignals()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "openai",
	Short: "An application for interacting with openai APIs",
	// Run:   func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	u, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(path.Join(u, ".config"))
	}
	viper.SetConfigName("openai")
	viper.AutomaticEnv()
	viper.EnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.ReadInConfig()

	c, err = openai.NewClient(viper.GetString("token"), openai.WithOrg(viper.GetString("org")))
	if err != nil {
		log.Fatal(err)
	}
}
