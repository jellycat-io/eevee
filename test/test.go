package test

import (
	"bytes"
)

// MakeInput takes code lines as strings or returns them in a single string separated by '\n'
func MakeInput(lines ...string) string {
	var out bytes.Buffer
	for _, line := range lines {
		out.WriteString(line)
		out.WriteString("\n")
	}
	return out.String()
}
