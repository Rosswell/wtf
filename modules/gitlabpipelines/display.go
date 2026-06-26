package gitlabpipelines

import (
	"fmt"

	glab "gitlab.com/gitlab-org/api/client-go"
)

func (widget *Widget) display() {
	widget.Redraw(widget.content)
}

func (widget *Widget) displayError() {
	title := fmt.Sprintf("%s - Error", widget.CommonSettings().Title)
	if widget.configError != nil {
		widget.Redraw(func() (string, string, bool) {
			return title, fmt.Sprintf("Error: \n [red]%v[white]", widget.configError), false
		})
		return
	}
	widget.Redraw(func() (string, string, bool) { return title, "Error", false })
}

func (widget *Widget) content() (string, string, bool) {
	project := widget.currentGitlabProject()
	if project == nil {
		return widget.CommonSettings().Title, " Gitlab pipeline data is unavailable ", true
	}

	// reset the selectable item collection on every render
	widget.Items = make([]string, 0)
	widget.SetItemCount(0)

	title := fmt.Sprintf("%s - [green]%s[white]", widget.CommonSettings().Title, project.path)

	_, _, width, _ := widget.View.GetRect()
	str := widget.settings.PaginationMarker(len(widget.GitlabProjects), widget.Idx, width) + "\n"
	str += widget.renderPipelines(project.Pipelines)

	return title, str, false
}

// statusColor maps a GitLab pipeline status to a display colour.
func statusColor(status string) string {
	switch status {
	case "success":
		return "green"
	case "failed":
		return "red"
	case "running", "pending":
		return "yellow"
	case "canceled", "skipped":
		return "grey"
	default:
		return "white"
	}
}

func (widget *Widget) renderPipelines(pipelines []*glab.PipelineInfo) string {
	if len(pipelines) == 0 {
		return " [grey]no pipelines[white]\n"
	}

	maxItems := widget.GetItemCount()

	str := ""
	for idx, p := range pipelines {
		when := ""
		if p.CreatedAt != nil {
			when = p.CreatedAt.Local().Format("01/02 15:04")
		}

		str += fmt.Sprintf(
			` ["%d"][%s]%-9s[white][grey]%s[white] %s`,
			maxItems+idx, statusColor(p.Status), p.Status, when, p.Ref,
		)
		str += "\n"
		widget.Items = append(widget.Items, p.WebURL)
	}
	widget.SetItemCount(maxItems + len(pipelines))

	return str
}
