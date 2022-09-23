package forth

import (
	"bufio"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
)

type any interface{}

// Forth contains the forth environment
type Forth struct {
	stack       stack
	returnStack stack // not used for the return addresses but for specific features like >r r> and loop parameters
	heap        []any // ints or words

	currentProgramWord    int // points to the current word being executed
	currentProgramAddress int // the current word in the heap Forth executes

	expectInput bool
	input       []rune
	output      string
	debug       bool
}

//go:embed bootstrap.fth
var bootstrap string

// NewForth creates the Forth environment
func NewForth(_debug bool) *Forth {
	var f = new(Forth)
	f.heap = make([]any, 10000)
	f.expectInput = true
	f.debug = _debug

	f.Init()
	f.Reset()

	result, err := f.Exec(bootstrap)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(result)
	}

	return f
}

func (f *Forth) next() {
	f.currentProgramAddress++
}

// Reset starts Forth by calling the word QUIT
func (f *Forth) Reset() {
	here := f.heap[HEREp].(int)
	f.heap[here] = pWord(f.Find("QUIT"))
	f.IncHere(1)
	f.currentProgramAddress = here // QUIT should never return
}

func errorHandler(r interface{}) error {
	switch r.(type) {
	case string:
		return errors.New(r.(string))
	case error:
		return r.(error)
	default:
		return errors.New("unknown error type")
	}
}

// Exec compiles the given command and runs it
func (f *Forth) Exec(command string) (result string, err error) {
	command += "\n"
	defer func() {
		if r := recover(); r != nil {
			if f.debug {
				debug.PrintStack()
			}
			result = f.output
			err = errorHandler(r)
			f.Reset()
		}
	}()

	// main loop, just exit in case of the blocking KEY word
	f.output = ""
	f.input = []rune(command)
	f.expectInput = false
	for !f.expectInput {
		if f.debug {
			f.PrintExecuteState()
		}
		f.currentProgramWord = int(f.heap[f.currentProgramAddress].(pWord))
		// indirect threading
		f.heap[f.currentProgramWord+1].(func())()
	}

	return f.output, nil
}

func (f *Forth) CmdLine() {
	fmt.Println("Forth Command Line Interpreter")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		result, err := f.Exec(text)
		fmt.Println(result)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			fmt.Println(f.State())
		}
	}
}
