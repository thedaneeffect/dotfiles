package ui

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
)

// Success prints a success message in green
func Success(format string, a ...interface{}) {
	green.Printf("✓ "+format+"\n", a...)
}

// Error prints an error message in red
func Error(format string, a ...interface{}) {
	red.Printf("✗ "+format+"\n", a...)
}

// Warning prints a warning message in yellow
func Warning(format string, a ...interface{}) {
	yellow.Printf("⊘ "+format+"\n", a...)
}

// Info prints an info message in yellow
func Info(format string, a ...interface{}) {
	yellow.Printf("→ "+format+"\n", a...)
}

// Plain prints without color or prefix
func Plain(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}
