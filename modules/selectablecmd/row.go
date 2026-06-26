package selectablecmd

import "strings"

// row is one parsed line of command output.
type row struct {
	cells []string
	url   string
}

// parseOutput splits raw command stdout into selectable rows using the
// configured separator and URL column. A header line is dropped when enabled.
func parseOutput(out string, settings *Settings) []row {
	rows := []row{}

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	for idx, line := range lines {
		if settings.header && idx == 0 {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}

		cells := strings.Split(line, settings.separator)

		urlIdx := settings.urlColumn
		if urlIdx < 0 {
			urlIdx = len(cells) - 1
		}

		url := ""
		if urlIdx >= 0 && urlIdx < len(cells) {
			url = strings.TrimSpace(cells[urlIdx])
		}

		rows = append(rows, row{cells: cells, url: url})
	}

	return rows
}

// display returns the cells to render for a row, honoring displayColumns when
// set, otherwise every cell except the URL column.
func (r row) display(settings *Settings) []string {
	urlIdx := settings.urlColumn
	if urlIdx < 0 {
		urlIdx = len(r.cells) - 1
	}

	if len(settings.displayColumns) > 0 {
		out := make([]string, 0, len(settings.displayColumns))
		for _, i := range settings.displayColumns {
			if i >= 0 && i < len(r.cells) {
				out = append(out, r.cells[i])
			}
		}
		return out
	}

	out := make([]string, 0, len(r.cells))
	for i, cell := range r.cells {
		if i == urlIdx {
			continue
		}
		out = append(out, cell)
	}
	return out
}
