package cmd

import (
	patternmanager "ae/patternManager"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var FilterPatternToAdd patternmanager.Pattern
var filterCmd = &cobra.Command{
	Use:   "filters",
	Short: "Management of filter patterns",

	Run: func(cmd *cobra.Command, args []string) {
		patternManager := patternmanager.NewPatternManager()
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}

		FilterPatternsToAdd, _ := cmd.Flags().GetStringArray("add")
		if len(FilterPatternsToAdd) != 0 {
			for _, entry := range FilterPatternsToAdd {
				parts := strings.Split(entry, "=")
				if len(parts) != 2 {
					fmt.Fprintf(os.Stderr, "invalid entry %q, expected name=regex", entry)
				} else {
					newFilter := patternmanager.Pattern{Name: parts[0], Regex: parts[1]}
					patternManager.AddFilter(newFilter)
					patternManager.SaveFilters()
					fmt.Fprintf(os.Stderr, "Added filter %q", newFilter)
				}
			}
		}

		list, _ := cmd.Flags().GetBool("list")
		if list {
			filterPatterns := patternManager.GetFilters()
			for _, filter := range filterPatterns {
				fmt.Println("Name: " + filter.Name + ", Pattern: " + filter.Regex)
			}
		}
		filterNameToDelete, _ := cmd.Flags().GetString("delete")
		if len(filterNameToDelete) != 0 {
			removedPattern, notFound := patternManager.RemoveFilterByName(filterNameToDelete)
			if notFound != nil {
				fmt.Fprintf(os.Stderr, "Filter with name %q not found", filterNameToDelete)
			} else {
				fmt.Fprintf(os.Stderr, "Removed filter %q", removedPattern)
				patternManager.SaveFilters()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(filterCmd)
	filterCmd.Flags().BoolP("list", "l", false, "List filter patterns")
	filterCmd.Flags().StringArrayP("add", "a", []string{}, "Add filter pattern in name=regex format")
	filterCmd.Flags().StringP("delete", "d", "", "Delete filter pattern by name")
}
