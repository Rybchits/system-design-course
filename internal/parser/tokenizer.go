package parser

import (
	"bufio"
	"fmt"
	"io"
	envsholder "shell/internal/envs_holder"
)

type TokenType int

type runeTokenClass int

type lexerState int

type Token struct {
	TokenType TokenType
	Value     string
}

func (a *Token) Equal(b *Token) bool {
	if a == nil || b == nil {
		return false
	}
	if a.TokenType != b.TokenType {
		return false
	}
	return a.Value == b.Value
}

const (
	spaceRunes            = " \t\r"
	escapingQuoteRunes    = `"`
	nonEscapingQuoteRunes = "'"
	escapeRunes           = `\`
	commentRunes          = "#"
	endLineRunes          = "\n"
	pipeRunes             = "|"
	envVarRunes           = "$"
)

const (
	unknownRuneClass runeTokenClass = iota
	spaceRuneClass
	escapingQuoteRuneClass
	nonEscapingQuoteRuneClass
	escapeRuneClass
	commentRuneClass
	endLineRuneClass
	pipeRuneClass
	eofRuneClass
	envVarClass
)

const (
	UnknownToken TokenType = iota
	WordToken
	CommentToken
	EndLineToken
	PipeToken
)

const (
	startState              lexerState = iota // еще не было символов
	inWordState                               // в процессе определения слова
	escapingState                             // экранирование, следующий символ должен быть литеральным
	escapingQuotedState                       // экранирование в заключенной в кавычки строке
	quotingEscapingState                      // внутри заключенной в кавычки строки, которая поддерживает экранирование
	quotingState                              // внутри строки, которая не поддерживает экранирование
	commentState                              // в пределах комментария
	pipeSymbolState                           // прошлый символ был pipe
	endLineState                              // прошлый символ был \n
	enviromentVariableState                   // внутри имени переменной окружения
)

type tokenClassifier map[rune]runeTokenClass

func (typeMap tokenClassifier) addRuneClass(runes string, tokenType runeTokenClass) {
	for _, runeChar := range runes {
		typeMap[runeChar] = tokenType
	}
}

func newDefaultClassifier() tokenClassifier {
	t := tokenClassifier{}
	t.addRuneClass(spaceRunes, spaceRuneClass)
	t.addRuneClass(escapingQuoteRunes, escapingQuoteRuneClass)
	t.addRuneClass(nonEscapingQuoteRunes, nonEscapingQuoteRuneClass)
	t.addRuneClass(escapeRunes, escapeRuneClass)
	t.addRuneClass(commentRunes, commentRuneClass)
	t.addRuneClass(endLineRunes, endLineRuneClass)
	t.addRuneClass(pipeRunes, pipeRuneClass)
	t.addRuneClass(envVarRunes, envVarClass)
	return t
}

func (t tokenClassifier) ClassifyRune(runeVal rune) runeTokenClass {
	return t[runeVal]
}

type Tokenizer struct {
	input             bufio.Reader
	classifier        tokenClassifier
	statesStack       lexerStateStack
	envsHolder        *envsholder.Env
	isEnded           bool
	currentTokenState *getTokenState
}

type getTokenState struct {
	tokenType    TokenType
	nextRuneType runeTokenClass
	nextRune     rune
	value        []rune
	envVarBuffer []rune
	err          error
}

func NewTokenizer(r io.Reader, vars *envsholder.Env) *Tokenizer {
	input := bufio.NewReader(r)

	return &Tokenizer{
		input:       *input,
		classifier:  newDefaultClassifier(),
		statesStack: *NewEmptyStack(),
		envsHolder:  vars,
		isEnded:     false,
	}
}

func (t *Tokenizer) handleInWordState() bool {
	nextRuneType := t.currentTokenState.nextRuneType
	value := &t.currentTokenState.value
	nextRune := t.currentTokenState.nextRune

	switch nextRuneType {
	case eofRuneClass:
		{
			t.isEnded = true
			return true
		}
	case spaceRuneClass:
		{
			t.statesStack.Pop()
			return true
		}
	case escapingQuoteRuneClass:
		{
			t.statesStack.Push(quotingEscapingState)
		}
	case nonEscapingQuoteRuneClass:
		{
			t.statesStack.Push(quotingState)
		}
	case escapeRuneClass:
		{
			t.statesStack.Push(escapingState)
		}
	case envVarClass:
		{
			t.statesStack.Push(enviromentVariableState)
		}
	case endLineRuneClass:
		{
			t.statesStack.Pop()
			t.statesStack.Push(endLineState)
			return true
		}
	case pipeRuneClass:
		{
			t.statesStack.Pop()
			t.statesStack.Push(pipeSymbolState)
			return true
		}
	default:
		{
			*value = append(*value, nextRune)
		}
	}
	return false
}

func (t *Tokenizer) handleEscapingState() bool {
	nextRuneType := t.currentTokenState.nextRuneType
	value := &t.currentTokenState.value
	nextRune := t.currentTokenState.nextRune

	switch nextRuneType {
	case eofRuneClass:
		{
			t.isEnded = true
			return true
		}
	default:
		{
			t.statesStack.Pop()
			*value = append(*value, nextRune)
		}
	}
	return false
}

func (t *Tokenizer) handleEscapingQuotedState() bool {
	nextRuneType := t.currentTokenState.nextRuneType
	value := &t.currentTokenState.value
	nextRune := t.currentTokenState.nextRune

	switch nextRuneType {
	case eofRuneClass:
		{
			t.isEnded = true
			return true
		}
	default:
		{
			t.statesStack.Pop()
			*value = append(*value, nextRune)
		}
	}
	return false
}

func (t *Tokenizer) handleQuotingEscapingState() bool {
	nextRuneType := t.currentTokenState.nextRuneType
	value := &t.currentTokenState.value
	nextRune := t.currentTokenState.nextRune

	switch nextRuneType {
	case eofRuneClass:
		{
			t.isEnded = true
			return true
		}
	case escapingQuoteRuneClass:
		{
			t.statesStack.Pop()
		}
	case escapeRuneClass:
		{
			t.statesStack.Push(escapingQuotedState)
		}
	case envVarClass:
		{
			t.statesStack.Push(enviromentVariableState)
		}
	default:
		{
			*value = append(*value, nextRune)
		}
	}
	return false
}

func (t *Tokenizer) handleQuotingState() bool {
	nextRuneType := t.currentTokenState.nextRuneType
	value := &t.currentTokenState.value
	nextRune := t.currentTokenState.nextRune

	switch nextRuneType {
	case eofRuneClass:
		{
			t.isEnded = true
			return true
		}
	case nonEscapingQuoteRuneClass:
		{
			t.statesStack.Pop()
		}
	default:
		{
			*value = append(*value, nextRune)
		}
	}
	return false
}

func (t *Tokenizer) handleCommentState() bool {
	nextRuneType := t.currentTokenState.nextRuneType
	value := &t.currentTokenState.value
	nextRune := t.currentTokenState.nextRune

	switch nextRuneType {
	case eofRuneClass:
		{
			t.isEnded = true
			return true
		}
	case endLineRuneClass:
		{
			t.statesStack.Pop()
			return true
		}
	case spaceRuneClass:
		{
			*value = append(*value, nextRune)
		}
	default:
		{
			*value = append(*value, nextRune)
		}
	}
	return false
}

func (t *Tokenizer) handleStartState() {
	tokenType := &t.currentTokenState.tokenType
	nextRuneType := t.currentTokenState.nextRuneType
	value := &t.currentTokenState.value
	nextRune := t.currentTokenState.nextRune

	switch nextRuneType {
	case eofRuneClass:
		{
			t.isEnded = true
		}
	case spaceRuneClass:
		{
		}
	case escapingQuoteRuneClass:
		{
			*tokenType = WordToken
			t.statesStack.Push(inWordState)
			t.statesStack.Push(quotingEscapingState)
		}
	case nonEscapingQuoteRuneClass:
		{
			*tokenType = WordToken
			t.statesStack.Push(inWordState)
			t.statesStack.Push(quotingState)
		}
	case escapeRuneClass:
		{
			*tokenType = WordToken
			t.statesStack.Push(inWordState)
			t.statesStack.Push(escapingState)
		}
	case envVarClass:
		{
			*tokenType = WordToken
			t.statesStack.Push(inWordState)
			t.statesStack.Push(enviromentVariableState)
		}
	case commentRuneClass:
		{
			*tokenType = CommentToken
			t.statesStack.Push(commentState)
		}
	case endLineRuneClass:
		{
			t.statesStack.Push(endLineState)
		}
	case pipeRuneClass:
		{
			t.statesStack.Push(pipeSymbolState)
		}
	default:
		{
			*tokenType = WordToken
			t.statesStack.Push(inWordState)
			*value = append(*value, nextRune)
		}
	}
}

func (t *Tokenizer) handleEnviromentVariableState() (*Token, error) {
	nextRuneType := t.currentTokenState.nextRuneType
	value := &t.currentTokenState.value
	envVarBuffer := &t.currentTokenState.envVarBuffer
	nextRune := t.currentTokenState.nextRune

	if nextRuneType != unknownRuneClass {
		if env, ok := t.envsHolder.Vars[string(*envVarBuffer)]; ok {
			*value = append(*value, []rune(env)...)
		}
		*envVarBuffer = []rune{}
		t.statesStack.Pop()
		return t.handleRune()
	} else {
		*envVarBuffer = append(*envVarBuffer, nextRune)
	}
	return nil, nil
}

func (t *Tokenizer) handleRune() (*Token, error) {
	tokenType := &t.currentTokenState.tokenType
	value := &t.currentTokenState.value

	switch t.statesStack.CurrentState() {
	case startState:
		{
			t.handleStartState()
			if t.isEnded {
				return nil, io.EOF
			}
		}
	case inWordState:
		{
			if t.handleInWordState() {
				var token *Token
				if len(*value) != 0 {
					token = &Token{
						TokenType: *tokenType,
						Value:     string(*value)}
				} else {
					token = nil
				}

				var err error = nil
				if t.isEnded {
					err = io.EOF
				}
				return token, err
			}
		}
	case escapingState:
		{
			if t.handleEscapingState() {
				token := &Token{
					TokenType: *tokenType,
					Value:     string(*value)}
				return token, fmt.Errorf("EOF found after escape character")
			}
		}
	case escapingQuotedState:
		{
			if t.handleEscapingQuotedState() {
				token := &Token{
					TokenType: *tokenType,
					Value:     string(*value)}
				return token, fmt.Errorf("EOF found after escape character")
			}
		}
	case quotingEscapingState:
		{
			if t.handleQuotingEscapingState() {
				token := &Token{
					TokenType: *tokenType,
					Value:     string(*value)}
				return token, fmt.Errorf("EOF found when expecting closing quote")
			}
		}
	case quotingState:
		{
			if t.handleQuotingState() {
				token := &Token{
					TokenType: *tokenType,
					Value:     string(*value)}
				return token, fmt.Errorf("EOF found when expecting closing quote")
			}
		}
	case commentState:
		{
			if !t.handleCommentState() {
				return nil, nil
			}

			if t.isEnded {
				token := &Token{
					TokenType: *tokenType,
					Value:     string(*value)}
				return token, io.EOF
			} else {
				return &Token{TokenType: EndLineToken, Value: endLineRunes}, nil
			}
		}
	case enviromentVariableState:
		{
			return t.handleEnviromentVariableState()
		}
	}
	return nil, nil
}

func (t *Tokenizer) scanStream() (*Token, error) {
	t.currentTokenState = &getTokenState{}

	if t.isEnded {
		return nil, io.EOF
	}

	for {
		state := t.statesStack.CurrentState()

		// Токен может быть получен на прошлой итерации, если так отдаем его
		if state == pipeSymbolState {
			t.statesStack.Pop()
			return &Token{TokenType: PipeToken, Value: pipeRunes}, nil

		} else if state == endLineState {
			t.statesStack.Pop()
			return &Token{TokenType: EndLineToken, Value: endLineRunes}, nil
		}

		// Читаем следующий символ и классифицируем его
		t.currentTokenState.nextRune, _, t.currentTokenState.err = t.input.ReadRune()
		t.currentTokenState.nextRuneType = t.classifier.ClassifyRune(t.currentTokenState.nextRune)

		// Если произошла ошибка при чтении, вернуть ошибку
		if t.currentTokenState.err == io.EOF {
			t.currentTokenState.nextRuneType = eofRuneClass
			t.currentTokenState.err = nil

		} else if t.currentTokenState.err != nil {
			return nil, t.currentTokenState.err
		}

		// Обработать текущий символ в контексте текущего состояни
		token, err := t.handleRune()

		if token != nil || err != nil {
			t.currentTokenState = nil
			return token, err
		}
	}
}

func (t *Tokenizer) Next() (*Token, error) {
	return t.scanStream()
}
