package tui

import (
	"fmt"
	"time"

	"maahsome/tool-notes/common"
	"maahsome/tool-notes/resource"

	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"
	"github.com/sirupsen/logrus"
)

type panels struct {
	currentPanel int
	panel        []panel
}

// vault resources
type resources struct {
	repositories   map[string]*repository
	repositoryRows map[int]string
	tools          map[string]*tool
	toolRows       map[int]string
	sections       map[string]*section
	sectionRows    map[int]string
	examples       map[string]*example
	exampleRows    map[int]string
	rowTracker     map[string]int
}

type state struct {
	panels     panels
	editPanels panels
	location   *location
	command    *commands
	grid       *tview.Flex
	edit       *tview.Flex
	info       *info
	resources  resources
	stopChans  map[string]chan int
}

const (
	enterRepository = iota
	enterTool
	enterSection
	enterExample
	applyFilter
)

const (
	allItems = iota
	readOnlyRepositories
	readWriteRepositories
)

const readWrite = "ReadWrite"
const readOnly = "ReadOnly"

func newState() *state {
	return &state{
		stopChans: make(map[string]chan int),
	}
}

// Tui - Structure that runs the application
type Tui struct {
	app    *tview.Application
	pages  *tview.Pages
	state  *state
	lang   *resource.Lang
	semVer string
}

// New create new Tui
func New(version string) *Tui {

	return &Tui{
		app:    tview.NewApplication(),
		state:  newState(),
		lang:   resource.NewLanguage(),
		semVer: version,
	}
}

func (t *Tui) repositoryPanel() *repositories {
	for _, panel := range t.state.panels.panel {
		if panel.name() == "repositories" {
			return panel.(*repositories)
		}
	}
	return nil
}

func (t *Tui) toolPanel() *tools {
	for _, panel := range t.state.panels.panel {
		if panel.name() == "tools" {
			return panel.(*tools)
		}
	}
	return nil
}

func (t *Tui) sectionPanel() *sections {
	for _, panel := range t.state.editPanels.panel {
		if panel.name() == "sections" {
			return panel.(*sections)
		}
	}
	return nil
}

func (t *Tui) examplePanel() *examples {
	for _, panel := range t.state.editPanels.panel {
		if panel.name() == "examples" {
			return panel.(*examples)
		}
	}
	return nil
}

