package main

import (
	"forthly/forth"
)

func main() {
	forth.SetDebug(true)
	f := forth.NewForth()
	f.CmdLine()
}
