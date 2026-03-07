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
	"strings"

	"github.com/spf13/cobra"
)

func mask(rawLine string, patternManager patternmanager.PatternManager, maskManager maskmanager.MaskManager) string {
	replaced_line := rawLine
	valuesToMaskMap, _ := patternManager.MapValuesToPatterns(rawLine)

	for _, match := range valuesToMaskMap {
		masks := maskManager.MapValuesToMasks(match)
		for sensitive_value, mask := range masks {
			replaced_line = strings.ReplaceAll(replaced_line, sensitive_value, mask)
		}
	}

	return replaced_line
}

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
		if len(filePath) > 0 {
			file, err := os.Open(filePath)
			check(err)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				replaced_line := mask(line, *patternManager, *maskManager)
				fmt.Println(replaced_line)
			}
		} else {
			if len(args) == 0 {
				stat, _ := os.Stdin.Stat()
				if (stat.Mode() & os.ModeCharDevice) == 0 {
					scanner := bufio.NewScanner(os.Stdin)
					for scanner.Scan() {
						line := scanner.Text()
						replaced_line := mask(line, *patternManager, *maskManager)
						fmt.Println(replaced_line)
					}
				}
			} else {
				for _, val := range strings.Split(args[0], "\n") {
					line := val
					replaced_line := mask(line, *patternManager, *maskManager)
					fmt.Println(replaced_line)
				}

			}
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
