package dumper

import (
	"database/sql"
	"fmt"
	"io"
	"strings"
	"time"
)

// WriteHeader is intended to be added at the beginning of a dump to manage database configuration.
// @todo, This header was taken from a mariadb mysqldump command. We need to determine if utf8mb4 should be configurable.
func (d *Client) WriteHeader(w io.Writer) error {
	header := `SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT;
SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS;
SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION;
SET NAMES utf8mb4;
SET @OLD_TIME_ZONE=@@TIME_ZONE;
SET TIME_ZONE='+00:00';
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0;`
	_, err := fmt.Fprintln(w, header)
	return err
}

// WriteFooter is intended to be added at the end of a dump to manage database configuration.
func (d *Client) WriteFooter(w io.Writer) error {
	footer := `SET TIME_ZONE=@OLD_TIME_ZONE;
SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT;
SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS;
SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION;
SET SQL_NOTES=@OLD_SQL_NOTES;`
	_, err := fmt.Fprintln(w, footer)
	return err
}

// WriteDumpCompleted is used for determining the freshness of a dump.
func (d *Client) WriteDumpCompleted(w io.Writer) error {
	_, err := fmt.Fprintln(w, "-- Dump completed on:", time.Now())
	return err
}

// WriteTableLockWrite to be used for a dump script.
func (d *Client) WriteTableLockWrite(w io.Writer, table string) {
	fmt.Fprintf(w, "LOCK TABLES `%s` WRITE;\n", table)
}

// WriteTableDisableKeys to be used for a dump script.
func (d *Client) WriteTableDisableKeys(w io.Writer, table string) {
	fmt.Fprintf(w, "ALTER TABLE `%s` DISABLE KEYS;\n", table)
}

// WriteTableEnableKeys to be used for a dump script.
func (d *Client) WriteTableEnableKeys(w io.Writer, table string) {
	fmt.Fprintf(w, "ALTER TABLE `%s` ENABLE KEYS;\n", table)
}

// WriteUnlockTables to be used for a dump script.
func (d *Client) WriteUnlockTables(w io.Writer) {
	fmt.Fprintln(w, "UNLOCK TABLES;")
}

// WriteCreateTable script used when dumping a database.
func (d *Client) WriteCreateTable(w io.Writer, table string) error {
	fmt.Fprintf(w, "\n--\n-- Structure for table `%s`\n--\n\n", table)
	fmt.Fprintf(w, "DROP TABLE IF EXISTS `%s`;\n", table)

	fmt.Fprintln(w, "SET @saved_cs_client = @@character_set_client;")
	fmt.Fprintln(w, "SET character_set_client = utf8;")

	row := d.DB.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`", table))

	var name, ddl string

	if err := row.Scan(&name, &ddl); err != nil {
		return err
	}

	fmt.Fprintf(w, "%s;\n", ddl)

	fmt.Fprintln(w, "SET character_set_client = @saved_cs_client;")

	return nil
}

// WriteTableHeader which contains debug information.
func (d *Client) WriteTableHeader(w io.Writer, table string) (uint64, error) {
	fmt.Fprintf(w, "\n--\n-- Data for table `%s`", table)

	count, err := d.GetRowCountForTable(table)
	if err != nil {
		return 0, err
	}

	fmt.Fprintf(w, " -- %d rows\n--\n\n", count)

	return count, nil
}

// WriteTableData for a specific table.
func (d *Client) WriteTableData(w io.Writer, table string) error {
	rows, columns, err := d.selectAllDataForTable(table)
	if err != nil {
		return err
	}

	defer rows.Close()

	values := make([]*sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	query := fmt.Sprintf("INSERT INTO `%s` VALUES", table)

	var data []string

	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			return err
		}

		var vals []string

		for _, col := range values {
			val := "NULL"

			if col != nil {
				val = fmt.Sprintf("'%s'", escape(string(*col)))
			}

			vals = append(vals, val)
		}

		data = append(data, fmt.Sprintf("( %s )", strings.Join(vals, ", ")))

		if d.ExtendedInsertRows == 0 {
			continue
		}

		if len(data) >= d.ExtendedInsertRows {
			fmt.Fprintf(w, "%s\n%s;\n", query, strings.Join(data, ",\n"))
			data = make([]string, 0)
		}
	}

	if len(data) > 0 {
		fmt.Fprintf(w, "%s\n%s;\n", query, strings.Join(data, ",\n"))
	}

	return nil
}

// WriteTables will create a script for all tables.
func (d *Client) writeTables(w io.Writer) error {
	tables, err := d.QueryTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		if err := d.writeTable(w, table); err != nil {
			return err
		}
	}

	return nil
}

// WriteTable allows for a single table dump script.
func (d *Client) writeTable(w io.Writer, table string) error {
	if d.FilterMap[strings.ToLower(table)] == OperationIgnore {
		return nil
	}

	skipData := d.FilterMap[strings.ToLower(table)] == OperationNoData
	if !skipData && d.UseTableLock {
		d.LockTableReading(table)
		d.FlushTable(table)
	}

	d.WriteCreateTable(w, table)

	if !skipData {
		cnt, err := d.WriteTableHeader(w, table)
		if err != nil {
			return err
		}

		if cnt > 0 {
			d.WriteTableLockWrite(w, table)
			d.WriteTableDisableKeys(w, table)
			d.WriteTableData(w, table)

			fmt.Fprintln(w)

			d.WriteTableEnableKeys(w, table)
			d.WriteUnlockTables(w)

			if d.UseTableLock {
				d.UnlockTables()
			}
		}
	}

	return nil
}
