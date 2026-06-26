package gitlab

import (
	"fmt"
)

func (widget *Widget) display() {
	widget.Redraw(widget.content)
}

func (widget *Widget) displayError() {
	widget.Redraw(widget.contentError)
}

func (widget *Widget) contentError() (string, string, bool) {

	title := fmt.Sprintf("%s - Error", widget.CommonSettings().Title)

	if widget.configError != nil {
		return title, fmt.Sprintf("Error: \n [red]%v[white]", widget.configError), false

	}
	return title, "Error", false
}

func (widget *Widget) content() (string, string, bool) {
	if len(widget.GitlabProjects) == 0 {
		return widget.CommonSettings().Title, " Gitlab project data is unavailable ", true
	}

	// reset the selectable item collection on every render
	widget.Items = make([]ContentItem, 0)
	widget.SetItemCount(0)

	title := widget.CommonSettings().Title

	str := ""
	for idx, project := range widget.GitlabProjects {
		if idx > 0 {
			str += "\n"
		}

		str += fmt.Sprintf(" [%s]%s[white]\n", widget.settings.Colors.Subheading, project.path)

		str += fmt.Sprintf("  [%s]To Review[white]\n", widget.settings.Colors.Subheading)
		str += widget.renderMergeRequests(project.myReviewMergeRequests())

		str += fmt.Sprintf("  [%s]Mine[white]\n", widget.settings.Colors.Subheading)
		str += widget.renderMergeRequests(project.myMergeRequests())
	}

	return title, str, false
}

// statusTag returns a colour and a fixed-width label describing where the MR
// stands, so the titles that follow stay aligned.
func statusTag(mr *MergeRequest) (color, label string) {
	switch {
	case mr.Draft:
		return "grey", "draft"
	case mr.Approved:
		return "green", "approved"
	default:
		return "yellow", "open"
	}
}

func (widget *Widget) renderMergeRequests(mrs []*MergeRequest) string {
	if len(mrs) == 0 {
		return "   [grey]none[white]\n"
	}

	maxItems := widget.GetItemCount()

	str := ""
	for idx, mr := range mrs {
		color, label := statusTag(mr)
		str += fmt.Sprintf(
			`   ["%d"][%s]%-8s[""][white] %s`,
			maxItems+idx, color, label, mr.Title,
		)
		str += "\n"
		widget.Items = append(widget.Items, ContentItem{URL: mr.WebURL})
	}
	widget.SetItemCount(maxItems + len(mrs))

	return str
}
