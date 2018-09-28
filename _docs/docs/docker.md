# Docker

## Description

The `docker` provider can be used to push image to any docker resgitry. You can login within `rocket`
or it can be done in a previous step in the pipeline.

Not that because of the client/server nature of `docker`, it cannot be directly embeded in `rocket`
so this provider is just a wrapper around the `docker` CLI which should be installed.

## Fields

| Field | Type | Default Value | Description |
| ----- | -----| ------------- |------------ |
| `username` | `string` | **$DOCKER_USERNAME** | The require docker username to login to the docker registry |
| `password` | `string` | **$DOCKER_PASSWORD** | The require docker username to login to the docker registry |
| `login` | `bool` | `true` | Whether to `docker login` or not. If set to false, the `docker login` command should be done before `rocket` usage |
| `images` | `[string]` | `[]` | The local docker images to publish|


## Example

```san
# .rocket.san
[docker]
username = "$MY_DOCKER_USERNAME"
password = "$MY_DOCKER_PASSWORD"
# images to push
images = [
  "astrocorp/rocket:lastest",
  "my-custom-registry/org/image:my-tag",
  "my-custom-registry/org/image:$VERSION", # we use env vars here
]
```
