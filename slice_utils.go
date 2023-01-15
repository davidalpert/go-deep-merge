package DeepMerge

// toStringSlice converts a slice of unknown types to a slice of their string representations
func toStringSlice(items []interface{}) []string {
	var ss = make([]string, len(items))
	for i, v := range items {
		ss[i] = v.(string)
	}
	return ss
}

// toAbstractSlice converts a slice of strings to a slice of unknown types
func toAbstractSlice(items []string) []interface{} {
	var ss = make([]interface{}, len(items))
	for i, v := range items {
		ss[i] = v
	}
	return ss
}
