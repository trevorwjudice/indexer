package generic

import "fmt"

// NOT THREAD SAFE
type Stack[T any] struct {
	stack []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(item T) {
	s.stack = append(s.stack, item)
}

func (s *Stack[T]) Pop() (t T, err error) {
	stackLen := len(s.stack)
	if stackLen == 0 {
		return t, fmt.Errorf("cannot pop 0 length stack")
	}

	t = s.stack[stackLen-1]
	if stackLen == 1 {
		s.stack = nil
		return t, nil
	}

	s.stack = s.stack[:stackLen-1]

	return t, nil
}

func (s *Stack[T]) Peek() (t T, err error) {
	stackLen := len(s.stack)
	if stackLen == 0 {
		return t, fmt.Errorf("cannot pop 0 length stack")
	}

	return s.stack[stackLen-1], nil
}

func (s *Stack[T]) Len() int {
	return len(s.stack)
}
