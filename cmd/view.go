/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/AlecAivazis/survey/v2"
	markdown "github.com/Maahsome/go-term-markdown"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/sirupsen/logrus"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:     "view",
	Aliases: []string{"show"},
	Short:   "Show the notes for the specified tool",
	Long: `EXAMPLE:
	
	> tool-notes view jq`,
	Run: func(cmd *cobra.Command, args []string) {
		tool := args[0]
		exists, yamlFile := findTool(tool)
		if exists {
			yamlData, err := ioutil.ReadFile(yamlFile[0])
			if err != nil {
				fmt.Printf("Error reading YAML file: %s\n", err)
				return
			}

			var toolNote ToolNote
			err = yaml.Unmarshal(yamlData, &toolNote)
			if err != nil {
				fmt.Printf("Error parsing YAML file: %s\n", err)
			}
			var sectionArray []string
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

			// Examples
			var exampleArray []string
			exampleData := make(map[string]Example)
			exampleArray = append(exampleArray, "< ALL EXAMPLES >")
			for i, v := range toolNote.Tool.Sections {
				if v.SectionName == sectionAnswers.SectionName {
					for _, example := range toolNote.Tool.Sections[i].Examples {
						exampleArray = append(exampleArray, example.Description)
						exampleData[example.Description] = example
					}
				}
			}
			// the questions to ask
			var exampleSurvey = []*survey.Question{
				{
					Name: "examplename",
					Prompt: &survey.Select{
						Message: "Choose an example:",
						Options: exampleArray,
					},
				},
			}

			// perform the questions
			if err = survey.Ask(exampleSurvey, exampleAnswers, opts); err != nil {
				logrus.Fatal("No section on the list")
			}
			fmt.Printf("Selected Example: %s\n", exampleAnswers.ExampleDescription)

			source := ""
			padBetween := ""
			source += fmt.Sprintf("# %s\n\n", sectionAnswers.SectionName)
			if exampleAnswers.ExampleDescription == "< ALL EXAMPLES >" {
				for i, v := range toolNote.Tool.Sections {
					if v.SectionName == sectionAnswers.SectionName {
						for _, example := range toolNote.Tool.Sections[i].Examples {
							source += fmt.Sprintf("%s## EXAMPLE - %s\n\n", padBetween, example.Description)
							if len(example.LongDescription) > 0 {
								source += fmt.Sprintf("%s\n\n", example.LongDescription)
							}
							source += fmt.Sprintf("```%s\n", example.Language)
							source += fmt.Sprintf("%s\n", example.Script)
							source += "```"
							padBetween = "\n\n"
						}

					}
				}
			} else {

				source += fmt.Sprintf("## EXAMPLE - %s\n\n", exampleAnswers.ExampleDescription)
				if len(exampleData[exampleAnswers.ExampleDescription].LongDescription) > 0 {
					source += fmt.Sprintf("%s\n\n", exampleData[exampleAnswers.ExampleDescription].LongDescription)
				}
				source += fmt.Sprintf("```%s\n", exampleData[exampleAnswers.ExampleDescription].Language)
				source += fmt.Sprintf("%s\n", exampleData[exampleAnswers.ExampleDescription].Script)
				source += "```"
			}
			w, _ := consolesize.GetConsoleSize()
			result := markdown.Render(source, w, 0)

			fmt.Println(string(result[:]))

			// path := "/Users/cmaahs/GDrive/src/Markdown/md-notebooks/Tools_Notes/jq/Table_Formatting.md"
			// mdsrc, err := ioutil.ReadFile(path)
			// if err != nil {
			// 	panic(err)
			// }

			// srcresult := markdown.Render(string(mdsrc), 80, 6)

			// // fmt.Printf("%v", srcresult)
			// fmt.Println(string(srcresult[:]))
		}
	},
}

func fetchToolFiles(dir_path string, tool_name string) []string {

	files := []string{}

	filepath.Walk(dir_path, func(path string, f os.FileInfo, err error) error {

		tool_file := fmt.Sprintf("%s/%s.yaml", path, tool_name)
		// logrus.Info(fmt.Sprintf("tool_file: %s", tool_file))
		if _, err := os.Stat(tool_file); err != nil {
			if os.IsNotExist(err) {
				return nil
			}
		}
		files = append(files, tool_file)
		return nil
	})

	return files
}

func findTool(t string) (bool, []string) {
	home, err := homedir.Dir()
	if err != nil {
		logrus.Error("Could not locate the HOME directory")
		return false, []string{}
	}
	toolPath := fmt.Sprintf("%s/.config/tool-notes/", home)
	tool_list := fetchToolFiles(toolPath, t)
	// toolFile := fmt.Sprintf("%s/.config/tool-notes/%s.yaml", home, t)
	// if _, err := os.Stat(toolFile); err != nil {
	// 	if os.IsNotExist(err) {
	// 		return false, ""
	// 	}
	// }
	return true, tool_list
}

func init() {
	rootCmd.AddCommand(viewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// viewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// viewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
