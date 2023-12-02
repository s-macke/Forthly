package forth

import (
	"fmt"
	"reflect"
	"strings"
)

type stack []any

// IsEmpty checks if stack is empty
func (s *stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *stack) Size() int {
	return len(*s)
}

// Push a new value onto the stack
func (s *stack) Push(str any) {
	if s.Size() > 1024 {
		panic("stack overflow")
	}

	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// Pop removes and return top element of stack.
func (s *stack) Pop() any {
	if s.IsEmpty() {
		panic("stack is empty")
	}
	index := len(*s) - 1   // Get the index of the top most element.
	element := (*s)[index] // Index into the slice and obtain the element.
	*s = (*s)[:index]      // Remove it from the stack by slicing it off.
	return element
}

// Get returns the n-th element in the Stack
func (s *stack) Get(idx int) any {
	return (*s)[len(*s)-idx-1]
}

func (s *stack) Clear() {
	*s = make([]any, 0)
}

func (s *stack) ToString() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i := 0; i < len(*s); i++ {
		if (*s)[i] == nil {
			sb.WriteString("nil")
		} else {
			sb.WriteString(fmt.Sprint((*s)[i]) + " (" + reflect.TypeOf((*s)[i]).Name() + ")")
		}
		if i < len(*s)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}
