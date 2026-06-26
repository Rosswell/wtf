package gitlab

import (
	glab "gitlab.com/gitlab-org/api/client-go"
)

type context struct {
	client *glab.Client
	user   *glab.User
}

func newContext(settings *Settings) (*context, error) {
	baseURL := settings.domain
	gitlabClient, _ := glab.NewClient(settings.apiKey, glab.WithBaseURL(baseURL))

	user, _, err := gitlabClient.Users.CurrentUser()

	if err != nil {
		return nil, err
	}

	ctx := &context{
		client: gitlabClient,
		user:   user,
	}

	return ctx, nil
}

// MergeRequest pairs a GitLab merge request with its computed approval status
// so the display can tag each row as draft / approved / open.
type MergeRequest struct {
	*glab.BasicMergeRequest
	Approved bool
}

type GitlabProject struct {
	context *context
	path    string

	ReviewMergeRequests   []*MergeRequest
	AuthoredMergeRequests []*MergeRequest
}

func NewGitlabProject(context *context, projectPath string) *GitlabProject {
	project := GitlabProject{
		context: context,
		path:    projectPath,
	}

	return &project
}

// Refresh reloads the gitlab data via the Gitlab API
func (project *GitlabProject) Refresh() {
	project.ReviewMergeRequests = project.withApproval(project.loadReviewMergeRequests())
	project.AuthoredMergeRequests = project.withApproval(project.loadAuthoredMergeRequests())
}

/* -------------------- Unexported Functions -------------------- */

// myReviewMergeRequests returns merge requests where the current user is a
// reviewer.
func (project *GitlabProject) myReviewMergeRequests() []*MergeRequest {
	return project.ReviewMergeRequests
}

// myMergeRequests returns merge requests authored by the current user.
func (project *GitlabProject) myMergeRequests() []*MergeRequest {
	return project.AuthoredMergeRequests
}

func (project *GitlabProject) loadReviewMergeRequests() ([]*glab.BasicMergeRequest, error) {
	state := "opened"
	opts := glab.ListProjectMergeRequestsOptions{
		State:      &state,
		ReviewerID: glab.ReviewerID(project.context.user.ID),
	}

	mrs, _, err := project.context.client.MergeRequests.ListProjectMergeRequests(project.path, &opts)

	if err != nil {
		return nil, err
	}

	return mrs, nil
}

func (project *GitlabProject) loadAuthoredMergeRequests() ([]*glab.BasicMergeRequest, error) {
	state := "opened"
	opts := glab.ListProjectMergeRequestsOptions{
		State:    &state,
		AuthorID: &project.context.user.ID,
	}

	mrs, _, err := project.context.client.MergeRequests.ListProjectMergeRequests(project.path, &opts)

	if err != nil {
		return nil, err
	}

	return mrs, nil
}

// withApproval wraps each merge request with its approval status. A load error
// yields an empty slice so the section renders as "none" rather than crashing.
func (project *GitlabProject) withApproval(mrs []*glab.BasicMergeRequest, err error) []*MergeRequest {
	if err != nil {
		return nil
	}

	wrapped := make([]*MergeRequest, 0, len(mrs))
	for _, mr := range mrs {
		wrapped = append(wrapped, &MergeRequest{
			BasicMergeRequest: mr,
			Approved:          project.isApproved(mr.IID),
		})
	}

	return wrapped
}

// isApproved reports whether every approval rule on the merge request is
// satisfied. An MR with no rules (no approval required) is not labelled
// approved.
func (project *GitlabProject) isApproved(iid int) bool {
	state, _, err := project.context.client.MergeRequestApprovals.GetApprovalState(project.path, iid)
	if err != nil || state == nil || len(state.Rules) == 0 {
		return false
	}

	for _, rule := range state.Rules {
		if !rule.Approved {
			return false
		}
	}

	return true
}