func (t *Tui) initPanels() {
	tools := newTools(t)
	repositories := newRepositories(t)
	command := newCommand(t)
	info := newInfo(t)
	location := newLocation()
	address := tview.NewTextView().SetTextColor(tcell.ColorWhite).
		SetText(fmt.Sprintf(" %s (%s)", "github", "maahsome"))
	sections := newSections(t)
	examples := newExamples(t)

	location.update("\n [white]/")
	go t.ClearIndicator(10)

	t.state.panels.panel = append(t.state.panels.panel, repositories)
	t.state.panels.panel = append(t.state.panels.panel, tools)
	t.state.editPanels.panel = append(t.state.editPanels.panel, sections)
	t.state.editPanels.panel = append(t.state.editPanels.panel, examples)
	t.state.command = command
	t.state.location = location
	t.state.info = info
	t.state.resources.rowTracker = make(map[string]int)

	// _, h := consolesize.GetConsoleSize()

	grid := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(info, 6, 1, false).
		AddItem(address, 2, 1, false).
		AddItem(command, 3, 1, false).
		// AddItem(repositories, 0, (h-4)/2, true).
		// AddItem(tools, 0, (h-4)/2, false).
		AddItem(repositories, 0, 1, true).
		AddItem(tools, 0, 2, false).
		AddItem(location, 3, 1, false)

	grid.ResizeItem(t.state.command, 0, 0)

	t.state.grid = grid

	// editGrid := tview.NewFlex().SetDirection(tview.FlexRow).
	// 	AddItem(info, 6, 1, false).
	// 	AddItem(address, 2, 1, false).
	// 	AddItem(command, 3, 1, false).
	// 	AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
	// 		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
	// 	        AddItem(sections, 0, 1, true).
	// 		    AddItem(tools, 0, 1, false), 0, 1, false), 0, 1, false).
	// 		AddItem(tview.NewBox().SetBorder(true).SetTitle("Editor"), 0, 1, false), 0, 10, true).
	// 	AddItem(location, 3, 1, false)

	editGrid := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(info, 6, 1, false).
		AddItem(address, 2, 1, false).
		AddItem(command, 3, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(sections, 0, 1, true).
				AddItem(examples, 0, 1, false), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Editor"), 0, 1, false), 0, 1, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("location"), 3, 1, false)

		// flex := tview.NewFlex().SetDirection(tview.FlexRow).
		// AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		// 	AddItem(tview.NewBox().SetBorder(true).SetTitle("Logo"), 20, 1, false).
		// 	AddItem(tview.NewBox().SetBorder(true).SetTitle("menu"), 0, 3, false).
		// 	AddItem(tview.NewBox().SetBorder(true).SetTitle("indicator"), 0, 2, false), 6, 1, false).
		// AddItem(address, 2, 1, false).
		// AddItem(tview.NewBox().SetBorder(true).SetTitle("command"), 3, 1, false).
		// AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
		// 	AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		// 		AddItem(tview.NewBox().SetBorder(true).SetTitle("Sections"), 0, 1, true).
		// 		AddItem(tview.NewBox().SetBorder(true).SetTitle("Examples"), 0, 1, false), 0, 1, false).
		// 	AddItem(tview.NewBox().SetBorder(true).SetTitle("Editor"), 0, 1, false), 0, 10, true).
		// AddItem(tview.NewBox().SetBorder(true).SetTitle("location"), 3, 1, false)

	editGrid.ResizeItem(t.state.command, 0, 0)

	t.state.edit = editGrid

	t.pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true).
		AddPage("edit", editGrid, true, false)

	t.app.SetRoot(t.pages, true)
	t.switchPanel("repositories")
	// TODO: test this line with NO repositories defined
	// We may need to seed at startup with the maahsome/toolnotes-base before TUI starts
	repositories.setEntries(t, enterRepository)
}

// Start start application
func (t *Tui) Start() error {
	t.initPanels()
	if err := t.app.Run(); err != nil {
		t.app.Stop()
		return err
	}

	return nil
}

// ClearIndicator - Clear the status message after s seconds
func (t *Tui) ClearIndicator(s time.Duration) {
	time.Sleep(s * time.Second)
	t.app.QueueUpdateDraw(func() {
		t.state.info.Status.display("")
	})
}

// Stop stop application
func (t *Tui) Stop() error {
	t.app.Stop()
	return nil
}

func (t *Tui) selectedRepository() *repository {
	if len(t.state.resources.repositories) == 0 {
		return nil
	}

	if t.repositoryPanel() != nil {
		row, _ := t.repositoryPanel().GetSelection()
		common.Logger.WithFields(logrus.Fields{
			"unit":     "tui",
			"function": "selection",
			"row":      row,
			"datamap":  t.state.resources.repositoryRows[row-1],
		}).Debug("Row for Selection")
		if len(t.state.resources.repositories) == 0 {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "tui",
				"function": "selection",
			}).Debug("Our repository resources is currenty EMPTY")
			return nil
		}
		if row-1 < 0 {
			// if we are at row 0, then return the folder for item 0
			return t.state.resources.repositories[t.state.resources.repositoryRows[0]]
		}
		return t.state.resources.repositories[t.state.resources.repositoryRows[row-1]]
	}
	return nil
}

