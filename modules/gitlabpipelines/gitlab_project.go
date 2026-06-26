package gitlabpipelines

import (
	glab "gitlab.com/gitlab-org/api/client-go"
)

type context struct {
	client *glab.Client
}

func newContext(settings *Settings) (*context, error) {
	gitlabClient, err := glab.NewClient(settings.apiKey, glab.WithBaseURL(settings.domain))
	if err != nil {
		return nil, err
	}

	return &context{client: gitlabClient}, nil
}

// GitlabProject holds the most recent pipelines for a single project.
type GitlabProject struct {
	context *context
	path    string
	count   int

	Pipelines []*glab.PipelineInfo
}

func NewGitlabProject(context *context, projectPath string, count int) *GitlabProject {
	return &GitlabProject{
		context: context,
		path:    projectPath,
		count:   count,
	}
}

// Refresh reloads the project's most recent pipelines via the GitLab API.
func (project *GitlabProject) Refresh() {
	project.Pipelines, _ = project.loadPipelines()
}

func (project *GitlabProject) loadPipelines() ([]*glab.PipelineInfo, error) {
	opts := glab.ListProjectPipelinesOptions{
		OrderBy:     glab.Ptr("id"),
		Sort:        glab.Ptr("desc"),
		ListOptions: glab.ListOptions{PerPage: project.count, Page: 1},
	}

	pipelines, _, err := project.context.client.Pipelines.ListProjectPipelines(project.path, &opts)
	if err != nil {
		return nil, err
	}

	return pipelines, nil
}
