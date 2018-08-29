# GitHub releases

## Fields

| Field             | Default Value | Description |
| --------------------- | -------| --- |
| `name` | **$ROCKET_LAST_TAG** | The release's name |
| `body` | `""` | The release's body | 
| `prerelease` | `false` | Identify the release as a prerelease |
| `repo` | **$ROCKET_GIT_REPO** | The GitHub repo to release |
| `api_key` | **$GITHUB_API_KEY** | The required GitHub API key |
| `assets` | `[]` | The assets to upload following the [`go` glob pattern](https://golang.org/pkg/path/filepath/#Match) |
| `tag` | **$ROCKET_LAST_TAG** | The `git` tag to release |


## Example

```toml
# .rocket.toml
[github_releases]
api_key = "$GITHUB_TOKEN"
assets = [
  "dist/*.zip",
  "dist/rocket_*_sha512sums.txt",
]
```
