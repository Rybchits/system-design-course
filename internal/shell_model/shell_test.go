package shellmodel

import (
	"bytes"
	"os"
	"testing"
)

func TestShellCommand(t *testing.T) {
	in_read, in_write, _ := os.Pipe()
	out_read, out_write, _ := os.Pipe()

	test_shell := NewShell()

	buf := make([]byte, 128)
	expected := []byte("42\n")

	go func(sh *Shell, in *os.File, out *os.File) {
		sh.ShellLoop(in, out, false)
	}(test_shell, in_read, out_write)

	in_write.WriteString("echo 42\n")
	in_write.Close()
	n, err := out_read.Read(buf)
	if err != nil {
		t.Fatal("Cant read pipe", err)
	}
	buf = buf[:n]

	if !bytes.Equal(buf, expected) {
		t.Fatalf(`Different outputs: %q != %q`, buf, expected)
	}
}

func TestPipe(t *testing.T) {
	in_read, in_write, _ := os.Pipe()
	out_read, out_write, _ := os.Pipe()

	test_shell := NewShell()

	buf := make([]byte, 128)
	expected := []byte("\t1\t1\t3\n")

	go func(sh *Shell, in *os.File, out *os.File) {
		sh.ShellLoop(in, out, false)
	}(test_shell, in_read, out_write)

	in_write.WriteString("echo 42 | wc\n")
	in_write.Close()
	n, err := out_read.Read(buf)
	if err != nil {
		t.Fatal("Cant read pipe", err)
	}
	buf = buf[:n]

	if !bytes.Equal(buf, expected) {
		t.Fatalf(`Different outputs: %q != %q`, buf, expected)
	}
}
