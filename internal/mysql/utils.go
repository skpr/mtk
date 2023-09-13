package mysql

import (
	"bytes"
	"fmt"
	"io"

	"github.com/asaskevich/govalidator"
)

func getValue(raw string) (string, error) {
	if raw == "" {
		return "''", nil
	}

	escaped, err := escape(raw)
	if err != nil {
		return "", err
	}

	if govalidator.IsInt(raw) {
		return escaped, nil
	}

	return fmt.Sprintf("'%s'", escaped), nil
}

func escape(str string) (string, error) {
	var (
		esc  string
		buf  bytes.Buffer
		last = 0
	)

	for i, c := range str {
		switch c {
		case 0:
			esc = `\0`
		case '\n':
			esc = `\n`
		case '\r':
			esc = `\r`
		case '\\':
			esc = `\\`
		case '\'':
			esc = `\'`
		case '"':
			esc = `\"`
		case '\032':
			esc = `\Z`
		default:
			continue
		}

		if _, err := io.WriteString(&buf, str[last:i]); err != nil {
			return "", err
		}

		if _, err := io.WriteString(&buf, esc); err != nil {
			return "", err
		}

		last = i + 1
	}

	if _, err := io.WriteString(&buf, str[last:]); err != nil {
		return "", err
	}

	return buf.String(), nil
}