func (t *Tui) selectedTool() *tool {
	if len(t.state.resources.tools) == 0 {
		return nil
	}

	if t.toolPanel() != nil {
		row, _ := t.toolPanel().GetSelection()
		common.Logger.WithFields(logrus.Fields{
			"unit":     "tui",
			"function": "selection",
			"row":      row,
			"datamap":  t.state.resources.toolRows[row-1],
		}).Debug("Row for Selection")
		if len(t.state.resources.tools) == 0 {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "tui",
				"function": "selection",
			}).Debug("Our tool resources is currenty EMPTY")
			return nil
		}
		if row-1 < 0 {
			// if we are at row 0, then return the folder for item 0
			common.Logger.WithFields(logrus.Fields{
				"unit":     "tui",
				"function": "selectionReturn",
				"row":      row,
				"datamap":  t.state.resources.toolRows[0],
			}).Debug("Return Selection 0")
			return t.state.resources.tools[t.state.resources.toolRows[0]]
		}
		common.Logger.WithFields(logrus.Fields{
			"unit":     "tui",
			"function": "selectionReturn",
			"row":      row,
			"datamap":  t.state.resources.toolRows[row-1],
		}).Debug("Return Selection row-1")
		return t.state.resources.tools[t.state.resources.toolRows[row-1]]
	}
	return nil
}

func (t *Tui) selectedSection() *section {
	if len(t.state.resources.sections) == 0 {
		return nil
	}

	if t.sectionPanel() != nil {
		row, _ := t.sectionPanel().GetSelection()
		common.Logger.WithFields(logrus.Fields{
			"unit":     "tui",
			"function": "selection",
			"row":      row,
			"datamap":  t.state.resources.sectionRows[row-1],
		}).Debug("Row for Section Selection")
		if len(t.state.resources.sections) == 0 {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "tui",
				"function": "selection",
			}).Debug("Our sections resources is currenty EMPTY")
			return nil
		}
		if row-1 < 0 {
			// if we are at row 0, then return the folder for item 0
			return t.state.resources.sections[t.state.resources.sectionRows[0]]
		}
		return t.state.resources.sections[t.state.resources.sectionRows[row-1]]
	}
	return nil
}

func (t *Tui) selectedExample() *example {
	if len(t.state.resources.examples) == 0 {
		return nil
	}

	if t.examplePanel() != nil {
		row, _ := t.examplePanel().GetSelection()
		common.Logger.WithFields(logrus.Fields{
			"unit":     "tui",
			"function": "selection",
			"row":      row,
			"datamap":  t.state.resources.exampleRows[row-1],
		}).Debug("Row for Example Selection")
		if len(t.state.resources.examples) == 0 {
			common.Logger.WithFields(logrus.Fields{
				"unit":     "tui",
				"function": "selection",
			}).Debug("Our examples resources is currenty EMPTY")
			return nil
		}
		if row-1 < 0 {
			// if we are at row 0, then return the folder for item 0
			return t.state.resources.examples[t.state.resources.exampleRows[0]]
		}
		return t.state.resources.examples[t.state.resources.exampleRows[row-1]]
	}
	return nil
}

func (t *Tui) switchPanel(panelName string) {
	for i, panel := range t.state.panels.panel {
		if panel.name() == panelName {
			t.state.info.Menu.display(panelName)
			panel.focus(t)
			t.state.panels.currentPanel = i
		} else {
			panel.unfocus()
		}
	}
	for i, panel := range t.state.editPanels.panel {
		if panel.name() == panelName {
			t.state.info.Menu.display(panelName)
			panel.focus(t)
			t.state.editPanels.currentPanel = i
		} else {
			panel.unfocus()
		}
	}
}

func (t *Tui) switchPage(pageName string) {
	currentPage, _ := t.pages.GetFrontPage()
	t.pages.HidePage(currentPage)
	t.pages.SwitchToPage(pageName)
	switch pageName {
	case "main":
		t.switchPanel("repositories")
	case "edit":
		t.switchPanel("sections")
		t.state.editPanels.panel[t.state.editPanels.currentPanel].setEntries(t, enterSection)
	}
}

func (t *Tui) closeAndSwitchPanel(removePanel, switchPanel string) {
	t.pages.RemovePage(removePanel).ShowPage("main")
	t.switchPanel(switchPanel)
}

func (t *Tui) modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func (t *Tui) currentPanel() panel {
	return t.state.panels.panel[t.state.panels.currentPanel]
}

func (t *Tui) hideCommand() {
	currentPage, _ := t.pages.GetFrontPage()
	switch currentPage {
	case "main":
		t.state.grid.ResizeItem(t.state.command, 0, 0)
	case "edit":
		t.state.edit.ResizeItem(t.state.command, 0, 0)
	}
}
