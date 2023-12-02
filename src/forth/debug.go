package forth

import (
	"fmt"
	"reflect"
	"strings"
)

func (f *Forth) SetDebug(value bool) {
	f.debug = value
}

func (f *Forth) PrintExecuteState() {
	fmt.Printf("Debug: 0x%04x: ", uint(f.currentProgramAddress))
	word, ok := f.heap[f.currentProgramAddress].(pWord)
	if !ok {
		fmt.Printf(" Invalid word of type %s and content %v\n", reflect.TypeOf(f.heap[f.currentProgramAddress]), f.heap[f.currentProgramAddress])
		panic("Invalid word")
	}
	for i := 0; i < f.returnStack.Size()*4; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("| Exec '%s'\n", f.heap[word].(*Word).name)
}

func (f *Forth) HeapDump() string {
	var sb strings.Builder

	// just estimate the size of the heap dump
	here := (f.heap[HEREp]).(int)
	sb.Grow(32 * here)

	sb.WriteString("\n========= HEAP DUMP START =========\n")
	for i := 0; i < len(f.heap); i++ {
		if f.heap[i] == nil {
			//sb.WriteString(fmt.Sprintf("0x%04x: nil\n", i))
			continue
		}
		sb.WriteString(fmt.Sprintf("0x%04x: %15s", uint(i), reflect.TypeOf(f.heap[i])))

		switch v := f.heap[i].(type) {
		case *Word:
			sb.WriteString(fmt.Sprintf(" '%s'  link: 0x%04x   imm: %v", v.name, uint(v.link), v.immediate))
		case int:
			sb.WriteString(fmt.Sprintf(" 0x%x", v))
		case pWord:
			if v == -1 {
				sb.WriteString(fmt.Sprintf(" -1"))
			} else {
				sb.WriteString(fmt.Sprintf(" 0x%04x %s", uint(v), f.heap[v].(*Word).name))
			}
		case func():
		default:
			sb.WriteString(fmt.Sprintf("Unknown type %v", v))
		}
		sb.WriteString(fmt.Sprintf("\n"))
	}
	sb.WriteString(fmt.Sprintf("========== HEAP DUMP END ==========\n"))
	return sb.String()
}

func (f *Forth) GetCallStack() string {
	var sb strings.Builder
	if f.currentProgramWord > 0 && f.currentProgramWord < pWord(len(f.heap)) {
		if w, ok := f.heap[f.currentProgramWord].(*Word); ok {
			sb.WriteString(fmt.Sprintf("    0x%04x: '%s'\n", uint(f.currentProgramWord), w.name))
		}
	}
	for i := len((*f).returnStack) - 1; i >= 0; i-- {
		if wp, ok1 := ((*f).returnStack[i]).(pWord); ok1 {
			sb.WriteString(fmt.Sprintf("    0x%04x: ", uint(wp)))
			if wp >= 0 && wp < pWord(len(f.heap)) {
				if w, ok2 := f.heap[uint(wp)].(*Word); ok2 {
					sb.WriteString(w.name)
				} else {
					sb.WriteString(" Unknown word")
				}
			} else {
				sb.WriteString(" out of heap range")
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// State returns a string representation of the current state of the Forth VM
// Be very careful, as this can be executed in an recover cycle
func (f *Forth) State() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Program address 0x%04x\n", f.currentProgramAddress))
	sb.WriteString("Stack: " + f.stack.ToString() + "\n")
	sb.WriteString("Return Stack: " + f.returnStack.ToString() + "\n")
	sb.WriteString("Stack Trace: \n" + f.GetCallStack() + "\n")
	/*
		slice := 50
		if len(f.input)-2 < slice {
			slice = len(f.input) - 2
		}
		sb.WriteString(fmt.Sprintf("input string:\n '%s'\n", string(f.input[:slice])))
	*/
	return sb.String()
}

func (f *Forth) Words() string {
	var sb strings.Builder
	latest := f.heap[LATESTp].(pWord)
	for latest != 0 {
		word := f.heap[latest].(*Word)
		sb.WriteString(word.name)
		latest = word.link
		if latest != 0 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("\n")
	return sb.String()
}
