package cmd

import (
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new objects in the Tool Notes data",
	Long: `EXAMPLE:
> tool-notes new tool <toolname>`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("new called")
	// },
}

func init() {
	rootCmd.AddCommand(newCmd)
}
