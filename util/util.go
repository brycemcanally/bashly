package util

import "errors"

// Stack is a stack data structure.
type Stack []interface{}

// Push pushes an item onto the stack.
func (stack *Stack) Push(item interface{}) {
	*stack = append(*stack, item)
}

// Pop pops an item off the stack.
func (stack *Stack) Pop() (interface{}, error) {
	if len(*stack) <= 0 {
		return "", errors.New("empty stack")
	}

	s := []interface{}(*stack)
	val := s[len(s)-1]
	s = s[:len(s)-1]
	*stack = s

	return val, nil
}

// Top returns the top item on the stack.
func (stack *Stack) Top() (interface{}, error) {
	if len(*stack) <= 0 {
		return "", errors.New("empty stack")
	}

	s := []interface{}(*stack)
	return s[len(s)-1], nil
}

// Top2 returns the top two items on the stack.
func (stack *Stack) Top2() (interface{}, interface{}, error) {
	if len(*stack) <= 1 {
		return nil, nil, errors.New("stack too small")
	}

	s := []interface{}(*stack)
	return s[len(s)-1], s[len(s)-2], nil
}
