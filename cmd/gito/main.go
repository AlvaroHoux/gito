package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
)

func main() {
	diffCmd := exec.Command("git", "diff", "--staged")
	diffOut, err := diffCmd.Output()

	if err != nil {
		fmt.Println("🐙 Gito: Erro", err)
		return
	}

	if err := copyToClipboard(string(diffOut)); err != nil {
		fmt.Println("🐙 Gito: Erro", err)
		return
	}

	fmt.Println("🐙 Gito: Copied to clipboard")
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("wl-copy")
	}
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}