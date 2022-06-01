package types

import "fmt"

// A data type for general purpose of parameters
type Parameter struct {
	// The name of the parameter
	Name *string

	// The value of the parameter
	Value *string
}

// A flag that accepts multiple strings as string array
type StringArrayFlag []string

// Output the value of flag
func (s *StringArrayFlag) String() string {
	return fmt.Sprint(*s)
}

// Append the value to array if the flag is used multiple times
func (s *StringArrayFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}
