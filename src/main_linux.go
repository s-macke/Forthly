package main

import (
	"forthly/forth"
)

func main() {
	f := forth.NewForth()
	f.SetDebug(true)
	f.CmdLine()
}
