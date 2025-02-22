package main

import "strings"

func isValidFile(filename string) bool {
	return strings.HasSuffix(filename, ".txt")
}
