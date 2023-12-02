package main

import (
	"fmt"
	"forthly/forth"
	"runtime/debug"
	"strings"
	"syscall/js"
)

func execFunc(this js.Value, args []js.Value) (out any) {
	var result string
	var err error
	defer func() {
		// someting went really wrong. Error in error handling
		if r := recover(); r != nil {
			debug.PrintStack()
			out = result + "\n\nUnrecoverable error. Check console log"
		}
	}()

	command := args[0].String()
	command = strings.ReplaceAll(command, "\r\n", "\n")
	f := forth.NewForth(false)
	result, err = f.ExecLoops(command, 10000)
	if err != nil {
		result += "\n" + "Error: " + err.Error() + "\n"
		result += f.State()
		// not necessary to reset forth. Forth is resetted at reach restart of this function
	}
	out = result
	return
}

func main() {
	fmt.Println("Init Forth environment")
	js.Global().Set("ExecFunc", js.FuncOf(execFunc))
	// Prevent main from exiting
	select {}
}
