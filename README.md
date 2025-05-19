# GoAgent - Your Local AI Programming Assistant

GoAgent is a command-line tool that acts as a local AI programming assistant. It leverages Large Language Models (LLMs) to understand and implement code changes based on natural language requests directly within your local Git repository.

## Features

- **Local Codebase Interaction:** Directly modifies files in your local Git repository based on your instructions.
- **Change Request Processing:** Accepts natural language change requests and translates them into code modifications.
- **Multiple LLM Support:** Supports various LLMs including Gemini, Groq, and OpenAI, allowing you to choose the best model for your needs.
- **Rate Limiting:** Implements rate limiting to manage API usage and prevent exceeding service limits.
- **Git Integration:** Automatically stages and commits changes with a generated commit message.
- **Interactive Mode:** Provides an interactive command-line interface for specifying change requests and asking questions about your codebase.
- **Configuration Options:** Allows customization of the LLM service, API keys, and rate limits through command-line flags and environment variables.

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Git installed and initialized in your project directory
- An API key for your chosen LLM service (Gemini, Groq, or OpenAI)

### Installation

Clone the repository to your local machine:

```bash
git clone <repository-url>
cd GoAgent
```

Build the GoAgent binary:

```bash
go build -o go-agent cmd/go-agent/main.go
```

Make the binary executable:

```bash
chmod +x go-agent
```

You might want to move the `go-agent` binary to a directory in your system's PATH (e.g., `/usr/local/bin` or `~/bin`) for easier access.

## Usage

GoAgent operates on a specified directory, which **should be a Git repository**.

### Running the Agent


To start GoAgent it's enough to run it from the repository you want to be working on (assuming the API key environment variables are set).
It will start with some sane defaults for the configurations.

```bash
go-agent
```

### Workflow

1.  **Enter Change Request:** GoAgent prompts you for a change request. Type your request in natural language (e.g., `Add a function to calculate the factorial of a number in math_utils.go`).
2.  **Processing:** GoAgent analyzes the request using its internal agents, identifies relevant files, plans the changes, and generates the necessary code modifications using the configured LLM.
3.  **Applying Changes:** GoAgent applies the generated changes (potentially including applying patches) to the files in your local repository.
4.  **Commit Confirmation:** After successfully applying changes, GoAgent prompts you:
    ```
    Do you want to commit the changes? [Y]es/[N]o/[A]llways
    ```
    -   Enter `Y` or `y` to stage all changes and commit them with a generated commit message.
    -   Enter `N` or `n` to leave the changes in your working directory without committing. You can review, modify, or discard them manually using Git commands.
    -   Enter `A` or `a` to commit the current changes *and* automatically commit all subsequent changes within the current GoAgent session without further prompting.
5.  **Error Handling:** If an error occurs while GoAgent is trying to implement the changes (e.g., the LLM produces invalid code, a patch fails to apply, or a file operation fails), it will prompt you:
    ```
    An error occurred... Do you want to reset to the last commit and discard the changes made by the agent? [Y]es/[N]o
    ```
    -   Enter `Y` or `y` to run `git reset --hard HEAD`, discarding the modifications made during the failed attempt and restoring your repository to the last committed state.
    -   Enter `N` or `n` to keep the (potentially broken) changes made by the agent in your working directory for manual inspection or recovery.
6.  **Repeat:** GoAgent waits for the next change request.

GoAgent also supports tab completion for file paths when entering change requests or `/ask` queries. Simply press the Tab key to activate file path completion based on the files tracked by Git in the target repository.

### Asking Questions

You can also ask questions about your codebase using the `/ask` command followed by your query:

```
/ask What is the purpose of the main function in cmd/go-agent/main.go?
```

GoAgent will respond with an answer based on its understanding of your code. The output of the `/ask` command is rendered using the configured Glamour style. This command does not modify files or trigger the commit workflow.

## Configuration Options

GoAgent can be configured using command-line flags, environment variables, and a configuration file.

### Command-Line Flags

