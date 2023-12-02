package forth

import "fmt"

// Representation of a word on the stack
type Word struct {
	link      pWord
	immediate bool
	hidden    bool
	name      string
}

type pWord int // a pointer to a Word on the heap array

func (pWord pWord) String() string {
	return fmt.Sprintf("0x%04x", uint(pWord))
}
