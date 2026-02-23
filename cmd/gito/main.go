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
	configData, err := config.LoadConfig()
	if err != nil {
		term.Error(err)
		return
	}

	var runModel string
	var ignoreAsk bool
	var onlyDiff bool

	flag.StringVar(&runModel, "m", configData.Model, "Ollama model name (shorthand)")
	flag.StringVar(&runModel, "model", configData.Model, "Ollama model name")
	flag.BoolVar(&ignoreAsk, "y", configData.SkipAsk, "Ignore asks")
	flag.BoolVar(&onlyDiff, "d", configData.OnlyDiff, "Copy only the staged diff to clipboard (shorthand)")
	flag.BoolVar(&onlyDiff, "diff", configData.OnlyDiff, "Copy only the staged diff to clipboard")


	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	var configModel string
	var configIgnoreAsk bool
	var configOnlyDiff bool

	configCmd.StringVar(&configModel, "m", configData.Model, "Model name to save (shorthand)")
	configCmd.StringVar(&configModel, "model", configData.Model, "Model name to save")
	configCmd.BoolVar(&configIgnoreAsk, "y", configData.SkipAsk, "Always ignore asks")
	configCmd.BoolVar(&configOnlyDiff, "d", configData.OnlyDiff, "Always copy only the staged diff to clipboard (shorthand)")
	configCmd.BoolVar(&configOnlyDiff, "diff", configData.OnlyDiff, "Always Ccpy only the staged diff to clipboard")

	flag.Parse()

	if len(os.Args) > 1 && os.Args[1] == "config" {
		if len(os.Args) == 2 {
			term.Log("Current configuration:")
		
			fmt.Printf("  - Model:      \033[1;36m%s\033[0m\n", configData.Model)
			fmt.Printf("  - Skip Ask:   \033[1;36m%v\033[0m\n", configData.SkipAsk)
			fmt.Printf("  - Only Diff:  \033[1;36m%v\033[0m\n\n", configData.OnlyDiff)
			
			fmt.Println("Run 'gito config -h' to see how to change these values.")
			return
		}

		configCmd.Parse(os.Args[2:])

		newConfig := config.GitoConfig{
			Model: configModel,
			SkipAsk: configIgnoreAsk,
			OnlyDiff: configOnlyDiff,
		}

		if err := config.SaveConfig(newConfig); err != nil {
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

		var fallbackText string
		var fallbackMsg string
		
		if onlyDiff {
			fallbackMsg = "Diff"
			fallbackText = diffOut	
		} else {
			fallbackMsg = "Prompt + Diff"
			fallbackText = ollama.GetSystemPrompt() + "\n\nDiff:\n" + diffOut
		}

		if err := clipboard.CopyToClipboard(fallbackText); err != nil {
			term.Error(err)
			return
		}
		term.Success(fallbackMsg, "copied to clipboard!")
	}
}
