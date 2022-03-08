package tui

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"

	// "github.com/maahsome/vault-view/clipboard"
	// "github.com/atotto/clipboard"
	"maahsome/tool-notes/common"
	"maahsome/tool-notes/resource"

	"github.com/sirupsen/logrus"
)

type tool struct {
	Path string
	Name string
}

type tools struct {
	*tview.Table
	filterWord string
	showTypes  int
	lang       *resource.Lang
}

func newTools(t *Tui) *tools {
	tools := &tools{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		lang:  t.lang,
	}

	common.Logger.WithFields(logrus.Fields{
		"unit":     "tools",
		"function": "new",
	}).Debug("New Tools Panel Created")
	tools.SetTitle(fmt.Sprintf(" [[ %s ]] ", t.lang.GetText("ui", "Tool"))).SetTitleAlign(tview.AlignLeft)
	tools.SetBorder(true)
	tools.SetBorderColor(tcell.ColorDeepSkyBlue)
	tools.setEntries(t, enterTool)
	tools.setKeybinding(t)
	return tools
}

func (i *tools) name() string {
	return "tools"
}

func (i *tools) setTitle() {

	if len(i.filterWord) > 0 {
		i.SetTitle(fmt.Sprintf(" [[ %s ]] - /%s/ ", i.lang.GetText("ui", "Tool"), i.filterWord)).SetTitleAlign(tview.AlignLeft)
	} else {
		i.SetTitle(fmt.Sprintf(" [[ %s ]] ", i.lang.GetText("ui", "Tool"))).SetTitleAlign(tview.AlignLeft)
	}

}

