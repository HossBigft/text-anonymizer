/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/lucasjones/reggen"
	"github.com/spf13/cobra"
	"os"
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
	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Execute() {
	IPV4_REGEX := `(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}`
	FQDN_REGEX := `(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)`
	patterns := make([]string, 3)
	patterns = append(patterns, IPV4_REGEX)
	patterns = append(patterns, FQDN_REGEX)

	err := rootCmd.Execute()
	file, err := os.Open("examples/nginx_access.log")
	check(err)
	defer file.Close()

	// r, _ := regexp.Compile(IPV4_REGEX)
	m := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		replaced_line := line
		for _, pattern := range patterns {
			regex, _ := regexp.Compile(pattern)
			sensitive_value := regex.FindString(line)
			// fmt.Println("Sensitive vlaue:" +sensitive_value)
			if len(sensitive_value) != 0 {
				mask, present := m[sensitive_value]
				if present {
					replaced_line = strings.ReplaceAll(replaced_line, sensitive_value, mask)
				} else {
					mask, _ = reggen.Generate(pattern, 1)
					m[sensitive_value] = mask
					replaced_line = strings.ReplaceAll(replaced_line, sensitive_value, mask)
				}
				// fmt.Println("Replaced with : "+mask)
				// fmt.Println("Original line: "+line)

			}

		}
		fmt.Println(replaced_line)
	}

	check(err)

	if err != nil {
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
