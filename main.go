package main

import (
	"os"
	"os/signal"
	shellmodel "shell/internal/shell_model"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sh := shellmodel.NewShell()

	go func() {
		sh.ShellLoop(os.Stdin, os.Stdout)
	}()

	<-sigChan
	sh.Terminate()
}
