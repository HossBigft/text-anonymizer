/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	maskManager "anonymizer/maskManager"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// masksCmd represents the masks command
var masksCmd = &cobra.Command{
	Use:   "masks",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}
		list, _ := cmd.Flags().GetBool("list")

		if list {
			maskPatterns := maskManager.NewMaskManager().GetMaskMap()
			for value, mask := range maskPatterns {
				fmt.Println(value + " => " + mask)
			}
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(masksCmd)

	masksCmd.Flags().BoolP("list", "l", false, "List masks")
}
