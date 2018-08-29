# Docker

## Fields

| Field | Type | Default Value | Description |
| ----- | -----| ------------- |------------ |
| `docker_username` | `string` | **$DOCKER_USERNAME** | The require docker username to login to the docker registry |
| `docker_password` | `string` | **$DOCKER_PASSWORD** | The require docker username to login to the docker registry |
| `login` | `bool` | `true` | Whether to `docker login` or not. If set to false, the `docker login` command should be done before `rocket` usage |
| `images` | `[string]` | `[]` | The local docker images to publish|


## Example

```toml
# .rocket.toml
[docker]
docker_username = "$MY_DOCKER_USERNAME"
docker_password = "$MY_DOCKER_PASSWORD"
# images to push
images = [
  "astrocorp/rocket:lastest",
  "my-custom-registry/org/image:my-tag",
  "my-custom-registry/org/image:$VERSION", # can we use env vars here ? to spec
]
```
