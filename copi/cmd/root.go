package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atakanozceviz/copi"
	"github.com/spf13/cobra"
)

var settingsFile string
var backupPath string
var keep int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "copi",
	Short: "Copy files and folders except specified in the settings",
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
		dst := args[1]

		if backupPath != "" && keep >= 1 {
			fmt.Printf("Backup: %s\n", dst)
			if err := copi.Backup(dst, backupPath, keep); err != nil {
				fmt.Printf("Cannot backup: %v\n", err)
				os.Exit(1)
			}
		}

		if err := copi.Copy(src, dst, settingsFile); err != nil {
			fmt.Printf("Cannot copy: %v\n", err)
			os.Exit(1)
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
	rootCmd.PersistentFlags().StringVarP(&settingsFile, "settings", "s", "", "filesystem path to settings file")
	rootCmd.PersistentFlags().StringVarP(&backupPath, "backup", "b", "", "filesystem path to backup folder")
	rootCmd.PersistentFlags().IntVarP(&keep, "keep", "k", 3, "number of backups to keep")
}
