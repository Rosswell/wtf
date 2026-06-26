package gitlabpipelines

import (
	"github.com/gdamore/tcell/v2"
)

func (widget *Widget) initializeKeyboardControls() {
	widget.InitializeHelpTextKeyboardControl(widget.ShowHelp)
	widget.InitializeRefreshKeyboardControl(widget.Refresh)

	widget.SetKeyboardChar("j", widget.Next, "Select next pipeline")
	widget.SetKeyboardChar("k", widget.Prev, "Select previous pipeline")
	widget.SetKeyboardChar("l", widget.NextSource, "Select next project")
	widget.SetKeyboardChar("h", widget.PrevSource, "Select previous project")
	widget.SetKeyboardChar("o", widget.openItemInBrowser, "Open pipeline in browser")

	widget.SetKeyboardKey(tcell.KeyDown, widget.Next, "Select next pipeline")
	widget.SetKeyboardKey(tcell.KeyUp, widget.Prev, "Select previous pipeline")
	widget.SetKeyboardKey(tcell.KeyRight, widget.NextSource, "Select next project")
	widget.SetKeyboardKey(tcell.KeyLeft, widget.PrevSource, "Select previous project")
	widget.SetKeyboardKey(tcell.KeyEnter, widget.openItemInBrowser, "Open pipeline in browser")
	widget.SetKeyboardKey(tcell.KeyEsc, widget.Unselect, "Clear selection")
}
