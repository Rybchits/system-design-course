package parser_test

import (
	"shell/internal/command_meta"
	envsholder "shell/internal/envs_holder"
	. "shell/internal/parser"
	"strings"
	"testing"
)

func TestComplexEnvs(t *testing.T) {
	s := "x=\"once upon\"=\"a\" y=\"a time\" bash -c 'echo $x $y'\n"
	vars := envsholder.Env{}
	tokenizer := NewTokenizer(strings.NewReader(s), &vars)
	parser := NewParser(tokenizer)
	commands, err := parser.Parse()
	if err != nil {
		t.Fail()
	}

	expected := []command_meta.CommandMeta{
		{
			Name: "bash",
			Args: []string{"-c", "echo $x $y"},
			Envs: envsholder.Env{
				Vars: map[string]string{"x": "once upon=a", "y": "a time"},
			},
		},
	}

	if len(expected) != len(commands) {
		t.Fail()
	}

	for i := range expected {
		if !commands[i].Equal(&expected[i]) {
			t.Fail()
		}
	}
}
