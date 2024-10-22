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
	input       bufio.Reader
	classifier  tokenClassifier
	statesStack LexerStateStack
	envsHolder  *envsholder.Env
	isEnded     bool
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

func (t *Tokenizer) handleInWordState(
	nextRuneType runeTokenClass,
	value *[]rune,
	nextRune rune,
) bool {
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

func (t *Tokenizer) handleEscapingState(
	nextRuneType runeTokenClass,
	value *[]rune,
	nextRune rune,
) bool {
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

func (t *Tokenizer) handleEscapingQuotedState(
	nextRuneType runeTokenClass,
	value *[]rune,
	nextRune rune,
) bool {
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

func (t *Tokenizer) handleQuotingEscapingState(
	nextRuneType runeTokenClass,
	value *[]rune,
	nextRune rune,
) bool {
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

func (t *Tokenizer) handleQuotingState(
	nextRuneType runeTokenClass,
	value *[]rune,
	nextRune rune,
) bool {
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

func (t *Tokenizer) handleCommentState(
	nextRuneType runeTokenClass,
	value *[]rune,
	nextRune rune,
) bool {
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

func (t *Tokenizer) handleStartState(
	tokenType *TokenType,
	nextRuneType runeTokenClass,
	value *[]rune,
	nextRune rune,
) {
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

func (t *Tokenizer) handleEnviromentVariableState(
	tokenType *TokenType,
	nextRuneType runeTokenClass,
	value *[]rune,
	envVarBuffer *[]rune,
	nextRune rune,
) (*Token, error) {
	if nextRuneType != unknownRuneClass {
		if env, ok := t.envsHolder.Vars[string(*envVarBuffer)]; ok {
			*value = append(*value, []rune(env)...)
		}
		*envVarBuffer = []rune{}
		t.statesStack.Pop()
		return t.handleRune(tokenType, nextRuneType, value, envVarBuffer, nextRune)
	} else {
		*envVarBuffer = append(*envVarBuffer, nextRune)
	}
	return nil, nil
}

func (t *Tokenizer) handleRune(
	tokenType *TokenType,
	nextRuneType runeTokenClass,
	value *[]rune,
	envVarBuffer *[]rune,
	nextRune rune,
) (*Token, error) {
	switch t.statesStack.CurrentState() {
	case startState:
		{
			t.handleStartState(tokenType, nextRuneType, value, nextRune)
			if t.isEnded {
				return nil, io.EOF
			}
		}
	case inWordState:
		{
			if t.handleInWordState(nextRuneType, value, nextRune) {
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
			if t.handleEscapingState(nextRuneType, value, nextRune) {
				token := &Token{
					TokenType: *tokenType,
					Value:     string(*value)}
				return token, fmt.Errorf("EOF found after escape character")
			}
		}
	case escapingQuotedState:
		{
			if t.handleEscapingQuotedState(nextRuneType, value, nextRune) {
				token := &Token{
					TokenType: *tokenType,
					Value:     string(*value)}
				return token, fmt.Errorf("EOF found after escape character")
			}
		}
	case quotingEscapingState:
		{
			if t.handleQuotingEscapingState(nextRuneType, value, nextRune) {
				token := &Token{
					TokenType: *tokenType,
					Value:     string(*value)}
				return token, fmt.Errorf("EOF found when expecting closing quote")
			}
		}
	case quotingState:
		{
			if t.handleQuotingState(nextRuneType, value, nextRune) {
				token := &Token{
					TokenType: *tokenType,
					Value:     string(*value)}
				return token, fmt.Errorf("EOF found when expecting closing quote")
			}
		}
	case commentState:
		{
			if !t.handleCommentState(nextRuneType, value, nextRune) {
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
			return t.handleEnviromentVariableState(tokenType, nextRuneType, value, envVarBuffer, nextRune)
		}
	}

	return nil, nil
}

func (t *Tokenizer) scanStream() (*Token, error) {
	var tokenType TokenType
	var value []rune
	var nextRune rune
	var envVarBuffer []rune
	var nextRuneType runeTokenClass
	var err error

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
		nextRune, _, err = t.input.ReadRune()
		nextRuneType = t.classifier.ClassifyRune(nextRune)

		// Если произошла ошибка при чтении, вернуть ошибку
		if err == io.EOF {
			nextRuneType = eofRuneClass
			err = nil

		} else if err != nil {
			return nil, err
		}

		// Обработать текущий символ в контексте текущего состояни
		token, err := t.handleRune(&tokenType, nextRuneType, &value, &envVarBuffer, nextRune)

		if token != nil || err != nil {
			return token, err
		}
	}
}

func (t *Tokenizer) Next() (*Token, error) {
	return t.scanStream()
}
