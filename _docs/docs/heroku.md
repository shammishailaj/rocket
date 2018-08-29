# Heroku

## Fields

| Field             | Default Value | Description |
| --------------------- | -------| --- |
| `api_key` | **$HEROKU_API_KEY** | The required Heroku API key |
| `app` | **$HEROKU_APP** | The Heroku app to deploy |
| `directory` | `"."` | The directory of your project (which will be tar gzipped and uploaded) |
| `version` | **$ROCKET_COMMIT_HASH** | The version of the app to release |


## Example

```toml
# .rocket.toml
[heroku]
app = "my-awesome-heroku-app"
api_key = "$HEROKU_TOKEN"
directory = "."
```