| Flag                 | Description                                                                                                                               | Default Value     | Environment Variable        |
|----------------------|-------------------------------------------------------------------------------------------------------------------------------------------|-------------------|-----------------------------|
| `-directory`         | Sets the project directory. **Must be a Git repository.**                                                                                   | Current directory | N/A                         |
| `-service`           | Sets the LLM service to use (`gemini`, `groq`, `openai`).                                                                                   | `gemini`          | N/A                         |
| `-gemini-api-key`    | Sets the Gemini API key. Required if `service` is `gemini`.                                                                               | ""                | `GEMINI_API_KEY`            |
| `-groq-api-key`      | Sets the Groq API key. Required if `service` is `groq`.                                                                                   | ""                | `GROQ_API_KEY`              |
| `-openai-api-key`    | Sets the OpenAI API key. Required if `service` is `openai`.                                                                               | ""                | `OPENAI_API_KEY`            |
| `-rate-limit`        | Sets the rate limit for API requests per minute. Prevents exceeding API usage limits. `0` means no limit.                                   | 0                 | N/A                         |
| `-log-level`         | Sets the logging level (`debug`, `info`, `warning`, `error`).                                                                                 | `info`            | N/A                         |
| `-glamour-style`     | Sets the Glamour style for Markdown rendering. Options: `ascii`, `auto`, `dark`, `dracula`, `tokyo-night`, `light`, `notty`, `pink`.        | `dracula`         | N/A                         |
| `-max-history-length`| Sets the maximum history length for LLM sessions.                                                                                         | 100               | N/A                         |
| `-max-process-loops` | Sets the maximum number of processing loops the agent will attempt for a single request.                                                    | 25                | N/A                         |

**Example Configuration:**

To use the Groq service with a specific API key, set a rate limit of 30 requests per minute, use the 'dark' output style, and limit process loops to 10:

```bash
./go-agent -directory /path/to/my/git/repo -service groq -groq-api-key <your-groq-api-key> -rate-limit 30 -glamour-style dark -max-process-loops 10
```

### Environment Variables

You can also set API keys using environment variables. GoAgent will automatically load these variables if the corresponding command-line flags are not provided.

- `GEMINI_API_KEY`
- `GROQ_API_KEY`
- `OPENAI_API_KEY`

**Example Environment Variable Setup (Linux/macOS):**

```bash
export GEMINI_API_KEY="your_gemini_api_key"
export GROQ_API_KEY="your_groq_api_key"
export OPENAI_API_KEY="your_openai_api_key"
```

## Supported LLM Services

- **Gemini:** Leverages the Gemini family of models for code generation and understanding.
- **Groq:** Utilizes the Groq API for fast and efficient LLM inference.
- **OpenAI:** Supports OpenAI models like GPT-4 and GPT-4.5.

### Configuration File (`.go-agent/config.yaml`)

GoAgent uses a configuration file located at `.go-agent/config.yaml` *relative to the directory where the `go-agent` command is executed*. This file allows customization of the LLM models used for specific internal tasks and can override `max_history_length` and `max_process_loops` set by flags.

**Purpose:** This file allows you to customize the specific LLM models used for different internal agent tasks within GoAgent, separately for each supported LLM service (Gemini, Groq, OpenAI). The currently configurable tasks are:
    - `instructions_model`: Used by the instruction agent for understanding the initial change request and structuring commands.
    - `generate_code_model`: Used by the code generation agent for creating or modifying file content.
    - `analysis_model`: Used by the analysis agent for analyzing code, planning changes, determining context files, and answering `/ask` queries.

You can also set `max_history_length` and `max_process_loops` in this file. Values set in the config file take precedence over command-line flags for these two options.

**Automatic Creation:** If the `.go-agent` directory or the `config.yaml` file does not exist in the *current working directory* when GoAgent starts, it will be automatically created with default model configurations and default values for `max_history_length` (100) and `max_process_loops` (5).

**Example Default `config.yaml`:**

```yaml
gemini:
  instructions_model: gemini-2.0-flash-exp
  generate_code_model: gemini-2.0-flash-thinking-exp-01-21
  analysis_model: gemini-2.0-flash-thinking-exp-01-21
groq:
  instructions_model: llama-3.3-70b-versatile
  generate_code_model: qwen-2.5-coder-32b
  analysis_model: llama-3.3-70b-versatile
openai:
  instructions_model: gpt-4.5-preview
  generate_code_model: gpt-4.5-preview
  analysis_model: gpt-4.5-preview
max_history_length: 100
max_process_loops: 5
```

Contributions to GoAgent are welcome! Please feel free to submit pull requests or open issues for bug reports and feature requests.

## License

[MIT License](LICENSE)