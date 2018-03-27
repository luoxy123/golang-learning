package polaris

// Bool is
func Bool(v bool) *bool {
	return &v
}

// BoolValue returns the value of the bool pointer passed in
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}

	return false
}

// String is
func String(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer
func StringValue(v *string) string {
	if v != nil {
		return *v
	}

	return ""
}
