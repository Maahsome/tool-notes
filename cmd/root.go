package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
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
	// Run: func(cmd *cobra.Command, args []string) {
	// },
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
