package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/AlecAivazis/survey/v2"
	markdown "github.com/Maahsome/go-term-markdown"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var cfgFile string

type (
	ToolNote struct {
		Tool struct {
			ToolName string `yaml:"tool_name"`
			Sections []struct {
				SectionName string    `yaml:"section_name"`
				Examples    []Example `yaml:"examples"`
			} `yaml:"sections"`
		} `yaml:"tool"`
	}

	Example struct {
		Description     string `yaml:"description"`
		LongDescription string `yaml:"long_description"`
		Language        string `yaml:"language"`
		Script          string `yaml:"script"`
	}
	sectionAnswer struct {
		SectionName string `survey:"sectionname"` // or you can tag fields to match a specific name
	}
	exampleAnswer struct {
		ExampleDescription string `survey:"examplename"` // or you can tag fields to match a specific name
	}
)

var (
	sectionAnswers = &sectionAnswer{}
	exampleAnswers = &exampleAnswer{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tool-notes",
	Short: "A cli tool to assist with remembering other cli tool command syntax",
	Long: `tool-notes is a place to capture those examples you learn from your web-searches,
	storing them locally and in a format that allows you to quickly modify and utilize them.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		tool := args[0]
		exists, yamlFile := toolExists(tool)
		if exists {
			yamlData, err := ioutil.ReadFile(yamlFile)
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
			source += fmt.Sprintf("# %s\n\n", sectionAnswers.SectionName)
			source += fmt.Sprintf("## EXAMPLE - %s\n\n", exampleAnswers.ExampleDescription)
			if len(exampleData[exampleAnswers.ExampleDescription].LongDescription) > 0 {
				source += fmt.Sprintf("%s\n\n", exampleData[exampleAnswers.ExampleDescription].LongDescription)
			}
			source += fmt.Sprintf("```%s\n", exampleData[exampleAnswers.ExampleDescription].Language)
			source += fmt.Sprintf("%s\n", exampleData[exampleAnswers.ExampleDescription].Script)
			source += "```"

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

func toolExists(t string) (bool, string) {
	home, err := homedir.Dir()
	if err != nil {
		logrus.Error("Could not locate the HOME directory")
		return false, ""
	}
	toolFile := fmt.Sprintf("%s/.config/tool-notes/%s.yaml", home, t)
	if _, err := os.Stat(toolFile); err != nil {
		if os.IsNotExist(err) {
			return false, ""
		}
	}
	return true, toolFile
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tool-notes.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tool-notes" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tool-notes")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
