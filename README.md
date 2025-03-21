# tbls-ask

`tbls-ask` is an external subcommand of tbls for asking LLM of the datasource.

## Usage

tbls-ask is provided as an external subcommand of [tbls](https://github.com/k1LoW/tbls).

### Ask OpenAI

``` console
$ tbls ask --dsn 'mysql://user:pass@localhost:3306/wordpress' which is blog comment table?
The blog comment table in the given DDL is `wp_comments`.
```

### Ask OpenAI for query

Use `--query` option.

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

## Ask Gemini

Add an option `--model` for asking Gemini.

```console
$ tbls ask --dsn 'mysql://user:pass@localhost:3306/wordpress' --model gemini-pro which is blog comment table?
```

## Ask Azure OpenAI
```console
export AZURE_OPENAI_API_KEY=your_api_key
export AZURE_OPENAI_ENDPOINT=your_endpoint
export AZURE_OPENAI_MODEL=your_model
export AZURE_OPENAI_API_VERSION=your_api_version
$ tbls ask --dsn 'mysql://user:pass@localhost:3306/wordpress' --model azure-openai which is blog comment table?
```

## Requirement

Either OpenAI or Gemini API key is required.

- `OPENAI_API_KEY` ... API Key of OpenAI.
- `GEMINI_API_KEY` ... API Key of Gemini.
