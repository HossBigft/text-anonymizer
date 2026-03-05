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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "anonymizer",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		maskManager := maskmanager.NewMaskManager()
		patternManager := patternmanager.NewPatternManager()
		maskPatterns := patternManager.GetPatterns()
		isMaskPatternsUpdated := false

		var isMasksUpdated bool
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

type (
	MaskPattern struct {
		Name  string `json:"name"`
		Regex string `json:"regex"`
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.anonymizer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
