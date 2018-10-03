package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atakanozceviz/copi"
	"github.com/spf13/cobra"
)

var settingsFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "copi",
	Short: "Copy files and folders except specified in settings",
	Long: `Usage:
copi [source] [destination]

Copies files and folders from [source] to [destination]
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		for i, v := range args {
			v = strings.Replace(v, "\\", "/", -1)
			if !strings.HasSuffix(v, "/") {
				v = v + "/"
			}
			args[i] = v
		}
		src := args[0]
		dest := args[1]

		if err := copi.Copy(src, dest, settingsFile); err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&settingsFile, "settings", "s", "copi.json", "path to settings file")
}
