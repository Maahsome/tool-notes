package tui

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	// "github.com/maahsome/vault-view/clipboard"
	// "github.com/atotto/clipboard"
	"maahsome/tool-notes/common"
	"maahsome/tool-notes/resource"

	"github.com/sirupsen/logrus"
)

type repository struct {
	Type      string
	Path      string
	Name      string
	CommitSHA string
}

type repositories struct {
	*tview.Table
	filterWord string
	showTypes  int
	lang       *resource.Lang
}

type ConfigData struct {
	Repositories ConfigDataRepositories `json:"repositories"`
}

type ConfigDataRepositories struct {
	Readonly  []string `json:"readonly"`
	Readwrite []string `json:"readwrite"`
}

const (
	chunksize int = 1024
)

func newRepositories(t *Tui) *repositories {
	repositories := &repositories{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		lang:  t.lang,
	}

	common.Logger.WithFields(logrus.Fields{
		"unit":     "repositories",
		"function": "new",
	}).Debug("New Repository Panel Created")
	repositories.SetTitle(fmt.Sprintf(" [[ %s ]] ", t.lang.GetText("ui", "Repository"))).SetTitleAlign(tview.AlignLeft)
	repositories.SetBorder(true)
	repositories.SetBorderColor(tcell.ColorDeepSkyBlue)
	repositories.setEntries(t, enterRepository)
	repositories.setKeybinding(t)
	return repositories
}

func (i *repositories) name() string {
	return "repositories"
}

func (i *repositories) setTitle() {

	itemsShown := ""
	if i.showTypes != allItems {
		switch i.showTypes {
		case readOnlyRepositories:
			itemsShown = fmt.Sprintf("[green](%s)[white]", i.lang.GetText("ui", "ReadOnly"))
		case readWriteRepositories:
			itemsShown = fmt.Sprintf("[green](%s)[white]", i.lang.GetText("ui", "ReadWrite"))
		}
	}
	if len(i.filterWord) > 0 {
		i.SetTitle(fmt.Sprintf(" [[ %s %s ]] - /%s/ ", i.lang.GetText("ui", "Repository"), itemsShown, i.filterWord)).SetTitleAlign(tview.AlignLeft)
	} else {
		i.SetTitle(fmt.Sprintf(" [[ %s %s]] ", i.lang.GetText("ui", "Respository"), itemsShown)).SetTitleAlign(tview.AlignLeft)
	}
}

func (i *repositories) setKeybinding(t *Tui) {
	i.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		t.setGlobalKeybinding(event)
		selectedRepository := t.selectedRepository()
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
			row, _ := t.repositoryPanel().GetSelection()

			common.Logger.WithFields(logrus.Fields{
				"unit":     "repositories",
				"function": "keystrokes",
			}).Debug(fmt.Sprintf("down repositories/row: %d/%d", len(t.state.resources.repositories), row))

			if row < len(t.state.resources.repositories) {
				tempRow := row + 1
				i.Select(tempRow, 0)
				t.toolPanel().setEntries(t, enterRepository)
				t.toolPanel().Select(0, 0)
				i.Select(row, 0)
			}
		case tcell.KeyUp:
			row, _ := t.repositoryPanel().GetSelection()
			common.Logger.WithFields(logrus.Fields{
				"unit":     "repositories",
				"function": "keystrokes",
			}).Debug(fmt.Sprintf("up repositories/row: %d/%d", len(t.state.resources.repositories), row))
			if row > 0 {
				tempRow := row - 1
				i.Select(tempRow, 0)
				t.toolPanel().setEntries(t, enterRepository)
				t.toolPanel().Select(0, 0)
				i.Select(row, 0)
			}
		case tcell.KeyRight:
			common.Logger.WithFields(logrus.Fields{
				"unit":     "repositories",
				"function": "keystrokes",
			}).Debug("KeyRight")

			if selectedRepository != nil {
				// 	row, _ := t.folderPanel().GetSelection()
				// 	common.Logger.WithFields(logrus.Fields{
				// 		"unit":        "repositories",
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
				"unit":     "repositories",
				"function": "keystrokes",
			}).Info("KeyLeft")
			// row, _ := t.folderPanel().GetSelection()
			// if selectedFolder != nil {
			// 	common.Logger.WithFields(logrus.Fields{
			// 		"unit":     "repositories",
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
				"unit":     "repositories",
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

