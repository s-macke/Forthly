package forth

import (
	"fmt"
	"strconv"
)

type dictEntry struct {
	immediate bool
	f         func()
}

type dictionary map[string]*dictEntry // dictionary

func NewDictEntry(f func()) *dictEntry {
	return &dictEntry{immediate: false, f: f}
}

func (f *Forth) fillDictionary() {
	f.dict["+"] = NewDictEntry(func() { f.stack.Push(f.stack.Pop() + f.stack.Pop()) })
	f.dict["*"] = NewDictEntry(func() { f.stack.Push(f.stack.Pop() * f.stack.Pop()) })
	f.dict["-"] = NewDictEntry(func() { temp := f.stack.Pop(); f.stack.Push(f.stack.Pop() - temp) })
	f.dict["/"] = NewDictEntry(func() { temp := f.stack.Pop(); f.stack.Push(f.stack.Pop() / temp) })

	f.dict[">"] = NewDictEntry(func() { temp := f.stack.Pop(); f.booleanPush(f.stack.Pop() > temp) })
	f.dict["<"] = NewDictEntry(func() { temp := f.stack.Pop(); f.booleanPush(f.stack.Pop() < temp) })
	f.dict["="] = NewDictEntry(func() { f.booleanPush(f.stack.Pop() == f.stack.Pop()) })
	f.dict["<>"] = NewDictEntry(func() { f.booleanPush(f.stack.Pop() == f.stack.Pop()) })

	f.dict["."] = NewDictEntry(func() { f.output += fmt.Sprintf("%d ", f.stack.Pop()) })
	f.dict["cr"] = NewDictEntry(func() { f.output += fmt.Sprintln() })
	f.dict["drop"] = NewDictEntry(func() { f.stack.Pop() })
	f.dict["dup"] = NewDictEntry(func() { temp := f.stack.Pop(); f.stack.Push(temp); f.stack.Push(temp) })
	f.dict["swap"] = NewDictEntry(func() { temp, temp2 := f.stack.Pop(), f.stack.Pop(); f.stack.Push(temp); f.stack.Push(temp2) })
	f.dict["@"] = NewDictEntry(func() { f.stack.Push(f.heap[f.stack.Pop()].(int)) })
	f.dict["!"] = NewDictEntry(func() { address := f.stack.Pop(); f.heap[address] = f.stack.Pop() })
	f.dict["depth"] = NewDictEntry(func() { f.stack.Push(len(f.stack)) })

	f.dict[">r"] = NewDictEntry(func() { f.returnStack.Push(f.stack.Pop()) })
	f.dict["r>"] = NewDictEntry(func() { f.stack.Push(f.returnStack.Pop()) })
	f.dict["i"] = NewDictEntry(func() { f.stack.Push(f.returnStack.Get(0)) })
	f.dict["do"] = NewDictEntry(func() { f.dict["swap"].f(); f.dict[">r"].f(); f.dict[">r"].f() })

	f.dict["bye"] = NewDictEntry(func() { panic("Program stopped") })

	f.dict["emit"] = NewDictEntry(func() {
		f.output += fmt.Sprintf("%c", f.stack.Pop())
	})

	f.dict["words"] = NewDictEntry(func() {
		for k := range f.dict {
			f.output += fmt.Sprint(k, " ")
		}
		f.output += fmt.Sprintln("")
	})

	f.dict[".s"] = NewDictEntry(func() {
		f.output += fmt.Sprintf("<%d> ", len(f.stack))
		for _, e := range f.stack {
			f.output += fmt.Sprint(e, " ")
		}
	})
}

func (f *Forth) getWordFromDictionary(w string) *dictEntry {
	if f.dict[w] == nil {
		i, err := strconv.Atoi(w)
		if err != nil {
			panic("Word '" + w + "' is neither a valid word nor a number")
		}
		return NewDictEntry(func() {
			f.stack.Push(i)
		})
	}
	return f.dict[w]
}
