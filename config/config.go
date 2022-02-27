package config

import (
	"fmt"
	"maahsome/tool-notes/objects"
	"strings"
)

type (
	Config struct {
		VersionDetail    objects.Version
		VersionJSON      string
		OutputFormat     string
		FormatOverridden bool
		NoHeaders        bool
		CACert           string
		CABundle         string
	}
	Outputtable interface {
		ToJSON() string
		ToYAML() string
		ToGRON() string
		ToTEXT(noHeaders bool) string
	}
)

func (c *Config) outputData(data Outputtable) string {
	switch strings.ToLower(c.OutputFormat) {
	case "json":
		return data.ToJSON()
	case "gron":
		return data.ToGRON()
	case "yaml":
		return data.ToYAML()
	case "text", "table":
		return data.ToTEXT(c.NoHeaders)
	default:
		return ""
	}
}

func (c *Config) OutputData(data Outputtable) {
	if !c.FormatOverridden {
		c.OutputFormat = "text"
	}

	fmt.Println(c.outputData(data))
}