func (i *repositories) toggleSelected(t *Tui, selectedRepository *repository) {
	common.Logger.WithFields(logrus.Fields{
		"unit":     "repositories",
		"function": "marking",
	}).Debug(fmt.Sprintf("MarkThis: %s", selectedRepository.Name))

	row, _ := t.repositoryPanel().GetSelection()
	rowColor := tcell.ColorLightBlue
	for col := 0; col <= t.repositoryPanel().GetColumnCount(); col++ {
		t.repositoryPanel().GetCell(row, col).SetTextColor(rowColor)
	}
}

// func (i *repositories) buildPanelData(t *Tui, operation int) {

// 	fetchPath := "/"
// 	selectedRepository := t.selectedRepository()
// 	common.Logger.WithFields(logrus.Fields{
// 		"unit":     "",
// 		"function": "data",
// 	}).Trace(fmt.Sprintf("Selected Repository: %#v", selectedFolder))

// 	// Determine the folder we will build data for
// 	switch operation {
// 	case enterFolder:
// 		if selectedFolder != nil {
// 			fetchPath = selectedFolder.FullPath
// 		}
// 	case enterParent:
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":     "repositories",
// 			"function": "data",
// 		}).Debug(fmt.Sprintf("NAV[repositories]: Selected Parent from Panel: %s", i.getParent()))
// 		fetchPath = fmt.Sprintf("%s/", filepath.Dir(strings.TrimSuffix(i.getParent(), "/")))
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":     "repositories",
// 			"function": "data",
// 		}).Debug(fmt.Sprintf("NAV[repositories]: Initial FetchPath: %s", fetchPath))
// 		if fetchPath == "//" || fetchPath == "./" {
// 			fetchPath = "/"
// 		}
// 	case applyFilter:
// 		fetchPath = i.getParent()
// 	}

// 	i.setShownPath(fetchPath)
// 	common.Logger.WithFields(logrus.Fields{
// 		"unit":     "repositories",
// 		"function": "data",
// 	}).Info(fmt.Sprintf("SET ShownPath: %s", fetchPath))

// 	var repositories map[string]vault.Paths
// 	var serr error

// 	if t.vaultCache.CachePathExists(fetchPath) {
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":     "repositories",
// 			"function": "cache",
// 		}).Info("Loading PATHS from Cache... wooo!")
// 		repositories = t.vaultCache.CachePaths[fetchPath].Paths
// 	} else {
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":     "repositories",
// 			"function": "data",
// 		}).Debug(fmt.Sprintf("Fetching Path Datas: %s", fetchPath))
// 		repositories, serr = t.vault.GetPaths(fetchPath)
// 		if serr != nil {
// 			common.Logger.WithFields(logrus.Fields{
// 				"unit":     "repositories",
// 				"function": "data",
// 			}).Error("failed to get paths")
// 		}
// 		if ok, err := t.vaultCache.UpdateCachePath(fetchPath, repositories); !ok {
// 			common.Logger.WithFields(logrus.Fields{
// 				"unit":     "repositories",
// 				"function": "cache",
// 			}).WithError(err).Error("Bad UpdateCachePath")
// 		}
// 	}

// 	t.state.resources.repositories = make(map[string]*folder, 0)
// 	t.state.resources.folderRows = make(map[int]string, 0)

// 	folderPaths := make([]string, 0, len(repositories))
// 	for k := range repositories {
// 		folderPaths = append(folderPaths, k)
// 	}
// 	sort.Strings(folderPaths)

// 	if len(repositories) > 0 {
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":     "repositories",
// 			"function": "data",
// 		}).Debug(fmt.Sprintf("NAV[repositories] SET Parent: %s", repositories[folderPaths[0]].Parent))
// 		i.setParent(repositories[folderPaths[0]].Parent)
// 		if t.state.location != nil {
// 			t.state.location.update(fmt.Sprintf("\n [white]%s", repositories[folderPaths[0]].Parent))
// 		}
// 	}

// 	rowCount := 0
// 	for _, sortedPath := range folderPaths {
// 		folderInfo := repositories[sortedPath]
// 		common.Logger.WithFields(logrus.Fields{
// 			"unit":          "repositories",
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
// 			// Folder, pass this along to pre-load KEYs in the repositories of this path
// 			go t.vaultCache.PreloadFolderPaths(folderInfo.FullPath)
// 		}

// 		t.state.resources.repositories[folderInfo.FullPath] = &folder{
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

