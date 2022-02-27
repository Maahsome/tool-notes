package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"maahsome/tool-notes/common"
	"maahsome/tool-notes/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	cfgFile   string
	semVer    string
	gitCommit string
	gitRef    string
	buildDate string
	// gitlabClient  *gitlab.Gitlab
	// gitlabToken   string
	// gitlabHost    string
	// repo          *git.Repository

	semVerReg = regexp.MustCompile(`(v[0-9]+\.[0-9]+\.[0-9]+).*`)
	// ecrRegex  = regexp.MustCompile(`^(?P<registry>\d+)\.dkr\.ecr.\w+-\w+-\d\.amazonaws\.com/(?P<repository>.+):(?P<tag>.+)$`)

	c = &config.Config{}
)

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
	fileAnswer struct {
		FileName string `survey:"filename"` // or you can tag fields to match a specific name
	}
	toolAnswer struct {
		ToolName string `survey:"toolname"` // or you can tag fields to match a specific name
	}
	sectionAnswer struct {
		SectionName string `survey:"sectionname"` // or you can tag fields to match a specific name
	}
	exampleAnswer struct {
		ExampleDescription string `survey:"examplename"` // or you can tag fields to match a specific name
	}
)

var (
	fileAnswers    = &fileAnswer{}
	toolAnswers    = &toolAnswer{}
	sectionAnswers = &sectionAnswer{}
	exampleAnswers = &exampleAnswer{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tool-notes",
	Short: "A cli tool to assist with remembering other cli tool command syntax",
	Long: `tool-notes is a place to capture those examples you learn from your web-searches,
	storing them locally and in a format that allows you to quickly modify and utilize them.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		c.VersionDetail.SemVer = semVer
		c.VersionDetail.BuildDate = buildDate
		c.VersionDetail.GitCommit = gitCommit
		c.VersionDetail.GitRef = gitRef
		c.VersionJSON = fmt.Sprintf("{\"SemVer\": \"%s\", \"BuildDate\": \"%s\", \"GitCommit\": \"%s\", \"GitRef\": \"%s\"}", semVer, buildDate, gitCommit, gitRef)
		if c.OutputFormat != "" {
			c.FormatOverridden = true
			c.NoHeaders = false
			c.OutputFormat = strings.ToLower(c.OutputFormat)
			switch c.OutputFormat {
			case "json", "gron", "yaml", "text", "table", "raw":
				break
			default:
				fmt.Println("Valid options for -o are [json|gron|text|table|yaml|raw]")
				os.Exit(1)
			}
		}
		if os.Args[1] != "version" && os.Args[1] != "config" {
			// Do PRE-SETUP Work here
			logFile, _ := cmd.Flags().GetString("log-file")
			logLevel, _ := cmd.Flags().GetString("log-level")
			ll := "Warning"
			switch strings.ToLower(logLevel) {
			case "trace":
				ll = "Trace"
			case "debug":
				ll = "Debug"
			case "info":
				ll = "Info"
			case "warning":
				ll = "Warning"
			case "error":
				ll = "Error"
			case "fatal":
				ll = "Fatal"
			}

			common.NewLogger(ll, logFile)

		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tool-notes.yaml)")
	rootCmd.PersistentFlags().StringVarP(&c.OutputFormat, "output", "o", "", "Set an output format: json, text, yaml, gron")
	rootCmd.PersistentFlags().StringP("log-file", "l", "", "Specify a log file to log events to, default to no logging")
	rootCmd.PersistentFlags().StringP("log-level", "v", "", "Specify a log level for logging, default to Warning (Trace, Debug, Info, Warning, Error, Fatal)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		workDir := fmt.Sprintf("%s/.config/tool-notes", home)
		if _, err := os.Stat(workDir); err != nil {
			if os.IsNotExist(err) {
				mkerr := os.MkdirAll(workDir, os.ModePerm)
				if mkerr != nil {
					logrus.Fatal("Error creating ~/.config/tool-notes directory", mkerr)
				}
			}
		}
		if stat, err := os.Stat(workDir); err == nil && stat.IsDir() {
			configFile := fmt.Sprintf("%s/%s", workDir, "config.yaml")
			createRestrictedConfigFile(configFile)
			viper.SetConfigFile(configFile)
		} else {
			logrus.Info("The ~/.config/tool-notes path is a file and not a directory, please remove the 'tool-notes' file.")
			os.Exit(1)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		logrus.Warn("Failed to read viper config file.")
	}
}

func createRestrictedConfigFile(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			file, ferr := os.Create(fileName)
			if ferr != nil {
				logrus.Info("Unable to create the configfile.")
				os.Exit(1)
			}
			mode := int(0600)
			if cherr := file.Chmod(os.FileMode(mode)); cherr != nil {
				logrus.Info("Chmod for config file failed, please set the mode to 0600.")
			}
		}
	}
}

// exists returns whether the given file or directory exists or not
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// remove directory
func RemoveDir(path string) bool {

	err := os.RemoveAll(path)

	if err != nil {
		fmt.Println(err)
		return false
	} else {
		return true
	}
}
