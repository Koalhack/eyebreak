package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eyebreak",
	Short: "A CLI for digital eye break reminder.",
	Long: `A Simple and Fast Eye Break utility built with love by Koalhack in Go.
This Tool respect the 20-20-20 rule.
(twenty minutes of work, twenty seconds of looking at something twenty feet away)`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Start()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of EyeBreak.",
	Long:  `All software has versions. This is EyeBreak`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("EyeBreak v0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
