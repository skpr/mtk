package dump

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/skpr/mtk/internal/dbutils"
	"github.com/skpr/mtk/pkg/config"
	"github.com/skpr/mtk/pkg/envar"
	"github.com/skpr/mysqlsuperdump/dumper"
)

type command struct {
	MaxConn  int
	Config   string
	Host     string
	User     string
	Password string
	Protocol string
	Port     string
	Database string
	Table    string
}

func (cmd *command) run(c *kingpin.ParseContext) error {
	conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", cmd.User, cmd.Password, cmd.Protocol, cmd.Host, cmd.Port, cmd.Database)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.SetMaxOpenConns(cmd.MaxConn)

	// Load the config.
	cfg, err := config.Load(cmd.Config)
	if err != nil {
		panic(err)
	}

	// Dump using the config.
	err = dump(os.Stdout, os.Stderr, db, cfg, cmd.Table)
	if err != nil {
		panic(err)
	}

	return nil
}

// Helper function to dump a database/table to stdout.
func dump(stdout, stderr io.Writer, db *sql.DB, cfg config.Rules, table string) error {
	logger := log.New(stderr, "", 0)

	d := dumper.NewMySQLDumper(db, logger)

	// Get a list of tables to nodata, passed through a globber.
	nodata, err := dbutils.ListTables(db, cfg.NoData)
	if err != nil {
		return err
	}

	// Get a list of tables to ignore, passed through a globber.
	ignore, err := dbutils.ListTables(db, cfg.Ignore)
	if err != nil {
		return err
	}

	// Assign nodata tables.
	d.FilterMap = make(map[string]string)
	for _, table := range nodata {
		d.FilterMap[table] = dumper.OperationNoData
	}

	// Assign ignore tables.
	for _, table := range ignore {
		d.FilterMap[table] = dumper.OperationIgnore
	}

	// Assign our sanitization rules to the dumper.
	d.SelectMap = cfg.SanitizeMap()

	// Assign conditional row rules to the dumper.
	d.WhereMap = cfg.WhereMap()

	if table != "" {
		return d.DumpTable(stdout, table)
	}

	return d.Dump(stdout)
}

// Command which dumps a database or table.
func Command(app *kingpin.Application) {
	cmd := new(command)

	command := app.Command("dump", "Dump a database or table.").Action(cmd.run)

	command.Flag("max-conn", "Sets the maximum number of open connections to the database").Default("50").Envar(envar.MaxConn).IntVar(&cmd.MaxConn)
	command.Flag("config", "Path to the configuration file which contains the rules").Default(".mtk/dump.yml").Envar(envar.Config).StringVar(&cmd.Config)
	command.Flag("host", "Host on which MySQL server is located").Short('h').Required().Envar(envar.Hostname).StringVar(&cmd.Host)
	command.Flag("user", "MySQL user name to use when connecting to server").Short('u').Required().Envar(envar.Username).StringVar(&cmd.User)
	command.Flag("password", "Password to use when connecting to server").Short('p').Required().Envar(envar.Password).StringVar(&cmd.Password)
	command.Flag("protocol", "Connection protocol to use").Default("tcp").Envar(envar.Protocol).StringVar(&cmd.Protocol)
	command.Flag("port", "TCP/IP port number for connection").Short('P').Default("3306").Envar(envar.Port).StringVar(&cmd.Port)

	command.Arg("database", "Name of the database to use when connecting to the server").Required().Envar(envar.Database).StringVar(&cmd.Database)
	command.Arg("table", "Name of the table of dump. All tables are dumped if omitted.").Envar(envar.Table).StringVar(&cmd.Table)
}
