package dumptables

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/skpr/mtk/cmd/mtk/dump"
	"github.com/skpr/mtk/internal/mysql"
	"github.com/skpr/mtk/pkg/config"
	"github.com/skpr/mtk/pkg/envar"
)

const cmdLong = `
  Dumps a sanitized output of a MySQL database into a directory, with one file created per table`

const cmdExample = `
  export MTK_HOSTNAME=localhost
  export MTK_USERNAME=test
  export MTK_PASSWORD=test

  # Dump all database tables.
  mtk dump-tables <database> <directory>

  # Dump all database tables using config file
  mtk dump-tables <database> --config <config file> <directory>`

// Options is the commandline options for 'dump' sub command
type Options struct {
	ConfigFile         string
	ExtendedInsertRows int
	SingleTransaction  bool
}

// NewOptions will return a new Options.
func NewOptions() Options {
	return Options{}
}

// NewCommand will return a new Cobra command.
func NewCommand(conn *mysql.Connection, provider, rdsRegion, rdsS3uri string) *cobra.Command {
	o := NewOptions()

	cmd := &cobra.Command{
		Use:                   "dump-tables",
		DisableFlagsInUseLine: true,
		Short:                 "Dumps a sanitized output of a MySQL database into a directory.",
		Args:                  cobra.ExactArgs(2),
		Long:                  cmdLong,
		Example:               cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				database  = args[0]
				directory = args[1]
			)

			if _, err := os.Stat(directory); os.IsNotExist(err) {
				// Fail if directory does not exist
				return fmt.Errorf("directory does not exist: %s", directory)
			}

			logger := log.New(os.Stderr, "", 0)

			cfg, err := config.Load(o.ConfigFile)
			if err != nil {
				return fmt.Errorf("failed to load config file: %w", err)
			}

			return o.Run(cmd.Context(), os.Stdout, logger, conn, database, directory, provider, rdsRegion, rdsS3uri, cfg)
		},
	}

	cmd.Flags().StringVar(&o.ConfigFile, "config", envar.GetStringWithFallback("", envar.Config), "Path to the configuration file which contains the rules")
	cmd.Flags().IntVar(&o.ExtendedInsertRows, "extended-insert-rows", envar.GetIntWithFallback(1000, envar.ExtendedInsertRows), "The number of rows to batch per INSERT statement")
	cmd.Flags().BoolVar(&o.SingleTransaction, "single-transaction", true, "No changes that occur to InnoDB tables during the dump will be included in the dump")

	return cmd
}

// Run will execute the dump command.
func (o *Options) Run(ctx context.Context, w io.Writer, logger *log.Logger, conn *mysql.Connection, database, directory, provider, region, uri string, cfg config.Rules) error {
	db, err := conn.Open(ctx, database, o.SingleTransaction)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	defer db.Close()

	client := mysql.NewClient(db, logger, provider, region, uri)

	return o.runDumpTables(ctx, w, directory, client, cfg)
}

// Helper function to dump all tables in a database.
func (o *Options) runDumpTables(ctx context.Context, w io.Writer, directory string, client *mysql.Client, cfg config.Rules) error {
	tables, err := client.QueryTables(ctx)
	if err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	for _, table := range tables {
		// Write to a specific file.
		file, err := os.Create(fmt.Sprintf("%s/%s.sql", directory, table))
		if err != nil {
			return fmt.Errorf("failed to create file for table %q: %w", table, err)
		}
		defer file.Close()

		if err := dump.RunDumpTable(ctx, file, client, table, o.ExtendedInsertRows, cfg); err != nil {
			return fmt.Errorf("failed to dump table %q: %w", table, err)
		}
	}

	return nil
}
