package common

func WithDefault(val string, defaultVal string) string {
	if val == "" {
		return defaultVal
	}
	return val
}
