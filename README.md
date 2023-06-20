# tbls-ask

`tbls-ask` is an external subcommand of tbls for asking OpenAI using the datasource.

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

## Requirement

- `OPENAI_API_KEY` ... API Key of OpenAI.
