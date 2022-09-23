package forth

import (
	"fmt"
	"strconv"
	"strings"
)

const LATESTp = 0
const HEREp = 3

func (f *Forth) IncHere(inc int) {
	f.heap[HEREp] = f.heap[HEREp].(int) + inc
}

func (f *Forth) NewPrimitiveWord(name string, _f func()) {
	here := f.heap[HEREp].(int)
	latest := f.heap[LATESTp].(int)
	f.heap[here] = &Word{link: latest, immediate: false, hidden: false, name: name}
	f.heap[here+1] = _f
	f.heap[LATESTp] = here
	f.IncHere(2)
}

func (f *Forth) NewWord(name string, words []any) {
	f.NewPrimitiveWord(name, func() { // DOCOL
		f.returnStack.Push(f.currentProgramAddress)
		f.currentProgramAddress = f.currentProgramWord + 1
		f.next()
	})
	words = append(words, "EXIT")

	for _, _e := range words {
		here := f.heap[HEREp].(int)
		switch e := _e.(type) {
		case string:
			w := f.Find(e)
			if w == -1 {
				panic("word '" + e + "' not found")
			} else {
				f.heap[here] = pWord(w)
				f.IncHere(1)
			}
		case int, func():
			f.heap[here] = e
			f.IncHere(1)
		default:
			panic("type unknown")
		}
	}
}

func (f *Forth) NewImmediateWord(name string, words []any) {
	f.NewWord(name, words)
	latest := f.heap[LATESTp].(int)
	f.heap[latest].(*Word).immediate = true
}

func (f *Forth) NewIntVariable(name string, x int) {
	here := f.heap[HEREp].(int)
	f.NewPrimitiveWord(name, func() { f.stack.Push(here + 2); f.next() })
	f.heap[here+2] = x
	f.IncHere(1)
}

func (f *Forth) booleanPush(result bool) {
	if result {
		f.stack.Push(-1)
	} else {
		f.stack.Push(0)
	}
}

// ParseCAddr parses a so called c-addr string.
// Given an address on the stack,
//
//	the first byte is the size of the string,
//	followed by the string.
func (f *Forth) ParseCAddr() string {
	buffer := f.stack.Pop()
	length := f.heap[buffer].(int)
	buffer++
	str := ""
	for i := 0; i < length; i++ {
		str += string(rune(f.heap[buffer+i].(int)))
	}
	return str
}

/*
func (f *Forth) Find(word string) *Word {
	latest := f.heap[LATESTp].(int)
	for latest != 0 {
		w := f.heap[latest].(*Word)
		if w.name == word {
			return &w
		}
		latest = w.link
	}
	return nil
}
*/

func (f *Forth) Find(word string) int {
	if f.debug {
		fmt.Println("Debug: FIND " + word)
	}

	latest := f.heap[LATESTp].(int)
	for latest != 0 {
		w := f.heap[latest].(*Word)
		if strings.ToLower(w.name) == strings.ToLower(word) {
			return latest
		}
		latest = w.link
	}
	return -1 // illegal reference if used
}

func (f *Forth) NewParserWord() {
	// TODO: According to https://forth-standard.org/standard/core/WORD word takes the delimiter as argument

	f.NewWord("WORD", []any{"KEY", "CODE", "NOP", "KEY", "CODE", "NOP"})

	latest := f.heap[LATESTp].(int)
	wordbuffer := f.heap[HEREp].(int)
	index := 0
	debugstr := ""

	// store of 32 bytes
	for i := 0; i < 32; i++ {
		here := f.heap[HEREp].(int)
		f.heap[here] = 0
		f.IncHere(1)
	}

	f.heap[latest+4] = func() {
		key := f.stack.Pop()
		if key == ' ' || key == '\n' || key == '\t' {
			f.currentProgramAddress -= 2 // Repeat the key word. the program address still points to the CODE word.
		} else {
			index = 1 // c-addr Syntax string. start at index 1, because index 0 is used for the length.
			f.heap[wordbuffer+index] = key
			debugstr = string(rune(key))
			index++
			f.next()
		}
	}

	f.heap[latest+7] = func() {
		key := f.stack.Pop()
		if key != ' ' && key != '\n' && key != '\t' {
			f.heap[wordbuffer+index] = key
			debugstr += string(rune(key))
			index++
			f.currentProgramAddress -= 2 // Repeat the key word.
		} else {
			if f.debug {
				fmt.Println("Debug: Read word '" + debugstr + "'")
			}
			f.heap[wordbuffer] = index - 1 // store length in first byte
			f.stack.Push(wordbuffer)
			f.next()
		}
	}
}

