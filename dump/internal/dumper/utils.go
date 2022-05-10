package dumper

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"io"
)

func getValue(raw string) string {
	if raw == "" {
		return "''"
	}

	if govalidator.IsInt(raw) {
		return fmt.Sprintf("%s", escape(raw))
	}

	return fmt.Sprintf("'%s'", escape(raw))
}

func escape(str string) string {
	var esc string
	var buf bytes.Buffer
	last := 0
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
		io.WriteString(&buf, str[last:i])
		io.WriteString(&buf, esc)
		last = i + 1
	}
	io.WriteString(&buf, str[last:])
	return buf.String()
}
