package tweetgo

// Bool converts a bool to a bool pointer so that it supports nil values
func Bool(input bool) *bool {
	return &input
}

// String converts a string to a string pointer so that it supports nil values
func String(input string) *string {
	return &input
}

// Int converts an int to an int pointer so that it supports nil values
func Int(input int) *int {
	return &input
}

// Int64 converts an int64 to an int64 pointer so that it supports nil values
func Int64(input int64) *int64 {
	return &input
}

// Float64 converts a float64 to a float64 pointer so that it supports nil values
func Float64(input float64) *float64 {
	return &input
}