func (f *Forth) NewInterpretWord() {
	f.NewWord("INTERPRET", []any{"[", "WORD", "FIND", "CODE", "NOP", "EXECUTE", "CODE", "NOP"})
	latest := f.heap[LATESTp].(int)

	f.heap[latest+6] = func() {
		STATE := f.heap[f.Find("STATE")+2].(int) // interpreter or compile mode
		status := f.stack.Pop()                  // retrieve results from the FIND word

		if STATE == 0 { // interpreter mode

			if status == 0 { // no word found
				w := f.ParseCAddr()
				number, err := strconv.Atoi(w)
				if err != nil {
					panic("Word '" + w + "' is neither a valid word nor a number")
				}
				f.stack.Push(number)
				f.currentProgramAddress -= 3 // Repeat the 'word' word.
			} else {
				f.next() // just execute the word on the stack
			}
		} else { // compile mode
			//panic("compile mode not implemented")
			here := f.heap[HEREp].(int)

			switch status {
			case 0: // no word found
				w := f.ParseCAddr()
				number, err := strconv.Atoi(w)
				if err != nil {
					panic("Word '" + w + "' is neither a valid word nor a number")
				}
				f.heap[here] = pWord(f.Find("LIT"))
				f.IncHere(1)
				f.heap[here+1] = number
				f.IncHere(1)
				f.currentProgramAddress -= 3 // Repeat the 'word' word.
			case 1: // IMMEDIATE word found
				f.next() // just execute the word on the stack
			case -1: // copy the word to the heap
				f.heap[here] = pWord(f.stack.Pop())
				f.IncHere(1)
				f.currentProgramAddress -= 3 // Repeat the 'word' word.
			default:
				panic("Unknown word type")
			}

		}
	}

	f.heap[latest+9] = func() {
		f.currentProgramAddress -= 6 // Repeat the 'word' word after execute
	}
}

