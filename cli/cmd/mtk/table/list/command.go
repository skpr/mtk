package list

import (
	"database/sql"
	"fmt"

	"github.com/alecthomas/kingpin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/skpr/mtk/internal/dbutils"
	"github.com/skpr/mtk/pkg/envar"
)

type command struct {
	MaxConn  int
	Host     string
	User     string
	Password string
	Protocol string
	Port     string
	Database string
	Globs    []string
}

func (cmd *command) run(c *kingpin.ParseContext) error {
	conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", cmd.User, cmd.Password, cmd.Protocol, cmd.Host, cmd.Port, cmd.Database)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	tables, err := dbutils.ListTables(db, cmd.Globs)
	if err != nil {
		return err
	}

	for _, table := range tables {
		fmt.Println(table)
	}

	return nil
}

// Command which dumps a database or table.
func Command(c *kingpin.CmdClause) {
	cmd := new(command)

	command := c.Command("list", "List tables using a list of globs").Action(cmd.run)

	command.Flag("max-conn", "Sets the maximum number of open connections to the database").Default("50").Envar(envar.MaxConn).IntVar(&cmd.MaxConn)
	command.Flag("host", "Host on which MySQL server is located").Short('h').Required().Envar(envar.Hostname).StringVar(&cmd.Host)
	command.Flag("user", "MySQL user name to use when connecting to server").Short('u').Required().Envar(envar.Username).StringVar(&cmd.User)
	command.Flag("password", "Password to use when connecting to server").Short('p').Required().Envar(envar.Password).StringVar(&cmd.Password)
	command.Flag("protocol", "Connection protocol to use").Default("tcp").Envar(envar.Protocol).StringVar(&cmd.Protocol)
	command.Flag("port", "TCP/IP port number for connection").Short('P').Default("3306").Envar(envar.Port).StringVar(&cmd.Port)

	command.Arg("database", "Name of the database to use when connecting to the server").Required().Envar(envar.Database).StringVar(&cmd.Database)
	command.Arg("glob", "A list of table name globs").Required().Envar(envar.Table).StringsVar(&cmd.Globs)
}
