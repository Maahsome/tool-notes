package tui

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"

	// "github.com/maahsome/vault-view/clipboard"
	// "github.com/atotto/clipboard"

	"maahsome/tool-notes/common"
	"maahsome/tool-notes/resource"

	"github.com/sirupsen/logrus"
)

type example struct {
	Name string
}

type examples struct {
	*tview.Table
	filterWord string
	showTypes  int
	lang       *resource.Lang
}

func newExamples(t *Tui) *examples {
	examples := &examples{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		lang:  t.lang,
	}

	common.Logger.WithFields(logrus.Fields{
		"unit":     "examples",
		"function": "new",
	}).Debug("New Examples Panel Created")
	examples.SetTitle(fmt.Sprintf(" [[ %s ]] ", t.lang.GetText("ui", "Examples"))).SetTitleAlign(tview.AlignLeft)
	examples.SetBorder(true)
	examples.SetBorderColor(tcell.ColorDeepSkyBlue)
	examples.setEntries(t, enterExample)
	examples.setKeybinding(t)
	return examples
}

func (i *examples) name() string {
	return "examples"
}

func (i *examples) setTitle() {

	if len(i.filterWord) > 0 {
		i.SetTitle(fmt.Sprintf(" [[ %s ]] - /%s/ ", i.lang.GetText("ui", "Examples"), i.filterWord)).SetTitleAlign(tview.AlignLeft)
	} else {
		i.SetTitle(fmt.Sprintf(" [[ %s ]] ", i.lang.GetText("ui", "Examples"))).SetTitleAlign(tview.AlignLeft)
	}

}

func (i *examples) setKeybinding(t *Tui) {
	i.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		selectedExample := t.selectedExample()
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
			row, _ := t.examplePanel().GetSelection()

			common.Logger.WithFields(logrus.Fields{
				"unit":     "examples",
				"function": "keystrokes",
			}).Debug(fmt.Sprintf("down examples/row: %d/%d", len(t.state.resources.examples), row))

			// if row < len(t.state.resources.tools) {
			// 	tempRow := row + 1
			// 	i.Select(tempRow, 0)
			// 	t.dataPanel().setEntries(t, enterFolder)
			// 	i.Select(row, 0)
			// }
		case tcell.KeyUp:
			row, _ := t.examplePanel().GetSelection()
			common.Logger.WithFields(logrus.Fields{
				"unit":     "examples",
				"function": "keystrokes",
			}).Debug(fmt.Sprintf("up examples/row: %d/%d", len(t.state.resources.examples), row))
			// if row > 0 {
			// 	tempRow := row - 1
			// 	i.Select(tempRow, 0)
			// 	t.dataPanel().setEntries(t, enterFolder)
			// 	i.Select(row, 0)
			// }
		case tcell.KeyRight:
			common.Logger.WithFields(logrus.Fields{
				"unit":     "examples",
				"function": "keystrokes",
			}).Debug("KeyRight")

			if selectedExample != nil {
				// 	row, _ := t.folderPanel().GetSelection()
				// 	common.Logger.WithFields(logrus.Fields{
				// 		"unit":        "examples",
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
				"unit":     "examples",
				"function": "keystrokes",
			}).Info("KeyLeft")
			// row, _ := t.folderPanel().GetSelection()
			// if selectedFolder != nil {
			// 	common.Logger.WithFields(logrus.Fields{
			// 		"unit":     "examples",
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
				"unit":     "examples",
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

func (i *examples) toggleSelected(t *Tui, selectedExample *example) {
	common.Logger.WithFields(logrus.Fields{
		"unit":     "examples",
		"function": "marking",
	}).Debug(fmt.Sprintf("MarkThis: %s", selectedExample.Name))

	row, _ := t.examplePanel().GetSelection()
	rowColor := tcell.ColorLightBlue
	for col := 0; col <= t.examplePanel().GetColumnCount(); col++ {
		t.examplePanel().GetCell(row, col).SetTextColor(rowColor)
	}
}

func (i *examples) buildPanelData(t *Tui, operation int) {

	common.Logger.WithFields(logrus.Fields{
		"unit":     "examples",
		"function": "data",
	}).Trace("Running buildPanelData for examples")

	selectedSection := t.selectedSection()
	common.Logger.WithFields(logrus.Fields{
		"unit":     "examples",
		"function": "data",
	}).Trace(fmt.Sprintf("Selected Section: %#v", selectedSection))
	if selectedSection != nil {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "examples",
			"function": "data",
		}).Trace(fmt.Sprintf("Selected Section: %#v", selectedSection))

		t.state.resources.examples = make(map[string]*example, 0)
		t.state.resources.exampleRows = make(map[int]string, 0)
		rowCount := 0

		for _, v := range selectedSection.Examples {
			t.state.resources.examples[v.Description] = &example{
				Name: v.Description,
			}
			// t.state.resources.sectionRows[rowCount] = fmt.Sprintf("%s/%s", selectedTool.Path, selectedTool.Name)
			t.state.resources.exampleRows[rowCount] = v.Description
			rowCount++
		}
	}
}

func (i *examples) setEntries(t *Tui, operation int) {
	i.buildPanelData(t, operation)
	table := i.Clear()

	headers := []string{
		i.lang.GetText("ui", "DESCRIPTION"),
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

	exampleNames := make([]string, 0, len(t.state.resources.examples))
	for k := range t.state.resources.examples {
		exampleNames = append(exampleNames, k)
	}
	sort.Strings(exampleNames)

	c := 0
	for _, sortedName := range exampleNames {
		example := t.state.resources.examples[sortedName]

		rowColor := tcell.ColorLightBlue
		table.SetCell(c+1, 0, tview.NewTableCell(example.Name).
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
	// if t.dataPanel() != nil {
	// 	common.Logger.WithFields(logrus.Fields{
	// 		"unit":     "sections",
	// 		"function": "tuibuild",
	// 	}).Debug("Setting Entries for Datas panel")
	// 	// TODO: This is the call to update TOOLS PANEL
	// 	// t.dataPanel().setEntries(t, enterRepository)
	// }

}

func (i *examples) updateEntries(t *Tui) {
	t.app.QueueUpdateDraw(func() {
		i.setEntries(t, enterRepository)
	})
}

func (i *examples) focus(t *Tui) {
	i.SetSelectable(true, false)
	t.app.SetFocus(i)
}

func (i *examples) unfocus() {
	i.SetSelectable(false, false)
}

func (i *examples) setFilterWord(word string) {
	i.filterWord = word
}

func (i *examples) setFilterType(which int) {
	i.showTypes = which
}
