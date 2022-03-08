package tui

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"
	"gopkg.in/yaml.v2"

	// "github.com/maahsome/vault-view/clipboard"
	// "github.com/atotto/clipboard"

	"maahsome/tool-notes/common"
	"maahsome/tool-notes/resource"

	"github.com/sirupsen/logrus"
)

type section struct {
	Name     string
	Examples []Example
}

type sections struct {
	*tview.Table
	filterWord string
	showTypes  int
	lang       *resource.Lang
}

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
)

func newSections(t *Tui) *sections {
	sections := &sections{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		lang:  t.lang,
	}

	common.Logger.WithFields(logrus.Fields{
		"unit":     "sections",
		"function": "new",
	}).Debug("New Sections Panel Created")
	sections.SetTitle(fmt.Sprintf(" [[ %s ]] ", t.lang.GetText("ui", "Sections"))).SetTitleAlign(tview.AlignLeft)
	sections.SetBorder(true)
	sections.SetBorderColor(tcell.ColorDeepSkyBlue)
	sections.setEntries(t, enterSection)
	sections.setKeybinding(t)
	return sections
}

func (i *sections) name() string {
	return "sections"
}

func (i *sections) setTitle() {

	if len(i.filterWord) > 0 {
		i.SetTitle(fmt.Sprintf(" [[ %s ]] - /%s/ ", i.lang.GetText("ui", "Sections"), i.filterWord)).SetTitleAlign(tview.AlignLeft)
	} else {
		i.SetTitle(fmt.Sprintf(" [[ %s ]] ", i.lang.GetText("ui", "Sections"))).SetTitleAlign(tview.AlignLeft)
	}

}

