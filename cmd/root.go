package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	maskmanager "github.com/HossBigft/ae/maskManager"
	patternmanager "github.com/HossBigft/ae/patternManager"

	"github.com/spf13/cobra"
)

func mask(rawLine string, patternManager patternmanager.PatternManager, maskManager maskmanager.MaskManager) string {
	replaced_line := rawLine
	valuesToMaskMap, _ := patternManager.MapValuesToPatterns(rawLine)
	for _, match := range valuesToMaskMap {
		masks := maskManager.MapValuesToMasks(match)
		for sensitive_value, mask := range masks {
			replaced_line = strings.ReplaceAll(strings.ToLower(replaced_line), sensitive_value, mask)
		}
	}

	return replaced_line
}
func unmask(rawLine string, maskManager maskmanager.MaskManager) string {
	line_to_replace := rawLine
	masks := maskManager.GetMasksToValuesMap()
	for sensitive_value, mask := range masks {
		line_to_replace = strings.ReplaceAll(strings.ToLower(line_to_replace), sensitive_value, mask)
	}

	return line_to_replace
}

var filePath string
var rootCmd = &cobra.Command{
	Use:   "ae",
	Short: "Mask sensitive values in text with lookalike text",
	Long: `ae is CLI program to mask sensitive values in given text and print masked text. It replaces values with randomly generated strings based on regex patterns used for extracting. Masks are persistent and saved on disk. Also it can unmask masked text back. 
	Usage:
	ae example.com
	cat example.txt | ae
	ae -f path/example.txt
	`,
	Run: func(cmd *cobra.Command, args []string) {
		maskManager := maskmanager.NewMaskManager()
		patternManager := patternmanager.NewPatternManager()
		decode, _ := cmd.Flags().GetBool("decode")

		if len(filePath) > 0 {
			file, err := os.Open(filePath)
			check(err)
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				var replaced_line string
				if decode {
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
						if decode {
							replaced_line = unmask(line, *maskManager)
						} else {
							replaced_line = mask(line, *patternManager, *maskManager)
						}
						fmt.Println(replaced_line)
					}
				} else {

					cmd.Help()
				}
			} else {
				for _, val := range strings.Split(args[0], "\n") {
					line := val
					var replaced_line string
					if decode {
						replaced_line = unmask(line, *maskManager)
					} else {
						replaced_line = mask(line, *patternManager, *maskManager)
					}
					fmt.Println(replaced_line)
				}

			}
		}

	},
	Args: cobra.ArbitraryArgs,
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
	rootCmd.Flags().BoolP("decode", "d", false, "Decode masked text")
}
