package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jwalton/gchalk"
	"github.com/spf13/cobra"

	"github.com/skpr/mtk/cmd/mtk/dump"
	"github.com/skpr/mtk/cmd/mtk/table"
	"github.com/skpr/mtk/internal/mysql"
	"github.com/skpr/mtk/pkg/envar"
)

var (
	conn         = new(mysql.Connection)
	providerName string
	rdsRegion    string
	rdsS3Uri     string
)

const cmdExample = `
  export MTK_HOSTNAME=localhost
  export MTK_USERNAME=test
  export MTK_PASSWORD=test
  export MTK_CONFIG=config.yml

  # Dump all database tables.
  mtk dump <database> > test.sql

  # List all database tables.
  mtk table list <database>
`

var cmd = &cobra.Command{
	Use:     "mtk",
	Short:   "Toolkit for exporting, sanitizing and packaging MySQL database.",
	Example: cmdExample,
	Long: `  __  __ _______ _  __
 |  \/  |__   __| |/ /
 | \  / |  | |  | ' / 
 | |\/| |  | |  |  <  
 | |  | |  | |  | . \ 
 |_|  |_|  |_|  |_|\_\
	
Toolkit for exporting, sanitizing and packaging MySQL databases.`,
}

func init() {
	cmd.PersistentFlags().StringVar(&conn.Hostname, "host", envar.GetStringWithFallback("", envar.Hostname, envar.MySQLHostname), "Host address of the MySQL instance")
	cmd.PersistentFlags().StringVar(&conn.Username, "user", envar.GetStringWithFallback("", envar.Username, envar.MySQLUsername), "Username used to connect to MySQL instance")
	cmd.PersistentFlags().StringVar(&conn.Password, "password", envar.GetStringWithFallback("", envar.Password, envar.MySQLPassword), "Password used to connect to MySQL instance")
	cmd.PersistentFlags().StringVar(&conn.Protocol, "protocol", envar.GetStringWithFallback("tcp", envar.Protocol, envar.MySQLProtocol), "Connection protocol to use when connecting to MySQL instance")
	cmd.PersistentFlags().Int32Var(&conn.Port, "port", int32(envar.GetIntWithFallback(3306, envar.Port, envar.MySQLPort)), "Port to connect to the MySQL instance on")
	cmd.PersistentFlags().IntVar(&conn.MaxConn, "max-conn", envar.GetIntWithFallback(50, envar.MaxConn), "Sets the maximum number of open connections to the database")

	cmd.PersistentFlags().StringVar(&providerName, "provider", envar.GetStringWithFallback("stdout", envar.Provider), "The provider to use (either 'stdout' or 'rds')")

	// RDS Provider Flags.
	cmd.PersistentFlags().StringVar(&rdsRegion, "rds-region", envar.GetStringWithFallback("", envar.RDSRegion), "The AWS region to use for S3 when connecting to the MySQL RDS instance")
	cmd.PersistentFlags().StringVar(&rdsS3Uri, "rds-s3-uri", envar.GetStringWithFallback("", envar.RDSS3Uri), "The S3 URI to use for exporting to S3 when exporting data from the MySQL RDS instance")
}

func main() {
	cobra.AddTemplateFunc("StyleHeading", styleHeading)
	usageTemplate := cmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Aliases:`, `{{StyleHeading "Aliases:"}}`,
		`Examples:`, `{{StyleHeading "Examples:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Global Flags:`, `{{StyleHeading "Global Flags:"}}`,
	).Replace(usageTemplate)

	re := regexp.MustCompile(`(?m)^Flags:\s*$`)
	usageTemplate = re.ReplaceAllLiteralString(usageTemplate, `{{StyleHeading "Flags:"}}`)
	cmd.SetUsageTemplate(usageTemplate)

	cmd.AddCommand(dump.NewCommand(conn, providerName, rdsRegion, rdsS3Uri))
	cmd.AddCommand(table.NewCommand(conn))

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Helper function for styling headings in the usage template.
func styleHeading(data string) string {
	return gchalk.WithHex("#ee5622").Bold(data)
}
