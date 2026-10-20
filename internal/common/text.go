// Package common holds small helpers used across the project. It imports nothing from the
// project, so any package can use it.
package common

// EqualFold compares ASCII strings case-insensitively (strings.EqualFold pulls the Unicode tables
// into RAM).
func EqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if lower(a[i]) != lower(b[i]) {
			return false
		}
	}
	return true
}

func lower(char byte) byte {
	if 'A' <= char && char <= 'Z' {
		return char + 'a' - 'A'
	}

	return char
}

// Flag formats a boolean as "1" or "0", like the factory firmware's log lines.
func Flag(value bool) string {
	if value {
		return "1"
	}

	return "0"
}
