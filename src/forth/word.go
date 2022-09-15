package forth

// Representation of a word on the stack
type Word struct {
	link      int
	immediate bool
	hidden    bool
	name      string
}

type pWord int // a pointer to a Word on the heap array
