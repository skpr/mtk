package config

// SanitizeMap list for sanitizing the MySQL dump.
func (r Rules) SanitizeMap() map[string]map[string]string {
	selectMap := make(map[string]map[string]string)

	for table, fields := range r.Rewrite {
		selectMap[table] = fields
	}

	return selectMap
}

// WhereMap list for conditional row exports in the MySQL dump.
func (r Rules) WhereMap() map[string]string {
	whereMap := make(map[string]string)

	for table, condition := range r.Where {
		whereMap[table] = condition
	}

	return whereMap
}
