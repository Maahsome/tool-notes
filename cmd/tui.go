package cmd

import (
	"os"

	"maahsome/tool-notes/common"
	"maahsome/tool-notes/tui"

	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"tui"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ret := startTUI()
		os.Exit(ret)
	},
}

func startTUI() int {

	tui := tui.New(semVer)

	if err := tui.Start(); err != nil {
		common.Logger.Errorf("cannot start tool-notes tui mode: %s", err)
		return 1
	}

	return 0
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
