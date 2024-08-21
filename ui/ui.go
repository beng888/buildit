package ui

func GetAttrValue(attr any) string {
	if value, ok := attr.(string); ok {
		return value
	}
	return ""
}
