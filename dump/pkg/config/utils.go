package config

// SanitizeMap list for sanitizing the MySQL dump.
func (r Rules) SanitizeMap() map[string]map[string]string {
	selectMap := make(map[string]map[string]string)

	for table, fields := range r.Sanitize {
		selectMap[table] = fields
	}

	return selectMap
}