func (i *tools) setKeybinding(t *Tui) {
	i.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		selectedTool := t.selectedTool()
		switch event.Key() {
		case tcell.KeyEnter:
			selectedRepository := t.selectedRepository()
			if selectedRepository != nil {
				if t.state.location != nil {
					t.state.location.update(fmt.Sprintf("\n [white]%s / %s", selectedRepository.Name, selectedTool.Name))
				}
			}
			t.switchPage("edit")
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
			row, _ := t.toolPanel().GetSelection()

			common.Logger.WithFields(logrus.Fields{
				"unit":     "tools",
				"function": "keystrokes",
			}).Debug(fmt.Sprintf("down tools/row: %d/%d", len(t.state.resources.tools), row))
			// if row < len(t.state.resources.tools) {
			// 	tempRow := row + 1
			// 	i.Select(tempRow, 0)
			// 	t.dataPanel().setEntries(t, enterFolder)
			// 	i.Select(row, 0)
			// }
		case tcell.KeyUp:
			row, _ := t.toolPanel().GetSelection()
			common.Logger.WithFields(logrus.Fields{
				"unit":     "tools",
				"function": "keystrokes",
			}).Debug(fmt.Sprintf("up tools/row: %d/%d", len(t.state.resources.tools), row))
			// if row > 0 {
			// 	tempRow := row - 1
			// 	i.Select(tempRow, 0)
			// 	t.dataPanel().setEntries(t, enterFolder)
			// 	i.Select(row, 0)
			// }
		case tcell.KeyRight:
			common.Logger.WithFields(logrus.Fields{
				"unit":     "tools",
				"function": "keystrokes",
			}).Debug("KeyRight")

			if selectedTool != nil {
				// 	row, _ := t.folderPanel().GetSelection()
				// 	common.Logger.WithFields(logrus.Fields{
				// 		"unit":        "tools",
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
				"unit":     "tools",
				"function": "keystrokes",
			}).Info("KeyLeft")
			// row, _ := t.folderPanel().GetSelection()
			// if selectedFolder != nil {
			// 	common.Logger.WithFields(logrus.Fields{
			// 		"unit":     "tools",
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
				"unit":     "tools",
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

func (i *tools) toggleSelected(t *Tui, selectedTool *tool) {
	common.Logger.WithFields(logrus.Fields{
		"unit":     "tools",
		"function": "marking",
	}).Debug(fmt.Sprintf("MarkThis: %s", selectedTool.Name))

	row, _ := t.toolPanel().GetSelection()
	rowColor := tcell.ColorLightBlue
	for col := 0; col <= t.toolPanel().GetColumnCount(); col++ {
		t.toolPanel().GetCell(row, col).SetTextColor(rowColor)
	}
}

// func (i *tools) buildPanelData(t *Tui, operation int) {

// 	fetchPath := "/"
// 	selectedTool := t.selectedTool()
// 	common.Logger.WithFields(logrus.Fields{
// 		"unit":     "",
// 		"function": "data",
// 	}).Trace(fmt.Sprintf("Selected tool: %#v", selectedFolder))

// 	// Determine the folder we will build data for
// 	switch operation {
// 	case enterFolder:
// 		if selectedFolder != nil {
// 			fetchPath = selectedFolder.FullPath
// 		}
// 	case enterParent:
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":     "tools",
// 			"function": "data",
// 		}).Debug(fmt.Sprintf("NAV[tools]: Selected Parent from Panel: %s", i.getParent()))
// 		fetchPath = fmt.Sprintf("%s/", filepath.Dir(strings.TrimSuffix(i.getParent(), "/")))
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":     "tools",
// 			"function": "data",
// 		}).Debug(fmt.Sprintf("NAV[tools]: Initial FetchPath: %s", fetchPath))
// 		if fetchPath == "//" || fetchPath == "./" {
// 			fetchPath = "/"
// 		}
// 	case applyFilter:
// 		fetchPath = i.getParent()
// 	}

// 	i.setShownPath(fetchPath)
// 	common.Logger.WithFields(logrus.Fields{
// 		"unit":     "tools",
// 		"function": "data",
// 	}).Info(fmt.Sprintf("SET ShownPath: %s", fetchPath))

// 	var tools map[string]vault.Paths
// 	var serr error

// 	if t.vaultCache.CachePathExists(fetchPath) {
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":     "tools",
// 			"function": "cache",
// 		}).Info("Loading PATHS from Cache... wooo!")
// 		tools = t.vaultCache.CachePaths[fetchPath].Paths
// 	} else {
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":     "tools",
// 			"function": "data",
// 		}).Debug(fmt.Sprintf("Fetching Path Datas: %s", fetchPath))
// 		tools, serr = t.vault.GetPaths(fetchPath)
// 		if serr != nil {
// 			common.Logger.WithFields(logrus.Fields{
// 				"unit":     "tools",
// 				"function": "data",
// 			}).Error("failed to get paths")
// 		}
// 		if ok, err := t.vaultCache.UpdateCachePath(fetchPath, tools); !ok {
// 			common.Logger.WithFields(logrus.Fields{
// 				"unit":     "tools",
// 				"function": "cache",
// 			}).WithError(err).Error("Bad UpdateCachePath")
// 		}
// 	}

// 	t.state.resources.tools = make(map[string]*folder, 0)
// 	t.state.resources.folderRows = make(map[int]string, 0)

// 	folderPaths := make([]string, 0, len(tools))
// 	for k := range tools {
// 		folderPaths = append(folderPaths, k)
// 	}
// 	sort.Strings(folderPaths)

// 	if len(tools) > 0 {
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":     "tools",
// 			"function": "data",
// 		}).Debug(fmt.Sprintf("NAV[tools] SET Parent: %s", tools[folderPaths[0]].Parent))
// 		i.setParent(tools[folderPaths[0]].Parent)
// 		if t.state.location != nil {
// 			t.state.location.update(fmt.Sprintf("\n [white]%s", tools[folderPaths[0]].Parent))
// 		}
// 	}

// 	rowCount := 0
// 	for _, sortedPath := range folderPaths {
// 		folderInfo := tools[sortedPath]
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":          "tools",
// 			"function":      "data",
// 			"folder_type":   folderInfo.Type,
// 			"folder_path":   folderInfo.Path,
// 			"folder_parent": folderInfo.Parent,
// 			"filter_word":   i.filterWord,
// 			"show_types":    i.showTypes,
// 		}).Debug(fmt.Sprintf("SortedKEY: [%s]", sortedPath))
// 		if strings.Index(folderInfo.Path, i.filterWord) == -1 {
// 			continue
// 		}
// 		if i.showTypes == dataItems && folderInfo.Type == vaultFolder {
// 			continue
// 		}
// 		if i.showTypes == folderItems && folderInfo.Type == vaultData {
// 			continue
// 		}

// 		var folderData vault.DataRecord
// 		// CACHE: Load the datas as we are building the KEY list.
// 		if folderInfo.Type == vaultData {
// 			// Check Cache
// 			if !t.vaultCache.CachePathExists(folderInfo.FullPath) {
// 				vaultPaths := make(map[string]vault.Paths)
// 				vaultPaths[folderInfo.FullPath] = vault.Paths{
// 					Type:    folderInfo.Type,
// 					Path:    folderInfo.Path,
// 					Parent:  folderInfo.Parent,
// 					Version: folderData.Data.Metadata.Version,
// 				}
// 				if len(vaultPaths) > 0 {
// 					// TODO: PERFORMANCE, this used to be called as a go routine
// 					// It makes the interface load MUCH faster, though the VERSION
// 					// for the "Data" types lags in the interface, need some way
// 					// to go back and populate them once the cache is complete.
// 					// There is a framework here for time based reloading that
// 					// is probably the deal...
// 					// go t.vaultCache.PreloadPaths(vaultPaths)
// 					t.vaultCache.PreloadPaths(vaultPaths)
// 				}
// 			}
// 			// Attempt to find VERSION in the Data cache, this will lag in the display
// 			if t.vaultCache.CacheDataExist(folderInfo.FullPath, expireMinutes) {
// 				folderData = t.vaultCache.GetCacheData(folderInfo.FullPath)
// 			}
// 		} else {
// 			// Folder, pass this along to pre-load KEYs in the tools of this path
// 			go t.vaultCache.PreloadFolderPaths(folderInfo.FullPath)
// 		}

// 		t.state.resources.tools[folderInfo.FullPath] = &folder{
// 			Type:     folderInfo.Type,
// 			Path:     folderInfo.Path,
// 			Parent:   folderInfo.Parent,
// 			FullPath: folderInfo.FullPath,
// 			Version:  folderData.Data.Metadata.Version,
// 		}
// 		t.state.resources.folderRows[rowCount] = folderInfo.FullPath

// 		// TODO: I was originally just changing them all to FALSE upon load, this may still be the desired state
// 		// I supposed I could loop through the array as my "work" book and generate scripts from ALL marked items
// 		// rather than just marked items on the page.  Could be interesting.
// 		// if t.state.resources.markedFolders != nil {
// 		// 	t.state.resources.markedFolders[fmt.Sprintf("%s%s", folderInfo.Parent, folderInfo.Path)] = false
// 		// }
// 		rowCount++
// 	}
// }

func (i *tools) buildPanelData(t *Tui, operation int) {

	selectedRepository := t.selectedRepository()
	if selectedRepository != nil {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "tools",
			"function": "data",
		}).Trace(fmt.Sprintf("Selected Repository: %#v", selectedRepository))

		t.state.resources.tools = make(map[string]*tool, 0)
		t.state.resources.toolRows = make(map[int]string, 0)
		rowCount := 0

		files, err := ioutil.ReadDir(fmt.Sprintf("/Users/christopher.maahs/.config/tool-notes/%s", selectedRepository.Name))
		if err != nil {
			common.Logger.Fatal(err)
		}

		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".yaml") {
				t.state.resources.tools[fmt.Sprintf("%s/%s", fmt.Sprintf("/Users/christopher.maahs/.config/tool-notes/%s", selectedRepository.Name), f.Name())] = &tool{
					Name: f.Name(),
					Path: fmt.Sprintf("%s/%s", fmt.Sprintf("/Users/christopher.maahs/.config/tool-notes/%s", selectedRepository.Name), f.Name()),
				}
				t.state.resources.toolRows[rowCount] = fmt.Sprintf("%s/%s", fmt.Sprintf("/Users/christopher.maahs/.config/tool-notes/%s", selectedRepository.Name), f.Name())
				rowCount++
			}
		}
		if t.state.location != nil {
			t.state.location.update(fmt.Sprintf("\n [white]%s", selectedRepository.Name))
		}
	}
}

