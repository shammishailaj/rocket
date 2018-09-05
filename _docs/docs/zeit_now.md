# Zeit Now

## Description

The `zeit_now` provider ease the deployment to Now.

It follows the below steps:
1. upload all files with the API ([https://zeit.co/api#endpoints/deployments/upload-deployment-files](https://zeit.co/api#endpoints/deployments/upload-deployment-files))
2. create a new deployment with the API ([https://zeit.co/api#endpoints/deployments/create-a-new-deployment](https://zeit.co/api#endpoints/deployments/create-a-new-deployment))

## Fields

| Field | Type | Default Value | Description |
| ----- | -----| ------------- |------------ |
| `token` | `string` | **$ZEIT_TOKEN** | The zeit token to use |
| `directory` | `string` | `"."` | The directory to upload |


## Example

```toml
# .rocket.toml
[zeit_now]
directory = "dist"
```
