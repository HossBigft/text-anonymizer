package cmd

import (
	maskmanager "ae/maskManager"
	patternmanager "ae/patternManager"
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
func unmask(rawLine string, maskManager maskmanager.MaskManager) string {
	replaced_line := rawLine
	masks := maskManager.GetMasksToValuesMap()
	for sensitive_value, mask := range masks {
		replaced_line = strings.ReplaceAll(replaced_line, sensitive_value, mask)
	}

	return replaced_line
}

var filePath string
var rootCmd = &cobra.Command{
	Use:   "ae",
	Short: "Mask sensitive values in text with lookalike text",
	Long: `Ae is CLI program to mask sensitive values in given text. It replaces values with randomly generated strings based on regex patterns used for extracting. Masks are persistent and savedon disk. Also it can unmask masked text back. 
	Usage:
	cat example.txt | ae
	ae -f path/example.txt
	`,
	Run: func(cmd *cobra.Command, args []string) {
		maskManager := maskmanager.NewMaskManager()
		patternManager := patternmanager.NewPatternManager()
		reverse, _ := cmd.Flags().GetBool("reverse")

		if len(filePath) > 0 {
			file, err := os.Open(filePath)
			check(err)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				var replaced_line string
				if reverse {
					replaced_line = unmask(line, *maskManager)
				} else {
					replaced_line = mask(line, *patternManager, *maskManager)
				}

				fmt.Println(replaced_line)
			}
		} else {
			if len(args) == 0 {
				stat, _ := os.Stdin.Stat()
				if (stat.Mode() & os.ModeCharDevice) == 0 {
					scanner := bufio.NewScanner(os.Stdin)
					for scanner.Scan() {
						line := scanner.Text()
						var replaced_line string
						if reverse {
							replaced_line = unmask(line, *maskManager)
						} else {
							replaced_line = mask(line, *patternManager, *maskManager)
						}
						fmt.Println(replaced_line)
					}
				}
			} else {
				for _, val := range strings.Split(args[0], "\n") {
					line := val
					var replaced_line string
					if reverse {
						replaced_line = unmask(line, *maskManager)
					} else {
						replaced_line = mask(line, *patternManager, *maskManager)
					}
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
	rootCmd.Flags().BoolP("reverse", "r", false, "Unmask masked text")
}
