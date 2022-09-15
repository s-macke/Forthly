package main

import (
	"fmt"
	"forthly/forth"
	"syscall/js"
)

func execFunc(this js.Value, args []js.Value) interface{} {
	command := args[0].String()
	forth.SetDebug(true)
	f := forth.NewForth()
	result, err := f.Exec(command)
	if err != nil {
		return "Error: " + err.Error()
	} else {
		return result
	}
}

func main() {
	fmt.Println("Init Forth environment")
	js.Global().Set("ExecFunc", js.FuncOf(execFunc))
	// Prevent main from exiting
	select {}
}
