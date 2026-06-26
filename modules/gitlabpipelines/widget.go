package gitlabpipelines

import (
	"strconv"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

type Widget struct {
	view.MultiSourceWidget
	view.TextWidget

	GitlabProjects []*GitlabProject

	context  *context
	settings *Settings
	Selected int
	maxItems int
	Items    []string // browser URL per rendered row

	configError error
}

// NewWidget creates a new instance of the widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	context, err := newContext(settings)

	widget := Widget{
		MultiSourceWidget: view.NewMultiSourceWidget(settings.Common, "project", "projects"),
		TextWidget:        view.NewTextWidget(tviewApp, redrawChan, pages, settings.Common),

		context:  context,
		settings: settings,

		configError: err,
	}

	widget.GitlabProjects = widget.buildProjectCollection(context, settings.projects)

	widget.initializeKeyboardControls()
	widget.View.SetRegions(true)
	widget.SetDisplayFunction(widget.display)

	widget.Unselect()

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	if widget.context == nil || widget.configError != nil {
		widget.displayError()
		return
	}

	for _, project := range widget.GitlabProjects {
		project.Refresh()
	}

	widget.display()
}

// SetItemCount sets the number of selectable rows rendered so far
func (widget *Widget) SetItemCount(items int) {
	widget.maxItems = items
}

// GetItemCount returns the number of selectable rows rendered so far
func (widget *Widget) GetItemCount() int {
	return widget.maxItems
}

// GetSelected returns the index of the currently highlighted item as an int
func (widget *Widget) GetSelected() int {
	if widget.Selected < 0 {
		return 0
	}
	return widget.Selected
}

// Next cycles the currently highlighted row down
func (widget *Widget) Next() {
	widget.Selected++
	if widget.Selected >= widget.maxItems {
		widget.Selected = 0
	}
	widget.View.Highlight(strconv.Itoa(widget.Selected))
	widget.View.ScrollToHighlight()
}

// Prev cycles the currently highlighted row up
func (widget *Widget) Prev() {
	widget.Selected--
	if widget.Selected < 0 {
		widget.Selected = widget.maxItems - 1
	}
	widget.View.Highlight(strconv.Itoa(widget.Selected))
	widget.View.ScrollToHighlight()
}

// Unselect stops highlighting the text and jumps the scroll position to the top
func (widget *Widget) Unselect() {
	widget.Selected = -1
	widget.View.Highlight()
	widget.View.ScrollToBeginning()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) buildProjectCollection(context *context, projectData []string) []*GitlabProject {
	gitlabProjects := []*GitlabProject{}

	for _, projectPath := range projectData {
		gitlabProjects = append(gitlabProjects, NewGitlabProject(context, projectPath, widget.settings.count))
	}

	return gitlabProjects
}

func (widget *Widget) currentGitlabProject() *GitlabProject {
	if widget.Idx < 0 || widget.Idx >= len(widget.GitlabProjects) {
		return nil
	}

	return widget.GitlabProjects[widget.Idx]
}

func (widget *Widget) openItemInBrowser() {
	currentSelection := widget.View.GetHighlights()
	if widget.Selected >= 0 && len(currentSelection) > 0 && currentSelection[0] != "" {
		if widget.Selected < len(widget.Items) && widget.Items[widget.Selected] != "" {
			utils.OpenFile(widget.Items[widget.Selected])
		}
	}
}
