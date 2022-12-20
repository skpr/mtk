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
	header := `/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;`
	_, err := fmt.Fprintln(w, header)
	return err
}

// WriteFooter is intended to be added at the end of a dump to manage database configuration.
func (d *Client) WriteFooter(w io.Writer) error {
	footer := `/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;`
	_, err := fmt.Fprintln(w, footer)
	return err
}

// WriteDumpCompleted is used for determining the freshness of a dump.
func (d *Client) WriteDumpCompleted(w io.Writer) error {
	_, err := fmt.Fprintln(w, "-- Dump completed on:", time.Now())
	return err
}

// WriteAutoCommitOff to be used for a dump script.
func (d *Client) WriteAutoCommitOff(w io.Writer) {
	fmt.Fprintln(w, "set autocommit=0;")
}

// WriteCommit to be used for a dump script.
func (d *Client) WriteCommit(w io.Writer) {
	fmt.Fprintln(w, "commit;")
}

// WriteTableLockWrite to be used for a dump script.
func (d *Client) WriteTableLockWrite(w io.Writer, table string) {
	fmt.Fprintf(w, "LOCK TABLES `%s` WRITE;\n", table)
}

// WriteTableDisableKeys to be used for a dump script.
func (d *Client) WriteTableDisableKeys(w io.Writer, table string) {
	fmt.Fprintf(w, "/*!40000 ALTER TABLE `%s` DISABLE KEYS */;\n", table)
}

// WriteTableEnableKeys to be used for a dump script.
func (d *Client) WriteTableEnableKeys(w io.Writer, table string) {
	fmt.Fprintf(w, "/*!40000 ALTER TABLE `%s` ENABLE KEYS */;\n", table)
}

// WriteUnlockTables to be used for a dump script.
func (d *Client) WriteUnlockTables(w io.Writer) {
	fmt.Fprintln(w, "UNLOCK TABLES;")
}

// WriteCreateTable script used when dumping a database.
func (d *Client) WriteCreateTable(w io.Writer, table string) error {
	d.Logger.Println("Dumping structure for table:", table)

	fmt.Fprintf(w, "\n--\n-- Structure for table `%s`\n--\n\n", table)
	fmt.Fprintf(w, "DROP TABLE IF EXISTS `%s`;\n", table)

	fmt.Fprintln(w, "/*!40101 SET @saved_cs_client     = @@character_set_client */;")
	fmt.Fprintln(w, "/*!40101 SET character_set_client = utf8 */;")

	row := d.DB.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`", table))

	var name, ddl string

	if err := row.Scan(&name, &ddl); err != nil {
		return err
	}

	fmt.Fprintf(w, "%s;\n", ddl)

	fmt.Fprintln(w, "/*!40101 SET character_set_client = @saved_cs_client */;")

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
	d.Logger.Println("Dumping data for table:", table)

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
				val = getValue(string(*col))
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

	if skipData {
		return nil
	}

	cnt, err := d.WriteTableHeader(w, table)
	if err != nil {
		return err
	}

	if cnt == 0 {
		return nil
	}

	d.WriteTableLockWrite(w, table)
	d.WriteTableDisableKeys(w, table)
	d.WriteAutoCommitOff(w)
	d.WriteTableData(w, table)
	d.WriteTableEnableKeys(w, table)
	d.WriteUnlockTables(w)
	d.WriteCommit(w)

	return nil
}
