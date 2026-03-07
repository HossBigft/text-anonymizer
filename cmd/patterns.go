/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	patternmanager "anonymizer/patternManager"
	"fmt"
	"github.com/spf13/cobra"
)

// patternsCmd represents the patterns command
var patternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		maskPatterns := patternmanager.NewPatternManager().GetPatterns()
		for _, pattern := range maskPatterns {
			fmt.Println("Name: " + pattern.Name + ", Pattern: " + pattern.Regex)
		}
	},
}

func init() {
	rootCmd.AddCommand(patternsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// patternsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// patternsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