func (f *Forth) setPrimitives() {
	f.NewPrimitiveWord("NOP", func() {})

	f.NewPrimitiveWord("@", func() { f.stack.Push(f.heap[f.stack.Pop()].(int)); f.next() })
	f.NewPrimitiveWord("!", func() { address := f.stack.Pop(); f.heap[address] = f.stack.Pop(); f.next() })
	f.NewPrimitiveWord("!w", func() { address := f.stack.Pop(); f.heap[address] = pWord(f.stack.Pop()); f.next() })

	f.NewPrimitiveWord("+", func() { f.stack.Push(f.stack.Pop() + f.stack.Pop()); f.next() })
	f.NewPrimitiveWord("*", func() { f.stack.Push(f.stack.Pop() * f.stack.Pop()); f.next() })
	f.NewPrimitiveWord("-", func() { temp := f.stack.Pop(); f.stack.Push(f.stack.Pop() - temp); f.next() })
	f.NewPrimitiveWord("/", func() { temp := f.stack.Pop(); f.stack.Push(f.stack.Pop() / temp); f.next() })
	f.NewPrimitiveWord("mod", func() { temp := f.stack.Pop(); f.stack.Push(f.stack.Pop() % temp); f.next() })

	f.NewPrimitiveWord("AND", func() { f.stack.Push(f.stack.Pop() & f.stack.Pop()); f.next() })
	f.NewPrimitiveWord("OR", func() { f.stack.Push(f.stack.Pop() | f.stack.Pop()); f.next() })
	f.NewPrimitiveWord("XOR", func() { f.stack.Push(f.stack.Pop() ^ f.stack.Pop()); f.next() })
	f.NewPrimitiveWord("INVERT", func() { f.stack.Push(^f.stack.Pop()); f.next() })

	f.NewPrimitiveWord("<", func() { temp := f.stack.Pop(); f.booleanPush(f.stack.Pop() < temp); f.next() })
	f.NewPrimitiveWord("=", func() { f.booleanPush(f.stack.Pop() == f.stack.Pop()); f.next() })

	f.NewPrimitiveWord(">R", func() { f.returnStack.Push(f.stack.Pop()); f.next() })
	f.NewPrimitiveWord("R>", func() { f.stack.Push(f.returnStack.Pop()); f.next() })
	f.NewPrimitiveWord("RPICK", func() { f.stack.Push(f.returnStack.Get(f.stack.Pop())); f.next() })
	f.NewPrimitiveWord("RDROP", func() { f.returnStack.Pop(); f.next() })

	f.NewPrimitiveWord("LIT", func() { f.stack.Push(f.heap[f.currentProgramAddress+1].(int)); f.next(); f.next() })
	f.NewPrimitiveWord("LITSTRING", func() {
		address := f.currentProgramAddress + 2
		f.stack.Push(address)
		length := f.heap[f.currentProgramAddress+1].(int)
		f.stack.Push(length)
		f.currentProgramAddress += length + 1
		f.next()
	})

	f.NewPrimitiveWord("BRANCH", func() {
		f.next()
		f.currentProgramAddress += f.heap[f.currentProgramAddress].(int)
	})

	f.NewPrimitiveWord("0BRANCH", func() { // Jump if stack is zero
		f.next()
		if f.stack.Pop() == 0 {
			f.currentProgramAddress += f.heap[f.currentProgramAddress].(int)
		} else {
			f.next()
		}
	})

	f.NewPrimitiveWord("EXIT", func() {
		f.currentProgramAddress = f.returnStack.Pop()
		f.next()
	})

	f.NewPrimitiveWord("IMMEDIATE", func() {
		latest := f.heap[LATESTp].(int)
		f.heap[latest].(*Word).immediate = true
		f.next()
	})
	f.heap[f.heap[LATESTp].(int)].(*Word).immediate = true

	f.NewPrimitiveWord("EXECUTE", func() {
		if f.debug {
			fmt.Printf("EXECUTE 0x%04x\n", f.stack.Get(0))
		}
		f.currentProgramWord = f.stack.Pop()
		f.heap[f.currentProgramWord+1].(func())()
	})

	f.NewIntVariable("STATE", 0)
	f.NewImmediateWord("[", []any{"LIT", 0, "STATE", "!"}) // Set STATE to 0, interpreter mode or immediate mode.
	f.NewWord("]", []any{"LIT", 1, "STATE", "!"})          // Set STATE to 1, compile mode

	f.NewIntVariable("BASE", 10)

	f.NewPrimitiveWord("BL", func() { f.stack.Push(0x20) })

	// execute arbitrary code in the next cell. The executed cell is responsible for the update of the program address.
	f.NewPrimitiveWord("CODE", func() { f.next(); f.heap[f.currentProgramAddress].(func())() })

	// Definition of word FIND: https://forth-standard.org/standard/core/FIND
	f.NewPrimitiveWord("FIND", func() {
		// save the c-addr in case of a miss
		buffer := f.stack.Get(0)

		word := f.ParseCAddr()
		wordp := f.Find(word)

		if wordp == -1 { // Word not found
			f.stack.Push(buffer)
			f.stack.Push(0)
		} else {
			f.stack.Push(wordp)
			if f.heap[wordp].(*Word).immediate {
				f.stack.Push(1)
			} else {
				f.stack.Push(-1)
			}
		}
		f.next()
	})

	// Because of the heap structure, simply do nothing here.
	// The word header has the same position on heap as the corresponding code
	f.NewPrimitiveWord(">CFA", func() {
		f.next()
	})

	f.NewPrimitiveWord(">DFA", func() {
		f.stack.Push(f.stack.Pop() + 1)
		f.next()
	})

	f.NewPrimitiveWord(".", func() { f.output += fmt.Sprintf("%d ", f.stack.Pop()); f.next() })

	f.NewPrimitiveWord("'", func() { // TICK word
		f.next()
		f.stack.Push(int(f.heap[f.currentProgramAddress].(pWord)))
		f.next()
	})

	f.NewPrimitiveWord("DROP", func() { f.stack.Pop(); f.next() })
	f.NewPrimitiveWord("SWAP", func() { temp, temp2 := f.stack.Pop(), f.stack.Pop(); f.stack.Push(temp); f.stack.Push(temp2); f.next() })
	f.NewPrimitiveWord("DUP", func() { temp := f.stack.Pop(); f.stack.Push(temp); f.stack.Push(temp); f.next() })

	f.NewPrimitiveWord(".S", func() {
		f.output += fmt.Sprintf("<%d> ", f.stack.Size())
		for _, e := range f.stack {
			f.output += fmt.Sprint(e, " ")
		}
		f.next()
	})

	f.NewPrimitiveWord("DUMP", func() {
		f.output += f.HeapDump()
		f.output += f.State()
		f.next()
	})

	f.NewPrimitiveWord("WORDS", func() {
		f.output += f.Words()
		f.next()
	})

	// The only blocking command here. Hence, have to be treated specially in the main loop
	f.NewPrimitiveWord("KEY", func() {
		if len(f.input) == 0 {
			f.expectInput = true
			return
		}
		f.stack.Push(int(f.input[0]))
		f.input = f.input[1:]
		f.next()
	})

	f.NewPrimitiveWord("EMIT", func() {
		char := f.stack.Pop()
		if f.debug {
			fmt.Printf("Debug: EMIT %c\n", char)
		}
		f.output += fmt.Sprintf("%c", char)
		f.next()
	})
}

