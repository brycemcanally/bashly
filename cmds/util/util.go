package util

import "errors"

// Stack is a stack data structure for strings.
type Stack []string

// Push pushes an element onto the stack.
func (stack *Stack) Push(elem string) {
	*stack = append(*stack, elem)
}

// Pop pops an element off the stack.
func (stack *Stack) Pop() (string, error) {
	if len(*stack) <= 0 {
		return "", errors.New("empty stack")
	}

	s := []string(*stack)
	val := s[len(s)-1]
	s = s[:len(s)-1]
	*stack = s

	return val, nil
}

// Top returns the top element on the stack.
func (stack *Stack) Top() (string, error) {
	if len(*stack) <= 0 {
		return "", errors.New("empty stack")
	}

	s := []string(*stack)
	return s[len(s)-1], nil
}

// Top2 returns the top two elements on the stack.
func (stack *Stack) Top2() ([]string, error) {
	if len(*stack) <= 1 {
		return nil, errors.New("stack too small")
	}

	s := []string(*stack)
	return s[len(s)-2:], nil
}
