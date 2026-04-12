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

func replaceAllCaseInsensitive(s, old, new string) string {
    lowerS := strings.ToLower(s)
    lowerOld := strings.ToLower(old)
    
    if lowerOld == "" {
        return s
    }
    
    var result strings.Builder
    for {
        idx := strings.Index(lowerS, lowerOld)
        if idx == -1 {
            result.WriteString(s)
            break
        }
        result.WriteString(s[:idx])
        result.WriteString(new)
        s = s[idx+len(old):]
        lowerS = lowerS[idx+len(lowerOld):]
    }
    return result.String()
}

func mask(rawLine string, patternManager patternmanager.PatternManager, maskManager maskmanager.MaskManager) string {
	line_to_replace := rawLine
	valuesToMaskMap, _ := patternManager.MapValuesToPatterns(rawLine)
	for _, match := range valuesToMaskMap {
		masks := maskManager.MapValuesToMasks(match)
		for sensitive_value, mask := range masks {
			line_to_replace = replaceAllCaseInsensitive(line_to_replace, sensitive_value, mask)
		}
	}

	return line_to_replace
}
func unmask(rawLine string, maskManager maskmanager.MaskManager) string {
	line_to_replace := rawLine
	masks := maskManager.GetMasksToValuesMap()
	for sensitive_value, mask := range masks {
		line_to_replace = replaceAllCaseInsensitive(line_to_replace, sensitive_value, mask)
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
