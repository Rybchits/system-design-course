package executor

import (
	"bytes"
	"os"
	"shell/internal/command_meta"
	"testing"
)

func TestExecutorEmpty(t *testing.T) {
	pf := NewPipelineFactory()
	p := pf.CreatePipeline(os.Stdin, os.Stdout, []command_meta.CommandMeta{})
  	if p != nil {
		t.Fatal("Empty pipeline is not nil")
	}
}

func TestExecutorEchoCatCat(t *testing.T) {
	expected := "Oh, hello"

	var args []string
	var meta command_meta.CommandMeta
	var metas []command_meta.CommandMeta

	args = []string{}
	args = append(args, expected)
	meta = command_meta.CommandMeta{Name: "echo", Args: args}
	metas = append(metas, meta)

	meta = command_meta.CommandMeta{Name: "cat", Args: []string{}}
	metas = append(metas, meta)

	meta = command_meta.CommandMeta{Name: "cat", Args: []string{}}
	metas = append(metas, meta)

	rp, wp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer rp.Close()

	pf := NewPipelineFactory()
	p := pf.CreatePipeline(nil, wp, metas)
	err = p.Execute()
	wp.Close()
	if err != nil {
		t.Fatal("Can't execute pipe", err)
	}

	buf := make([]byte, 128)
	n, err := rp.Read(buf)
	if err != nil {
		t.Fatal("Can't read pipe", err)
	}
	if n > 0 && !bytes.Equal(buf[:n-1], []byte(expected)) {
		t.Fatalf(`Different outputs: %q != %q`, buf[:n-1], []byte(expected))
	}
}

func TestExecutorEchoCatWc(t *testing.T) {
	input := "12345"
	expected := "\t1\t1\t6"

	var args []string
	var meta command_meta.CommandMeta
	var metas []command_meta.CommandMeta

	args = []string{}
	args = append(args, input)
	meta = command_meta.CommandMeta{Name: "echo", Args: args}
	metas = append(metas, meta)

	meta = command_meta.CommandMeta{Name: "cat", Args: []string{}}
	metas = append(metas, meta)

	meta = command_meta.CommandMeta{Name: "wc", Args: []string{}}
	metas = append(metas, meta)

	rp, wp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer rp.Close()

	pf := NewPipelineFactory()
	p := pf.CreatePipeline(nil, wp, metas)
	err = p.Execute()
	wp.Close()
	if err != nil {
		t.Fatal("Can't execute pipe", err)
	}

	buf := make([]byte, 128)
	n, err := rp.Read(buf)
	if err != nil {
		t.Fatal("Can't read pipe", err)
	}
	if n > 0 && !bytes.Equal(buf[:n-1], []byte(expected)) {
		t.Fatalf(`Different outputs: %q != %q`, buf[:n-1], []byte(expected))
	}
}