func (i *repositories) buildPanelData(t *Tui, operation int) {

	selectedRepository := t.selectedRepository()
	// if selectedRepository != nil {
	common.Logger.WithFields(logrus.Fields{
		"unit":     "",
		"function": "data",
	}).Trace(fmt.Sprintf("Selected Repository: %#v", selectedRepository))

	t.state.resources.repositories = make(map[string]*repository, 0)
	t.state.resources.repositoryRows = make(map[int]string, 0)
	rowCount := 0

	_, buffer := openFile(viper.ConfigFileUsed())

	original := buffer.Bytes()
	repoList := ConfigData{}
	if err := yaml.Unmarshal(original, &repoList); err != nil {
		logrus.Info("DEBUG: failed to reparse our base structure")
	}

	for _, f := range repoList.Repositories.Readwrite {
		repoSplit := strings.Split(f, "/")
		repoName := strings.TrimSuffix(repoSplit[len(repoSplit)-2]+"/"+repoSplit[len(repoSplit)-1], ".git")
		t.state.resources.repositories[f] = &repository{
			Type:      readWrite,
			Path:      f,
			Name:      repoName,
			CommitSHA: "unknown",
		}
		t.state.resources.repositoryRows[rowCount] = f
		rowCount++
	}
	for _, f := range repoList.Repositories.Readonly {
		repoSplit := strings.Split(f, "/")
		repoName := strings.TrimSuffix(repoSplit[len(repoSplit)-2]+"/"+repoSplit[len(repoSplit)-1], ".git")
		t.state.resources.repositories[f] = &repository{
			Type:      readOnly,
			Path:      f,
			Name:      repoName,
			CommitSHA: "unknown",
		}
		t.state.resources.repositoryRows[rowCount] = f
		rowCount++
	}
	// }
}

func (i *repositories) setEntries(t *Tui, operation int) {
	i.buildPanelData(t, operation)
	table := i.Clear()

	headers := []string{
		i.lang.GetText("ui", "TYPE"),
		i.lang.GetText("ui", "NAME"),
		i.lang.GetText("ui", "COMMIT"),
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

	repositoryNames := make([]string, 0, len(t.state.resources.repositories))
	for k := range t.state.resources.repositories {
		repositoryNames = append(repositoryNames, k)
	}
	sort.Strings(repositoryNames)

	c := 0
	for _, sortedName := range repositoryNames {
		repository := t.state.resources.repositories[sortedName]

		rowColor := tcell.ColorLightBlue
		if repository.Type == readOnly {
			rowColor = tcell.ColorMediumSeaGreen
		}
		table.SetCell(c+1, 0, tview.NewTableCell(i.lang.GetText("ui", repository.Type)).
			SetTextColor(rowColor).
			SetMaxWidth(10).
			SetExpansion(0))

		table.SetCell(c+1, 1, tview.NewTableCell(repository.Name).
			SetTextColor(rowColor).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(c+1, 2, tview.NewTableCell(repository.CommitSHA).
			SetTextColor(rowColor).
			SetMaxWidth(1).
			SetExpansion(1))
		c++
	}

	// lastRow := 0
	// if len(folderPaths) > 0 {
	// 	lastRow = t.state.resources.rowTracker[t.state.resources.repositories[[0]].Parent]
	// }
	// if lastRow <= c {
	// 	table.Select(lastRow, 0)
	// } else {
	// 	table.Select(0, 0)
	// }
	// 	table.Select(0, 0)
	i.ScrollToBeginning()

	common.Logger.WithFields(logrus.Fields{
		"unit":     "repositories",
		"function": "tuibuild",
	}).Trace(fmt.Sprintf("Checking toolPanel %#v", t.toolPanel()))
	if t.toolPanel() != nil {
		common.Logger.WithFields(logrus.Fields{
			"unit":     "repositories",
			"function": "tuibuild",
		}).Debug("Setting Entries for Tools panel")
		t.toolPanel().setEntries(t, enterRepository)
	}

}

func (i *repositories) updateEntries(t *Tui) {
	t.app.QueueUpdateDraw(func() {
		i.setEntries(t, enterRepository)
	})
}

func (i *repositories) focus(t *Tui) {
	i.SetSelectable(true, false)
	t.app.SetFocus(i)
}

func (i *repositories) unfocus() {
	i.SetSelectable(false, false)
}

func (i *repositories) setFilterWord(word string) {
	i.filterWord = word
}

func (i *repositories) setFilterType(which int) {
	i.showTypes = which
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
