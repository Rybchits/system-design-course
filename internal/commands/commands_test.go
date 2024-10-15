package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"shell/internal/command_meta"
	"testing"
)

func TestWcExecuteSimple(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		t.Fatal("Can't create temp file", err)
	}
	defer os.Remove(file.Name())

	expected := []byte(fmt.Sprintf("\t%d\t%d\t%d\t%s\n", 2, 2, 12, file.Name()))
	file.Write([]byte("Hello\nworld"))
	file.Close()

	args := make([]string, 0)
	args = append(args, file.Name())
	meta := command_meta.CommandMeta{Name: "wc", Args: args}
	rp, wp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer rp.Close()

	r := bufio.NewReader(rp)
	buf := make([]byte, 0, 128)
	cmd := WcCommand{nil, wp, meta}
	go func(cmd WcCommand, wp *os.File) {
		defer wp.Close()
		cmd.Execute()
	}(cmd, wp)

	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			t.Fatal("Can't read buffer", err)
		}

		if !bytes.Equal(buf, expected) {
			t.Fatalf(`Different outputs: %q != %q`, buf, expected)
		}
		break
	}
}

func TestWcExecuteEmpty(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		t.Fatal("Can't create temp file", err)
	}
	defer os.Remove(file.Name())

	expected := []byte(fmt.Sprintf("\t%d\t%d\t%d\t%s\n", 0, 0, 0, file.Name()))
	file.Write([]byte(""))
	file.Close()

	args := make([]string, 0)
	args = append(args, file.Name())
	meta := command_meta.CommandMeta{Name: "wc", Args: args}
	rp, wp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer rp.Close()

	r := bufio.NewReader(rp)
	buf := make([]byte, 0, 128)
	cmd := WcCommand{nil, wp, meta}
	go func(cmd WcCommand, wp *os.File) {
		defer wp.Close()
		cmd.Execute()
	}(cmd, wp)

	{
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]

		if err != nil {
			t.Fatal("Can't read buffer", err)
		}

		if !bytes.Equal(buf, expected) {
			t.Fatalf(`Different outputs: %q != %q`, buf, expected)
		}
	}
}

//////////////////////////////////

func TestCatExecuteSimple(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		t.Fatal("Can't create temp file", err)
	}
	defer os.Remove(file.Name())

	expected := []byte("Hello world")
	file.Write([]byte("Hello world"))
	file.Close()

	args := make([]string, 0)
	args = append(args, file.Name())
	meta := command_meta.CommandMeta{Name: "cat", Args: args}
	rp, wp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer rp.Close()

	r := bufio.NewReader(rp)
	buf := make([]byte, 0, 128)
	cmd := CatCommand{nil, wp, meta}
	go func(cmd CatCommand, wp *os.File) {
		defer wp.Close()
		cmd.Execute()
	}(cmd, wp)

	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			t.Fatal("Can't read buffer", err)
		}

		if !bytes.Equal(buf, expected) {
			t.Fatalf(`Different outputs: %q != %q`, buf, expected)
		}
		break
	}
}

func TestCatExecuteEmpty(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		t.Fatal("Can't create temp file", err)
	}
	defer os.Remove(file.Name())

	expected := []byte("")
	file.Write([]byte(""))
	file.Close()

	args := make([]string, 0)
	args = append(args, file.Name())
	meta := command_meta.CommandMeta{Name: "cat", Args: args}
	rp, wp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer rp.Close()

	r := bufio.NewReader(rp)
	buf := make([]byte, 0, 128)
	cmd := CatCommand{nil, wp, meta}
	go func(cmd CatCommand, wp *os.File) {
		defer wp.Close()
		cmd.Execute()
	}(cmd, wp)

	{
		n, _ := r.Read(buf[:cap(buf)])
		buf = buf[:n]

		if !bytes.Equal(buf, expected) {
			t.Fatalf(`Different outputs: %q != %q`, buf, expected)
		}
	}
}

//////////////////////////////////

func TestEchoExecuteSimple(t *testing.T) {
	expected := []byte("1 2 3\n")

	args := make([]string, 0)
	args = append(args, "1")
	args = append(args, "2")
	args = append(args, "3")
	meta := command_meta.CommandMeta{Name: "echo", Args: args}
	rp, wp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer rp.Close()

	r := bufio.NewReader(rp)
	buf := make([]byte, 0, 128)
	cmd := EchoCommand{nil, wp, meta}
	go func(cmd EchoCommand, wp *os.File) {
		defer wp.Close()
		cmd.Execute()
	}(cmd, wp)

	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			t.Fatal("Can't read buffer", err)
		}

		if !bytes.Equal(buf, expected) {
			t.Fatalf(`Different outputs: %q != %q`, buf, expected)
		}
		break
	}
}

func TestEchoExecuteEmpty(t *testing.T) {
	expected := []byte("\n")

	args := make([]string, 0)
	meta := command_meta.CommandMeta{Name: "echo", Args: args}
	rp, wp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer rp.Close()

	r := bufio.NewReader(rp)
	buf := make([]byte, 0, 128)
	cmd := EchoCommand{nil, wp, meta}
	go func(cmd EchoCommand, wp *os.File) {
		defer wp.Close()
		cmd.Execute()
	}(cmd, wp)

	{
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]

		if err != nil {
			t.Fatal("Can't read buffer", err)
		}

		if !bytes.Equal(buf, expected) {
			t.Fatalf(`Different outputs: %q != %q`, buf, expected)
		}
	}
}

//////////////////////////////////

func TestPwdExecuteSimple(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to read current directory path with err: %s\n", err)
	}
	expected := []byte(dir)

	args := make([]string, 0)
	meta := command_meta.CommandMeta{Name: "pwd", Args: args}
	rp, wp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer rp.Close()

	r := bufio.NewReader(rp)
	buf := make([]byte, 0, 128)
	cmd := PwdCommand{nil, wp, meta}
	go func(cmd PwdCommand, wp *os.File) {
		defer wp.Close()
		cmd.Execute()
	}(cmd, wp)

	{
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]

		if err != nil {
			t.Fatal("Can't read buffer", err)
		}

		if !bytes.Equal(buf, expected) {
			t.Fatalf(`Different outputs: %q != %q`, buf, expected)
		}
	}
}

//////////////////////////////////

func TestProcessExecuteSimple(t *testing.T) {
	exrp, exwp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer exwp.Close()
	defer exrp.Close()

	process := exec.Command("ls")
	process.Stdout = exwp
	go process.Run()

	exr := bufio.NewReader(exrp)
	expected := make([]byte, 0, 1024)
	n, err := exr.Read(expected[:cap(expected)])
	expected = expected[:n]

	if err != nil {
		t.Fatal("Can't read buffer", err)
	}

	args := make([]string, 0)
	meta := command_meta.CommandMeta{Name: "ls", Args: args}
	rp, wp, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer rp.Close()

	r := bufio.NewReader(rp)
	buf := make([]byte, 0, 1024)
	cmd := ProcessCommand{nil, wp, meta}
	go func(cmd ProcessCommand, wp *os.File) {
		defer wp.Close()
		cmd.Execute()
	}(cmd, wp)

	{
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]

		if err != nil {
			t.Fatal("Can't read buffer", err)
		}

		if !bytes.Equal(buf, expected) {
			t.Fatalf(`Different outputs: %q != %q`, buf, expected)
		}
	}
}
