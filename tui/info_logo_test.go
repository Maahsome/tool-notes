package tui

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

// TestNewInfoLogo - testing newInfoLogo
func TestNewInfoLogo(t *testing.T) {

	tv := newInfoLogo("v0.0.1")

	logo := `_____
__  /_ %s
_  __/_______
/ /_  __  __ \
\__/  _  / / /
      /_/ /_/
`

	assert.Equal(t, tv.GetText(true), fmt.Sprintf(logo, "v0.0.1"))

}
