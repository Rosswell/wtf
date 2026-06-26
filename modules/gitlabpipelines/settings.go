package gitlabpipelines

import (
	"os"

	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
	"github.com/wtfutil/wtf/utils"
)

const (
	defaultFocusable = true
	defaultTitle     = "GitLab Pipelines"
)

// Settings defines the configuration properties for this module
type Settings struct {
	*cfg.Common

	apiKey   string   `help:"A GitLab personal access token. Requires at least read_api access."`
	domain   string   `help:"Your GitLab corporate domain."`
	projects []string `help:"A list of GitLab projects to fetch pipeline status for. Each is shown on its own page."`
	count    int      `help:"How many of the most recent pipelines to show per project." optional:"true"`
}

// NewSettingsFromYAML creates a new settings instance from a YAML config block
func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		apiKey: ymlConfig.UString("apiKey", ymlConfig.UString("apikey", os.Getenv("WTF_GITLAB_TOKEN"))),
		domain: ymlConfig.UString("domain", "https://gitlab.com"),
		count:  ymlConfig.UInt("count", 5),
	}

	cfg.ModuleSecret(name, globalConfig, &settings.apiKey).
		Service(settings.domain).Load()

	settings.projects = cfg.ParseAsMapOrList(ymlConfig, "projects")

	return &settings
}

func (widget *Widget) ConfigText() string {
	return utils.HelpFromInterface(Settings{})
}
