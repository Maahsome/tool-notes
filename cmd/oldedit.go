package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	// "k8s.io/kubectl/pkg/cmd/util/editor"
	// "k8s.io/kubectl/pkg/cmd/util/editor/crlf"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util/editor"
)

const (
	chunksize int = 1024
)

// oldeditCmd represents the edit command
var oldeditCmd = &cobra.Command{
	Use:   "oldedit",
	Short: "Edit a tool notes yaml file",
	Long: `EXAMPLE

	To edit the notes for tool 'jq'

	> tool-notes oldedit jq`,
	Run: func(cmd *cobra.Command, args []string) {
		tool := args[0]
		exists, yamlFile := findTool(tool)
		if exists {
			editTool(yamlFile)
		} else {
			logrus.Warn("Could not find tool ", tool)
		}
	},
}

func editTool(file []string) {

	edit := editor.NewDefaultEditor([]string{
		"TOOLNOTES_EDITOR",
		"EDITOR",
	})

	fileToOpen := file[0]
	if len(file) > 1 {
		// the questions to ask
		var fileSurvey = []*survey.Question{
			{
				Name: "filename",
				Prompt: &survey.Select{
					Message: "Choose a file:",
					Options: file,
				},
			},
		}

		opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

		// perform the questions
		if err := survey.Ask(fileSurvey, fileAnswers, opts); err != nil {
			logrus.Fatal("No section on the list")
		}
		fmt.Printf("Selected Section: %s\n", fileAnswers.FileName)
		fileToOpen = fileAnswers.FileName
	}
	_, buffer := openFile(fileToOpen)

	original := buffer.Bytes()

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
	// logrus.Warn(string(edited[:]))
}

func openFile(name string) (byteCount int, buffer *bytes.Buffer) {

	var (
		data  *os.File
		part  []byte
		err   error
		count int
	)

	data, err = os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	reader := bufio.NewReader(data)
	buffer = bytes.NewBuffer(make([]byte, 0))
	part = make([]byte, chunksize)

	for {
		if count, err = reader.Read(part); err != nil {
			break
		}
		buffer.Write(part[:count])
	}
	if err != io.EOF {
		log.Fatal("Error Reading ", name, ": ", err)
	} else {
		err = nil
	}

	byteCount = buffer.Len()
	return
}

func init() {
	rootCmd.AddCommand(oldeditCmd)
}
