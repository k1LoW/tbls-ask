# tbls-ask

`tbls-ask` is an external subcommand of tbls for asking LLM of the datasource.

## Requirements

- [`tbls`](https://github.com/k1LoW/tbls) ... `tbls-ask` is an external subcommand of `tbls`. You need `tbls` installed first.
- `OPENAI_API_KEY` ... API Key for OpenAI or compatible API.
- `OPENAI_BASE_URL` ... (Optional) Base URL for OpenAI-compatible endpoints (e.g. OpenRouter, LiteLLM, Ollama, Gemini).

## Install

First, install [`tbls`](https://github.com/k1LoW/tbls) if you haven't already:

``` console
$ brew install tbls
# or
$ go install github.com/k1LoW/tbls@latest
```

Then install `tbls-ask`:

### Homebrew

``` console
$ brew install k1LoW/tap/tbls-ask
```

### aqua

``` console
$ aqua g -i k1LoW/tbls-ask
```

### `go install`

``` console
$ go install github.com/k1LoW/tbls-ask@latest
```

Make sure the installed binary (`tbls-ask`) is in your `$PATH`.

## Usage

Once installed, `tbls-ask` is executed via `tbls ask`. By default, it uses OpenAI (`chat-latest`).

### Ask questions about database

``` console
$ tbls ask --dsn 'mysql://user:pass@localhost:3306/wordpress' which is blog comment table?
The blog comment table in the given DDL is `wp_comments`.
```

### Generate SQL query (`--query` / `-q`)

Use `--query` (or `-q`) to output only the executable SQL query.

``` console
$ tbls ask --query --dsn 'json://path/to/wordpress/schema.json' count blog posts per user per month
SELECT
    YEAR(p.post_date) AS `Year`,
    MONTH(p.post_date) AS `Month`,
    u.display_name AS `User`,
    COUNT(p.ID) AS `Post Count`
FROM wp_posts p
INNER JOIN wp_users u ON p.post_author = u.ID
WHERE p.post_type = 'post' AND p.post_status = 'publish'
GROUP BY `Year`, `Month`, `User`
ORDER BY `Year` DESC, `Month` DESC, `User` ASC
```

## Using Other Models

You can switch models using the `--model` (or `-m`) option and `OPENAI_BASE_URL`.

### Gemini

#### Google AI Studio

You can use Gemini via the Google AI Studio OpenAI-compatible endpoint by setting `OPENAI_BASE_URL`:

```console
export OPENAI_BASE_URL="https://generativelanguage.googleapis.com/v1beta/openai/"
export OPENAI_API_KEY=your_gemini_api_key
$ tbls ask --dsn 'mysql://user:pass@localhost:3306/wordpress' --model gemini-flash-latest which is blog comment table?
```

#### Vertex AI (Google Cloud)

You can also use Gemini on Vertex AI via its OpenAI-compatible endpoint with a Google Cloud access token:

```console
export PROJECT_ID="your-gcp-project-id"
export LOCATION="us-central1"
export OPENAI_BASE_URL="https://${LOCATION}-aiplatform.googleapis.com/v1/projects/${PROJECT_ID}/locations/${LOCATION}/endpoints/openapi"
export OPENAI_API_KEY=$(gcloud auth print-access-token)
$ tbls ask --dsn 'mysql://user:pass@localhost:3306/wordpress' --model google/gemini-flash-latest which is blog comment table?
```

### Azure OpenAI

```console
export AZURE_OPENAI_KEY=your_api_key  # or AZURE_OPENAI_API_KEY
export AZURE_OPENAI_ENDPOINT=your_endpoint
export AZURE_OPENAI_MODEL=your_deployment_model
export AZURE_OPENAI_API_VERSION=your_api_version
$ tbls ask --dsn 'mysql://user:pass@localhost:3306/wordpress' --model azure-openai which is blog comment table?
```

Azure OpenAI uses dedicated environment variables (`AZURE_OPENAI_*`) instead of `OPENAI_API_KEY` and `OPENAI_BASE_URL`:

- `AZURE_OPENAI_KEY` or `AZURE_OPENAI_API_KEY` ... (Required) API Key for Azure OpenAI.
- `AZURE_OPENAI_ENDPOINT` ... (Required) Endpoint URL (e.g. `https://<resource>.openai.azure.com`).
- `AZURE_OPENAI_MODEL` ... (Optional) Deployment model name on Azure.
- `AZURE_OPENAI_API_VERSION` ... (Optional) API version (e.g. `2024-02-15-preview`).

### Others (OpenRouter, LiteLLM, Ollama, etc.)

`tbls-ask` supports any OpenAI-compatible API by setting `OPENAI_BASE_URL`.

#### OpenRouter / LiteLLM

```console
export OPENAI_BASE_URL="https://openrouter.ai/api/v1"
export OPENAI_API_KEY=your_openrouter_api_key
$ tbls ask --dsn 'mysql://user:pass@localhost:3306/wordpress' --model anthropic/claude-sonnet-latest which is blog comment table?
```

#### Ollama (Local LLM)

```console
export OPENAI_BASE_URL="http://localhost:11434/v1"
export OPENAI_API_KEY="ollama"
$ tbls ask --dsn 'mysql://user:pass@localhost:3306/wordpress' --model llama3.3 which is blog comment table?
```
