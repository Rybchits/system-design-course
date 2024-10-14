package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"shell/internal/command_meta"
	"testing"
)

func TestWcExecute(t *testing.T) {
	file, err := ioutil.TempFile(".", "test")
	if err != nil {
		t.Fatal("Can't create temp file", err)
	}
	defer os.Remove(file.Name())

	expected := []byte(fmt.Sprintf("%d %d %d %s\n", 1, 2, 11, file.Name()))
	file.Write([]byte("Hello world"))

	args := make([]string, 1)
	args = append(args, file.Name())
	meta := command_meta.CommandMeta{Name: "wc", Args: args}
	in, out, err := os.Pipe()
	if err != nil {
		t.Fatal("Can't create pipe", err)
	}
	defer in.Close()
	defer out.Close()

	r := bufio.NewReader(in)
	buf := make([]byte, 0, 128)
	cmd := WcCommand{nil, out, meta}
	go cmd.Execute()

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
