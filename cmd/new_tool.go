package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util/editor"
)

// newToolCmd represents the tool command
var newToolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Create a new tool entry",
	Long: `EXAMPLE:
> tool-notes new tool gh
`,
	Run: func(cmd *cobra.Command, args []string) {
		tool := args[0]
		exists, _ := findTool(tool)
		if exists {
			logrus.Fatal("This tool already exists, please use 'edit'", tool)
		} else {
			newTool(tool)
		}
	},
}

func newTool(name string) {

	var buffer *bytes.Buffer
	toolName := ""
	toolFile := ""
	edit := editor.NewDefaultEditor([]string{
		"TOOLNOTES_EDITOR",
		"EDITOR",
	})

	// Find REPOS, prompt if more than 1, else select one
	// TODO: get this done, hard code to main repo for now
	repo := "/Users/christopher.maahs/.config/tool-notes/maahsome/toolnotes-base/"

	if strings.HasSuffix(name, ".yaml") {
		toolName = strings.TrimSuffix(name, ".yaml")
		toolFile = name
	} else {
		toolName = name
		toolFile = fmt.Sprintf("%s.yaml", name)
	}
	fileToOpen := fmt.Sprintf("%s%s", repo, toolFile)
	// touch the file
	// file, err := os.Create(fileToOpen)
	// if err != nil {
	// 	logrus.WithError(err).Fatal("Cannot create empty file.")
	// }
	template := `tool:
  tool_name: NAME
  sections:
  - section_name: General Commands
    examples:
    - description: What does your example accomplish
      language: bash
      script: |
        # script lines
`
	// fileContent := []byte(strings.Replace(template, "NAME", name, 1))
	// _, werr := file.Write(fileContent)
	// if werr != nil {
	// 	logrus.WithError(werr).Fatal("Failed to write template to file")
	// }
	// file.Close()

	// _, buffer := openFile(fileToOpen)

	// original := buffer.Bytes()
	original := []byte(strings.Replace(template, "NAME", toolName, 1))
	buffer = bytes.NewBuffer(make([]byte, 0))
	buffer.Write(original[:])
	edited, _, err := edit.LaunchTempFile("tn-edit", ".yaml", buffer)
	if err != nil {
		logrus.WithError(err).Error("Bad, bad")
	}

	if bytes.Equal(edited, original) {
		logrus.Info("Apply was skipped: no changes detected.")
	} else {
		logrus.Info("Applied: changes detected.")
		err := os.WriteFile(fileToOpen, edited, 0644)
		if err != nil {
			logrus.WithError(err).Error("Failed to write changes")
		}
	}
	logrus.Warn(string(edited[:]))
}

func init() {
	newCmd.AddCommand(newToolCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toolCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// toolCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
