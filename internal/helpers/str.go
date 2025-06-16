package helpers

func TruncateString(str string, max int) string {
	if len(str) <= max {
		return str
	}
	return str[:max] + "..."
}
