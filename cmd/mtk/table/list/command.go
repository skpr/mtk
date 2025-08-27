package list

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/skpr/mtk/internal/mysql"
	"github.com/skpr/mtk/pkg/config"
	"github.com/skpr/mtk/pkg/envar"
)

const cmdLong = `
  List all tables in a database.

  Useful when chaining with the "dump" command for a "dump per table" model.`

const cmdExample = `
  export MTK_HOSTNAME=localhost
  export MTK_USERNAME=test
  export MTK_PASSWORD=test

  # List all database tables.
  mtk table list <database>

  # List all database tables using config file to skip certain tables.
  mtk table list <database> --config <config file>

  # List all database tables and dump each table to a file.
  mtk table list <database> | xargs -I {} sh -c "mtk dump <database> '{}' > '{}.sql'"`

// Options is the commandline options for 'list' sub command
type Options struct {
	ConfigFile        string
	SingleTransaction bool
}

// NewOptions returns a new Options struct.
func NewOptions() Options {
	return Options{}
}

// NewCommand returns a new cobra command.
func NewCommand(conn *mysql.Connection) *cobra.Command {
	o := NewOptions()

	cmd := &cobra.Command{
		Use:                   "list <database>",
		DisableFlagsInUseLine: true,
		Short:                 "List all tables in a database.",
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.New(os.Stderr, "", 0)

			cfg, err := config.Load(o.ConfigFile)
			if err != nil {
				return err
			}

			for _, database := range args {
				if err := o.Run(cmd.Context(), logger, conn, database, cfg.Ignore); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.ConfigFile, "config", envar.GetStringWithFallback("", envar.Config), "Path to the configuration file which contains the rules")
	cmd.Flags().BoolVar(&o.SingleTransaction, "single-transaction", true, "No changes that occur to InnoDB tables during the dump will be included in the dump")

	return cmd
}

// Run the command which will list all tables.
func (o *Options) Run(ctx context.Context, logger *log.Logger, conn *mysql.Connection, database string, exclude []string) error {
	db, err := conn.Open(ctx, database, o.SingleTransaction)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	defer db.Close()

	client := mysql.NewClient(db, logger, "", "", "")

	tables, err := client.QueryTables(ctx)
	if err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	skip, err := client.ListTablesByGlob(ctx, exclude)
	if err != nil {
		return fmt.Errorf("failed to list tables to skip: %w", err)
	}

	for _, table := range tables {
		if slices.Contains(skip, table) {
			continue
		}

		fmt.Println(table)
	}

	return nil
}
