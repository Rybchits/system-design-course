package parser

type lexerStateStack struct {
	items []lexerState
}

func NewEmptyStack() *lexerStateStack {
	return &lexerStateStack{
		items: nil,
	}
}

func (stack *lexerStateStack) Push(item lexerState) {
	stack.items = append(stack.items, item)
}

func (stack *lexerStateStack) Pop() lexerState {
	if len(stack.items) == 0 {
		panic("Empty stack")
	}

	lastItem := stack.items[len(stack.items)-1]
	stack.items = stack.items[:len(stack.items)-1]

	return lastItem
}

func (stack *lexerStateStack) IsEmpty() bool {
	return len(stack.items) == 0
}

func (stack *lexerStateStack) CurrentState() lexerState {
	if stack.IsEmpty() {
		return startState
	}
	return stack.items[len(stack.items)-1]
}
