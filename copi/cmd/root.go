package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/atakanozceviz/copi"
	"github.com/spf13/cobra"
)

var ignoreList string
var backupPath string
var keep int
var remove bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "copi",
	Short: "Copy files and folders except specified in the ignore list",
	Long: `Usage:
copi [source] [destination]

Copies files and folders from [source] to [destination]
`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		for i, arg := range args {
			if !path.IsAbs(arg) {
				arg = path.Clean(path.Join(wd, arg))
			}
			arg = strings.Replace(arg, "\\", "/", -1)
			if !strings.HasSuffix(arg, "/") {
				arg = arg + "/"
			}
			args[i] = arg
		}

		list, err := copi.ParseIgnore(ignoreList)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		src := args[0]
		dst := args[1]

		if backupPath != "" {
			backupPath = strings.Replace(backupPath, "\\", "/", -1)
			if !path.IsAbs(backupPath) {
				backupPath = path.Clean(path.Join(wd, backupPath))
			}
			err = copi.Backup(dst, backupPath, keep)
			if err != nil {
				fmt.Printf("Cannot backup: %v\n", err)
				os.Exit(1)
			}
		}

		if remove {
			err = copi.RemoveContentsExcept(dst, list)
			if err != nil {
				fmt.Printf("Cannot remove contents: %v\n", err)
				os.Exit(1)
			}
		}

		err = copi.CopyContentsExcept(src, dst, list)
		if err != nil {
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
	rootCmd.PersistentFlags().StringVarP(&ignoreList, "ignore", "s", "", "filesystem path to list of files and folders to ignore") // TODO: Change shorthand form s to i.(BREAKS EVERYTHING!)
	rootCmd.PersistentFlags().StringVarP(&backupPath, "backup", "b", "", "filesystem path to backup folder")
	rootCmd.PersistentFlags().IntVarP(&keep, "keep", "k", 3, "number of backups to keep")
	rootCmd.PersistentFlags().BoolVarP(&remove, "remove", "r", true, "remove destination contents")
}
