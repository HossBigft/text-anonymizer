/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	maskmanager "anonymizer/maskManager"
	patternmanager "anonymizer/patternManager"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/lucasjones/reggen"
	"github.com/spf13/cobra"
)

var filePath string
var rootCmd = &cobra.Command{
	Use:   "anonymizer",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		maskManager := maskmanager.NewMaskManager()
		patternManager := patternmanager.NewPatternManager()
		maskPatterns := patternManager.GetPatterns()
		isMaskPatternsUpdated := false
		var isMasksUpdated bool
		if len(filePath) > 0 {
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
		} else {
			if len(args) == 0 {
				scanner := bufio.NewScanner(os.Stdin)
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
								isMasksUpdated = true
							}
						}

					}
					fmt.Println(replaced_line)
				}
			} else {
				for _, val := range strings.Split(args[0], "\n") {
					line := val
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
								isMasksUpdated = true
							}
						}

					}
					fmt.Println(replaced_line)
				}

			}
		}

		if isMasksUpdated {
			err := maskManager.SaveMasks()
			if err != nil {
				fmt.Println(err)
			}
		}

		if isMaskPatternsUpdated {
			patternManager.SavePatterns()
		}
	},
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "File to process")
}
