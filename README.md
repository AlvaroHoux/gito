# 🐙 Gito

> A Go-based CLI tool designed to streamline AI-assisted commits. Gito captures your staged changes to provide instant context for LLMs, making it effortless to generate accurate, Conventional Commit messages.

Gito acts as a bridge between your code and AI. It natively integrates with [Ollama](https://ollama.com/) to generate commits automatically on your local machine. If Ollama is offline, Gito smartly falls back to copying your diff and a specialized system prompt directly to your clipboard, so you can paste it into ChatGPT, Gemini, or Claude.

## ✨ Features

- **Local AI Integration:** Connects directly to Ollama to generate commit messages without leaving the terminal.
- **Smart Fallback (Clipboard):** If Ollama isn't running, it copies the `git diff` + a meticulously crafted prompt directly to your clipboard (`wl-copy`, `pbcopy`, or `clip`).
- **Conventional Commits:** The built-in prompt ensures your commits always follow standard formatting (e.g., `feat:`, `fix:`, `refactor:`).
- **Interactive or Silent:** Prompts for confirmation before applying the commit, or accepts the `-y` flag to skip the prompt.
- **Cross-Platform:** Works seamlessly on Linux (Wayland/X11), macOS, and Windows.

## 🚀 Installation

Ensure you have [Go](https://go.dev/) installed on your machine.

```bash
go install github.com/AlvaroHoux/gito/cmd/gito@latest
```

Make sure your Go bin directory is in your system's $PATH.

## 🛠️ Prerequisites
- **Git**: Must be installed and initialized in your repository.
- **Ollama** (Optional but recommended): For local AI generation.
  - Install from ollama.com
  - Pull a model, e.g., ollama pull granite3.3:2b

## 💻 Usage

First, stage your changes as you normally would:
```bash
git add .
```

Then, just run:
```bash
gito
```

### Flags and Options
**Skip Confirmation**: Apply the generated commit immediately without asking.
```bash
gito -y
```
**Specify a Model**: Temporarily use a different Ollama model for this commit.
```bash
gito -m llama3
```

### Configuration

You can set a default model so you don't have to specify it every time. Gito saves this safely in your OS's native config directory (`~/.config/gito/config.json` on Linux).

```bash
gito config -m granite3.3:2b
```

### 📝 The Prompt & Customization

Gito comes with a carefully crafted, built-in system prompt (you can check it out in the [`prompt.txt`](./prompt.txt) file in this repository). If you prefer using a web interface, you can also use this exact prompt via our **[Official Gito GEM](https://gemini.google.com/gem/1qXqsxUgAovbu6QgAOtAzQcMCqcmwMLnc?usp=sharing)**.

However, if you want to enforce specific rules for your team locally (e.g., "always use emojis", "write in Portuguese", or "reference Jira tickets"), you can easily override the default prompt.

Simply create a `prompt.txt` file in Gito's configuration directory:

- **Linux:** `~/.config/gito/prompt.txt`
- **Windows:** `%APPDATA%\gito\prompt.txt`
- **macOS:** `~/Library/Application Support/gito/prompt.txt`

If this file exists, Gito will automatically use its contents instead of the default prompt!

## 🧠 How the Fallback Works

If you run `gito` and the Ollama server is **not** active, the CLI won't crash. Instead, it captures your `git diff --staged` and copies the entire block to your clipboard. 

You will see: `🐙 Gito: Copied to clipboard!`. 
Just `Ctrl+V` into your favorite web AI, and it will give you the perfect commit message.

✨ **For Gemini Users:** We have an official **[Gito GEM](https://gemini.google.com/gem/1qXqsxUgAovbu6QgAOtAzQcMCqcmwMLnc?usp=sharing)** pre-configured with our exact system prompt. If you are using the fallback, just open the GEM and paste your diff there!

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
