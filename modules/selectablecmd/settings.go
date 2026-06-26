package selectablecmd

import (
	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
	"github.com/wtfutil/wtf/utils"
)

const (
	defaultFocusable = true
	defaultTitle     = "SelectableCmd"
)

// Settings for the selectablecmd widget.
//
// The configured command must emit one row per line, with columns separated by
// the `separator` (default: tab). Rows are rendered as selectable lines; press
// Enter/o on a selected row to open the URL found in `urlColumn` in the browser.
type Settings struct {
	*cfg.Common

	cmd            string   `help:"The command to run. Its stdout is parsed into selectable rows."`
	args           []string `help:"Arguments to the command, one element per array item."`
	workingDir     string   `help:"Working directory the command runs in." optional:"true"`
	separator      string   `help:"Column separator in the command output. Default: tab." optional:"true"`
	urlColumn      int      `help:"Zero-based index of the column holding the open-in-browser URL. -1 means the last column. Default: -1." optional:"true"`
	displayColumns []int    `help:"Zero-based indexes of columns to display, in order. Empty means all columns except the URL column." optional:"true"`
	header         bool     `help:"Treat the first output line as a non-selectable header row. Default: false." optional:"true"`
}

// NewSettingsFromYAML loads the selectablecmd portion of the WTF config.
func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		cmd:        ymlConfig.UString("cmd"),
		args:       utils.ToStrs(ymlConfig.UList("args")),
		workingDir: ymlConfig.UString("workingDir", "."),
		separator:  ymlConfig.UString("separator", "\t"),
		urlColumn:  ymlConfig.UInt("urlColumn", -1),
		header:     ymlConfig.UBool("header", false),
	}

	for _, col := range ymlConfig.UList("displayColumns") {
		if i, ok := col.(int); ok {
			settings.displayColumns = append(settings.displayColumns, i)
		}
	}

	return &settings
}

func (widget *Widget) ConfigText() string {
	return utils.HelpFromInterface(Settings{})
}
