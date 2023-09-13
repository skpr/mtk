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

var conn = new(mysql.Connection)

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
	cmd.PersistentFlags().StringVar(&conn.Hostname, "host", envar.GetStringWithFallback(envar.Hostname, ""), "Host address of the MySQL instance")
	cmd.PersistentFlags().StringVar(&conn.Username, "user", envar.GetStringWithFallback(envar.Username, ""), "Username used to connect to MySQL instance")
	cmd.PersistentFlags().StringVar(&conn.Password, "password", envar.GetStringWithFallback(envar.Password, ""), "Password used to connect to MySQL instance")
	cmd.PersistentFlags().StringVar(&conn.Protocol, "protocol", envar.GetStringWithFallback(envar.Protocol, "tcp"), "Connection protocol to use when connecting to MySQL instance")
	cmd.PersistentFlags().Int32Var(&conn.Port, "port", int32(envar.GetIntWithFallback(envar.Port, 3306)), "Port to connect to the MySQL instance on")
	cmd.PersistentFlags().IntVar(&conn.MaxConn, "max-conn", envar.GetIntWithFallback(envar.MaxConn, 50), "Sets the maximum number of open connections to the database")
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

	cmd.AddCommand(dump.NewCommand(conn))
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
