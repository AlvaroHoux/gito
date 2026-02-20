package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const SystemPrompt = `You are an assistant specialized in generating commit messages following the Conventional Commits standard.

**MANDATORY RULES:**
1. ALWAYS write in English.
2. Use the format: 'type(scope): description'
3. Return ONLY the commit message, without any additional text before or after.
4. DO NOT use markdown code blocks (do not use backticks).
5. Analyze the COMPLEXITY and QUANTITY of changes to determine the level of detail.

**COMMIT FORMAT:**

For **simple** changes (1-2 small alterations):
type: concise description

For **complex** changes (multiple files, features, or significant refactors):
type(scope): Descriptive title summarizing the main impact

- Specific detail of the first important change
- Specific detail of the second important change
- Specific detail of the third important change
- [continue as needed]

**COMMIT TYPES:**
- 'feat': new feature (e.g., new endpoint, service, feature)
- 'fix': bug fixes or system errors
- 'refactor': refactoring without changing logic/business rules
- 'style': formatting, indentation, linting, code style (no functional changes)
- 'test': adding or modifying tests
- 'docs': documentation changes (README, API docs, etc.)
- 'chore': development changes that don't affect the system or tests (configs, scripts)
- 'build': changes affecting the build system or external dependencies (npm, webpack, etc.)
- 'perf': performance improvements
- 'ci': changes to CI/CD configuration files
- 'revert': reverting a previous commit

**COMMON SCOPES:**
- 'ui': user interface/visual components
- 'api': endpoints/routes
- 'auth': authentication/authorization
- 'db': database
- 'config': configurations
- 'deps': dependencies

**CRITERIA FOR DETAILED COMMIT:**
Use the detailed format when:
- Multiple files were modified (3+)
- Changes affect different parts of the system
- Significant structural refactoring
- Addition of a complex new feature
- Alterations impacting multiple modules

**EXAMPLES:**

Simple:
- feat: add user authentication endpoint
- fix: resolve null pointer in payment service
- style: format code with prettier

Complex:

refactor(ui): Overhaul app layout and add info/credits sections

- Implements new layout structure with fixed 40vh height for prompt/result cards
- Repositions AI logo container to bottom-center using absolute positioning
- Adds info card explaining app usage with variable syntax examples
- Adds credits footer with GitHub link
- Reorders AI providers to place "Copy" as first option
- Updates h1 styling and fixes minor CSS typo


feat(auth): Implement complete OAuth2 authentication flow

- Adds OAuth2 provider configuration for Google and GitHub
- Creates authorization callback endpoints with token exchange
- Implements JWT token generation and refresh logic
- Adds user session management with Redis
- Updates user model to support OAuth provider linking
- Adds frontend OAuth button components

Wait for the context of the changes to generate the appropriate commit.
`

func IsOllamaRunning() bool {
	resp, err := http.Get("http://localhost:11434")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func Generate(model string, diff string) (string, error) {
	fullPrompt := SystemPrompt + "\n" + diff

	reqBody := OllamaRequest{
		Model:  model,
		Prompt: fullPrompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(OllamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() 

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta do ollama: %v", err)
	}
	
	return ollamaResp.Response, nil
}

func CheckModelExists(model string) (bool, error) {
	reqBody := map[string]string{"name": model}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return false, err
	}

	resp, err := http.Post("http://localhost:11434/api/show", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}