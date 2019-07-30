package config

import (
	"fmt"
)

// Map list for sanitizing the MySQL dump.
func (s Sanitize) Map() map[string]map[string]string {
	selectMap := make(map[string]map[string]string)

	for _, table := range s.Tables {
		tableMap := make(map[string]string, len(table.Fields))

		for _, field := range table.Fields {
			if field.Value == "" {
				field.Value = DefaultPlaceholder
			}

			tableMap[field.Name] = fmt.Sprintf("'%s'", field.Value)
		}

		selectMap[table.Name] = tableMap
	}

	return selectMap
}
