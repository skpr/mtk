package dump

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/skpr/mtk/internal/mysql"
	"github.com/skpr/mtk/pkg/config"
	"github.com/skpr/mtk/pkg/envar"
)

const cmdLong = `
  Dumps a sanitized output of a MySQL database.`

const cmdExample = `
  export MTK_HOSTNAME=localhost
  export MTK_USERNAME=test
  export MTK_PASSWORD=test

  # Dump all database tables.
  mtk dump <database> > test.sql

  # Dump all database tables using config file
  mtk dump <database> --config <config file> > test.sql

  # List all database tables and dump each table to a file.
  mtk table list <database> | xargs -I {} sh -c "mtk dump <database> '{}' > '{}.sql'"`

// Options is the commandline options for 'config' sub command
type Options struct {
	ConfigFile         string
	ExtendedInsertRows int
}

func NewOptions() Options {
	return Options{}
}

func NewCommand(conn *mysql.Connection) *cobra.Command {
	o := NewOptions()

	cmd := &cobra.Command{
		Use:                   "dump > test.sql",
		DisableFlagsInUseLine: true,
		Short:                 "Dumps a sanitized output of a MySQL database.",
		Args:                  cobra.RangeArgs(1, 2),
		Long:                  cmdLong,
		Example:               cmdExample,
		Run: func(cmd *cobra.Command, args []string) {
			var (
				database = args[0]
				table    = ""
			)

			if len(args) == 2 {
				table = args[1]
			}

			logger := log.New(os.Stderr, "", 0)

			cfg, err := config.Load(o.ConfigFile)
			if err != nil {
				panic(err)
			}

			if err := o.Run(os.Stdout, logger, conn, database, table, cfg); err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringVar(&o.ConfigFile, "config", envar.GetStringWithFallback("", envar.Config), "Path to the configuration file which contains the rules")
	cmd.Flags().IntVar(&o.ExtendedInsertRows, "extended-insert-rows", envar.GetIntWithFallback(1000, envar.ExtendedInsertRows), "The number of rows to batch per INSERT statement")

	return cmd
}

func (o *Options) Run(w io.Writer, logger *log.Logger, conn *mysql.Connection, database, table string, cfg config.Rules) error {
	db, err := conn.Open(database)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	defer db.Close()

	client := mysql.NewClient(db, logger)

	// Get a list of tables to nodata, passed through a globber.
	nodata, err := client.ListTablesByGlob(cfg.NoData)
	if err != nil {
		return err
	}

	// Get a list of tables to ignore, passed through a globber.
	ignore, err := client.ListTablesByGlob(cfg.Ignore)
	if err != nil {
		return err
	}

	params := mysql.DumpParams{
		ExtendedInsertRows: o.ExtendedInsertRows,
	}

	// Assign nodata tables.
	params.FilterMap = make(map[string]string)
	for _, table := range nodata {
		// @todo, Needs to be const values in mysqlsuperdump
		params.FilterMap[table] = "nodata"
	}

	// Assign ignore tables.
	for _, table := range ignore {
		// @todo, Needs to be const values in mysqlsuperdump
		params.FilterMap[table] = "ignore"
	}

	// Assign our sanitization rules to the dumper.
	params.SelectMap = cfg.SanitizeMap()

	// Assign conditional row rules to the dumper, passed through a globber.
	where := make(map[string]string, 0)

	for glob, condition := range cfg.WhereMap() {
		tables, err := client.ListTablesByGlob([]string{glob})
		if err != nil {
			return err
		}

		for _, table := range tables {
			where[table] = condition
		}
	}

	params.WhereMap = where

	if table != "" {
		return client.DumpTable(w, table, params)
	}

	return client.DumpTables(w, params)
}
