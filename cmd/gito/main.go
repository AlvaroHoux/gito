package main

import (
	"flag"
	"orbit/internal/config"
	"orbit/internal/ollama"
	"orbit/internal/clipboard"
	"orbit/internal/term"
	"os"
	"os/exec"
)

func main() {
	configCmd := flag.NewFlagSet("config", flag.ContinueOnError)
	modelArg := flag.String("model", "default", "Ollama default model")

	flag.Parse()

	config, err := config.LoadConfig()
	if err != nil {
		term.Error(err)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "config" {
		configCmd.Parse(os.Args[2:])
	}

	model := *modelArg
	if *modelArg == "default" {
		model = config.Model
	}

	diffCmd := exec.Command("git", "diff", "--staged")
	diffOut, err := diffCmd.Output()

	if err != nil {
		term.Error(err)
		return
	}

	if running := ollama.IsOllamaRunning(); running == false {
		term.Warn("Ollama is not running, use ollama serve in terminal")
		if err := clipboard.CopyToClipboard(string(diffOut)); err != nil {
			term.Error(err)
			return
		}
		term.Log("Copied to clipboard!")
	} else {
		term.Log("Using model", model)
		output, err := ollama.Generate(model, string(diffOut))
		if err != nil {
			term.Error(err)
		}
		term.Log(output)
	}
}