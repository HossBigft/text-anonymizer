/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"bufio"
	"os"
	"regexp"
	"strings"
	"encoding/json"
	"github.com/spf13/cobra"
		"github.com/lucasjones/reggen"
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
		IPV4_REGEX := `(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}`
		FQDN_REGEX := `(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)`
		// configDir := filepath.Join(os.Getenv("HOME") + "/.config/anonymizer/")
		maskedValuesFilePath := os.Getenv("HOME") + "/.config/anonymizer/map.json"
		masksPatternsFilePath := os.Getenv("HOME") + "/.config/anonymizer/maskPatterns.json"
		maskPatterns, err := loadPatterns(masksPatternsFilePath)
		isMaskPatternsUpdated := false
		if err != nil {
			maskPatterns = append(maskPatterns, MaskPattern{Name: "ipv4", Regex: IPV4_REGEX})
			maskPatterns = append(maskPatterns, MaskPattern{Name: "fqdn", Regex: FQDN_REGEX})
			isMaskPatternsUpdated = true
		}
		fmt.Println("file called")
		filePath := args[0]
		file, err := os.Open(filePath)
		check(err)
		defer file.Close()
				valuesToMasks := make(map[string]string)
		maskedValuesFileHandle, err := os.ReadFile(maskedValuesFilePath)
		err = json.Unmarshal(maskedValuesFileHandle, &valuesToMasks)
		scanner := bufio.NewScanner(file)
		var isMasksUpdated bool
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
			fmt.Println(isMaskPatternsUpdated, isMasksUpdated)
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
