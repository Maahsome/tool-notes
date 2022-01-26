package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/miracl/conflate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	// "bytes"
	// "fmt"
	// "os"
	// "strings"
	// "github.com/sirupsen/logrus"
	// "github.com/spf13/cobra"
	// "k8s.io/kubectl/pkg/cmd/util/editor"
)

// newExampleCmd represents the tool command
var newExampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Create a new example entry for a tool",
	Long: `EXAMPLE:
> tool-notes new example gh
`,
	Run: func(cmd *cobra.Command, args []string) {
		// tool := args[0]
		// exists, _ := findTool(tool)
		// if exists {
		// 	newExample(tool)
		// } else {
		// 	logrus.Fatal("This tool does not exist, please use 'new tool <toolname>'", tool)
		// }
		var err error
		var tool string

		if len(args) == 0 {
			// survey for tool name
			_, toolArray := getToolList()
			// the questions to ask
			var toolSurvey = []*survey.Question{
				{
					Name: "toolname",
					Prompt: &survey.Select{
						Message: "Choose a Tool:",
						Options: toolArray,
					},
				},
			}

			opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

			// perform the questions
			if err = survey.Ask(toolSurvey, toolAnswers, opts); err != nil {
				logrus.Fatal("No section on the list")
			}
			fmt.Printf("Selected Tool: %s\n", toolAnswers.ToolName)
			tool = toolAnswers.ToolName
		} else {
			tool = args[0]
		}

		exists, yamlFile := findTool(tool)
		if exists {
			var toolNote ToolNote

			mergedYaml, err := conflate.FromFiles(yamlFile...)
			if err != nil {
				logrus.WithError(err).Error("Failed to merge YAML Files")
				return
			}

			rawYaml, err := mergedYaml.MarshalYAML()
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := yaml.Unmarshal(rawYaml, &toolNote); err != nil {
				logrus.Info("DEBUG: failed to reparse our base structure")
			}

			var sectionArray []string
			sectionArray = append(sectionArray, "< NEW SECTION >")
			for _, v := range toolNote.Tool.Sections {
				sectionArray = append(sectionArray, v.SectionName)
			}
			// the questions to ask
			var sectionSurvey = []*survey.Question{
				{
					Name: "sectionname",
					Prompt: &survey.Select{
						Message: "Choose a section:",
						Options: sectionArray,
					},
				},
			}

			opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

			// perform the questions
			if err = survey.Ask(sectionSurvey, sectionAnswers, opts); err != nil {
				logrus.Fatal("No section on the list")
			}
			fmt.Printf("Selected Section: %s\n", sectionAnswers.SectionName)

		}
	},
}

// func newExample(name string) {

// }

func init() {
	newCmd.AddCommand(newExampleCmd)

	conflate.Unmarshallers = conflate.UnmarshallerMap{
		".json": {conflate.JSONUnmarshal},
		".jsn":  {conflate.JSONUnmarshal},
		".yaml": {conflate.YAMLUnmarshal},
		".yml":  {conflate.YAMLUnmarshal},
		".toml": {conflate.TOMLUnmarshal},
		".tml":  {conflate.TOMLUnmarshal},
		"":      {conflate.JSONUnmarshal, conflate.YAMLUnmarshal, conflate.TOMLUnmarshal},
	}
}
