package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"maahsome/tool-notes/objects"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Express the 'version' of compare-repo-members.",
	Run: func(cmd *cobra.Command, args []string) {
		var verData objects.Version
		err := json.Unmarshal([]byte(c.VersionJSON), &verData)
		if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal JSON")
		}

		if !c.FormatOverridden {
			c.OutputFormat = "raw"
		}
		fmt.Println(versionDataToString(verData))

	},
}

func versionDataToString(versionData objects.Version) string {

	switch strings.ToLower(c.OutputFormat) {
	case "raw":
		return c.VersionJSON
	case "json":
		return versionData.ToJSON()
	case "gron":
		return versionData.ToGRON()
	case "yaml":
		return versionData.ToYAML()
	case "text", "table":
		return versionData.ToTEXT(c.NoHeaders)
	default:
		return versionData.ToTEXT(c.NoHeaders)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
