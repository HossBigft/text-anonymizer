package cmd

import (
	maskManager "ae/maskManager"
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
				fmt.Printf("Removed mask %q", removedMask)
				maskManager.SaveMasks()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(masksCmd)

	masksCmd.Flags().BoolP("list", "l", false, "List masks")
	masksCmd.Flags().StringP("delete", "d", "", "Delete mask by value")
}
