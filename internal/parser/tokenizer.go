package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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
)

const (
	UnknownToken TokenType = iota
	WordToken
	SpaceToken
	CommentToken
	EndLineToken
	PipeToken
)

const (
	startState           lexerState = iota // еще не было символов
	inWordState                            // в процессе определения слова
	escapingState                          // экранирование, следующий символ должен быть литеральным
	escapingQuotedState                    // экранирование в заключенной в кавычки строке
	quotingEscapingState                   // внутри заключенной в кавычки строки, которая поддерживает экранирование
	quotingState                           // внутри строки, которая не поддерживает экранирование
	commentState                           // в пределах комментария
	pipeSymbolState                        // прошлый символ был pipe
	endLineState                           // прошлый символ был \n
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
	return t
}

func (t tokenClassifier) ClassifyRune(runeVal rune) runeTokenClass {
	return t[runeVal]
}

type Tokenizer struct {
	input       bufio.Reader
	classifier  tokenClassifier
	statesStack LexerStateStack
	isEnded     bool
}

func NewTokenizer(r io.Reader) *Tokenizer {
	input := bufio.NewReader(r)
	classifier := newDefaultClassifier()

	return &Tokenizer{
		input:       *input,
		classifier:  classifier,
		statesStack: *NewEmptyStack(),
		isEnded:     false,
	}
}

func (t *Tokenizer) scanStream() (*Token, error) {
	var tokenType TokenType
	var value []rune
	var nextRune rune
	var nextRuneType runeTokenClass
	var err error

	if t.isEnded {
		return nil, io.EOF
	}

	var state lexerState

	for {
		if t.statesStack.IsEmpty() {
			state = startState
		} else {
			state = t.statesStack.Top()
		}

		if state == pipeSymbolState {
			t.statesStack.Pop()
			return &Token{TokenType: PipeToken, Value: pipeRunes}, nil

		} else if state == endLineState {
			t.statesStack.Pop()
			return &Token{TokenType: EndLineToken, Value: endLineRunes}, nil
		}

		nextRune, _, err = t.input.ReadRune()
		nextRuneType = t.classifier.ClassifyRune(nextRune)

		if err == io.EOF {
			nextRuneType = eofRuneClass
			err = nil

		} else if err != nil {
			return nil, err
		}

		switch state {
		case startState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						t.isEnded = true
						return nil, io.EOF
					}
				case spaceRuneClass:
					{
					}
				case escapingQuoteRuneClass:
					{
						tokenType = WordToken
						t.statesStack.Push(inWordState)
						t.statesStack.Push(quotingEscapingState)
					}
				case nonEscapingQuoteRuneClass:
					{
						tokenType = WordToken
						t.statesStack.Push(inWordState)
						t.statesStack.Push(quotingState)
					}
				case escapeRuneClass:
					{
						tokenType = WordToken
						t.statesStack.Push(inWordState)
						t.statesStack.Push(escapingState)
					}
				case commentRuneClass:
					{
						tokenType = CommentToken
						t.statesStack.Push(commentState)
					}
				case endLineRuneClass:
					{
						token := &Token{
							TokenType: EndLineToken,
							Value:     string(nextRune)}
						return token, nil
					}
				case pipeRuneClass:
					{
						token := &Token{
							TokenType: PipeToken,
							Value:     string(nextRune)}
						return token, nil
					}
				default:
					{
						tokenType = WordToken
						t.statesStack.Push(inWordState)
						value = append(value, nextRune)
					}
				}
			}
		case inWordState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						t.isEnded = true
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, err
					}
				case spaceRuneClass:
					{
						t.statesStack.Pop()
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, err
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
				case endLineRuneClass:
					{
						t.statesStack.Pop()
						t.statesStack.Push(endLineState)
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, nil
					}
				case pipeRuneClass:
					{
						t.statesStack.Pop()
						t.statesStack.Push(pipeSymbolState)
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, nil
					}
				default:
					{
						value = append(value, nextRune)
					}
				}
			}
		case escapingState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						t.isEnded = true
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, fmt.Errorf("EOF found after escape character")
					}
				default:
					{
						t.statesStack.Pop()
						value = append(value, nextRune)
					}
				}
			}
		case escapingQuotedState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						t.isEnded = true
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, fmt.Errorf("EOF found after escape character")
					}
				default:
					{
						t.statesStack.Pop()
						value = append(value, nextRune)
					}
				}
			}
		case quotingEscapingState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, fmt.Errorf("EOF found when expecting closing quote")
					}
				case escapingQuoteRuneClass:
					{
						t.statesStack.Pop()
					}
				case escapeRuneClass:
					{
						t.statesStack.Push(escapingQuotedState)
					}
				default:
					{
						value = append(value, nextRune)
					}
				}
			}
		case quotingState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						t.isEnded = true
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, fmt.Errorf("EOF found when expecting closing quote")
					}
				case nonEscapingQuoteRuneClass:
					{
						t.statesStack.Pop()
					}
				default:
					{
						value = append(value, nextRune)
					}
				}
			}
		case commentState:
			{
				switch nextRuneType {
				case eofRuneClass:
					{
						t.isEnded = true
						token := &Token{
							TokenType: tokenType,
							Value:     string(value)}
						return token, err
					}
				case endLineRuneClass:
					{
						t.statesStack.Pop()
						return &Token{TokenType: EndLineToken, Value: endLineRunes}, nil
					}
				case spaceRuneClass:
					{
						value = append(value, nextRune)
					}
				default:
					{
						value = append(value, nextRune)
					}
				}
			}
		default:
			{
				return nil, fmt.Errorf("Unexpected state: %v", state)
			}
		}
	}
}

func (t *Tokenizer) Next() (*Token, error) {
	return t.scanStream()
}

func SplitOnTokens(s string) ([]Token, error) {
	tokenizer := NewTokenizer(strings.NewReader(s))
	tokens := make([]Token, 0)
	for {
		token, err := tokenizer.Next()
		if err != nil {
			if err == io.EOF {
				return tokens, nil
			}
			return []Token{}, err
		}
		tokens = append(tokens, *token)
	}
}
