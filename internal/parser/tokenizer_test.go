package parser_test

import (
	"fmt"
	"io"
	envsholder "shell/internal/envs_holder"
	. "shell/internal/parser"
	"strings"
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

func splitOnTokens(s string, vars map[string]string) ([]Token, error) {
	envs := envsholder.Env{}
	envs.Init()
	for k, v := range vars {
		envs.Set(k, v)
	}
	tokenizer := NewTokenizer(strings.NewReader(s), &envs)

	tokens := make([]Token, 0)
	for {
		token, err := tokenizer.Next()
		if err != nil {
			if err == io.EOF {
				if token != nil {
					tokens = append(tokens, *token)
				}
				return tokens, nil
			}
			return []Token{}, err
		}
		tokens = append(tokens, *token)
	}
}

func TestSingleWordTokenizer(t *testing.T) {
	s := "bober"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "bober"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestSimpleTokenizer(t *testing.T) {
	s := "echo hello world"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "echo"},
		{TokenType: WordToken, Value: "hello"},
		{TokenType: WordToken, Value: "world"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestPipeTokenizer1(t *testing.T) {
	s := "cat|echo"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		fmt.Println(tokens)
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "cat"},
		{TokenType: PipeToken, Value: "|"},
		{TokenType: WordToken, Value: "echo"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestEmptyInput(t *testing.T) {
	s := ""
	tokens, err := splitOnTokens(s, map[string]string{})

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
	s := " \t\t\t  \t"
	tokens, err := splitOnTokens(s, map[string]string{})

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
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "echo"},
		{TokenType: WordToken, Value: "hello"},
		{TokenType: EndLineToken, Value: "\n"},
		{TokenType: WordToken, Value: "world"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestPipeTokenizer2(t *testing.T) {
	s := "cat | echo"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "cat"},
		{TokenType: PipeToken, Value: "|"},
		{TokenType: WordToken, Value: "echo"},
	})

	if !result {
		t.Fail()
	}
}

func TestStringWithEndLine(t *testing.T) {
	s := "cat echo\n"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "cat"},
		{TokenType: WordToken, Value: "echo"},
		{TokenType: EndLineToken, Value: "\n"},
	})

	if !result {
		t.Fail()
	}
}

func TestMultipleSpace(t *testing.T) {
	s := "cat     echo     \n\n"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "cat"},
		{TokenType: WordToken, Value: "echo"},
		{TokenType: EndLineToken, Value: "\n"},
		{TokenType: EndLineToken, Value: "\n"},
	})

	if !result {
		t.Fail()
	}
}

func TestSingleCommand(t *testing.T) {
	s := "'cat echo'"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "cat echo"},
	})

	if !result {
		t.Fail()
	}
}

func TestNewlineInString(t *testing.T) {
	s := "'cat\necho'"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "cat\necho"},
	})

	if !result {
		t.Fail()
	}
}

func TestEscapingInString(t *testing.T) {
	s := "'cat\\necho'"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "cat\\necho"},
	})

	if !result {
		t.Fail()
	}
}

func TestNonEscapingQuoteStirng(t *testing.T) {
	s := "echo 'abc'"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "echo"},
		{TokenType: WordToken, Value: "abc"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestNonEscapingQuote2Stirng(t *testing.T) {
	s := "echo'abc'"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "echoabc"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestManyPipesStirng(t *testing.T) {
	s := "|| a |"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: PipeToken, Value: "|"},
		{TokenType: PipeToken, Value: "|"},
		{TokenType: WordToken, Value: "a"},
		{TokenType: PipeToken, Value: "|"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestSpecialSimbols(t *testing.T) {
	s := "echo '(´･_･`)'"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "echo"},
		{TokenType: WordToken, Value: "(´･_･`)"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestSimpleEnvs(t *testing.T) {
	s := "x=\"once upon \" y=\"a time\" bash -c 'echo $x $y'"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "x=once upon "},
		{TokenType: WordToken, Value: "y=a time"},
		{TokenType: WordToken, Value: "bash"},
		{TokenType: WordToken, Value: "-c"},
		{TokenType: WordToken, Value: "echo $x $y"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestComplexEnvsWithStrings(t *testing.T) {
	s := "x=\"once upon\"=\"a\" y=\"a time\" bash -c 'echo $x $y'"
	tokens, err := splitOnTokens(s, map[string]string{})

	if err != nil {
		t.Fail()
	}

	result := compareTwoTokensArray(tokens, []Token{
		{TokenType: WordToken, Value: "x=once upon=a"},
		{TokenType: WordToken, Value: "y=a time"},
		{TokenType: WordToken, Value: "bash"},
		{TokenType: WordToken, Value: "-c"},
		{TokenType: WordToken, Value: "echo $x $y"},
	})

	if !result {
		fmt.Println(tokens)
		t.Fail()
	}
}

func TestEnvVarTestStirng(t *testing.T) {
	list := []string{
		"echo $a",
		"echo $aa",
		"echo '$aa'",
		"echo \"$aa\"",
		"echo \"$a\"",
		"$a$b",
		"echo $a$b",
		"echo $a|$b",
		"echo \"$a|$b\"",
	}

	vars := map[string]string{
		"a": "ec",
		"b": "ho",
	}

	expected := [][]Token{
		{
			{TokenType: WordToken, Value: "echo"},
			{TokenType: WordToken, Value: "ec"},
		},
		{
			{TokenType: WordToken, Value: "echo"},
		},
		{
			{TokenType: WordToken, Value: "echo"},
			{TokenType: WordToken, Value: "$aa"},
		},
		{
			{TokenType: WordToken, Value: "echo"},
		},
		{
			{TokenType: WordToken, Value: "echo"},
			{TokenType: WordToken, Value: "ec"},
		},
		{
			{TokenType: WordToken, Value: "echo"},
		},
		{
			{TokenType: WordToken, Value: "echo"},
			{TokenType: WordToken, Value: "echo"},
		},
		{
			{TokenType: WordToken, Value: "echo"},
			{TokenType: WordToken, Value: "ec"},
			{TokenType: PipeToken, Value: "|"},
			{TokenType: WordToken, Value: "ho"},
		},
		{
			{TokenType: WordToken, Value: "echo"},
			{TokenType: WordToken, Value: "ec|ho"},
		},
	}

	for i, str := range list {
		tokens, err := splitOnTokens(str, vars)

		if err != nil {
			t.Fail()
		}

		result := compareTwoTokensArray(tokens, expected[i])

		if !result {
			fmt.Println(i, " ", tokens)
			t.Fail()
		}
	}
}
