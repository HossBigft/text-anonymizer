/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"anonymizer/patternManager"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/lucasjones/reggen"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

		configDir := filepath.Join(os.Getenv("HOME") + "/.config/anonymizer/")
		maskedValuesFilePath := os.Getenv("HOME") + "/.config/anonymizer/map.json"
		maskManager := patternmanager.NewMaskManager()
		maskPatterns := maskManager.GetPatterns()
		isMaskPatternsUpdated := false
		file, err := os.Open("examples/nginx_access.log")
		check(err)
		defer file.Close()

		valuesToMasks := make(map[string]string)
		maskedValuesFileHandle, err := os.ReadFile(maskedValuesFilePath)
		err = json.Unmarshal(maskedValuesFileHandle, &valuesToMasks)

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
						mask, present := valuesToMasks[sensitive_value]
						if present {
							replaced_line = strings.ReplaceAll(replaced_line, sensitive_value, mask)
						} else {
							mask, _ = reggen.Generate(pattern.Regex, 7)
							valuesToMasks[sensitive_value] = mask
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
						mask, present := valuesToMasks[sensitive_value]
						if present {
							replaced_line = strings.ReplaceAll(replaced_line, sensitive_value, mask)
						} else {
							mask, _ = reggen.Generate(pattern.Regex, 7)
							valuesToMasks[sensitive_value] = mask
							replaced_line = strings.ReplaceAll(replaced_line, sensitive_value, mask)
							isMasksUpdated = true
						}
					}

				}
				fmt.Println(replaced_line)
			}

		}

		if isMasksUpdated {
			err = os.MkdirAll(configDir, 0755)
			valueMapJson, err := json.Marshal(valuesToMasks)
			check(err)
			err = os.WriteFile(maskedValuesFilePath, valueMapJson, 0644)
			check(err)
		}

		if isMaskPatternsUpdated {
			maskManager.SavePatterns()
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
