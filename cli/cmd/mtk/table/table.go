package table

import (
	"github.com/alecthomas/kingpin"

	"github.com/skpr/mtk/cmd/mtk/table/list"
)

// Command initializes the config commands.
func Command(app *kingpin.Application) {
	cmd := app.Command("table", "Config operations")
	list.Command(cmd)
}