func (i *sections) setKeybinding(t *Tui) {
	i.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		selectedSection := t.selectedSection()
		switch event.Key() {
		case tcell.KeyEnter:
			// if selectedFolder != nil {
			// 	if selectedFolder.Type == vaultFolder {
			// 		i.filterWord = ""
			// 		i.setTitle()
			// 		i.setEntries(t, enterFolder)
			// 	} else {
			// 		if len(t.state.resources.datas) > 0 {
			// 			t.dataPanel().setEntries(t, enterFolder)
			// 			t.switchPanel("datas")
			// 		}
			// 	}
			// }
		// KeyDown/KeyUp happen BEFORE the actual selection changes
		// So we temporarily increase by one, then decrease and allow the
		// default behavior to increase it after we display the path datas.
		case tcell.KeyCtrlSpace:
			// if selectedFolder != nil {
			// 	i.toggleSelected(t, selectedFolder)
			// }
		case tcell.KeyDown:
			row, _ := t.sectionPanel().GetSelection()
			common.Logger.WithFields(logrus.Fields{
				"unit":     "sections",
				"function": "keystrokes",
			}).Debug(fmt.Sprintf("down sections/row: %d/%d", len(t.state.resources.sections), row))
			if row < len(t.state.resources.sections) {
				tempRow := row + 1
				i.Select(tempRow, 0)
				t.examplePanel().setEntries(t, enterExample)
				t.examplePanel().Select(0, 0)
				i.Select(row, 0)
				t.state.example.SetText(t.state.resources.examples[t.state.resources.exampleRows[0]].Name)
			}
			// selectedRepository := t.selectedRepository()
			// selectedTool := t.selectedTool()
			// if selectedRepository != nil && selectedTool != nil && selectedSection != nil {
			// 	if t.state.location != nil {
			// 		t.state.location.update(fmt.Sprintf("\n [white]%s / %s (.Section=%s)", selectedRepository.Name, selectedTool.Name, selectedSection.Name))
			// 	}
			// }
		case tcell.KeyUp:
			row, _ := t.sectionPanel().GetSelection()
			common.Logger.WithFields(logrus.Fields{
				"unit":     "sections",
				"function": "keystrokes",
			}).Debug(fmt.Sprintf("up sections/row: %d/%d", len(t.state.resources.sections), row))
			if row > 0 {
				tempRow := row - 1
				i.Select(tempRow, 0)
				t.examplePanel().setEntries(t, enterExample)
				t.examplePanel().Select(0, 0)
				i.Select(row, 0)
				t.state.example.SetText(t.state.resources.examples[t.state.resources.exampleRows[0]].Name)
			}
			// selectedRepository := t.selectedRepository()
			// selectedTool := t.selectedTool()
			// currentSelectedSection := t.selectedSection()
			// if selectedRepository != nil && selectedTool != nil && currentSelectedSection != nil {
			// 	if t.state.location != nil {
			// 		t.state.location.update(fmt.Sprintf("\n [white]%s / %s (.Section=%s)", selectedRepository.Name, selectedTool.Name, currentSelectedSection.Name))
			// 	}
			// }
		case tcell.KeyRight:
			common.Logger.WithFields(logrus.Fields{
				"unit":     "sections",
				"function": "keystrokes",
			}).Debug("KeyRight")

			if selectedSection != nil {
				// 	row, _ := t.folderPanel().GetSelection()
				// 	common.Logger.WithFields(logrus.Fields{
				// 		"unit":        "sections",
				// 		"function":    "keystrokes",
				// 		"row":         row,
				// 		"parent":      selectedFolder.Parent,
				// 		"folder_type": selectedFolder.Type,
				// 	}).Info("RIGHT: Remember Where I was")
				// 	t.state.resources.rowTracker[selectedFolder.Parent] = row
				// 	if selectedFolder.Type == vaultFolder {
				// 		i.filterWord = ""
				// 		i.setTitle()
				// 		i.setEntries(t, enterFolder)
				// 	} else {
				// 		t.dataPanel().setEntries(t, enterFolder)
				// 	}
			}
		case tcell.KeyLeft, tcell.KeyEscape:
			common.Logger.WithFields(logrus.Fields{
				"unit":     "sections",
				"function": "keystrokes",
			}).Info("KeyLeft")
			// row, _ := t.folderPanel().GetSelection()
			// if selectedFolder != nil {
			// 	common.Logger.WithFields(logrus.Fields{
			// 		"unit":     "sections",
			// 		"function": "keystrokes",
			// 		"row":      row,
			// 		"parent":   selectedFolder.Parent,
			// 	}).Info("LEFT: Remember Where I was")
			// 	t.state.resources.rowTracker[selectedFolder.Parent] = row
			// }
			// i.filterWord = ""
			// i.setTitle()
			// i.setEntries(t, enterParent)
		}

		var showValuesRune rune
		var runeerr error
		showValuesRune, runeerr = i.lang.GetRune("kbd", "s")
		if runeerr != nil {
			common.Logger.WithError(runeerr).WithFields(logrus.Fields{
				"unit":     "sections",
				"function": "keystrokes",
				"rune":     "s",
			}).Error("Rune undefined")
		}
		switch event.Rune() {
		case showValuesRune:
			// obscured = !obscured
			// selectedFolder := t.selectedFolder()
			// if selectedFolder.Type == vaultData {
			// 	t.dataPanel().setEntries(t, enterFolder)
			// }
		case ' ':
			// i.toggleSelected(t, selectedFolder)
		case 'c':
			// if selectedFolder.Type == vaultData {
			// 	i.selectedCmdToClipboard(t, selectedFolder)
			// } else {
			// 	origText := t.state.info.Status.GetText(false)
			// 	t.state.info.Status.SetText(fmt.Sprintf("%s%s", origText, t.lang.GetText("ui", "Only Valid for Data Type")))
			// 	go t.ClearIndicator(1)
			// }
		case 'C':
			// i.markedCmdToClipboard(t)
		}

		return event
	})
}

func (i *sections) toggleSelected(t *Tui, selectedSection *section) {
	common.Logger.WithFields(logrus.Fields{
		"unit":     "sections",
		"function": "marking",
	}).Debug(fmt.Sprintf("MarkThis: %s", selectedSection.Name))

	row, _ := t.sectionPanel().GetSelection()
	rowColor := tcell.ColorLightBlue
	for col := 0; col <= t.sectionPanel().GetColumnCount(); col++ {
		t.sectionPanel().GetCell(row, col).SetTextColor(rowColor)
	}
}

