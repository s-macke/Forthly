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
	fmt.Printf("Debug: 0x%04x: ", f.currentProgramAddress)
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
			break
		}
		sb.WriteString(fmt.Sprintf("0x%04x: %15s", i, reflect.TypeOf(f.heap[i])))

		switch v := f.heap[i].(type) {
		case *Word:
			sb.WriteString(fmt.Sprintf(" '%s'  link: 0x%04x   immediate: %v", v.name, v.link, v.immediate))
		case int:
			sb.WriteString(fmt.Sprintf(" 0x%x", v))
		case pWord:
			sb.WriteString(fmt.Sprintf(" 0x%04x %s", v, f.heap[v].(*Word).name))
		case func():
		default:
			sb.WriteString(fmt.Sprintf("Unknown type %v", v))

		}
		sb.WriteString(fmt.Sprintf("\n"))
	}
	sb.WriteString(fmt.Sprintf("========== HEAP DUMP END ==========\n"))
	return sb.String()
}

func (f *Forth) State() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("program address at 0x%04x\n", f.currentProgramAddress))
	sb.WriteString(fmt.Sprintf("Stack: %v\n", f.stack))
	sb.WriteString(fmt.Sprintf("Return Stack: %v\n", f.returnStack))
	return sb.String()
}

func (f *Forth) Words() string {
	var sb strings.Builder
	latest := f.heap[LATESTp].(int)
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
