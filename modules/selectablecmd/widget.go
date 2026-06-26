package selectablecmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

// Widget runs a command and renders its tab-separated output as selectable
// rows. Pressing Enter/o on a row opens that row's URL column in the browser.
type Widget struct {
	view.ScrollableWidget

	settings *Settings
	rows     []row
	err      error
}

// NewWidget creates a new instance of the widget.
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.Common),

		settings: settings,
	}

	widget.SetRenderFunction(widget.Render)
	widget.initializeKeyboardControls()

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh runs the command and reparses its output.
func (widget *Widget) Refresh() {
	out, err := widget.runCommand()
	if err != nil {
		widget.err = err
		widget.rows = nil
		widget.SetItemCount(0)
	} else {
		widget.err = nil
		widget.rows = parseOutput(out, widget.settings)
		widget.SetItemCount(len(widget.rows))
	}
	widget.Render()
}

func (widget *Widget) Render() {
	widget.Redraw(widget.content)
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) runCommand() (string, error) {
	cmd := exec.Command(widget.settings.cmd, widget.settings.args...)
	cmd.Dir = widget.settings.workingDir

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	return stdout.String(), err
}

func (widget *Widget) openItem() {
	sel := widget.GetSelected()
	if sel >= 0 && sel < len(widget.rows) {
		url := widget.rows[sel].url
		if url != "" {
			utils.OpenFile(url)
		}
	}
}

func (widget *Widget) content() (string, string, bool) {
	title := widget.CommonSettings().Title

	if widget.err != nil {
		return title, widget.err.Error(), true
	}
	if len(widget.rows) == 0 {
		return title, "No results to display", false
	}

	widths := widget.columnWidths()

	var str strings.Builder
	for idx, r := range widget.rows {
		cells := r.display(widget.settings)

		var line strings.Builder
		fmt.Fprintf(&line, "[%s]", widget.RowColor(idx))
		for i, cell := range cells {
			pad := 0
			if i < len(widths) {
				pad = widths[i] + 1
			}
			fmt.Fprintf(&line, "%-*s", pad, tview.Escape(cell))
		}

		rendered := line.String()
		str.WriteString(utils.HighlightableHelper(widget.View, rendered, idx, len(rendered)))
	}

	return title, str.String(), false
}

// columnWidths computes the max width of each displayed column for alignment.
func (widget *Widget) columnWidths() []int {
	widths := []int{}
	for _, r := range widget.rows {
		for i, cell := range r.display(widget.settings) {
			if i >= len(widths) {
				widths = append(widths, 0)
			}
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	return widths
}
