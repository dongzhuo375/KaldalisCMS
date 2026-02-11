package utils

import "strconv"

// ParseInt64 parses base-10 int64; returns 0 on empty/invalid.
func ParseInt64(s string) int64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// ParseInt parses base-10 int; returns 0 on empty/invalid.
func ParseInt(s string) int {
	if s == "" {
		return 0
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}
