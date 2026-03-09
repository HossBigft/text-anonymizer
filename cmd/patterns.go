package cmd

import (
	patternmanager "ae/patternManager"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var MaskPatternToAdd patternmanager.MaskPattern
var patternsCmd = &cobra.Command{
	Use:   "patterns",
	Short: "Management of mask patterns",

	Run: func(cmd *cobra.Command, args []string) {
		patternManager := patternmanager.NewPatternManager()
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}

		MaskPatternsToAdd, _ := cmd.Flags().GetStringArray("add")
		if len(MaskPatternsToAdd) != 0 {
			for _, entry := range MaskPatternsToAdd {
				parts := strings.Split(entry, "=")
				if len(parts) != 2 {
					fmt.Fprintf(os.Stderr, "invalid entry %q, expected name=regex", entry)
				} else {
					newPattern := patternmanager.MaskPattern{Name: parts[0], Regex: parts[1]}
					patternManager.AddPattern(newPattern)
					patternManager.SavePatterns()
					fmt.Printf("Added pattern %q", newPattern)
				}
			}
		}

		list, _ := cmd.Flags().GetBool("list")
		if list {
			maskPatterns := patternManager.GetPatterns()
			for _, pattern := range maskPatterns {
				fmt.Println("Name: " + pattern.Name + ", Pattern: " + pattern.Regex)
			}
		}
		patternNameToDelete, _ := cmd.Flags().GetString("delete")
		if len(patternNameToDelete) != 0 {
			removedPattern, notFound := patternManager.RemovePatternByName(patternNameToDelete)
			if notFound != nil {
				fmt.Fprintf(os.Stderr, "Pattern with name %q not found", patternNameToDelete)
			} else {
				fmt.Printf("Removed pattern %q", removedPattern)
				patternManager.SavePatterns()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(patternsCmd)
	patternsCmd.Flags().BoolP("list", "l", false, "List mask patterns")
	patternsCmd.Flags().StringArrayP("add", "a", []string{}, "Add mask pattern in name=regex format")
	patternsCmd.Flags().StringP("delete", "d", "", "Delete mask patter by name")
}
