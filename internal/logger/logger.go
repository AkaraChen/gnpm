package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var (
	// Output writers
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr

	// Colors - minimal usage
	dim    = color.New(color.FgHiBlack).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

// Command prints a command that will be executed (to stderr)
func Command(cmd string) {
	fmt.Fprintf(stderr, "%s %s\n", dim("$"), cmd)
}

// DryRun prints a command in dry-run mode (command to stdout, context to stderr)
func DryRun(cmd string, workDir string) {
	fmt.Fprintln(stdout, cmd)
	fmt.Fprintf(stderr, "%s %s\n", yellow("dry-run"), dim(workDir))
}

// Info prints an info message (to stderr)
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(stderr, msg)
}

// Success prints a success message (to stderr)
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(stderr, "%s %s\n", green("✓"), msg)
}

// Warn prints a warning message (to stderr)
func Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(stderr, "%s %s\n", yellow("!"), msg)
}

// Error prints an error message (to stderr)
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(stderr, "%s %s\n", red("✗"), msg)
}

// List prints a list item (to stderr)
func List(item string) {
	fmt.Fprintf(stderr, "  %s %s\n", dim("•"), item)
}

// Header prints a section header (to stderr)
func Header(title string) {
	fmt.Fprintf(stderr, "%s\n", title)
}

// Dim prints dimmed text (to stderr)
func Dim(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(stderr, dim(msg))
}

// Plain prints plain text to stdout
func Plain(format string, args ...interface{}) {
	fmt.Fprintf(stdout, format, args...)
}

// Plainln prints plain text to stdout with newline
func Plainln(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintln(stdout, msg)
}
