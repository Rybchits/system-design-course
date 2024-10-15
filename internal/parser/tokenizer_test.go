package parser

import (
	"fmt"
	"testing"
)

func compareTwoTokensArray(lhs []Token, rhs []Token) bool {
	if len(lhs) != len(rhs) {
		return false
	}
	for index := range lhs {
		if !lhs[index].Equal(&rhs[index]) {
			return false
		}
	}
	return true
}

func TestSimpleTokenizer(t *testing.T) {
	s := "echo \"hello world\""
	tokens, err := SplitOnTokens(s)

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		Token{TokenType: WordToken, Value: "echo"},
		Token{TokenType: WordToken, Value: "hello world"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestPipeTokenizer1(t *testing.T) {
	s := "cat|echo"
	tokens, err := SplitOnTokens(s)

	if err != nil {
		fmt.Println(tokens)
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		Token{TokenType: WordToken, Value: "cat"},
		Token{TokenType: PipeToken, Value: "|"},
		Token{TokenType: WordToken, Value: "echo"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestEmptyInput(t *testing.T) {
	s := ""
	tokens, err := SplitOnTokens(s)

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestOnlySpaceSymbols(t *testing.T) {
	s := " \t \t \t  \t"
	tokens, err := SplitOnTokens(s)

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestEndLineInMiddleWorld(t *testing.T) {
	s := "echo hello\nworld"
	tokens, err := SplitOnTokens(s)

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		Token{TokenType: WordToken, Value: "echo"},
		Token{TokenType: WordToken, Value: "hello"},
		Token{TokenType: EndLineToken, Value: "\n"},
		Token{TokenType: WordToken, Value: "world"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestPipeTokenizer2(t *testing.T) {
	s := "cat | echo"
	tokens, err := SplitOnTokens(s)

	if err != nil {
		t.Fail()
	}

	fmt.Println(tokens)
	result := compareTwoTokensArray(tokens, []Token{
		Token{TokenType: WordToken, Value: "cat"},
		Token{TokenType: PipeToken, Value: "|"},
		Token{TokenType: WordToken, Value: "echo"},
	})

	if !result {
		t.Fail()
	}
}

func TestStringWithEndLine(t *testing.T) {
	s := "cat echo\n"
	tokens, err := SplitOnTokens(s)

	if err != nil {
		t.Fail()
	}

	fmt.Println(tokens)
	result := compareTwoTokensArray(tokens, []Token{
		Token{TokenType: WordToken, Value: "cat"},
		Token{TokenType: WordToken, Value: "echo"},
		Token{TokenType: EndLineToken, Value: "\n"},
	})

	if !result {
		t.Fail()
	}
}

func TestMultipleSpace(t *testing.T) {
	s := "cat     echo     \n\n"
	tokens, err := SplitOnTokens(s)

	if err != nil {
		t.Fail()
	}

	fmt.Println(tokens)
	result := compareTwoTokensArray(tokens, []Token{
		Token{TokenType: WordToken, Value: "cat"},
		Token{TokenType: WordToken, Value: "echo"},
		Token{TokenType: EndLineToken, Value: "\n"},
		Token{TokenType: EndLineToken, Value: "\n"},
	})

	if !result {
		t.Fail()
	}
}
