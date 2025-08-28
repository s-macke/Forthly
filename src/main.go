//go:build !js && !wasm
package main

import (
	"forthly/forth"
)

func main() {
	f := forth.NewForth(false)
	f.CmdLine()
}
