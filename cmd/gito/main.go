package main

import (
	"bytes"
	"flag"
	"fmt"
	"orbit/internal/config"
	"orbit/internal/ollama"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	configCmd := flag.NewFlagSet("config", flag.ContinueOnError)
	model := configCmd.String("model", "granite3.3:2b", "Ollama default model")

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("🐙 Gito: Cannot load config data >", err)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "config" {
		configCmd.Parse(os.Args[2:])
		fmt.Println(config.Model, *model)
	}

	diffCmd := exec.Command("git", "diff", "--staged")
	diffOut, err := diffCmd.Output()

	if err != nil {
		fmt.Println("🐙 Gito: Erro", err)
		return
	}

	if running := ollama.IsOllamaRunning(); running == false {
		fmt.Println("🐙 Gito: Ollama is not running, use ollama serve in terminal")
		if err := copyToClipboard(string(diffOut)); err != nil {
			fmt.Println("🐙 Gito: Erro", err)
			return
		}
		fmt.Println("🐙 Gito: Copied to clipboard")
	} else {
		output, err := ollama.Generate(config.Model, string(diffOut))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(output)
	}
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("wl-copy")
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("Not suported system")
	}
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}
