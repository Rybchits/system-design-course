package parser

type LexerStateStack struct {
	items []lexerState
}

func NewEmptyStack() *LexerStateStack {
	return &LexerStateStack{
		items: nil,
	}
}

func (stack *LexerStateStack) Push(item lexerState) {
	stack.items = append(stack.items, item)
}

func (stack *LexerStateStack) Pop() lexerState {
	if len(stack.items) == 0 {
		panic("Empty stack")
	}

	lastItem := stack.items[len(stack.items)-1]
	stack.items = stack.items[:len(stack.items)-1]

	return lastItem
}

func (stack *LexerStateStack) IsEmpty() bool {
	return len(stack.items) == 0
}

func (stack *LexerStateStack) CurrentState() lexerState {
	if stack.IsEmpty() {
		return startState
	}
	return stack.items[len(stack.items)-1]
}
