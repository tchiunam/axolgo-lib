package utility

// Return an empty string when string pointer is nil so
// it's easier to do printing.
func HushedStringPtr(s *string) *string {
	if s == nil {
		return new(string)
	} else {
		return s
	}
}
