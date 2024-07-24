package table

import (
	"github.com/spf13/cobra"

	"github.com/skpr/mtk/cmd/mtk/table/list"
	"github.com/skpr/mtk/internal/mysql"
)

const cmdLong = `
  Perform MySQL tasks related to database tables.`

const cmdExample = `
  # List all database tables.
  mtk table list <database>`

// NewCommand will execute the table command.
func NewCommand(conn *mysql.Connection) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "table",
		DisableFlagsInUseLine: true,
		Short:                 "Perform MySQL tasks related to database tables.",
		Long:                  cmdLong,
		Example:               cmdExample,
	}

	cmd.AddCommand(list.NewCommand(conn))

	return cmd
}
