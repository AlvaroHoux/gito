package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"orbit/internal/clipboard"
	"orbit/internal/config"
	"orbit/internal/git"
	"orbit/internal/ollama"
	"orbit/internal/term"
)

func askConfirmation() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("🐙 Gito: Would you like to commit this message? [Y/n]: ")
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.ToLower(strings.TrimSpace(response))

	switch response {
	case "yes", "y", "":
		return true, nil
	case "no", "n":
		return false, nil
	}

	return true, nil
}

func main() {
	var runModel string
	var ignoreAsk bool

	flag.StringVar(&runModel, "m", "", "Ollama model name (shorthand)")
	flag.StringVar(&runModel, "model", "", "Ollama model name")
	flag.BoolVar(&ignoreAsk, "y", false, "Ignore asks")

	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	var configModel string

	configCmd.StringVar(&configModel, "m", "", "Model name to save (shorthand)")
	configCmd.StringVar(&configModel, "model", "", "Model name to save")

	flag.Parse()

	configData, err := config.LoadConfig()
	if err != nil {
		term.Error(err)
		return
	}

	if runModel == "" {
		runModel = configData.Model
	}

	if len(os.Args) > 1 && os.Args[1] == "config" {
		configCmd.Parse(os.Args[2:])

		if configModel == "" {
			term.Warning("Usage: gito config -m <model>")
			return
		}

		if err := config.SaveConfig(configModel); err != nil {
			term.Error(err)
			return
		}
		term.Log("Config saved")
		return
	}

	diffOut, err := git.GetDiff()
	if err != nil {
		term.Error(err)
		return
	}

	if strings.TrimSpace(diffOut) == "" {
		term.Warning("No staged changes found. Did you forget to run 'git add'?")
		return
	}

	if ollama.IsOllamaRunning() {
		exists, err := ollama.CheckModelExists(runModel)
		if err != nil {
			term.Error(err)
			return
		}

		if !exists {
			term.Error(fmt.Errorf("model '%s' not found. Run 'ollama pull %s'", runModel, runModel))
			return
		}

		term.Log("Using model", runModel)

		output, err := ollama.Generate(runModel, diffOut)
		if err != nil {
			term.Error(err)
			return
		}

		term.Log("Generated commit\n")
		fmt.Println(output, "\n")

		var confirm bool
		if ignoreAsk {
			confirm = true
		} else {
			confirm, err = askConfirmation()
			if err != nil {
				term.Error(err)
				return
			}
		}

		if confirm {
			if err := git.Commit(output); err != nil {
				term.Error(err)
				return
			} else {
				term.Success("Message commited!")
			}
		} else {
			term.Warning("Commit canceled")
		}

	} else {
		term.Warning("Ollama is not running, use ollama serve in terminal")
		if err := clipboard.CopyToClipboard(diffOut); err != nil {
			term.Error(err)
			return
		}
		term.Success("Copied to clipboard!")
	}
}
