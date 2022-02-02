package main

import (
	"os"

	"github.com/alecthomas/kingpin"

	"github.com/skpr/mtk/cmd/mtk/dump"
	"github.com/skpr/mtk/cmd/mtk/table"
)

func main() {
	app := kingpin.New("mtk", "MySQL Toolkit")

	dump.Command(app)
	table.Command(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