func (i *tools) setEntries(t *Tui, operation int) {
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

	toolNames := make([]string, 0, len(t.state.resources.tools))
	for k := range t.state.resources.tools {
		toolNames = append(toolNames, k)
	}
	sort.Strings(toolNames)

	c := 0
	for _, sortedName := range toolNames {
		tool := t.state.resources.tools[sortedName]

		rowColor := tcell.ColorLightBlue
		table.SetCell(c+1, 0, tview.NewTableCell(tool.Name).
			SetTextColor(rowColor).
			SetMaxWidth(30).
			SetExpansion(0))

		c++
	}

	// lastRow := 0
	// if len(folderPaths) > 0 {
	// 	lastRow = t.state.resources.rowTracker[t.state.resources.tools[[0]].Parent]
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
	// 		"unit":     "tools",
	// 		"function": "tuibuild",
	// 	}).Debug("Setting Entries for Datas panel")
	// 	// TODO: This is the call to update TOOLS PANEL
	// 	// t.dataPanel().setEntries(t, enterRepository)
	// }

}

func (i *tools) updateEntries(t *Tui) {
	t.app.QueueUpdateDraw(func() {
		i.setEntries(t, enterRepository)
	})
}

func (i *tools) focus(t *Tui) {
	i.SetSelectable(true, false)
	t.app.SetFocus(i)
}

func (i *tools) unfocus() {
	i.SetSelectable(false, false)
}

func (i *tools) setFilterWord(word string) {
	i.filterWord = word
}

func (i *tools) setFilterType(which int) {
	i.showTypes = which
}