func (f *Forth) InitHeao() {

	// Set two of the most important variables, LATEST and HERE
	f.heap[LATESTp] = 4 // LATEST: Points to the latest (most recently defined) word in the dictionary.
	f.heap[1] = &Word{link: 0, immediate: false, hidden: false, name: "&LATEST"}
	f.heap[2] = func() { f.stack.Push(LATESTp); f.next() }

	f.heap[HEREp] = 6 // HERE: Points to the next free byte of memory.  When compiling, compiled words go here.
	f.heap[4] = &Word{link: 1, immediate: false, hidden: false, name: "&HERE"}
	f.heap[5] = func() { f.stack.Push(HEREp); f.next() }

	// Now we can use the easier NewPrimitiveWord and NewIntVariable
}

func (f *Forth) Init() {
	f.InitHeao()

	f.setPrimitives()

	f.NewParserWord()

	f.NewWord("CREATE", []any{"WORD", "CODE", func() {
		name := f.ParseCAddr()
		here := f.heap[HEREp].(int)
		// create new word with initial behavior
		f.NewPrimitiveWord(name, func() {
			f.stack.Push(here + 2)
			f.next()
		})
		f.next()
	}})

	// this is a crazy function
	f.NewPrimitiveWord("DOES>", func() {
		latest := f.heap[LATESTp].(int)
		currentAddress := f.currentProgramAddress // the address of DOES>
		f.heap[latest+1] = func() {
			// push current address to the return stack, so that the EXIT statement can return to the correct address
			f.returnStack.Push(f.currentProgramAddress)
			// first replicate the CREATE functionality
			f.stack.Push(latest + 2)
			// then jump to the first word after the >DOES word
			f.currentProgramAddress = currentAddress + 1
		}
		// exit the function here, everything after >DOES is ignored
		f.currentProgramAddress = f.returnStack.Pop()
		f.next()
	})

	// TODO HIDE
	f.NewWord(":", []any{"CREATE", "CODE", func() {
		// overwrite default behavior of CREATE with DOCOL
		latest := f.heap[LATESTp].(int)
		f.heap[latest+1] = func() {
			// DOCOL
			f.returnStack.Push(f.currentProgramAddress)
			f.currentProgramAddress = f.currentProgramWord + 1
			f.next()
		}
		f.next()
	}, "]"})

	// define commentary
	f.NewImmediateWord("\\", []any{"KEY", "CODE", func() {
		key := f.stack.Pop()
		if key == 0x0a { // if is new line
			f.next()
			return
		}
		f.currentProgramAddress -= 2 // back to KEY word
	}})

	f.NewImmediateWord(";", []any{"CODE", func() {
		// set EXIT at the end of the function
		here := f.heap[HEREp].(int)
		f.heap[here] = pWord(f.Find("EXIT"))
		f.IncHere(1)
		f.next()
	}, "["})

	f.NewInterpretWord()

	// QUIT is the first word executed
	// Reset return stack and parameter stack
	// call Interpret
	// endless loop. e.g. BRANCH with 0 for example
	f.NewWord("QUIT", []any{"CODE", func() { f.stack.Clear(); f.returnStack.Clear(); f.next() }, "INTERPRET", "BRANCH", 0})
}
