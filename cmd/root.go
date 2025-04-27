/*
Copyright © 2025 ldinc <drogunov.igor@gmail.com>
*/

package cmd

import (
	"goft/internal/manager"
	"os"

	"github.com/spf13/cobra"
)

var (
	requestedInfo  bool
	requestedBuild bool

	targetPath string
	outputPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goft",
	Short: "Go Factorio tool",
	Long:  `Small set usefull factorio commands for mod developers written with Go.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		if targetPath == "" {
			var err error

			targetPath, err = os.Getwd()

			if err != nil {
				panic(err)
			}
		}

		if outputPath == "" {
			outputPath = targetPath
		}

		manager := manager.New(targetPath, outputPath, args)

		if requestedBuild {
			manager.BuildAll()
		}

	},
}

func init() {
	rootCmd.Flags().BoolVarP(&requestedBuild, "build", "b", false, "Build zip archive with mod")
	rootCmd.Flags().StringVarP(&targetPath, "dir", "d", "", "Path to mod (by default: current dir)")
	rootCmd.Flags().StringVarP(&outputPath, "out", "o", "", "Path to zip (by default: current dir)")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
