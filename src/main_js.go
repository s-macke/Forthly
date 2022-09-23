package main

import (
	"fmt"
	"forthly/forth"
	"syscall/js"
)

func execFunc(this js.Value, args []js.Value) interface{} {
	command := args[0].String()
	f := forth.NewForth(false)
	result, err := f.Exec(command)
	if err != nil {
		return result + "\n" + "Error: " + err.Error() + "\n" + f.State()
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
