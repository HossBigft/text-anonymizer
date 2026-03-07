/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	patternmanager "anonymizer/patternManager"

	"fmt"
	"github.com/spf13/cobra"
	"os"
)


var patternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		list, _ := cmd.Flags().GetBool("list")

		if list {
			maskPatterns := patternmanager.NewPatternManager().GetPatterns()
			for _, pattern := range maskPatterns {
				fmt.Println("Name: " + pattern.Name + ", Pattern: " + pattern.Regex)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(patternsCmd)

	patternsCmd.Flags().BoolP("list", "l", false, "List mask patterns")
}
