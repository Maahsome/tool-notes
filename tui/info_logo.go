package tui

import (
	"fmt"
	"runtime"

	"github.com/gdamore/tcell/v2"
	"github.com/maahsome/tview"
)

type infoLogo struct {
	*tview.TextView
	Host   *hostInfo
	AppVer *appVer
}

type appVer struct {
	Name    string
	Version string
}

type hostInfo struct {
	OSType       string
	Architecture string
}

func newInfoLogo(semVer string) *infoLogo {

	i := &infoLogo{
		TextView: tview.NewTextView(),
		Host:     newHostInfo(),
		AppVer:   newAppVerInfo(semVer),
	}

	i.display()

	return i
}

func newAppVerInfo(semVer string) *appVer {
	return &appVer{
		Name:    "tool-notes",
		Version: semVer,
	}
}

func newHostInfo() *hostInfo {
	return &hostInfo{
		OSType:       runtime.GOOS,
		Architecture: runtime.GOARCH,
	}
}

func (i *infoLogo) display() {

	appVersion := fmt.Sprintf("%s", i.AppVer.Version)

	i.SetDynamicColors(true).SetTextColor(tcell.ColorDarkOrange)

	logo := `_____
__  /_ %s
_  __/_______
/ /_  __  __ \
\__/  _  / / /
      /_/ /_/
`

	i.SetText(fmt.Sprintf(logo, appVersion))
}
