package log

import (
	"fmt"

	"github.com/fatih/color"
)

var infoColor = color.New(color.FgCyan).PrintfFunc()

// Infof is Printf prefixed with [INFO]
func Infof(format string, a ...interface{}) {
	infoColor("[INFO] ")
	fmt.Printf(format+"\n", a...)
}

var loadColor = color.New(color.FgGreen).PrintfFunc()

// Loadf is Printf prefixed with [LOAD]
func Loadf(format string, a ...interface{}) {
	loadColor("[LOAD] ")
	fmt.Printf(format+"\n", a...)
}

var warnColor = color.New(color.FgYellow).PrintfFunc()

// Warnf is Fprintf(os.Stderr) prefixed with [WARN]
func Warnf(format string, a ...interface{}) {
	warnColor("[WARN] ")
	fmt.Printf(format+"\n", a...)
}

var errorColor = color.New(color.FgRed).PrintfFunc()

// Errorf is Fprintf(os.Stderr) prefixed with [ERRO]
func Errorf(format string, a ...interface{}) {
	errorColor("[ERRO] ")
	fmt.Printf(format+"\n", a...)
}

var verboseColor = color.New(color.FgMagenta).PrintfFunc()

// Verbosef is Printf prefixed with [VERB]
func Verbosef(format string, a ...interface{}) {
	verboseColor("[VERB] ")
	fmt.Printf(format+"\n", a...)
}
