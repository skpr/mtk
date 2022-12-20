package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/skpr/mtk/dump/internal/dbutils"
	"github.com/skpr/mtk/dump/internal/dumper"
	"github.com/skpr/mtk/dump/pkg/config"
	"github.com/skpr/mtk/dump/pkg/envar"
)

var (
	cliMaxConn  = kingpin.Flag("max-conn", "Sets the maximum number of open connections to the database").Default("50").Envar(envar.MaxConn).Int()
	cliConfig   = kingpin.Flag("config", "Path to the configuration file which contains the rules").Default(".mtk/dump.yml").Envar(envar.Config).String()
	cliHostname = kingpin.Flag("host", "Host on which MySQL server is located").Short('h').Required().Envar(envar.Hostname).String()
	cliUsername = kingpin.Flag("user", "MySQL user name to use when connecting to server").Short('u').Required().Envar(envar.Username).String()
	cliPassword = kingpin.Flag("password", "Password to use when connecting to server").Short('p').Required().Envar(envar.Password).String()
	cliProtocol = kingpin.Flag("protocol", "Connection protocol to use").Default("tcp").Envar(envar.Protocol).String()
	cliPort     = kingpin.Flag("port", "TCP/IP port number for connection").Short('P').Default("3306").Envar(envar.Port).String()
	cliDatabase = kingpin.Arg("name", "Name of the database to use when connecting to the server").Required().Envar(envar.Database).String()
)

func main() {
	kingpin.Parse()

	conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", *cliUsername, *cliPassword, *cliProtocol, *cliHostname, *cliPort, *cliDatabase)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.SetMaxOpenConns(*cliMaxConn)

	// Load the config.
	cfg, err := config.Load(*cliConfig)
	if err != nil {
		panic(err)
	}

	// Dump using the config.
	err = dump(os.Stdout, os.Stderr, db, cfg)
	if err != nil {
		panic(err)
	}
}

//
func dump(stdout, stderr io.Writer, db *sql.DB, cfg config.Rules) error {
	d := dumper.NewClient(db, log.New(stderr, "", 0))

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
		// @todo, Needs to be const values in mysqlsuperdump
		d.FilterMap[table] = "nodata"
	}

	// Assign ignore tables.
	for _, table := range ignore {
		// @todo, Needs to be const values in mysqlsuperdump
		d.FilterMap[table] = "ignore"
	}

	// Assign our sanitization rules to the dumper.
	d.SelectMap = cfg.SanitizeMap()

	// Assign conditional row rules to the dumper.
	d.WhereMap = cfg.WhereMap()

	return d.DumpTables(stdout)
}
