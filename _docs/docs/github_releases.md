# GitHub releases

## Fields

| Field             | Type |Default Value | Description |
| ------------------| ---- | ------------ | ----------- |
| `name` | `string` | **$ROCKET_LAST_TAG** | The release's name |
| `body` | `string` | `""` | The release's body | 
| `prerelease` | `bool` | `false` | Identify the release as a prerelease |
| `repo` | `string` | **$ROCKET_GIT_REPO** | The GitHub repo to release |
| `api_key` | `string` | **$GITHUB_API_KEY** | The required GitHub API key |
| `assets` | `[string]` | `[]` | The assets to upload following the [`go` glob pattern](https://golang.org/pkg/path/filepath/#Match) |
| `tag` | `string` | **$ROCKET_LAST_TAG** | The `git` tag to release |


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
