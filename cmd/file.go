/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	maskmanager "anonymizer/maskManager"
	patternmanager "anonymizer/patternManager"
	"bufio"
	"fmt"
	"github.com/lucasjones/reggen"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		maskManager := maskmanager.NewMaskManager()
		patternManager := patternmanager.NewPatternManager()
		maskPatterns := patternManager.GetPatterns()

		isMaskPatternsUpdated := false
		filePath := args[0]
		file, err := os.Open(filePath)
		check(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			replaced_line := line
			for _, pattern := range maskPatterns {
				regex, _ := regexp.Compile(pattern.Regex)
				sensitive_value := regex.FindString(line)
				if len(sensitive_value) != 0 {
					mask, present := maskManager.GetMask(sensitive_value)
					if present {
						replaced_line = strings.ReplaceAll(replaced_line, sensitive_value, mask)
					} else {
						mask, _ = reggen.Generate(pattern.Regex, 7)
						maskManager.UpdateMask(sensitive_value, mask)
						replaced_line = strings.ReplaceAll(replaced_line, sensitive_value, mask)
						isMaskPatternsUpdated = true
					}
				}

			}
			fmt.Println(replaced_line)
			if isMaskPatternsUpdated {
				patternManager.SavePatterns()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