func (i *sections) buildPanelData(t *Tui, operation int) {

	common.Logger.WithFields(logrus.Fields{
		"unit":     "sections",
		"function": "data",
	}).Trace("Running buildPanelData for sections")

	selectedTool := t.selectedTool()
	common.Logger.WithFields(logrus.Fields{
		"unit":     "sections",
		"function": "data",
	}).Trace(fmt.Sprintf("Selected Tool: %#v", selectedTool))
	if selectedTool != nil {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "sections",
			"function": "data",
		}).Trace(fmt.Sprintf("Selected Tool: %#v", selectedTool))

		t.state.resources.sections = make(map[string]*section, 0)
		t.state.resources.sectionRows = make(map[int]string, 0)
		rowCount := 0

		_, buffer := openFile(selectedTool.Path)
		yamlData := buffer.Bytes()

		var toolNote ToolNote

		if err := yaml.Unmarshal(yamlData, &toolNote); err != nil {
			logrus.Info("DEBUG: failed to reparse our base structure")
		}

		toolNames := make([]string, 0, len(toolNote.Tool.Sections))
		toolList := make(map[string]*section)
		for _, k := range toolNote.Tool.Sections {
			toolNames = append(toolNames, k.SectionName)
			toolList[k.SectionName] = &section{
				Name:     k.SectionName,
				Examples: k.Examples,
			}
		}
		sort.Strings(toolNames)

		// var sectionArray []string
		// for _, v := range toolNote.Tool.Sections {
		for _, v := range toolNames {
			// sectionArray = append(sectionArray, v.SectionName)
			t.state.resources.sections[v] = toolList[v]
			// t.state.resources.sectionRows[rowCount] = fmt.Sprintf("%s/%s", selectedTool.Path, selectedTool.Name)
			t.state.resources.sectionRows[rowCount] = v
			rowCount++
		}
	}
}

func (i *sections) setEntries(t *Tui, operation int) {
	i.buildPanelData(t, operation)
	table := i.Clear()

	headers := []string{
		i.lang.GetText("ui", "NAME"),
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorLightYellow,
			BackgroundColor: tcell.ColorDefault,
		})
	}

	sectionNames := make([]string, 0, len(t.state.resources.sections))
	for k := range t.state.resources.sections {
		sectionNames = append(sectionNames, k)
	}
	sort.Strings(sectionNames)

	c := 0
	for _, sortedName := range sectionNames {
		section := t.state.resources.sections[sortedName]

		rowColor := tcell.ColorLightBlue
		table.SetCell(c+1, 0, tview.NewTableCell(section.Name).
			SetTextColor(rowColor).
			SetMaxWidth(60).
			SetExpansion(0))
		c++
	}

	// lastRow := 0
	// if len(folderPaths) > 0 {
	// 	lastRow = t.state.resources.rowTracker[t.state.resources.sections[[0]].Parent]
	// }
	// if lastRow <= c {
	// 	table.Select(lastRow, 0)
	// } else {
	// 	table.Select(0, 0)
	// }
	// 	table.Select(0, 0)
	i.ScrollToBeginning()
	// TODO: update the example panel on re-load
	// if operation == enterSection {
	// 	t.examplePanel().setEntries(t, enterExample)
	// }
	if t.examplePanel() != nil {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "sections",
			"function": "tuibuild",
		}).Debug("Setting Entries for Examples panel")
		// TODO: This is the call to update EXAMPLES PANEL
		t.examplePanel().setEntries(t, enterSection)
	}

}

func (i *sections) updateEntries(t *Tui) {
	t.app.QueueUpdateDraw(func() {
		i.setEntries(t, enterRepository)
	})
}

func (i *sections) focus(t *Tui) {
	i.SetSelectable(true, false)
	t.app.SetFocus(i)
}

func (i *sections) unfocus() {
	i.SetSelectable(false, false)
}

func (i *sections) setFilterWord(word string) {
	i.filterWord = word
}

func (i *sections) setFilterType(which int) {
	i.showTypes = which
}
