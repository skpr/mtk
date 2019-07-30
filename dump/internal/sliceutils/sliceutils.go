package sliceutils

// AppendIfMissing from the slice.
func AppendIfMissing(slice []string, add string) []string {
	for _, existing := range slice {
		if existing == add {
			return slice
		}
	}

	return append(slice, add)
}
