package cmd

import (
	"fmt"
	"os"
	"strings"

	maskManager "github.com/HossBigft/ae/maskManager"
	maskmanager "github.com/HossBigft/ae/maskManager"

	"github.com/spf13/cobra"
)

var masksCmd = &cobra.Command{
	Use:   "masks",
	Short: "Management of created masks for sensitive values",

	Run: func(cmd *cobra.Command, args []string) {
		maskManager := maskManager.NewMaskManager()
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			os.Exit(0)
		}

		list, _ := cmd.Flags().GetBool("list")
		if list {
			maskPatterns := maskManager.GetMaskMap()
			for value, mask := range maskPatterns {
				fmt.Println(value + " => " + mask)
			}
			os.Exit(0)
		}

		maskValueToDelete, _ := cmd.Flags().GetString("delete")
		if len(maskValueToDelete) != 0 {
			removedMask, notFound := maskManager.RemoveMaskByValue(maskValueToDelete)
			if notFound != nil {
				fmt.Fprintf(os.Stderr, "Mask with value %q not found", maskValueToDelete)
			} else {
				fmt.Fprintf(os.Stderr, "Removed mask %q", removedMask)
				maskManager.SaveMasks()
			}
		}
		MasksToAdd, _ := cmd.Flags().GetStringArray("add")
		if len(MasksToAdd) != 0 {
			for _, entry := range MasksToAdd {
				parts := strings.Split(entry, "=")
				if len(parts) != 2 {
					fmt.Fprintf(os.Stderr, "invalid entry %q, expected value=regex", entry)
				} else {
					newMask := maskmanager.ValueMask{Value: parts[0], Mask: parts[1]}
					maskManager.AddPattern(newMask)
					maskManager.SaveMasks()
					fmt.Fprintf(os.Stderr, "Added mask %q", newMask)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(masksCmd)

	masksCmd.Flags().BoolP("list", "l", false, "List masks")
	masksCmd.Flags().StringP("delete", "d", "", "Delete mask by value")
	masksCmd.Flags().StringArrayP("add", "a", []string{}, "Add mask in value=mask format")
}
