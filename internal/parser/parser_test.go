package parser_test

import (
	"fmt"
	. "shell/internal/parser"
	"shell/internal/command_meta"
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

func TestComplexEnvs(t *testing.T) {
	s := "x=\"once upon\"=\"a\" y=\"a time\" bash -c 'echo $x $y'"
	parser := NewParser(strings.NewReader(s))
	commands, err := parser.Parse()
	if err != nil {
		t.Fail()
	}
	
	expected := []command_meta.CommandMeta{
		command_meta.CommandMeta{
			Name: "bash", 
			Args: []string{"-c", "echo $x $y"},
		},
	}

	for i := range commands {
		if !commands[i].Equal(&expected[i]) {
			t.Fail()
		}
	}
}