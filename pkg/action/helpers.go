package action

func pointerString(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

func pointerBool(val *bool) bool {
	if val == nil {
		return false
	}
	return *val
}
