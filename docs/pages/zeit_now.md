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
| `env` | `map[string]string` | `{}` | The environment for the deployment |
| `public` | `bool` | `false` | Whether the deployment is public or not |
| `deployment_type` | `string` | `"NPM"` | see the zeit API [documentation](https://zeit.co/api#endpoints/deployments/create-a-new-deployment) |
| `name` | `string` | **$ZEIT_NOW_NAME** | see the zeit API [documentation](https://zeit.co/api#endpoints/deployments/create-a-new-deployment) |
| `force_new` | `bool` | `true` | see the zeit API [documentation](https://zeit.co/api#endpoints/deployments/create-a-new-deployment) |
| `engines` | `map[string]string` | `{}` | see the zeit API [documentation](https://zeit.co/api#endpoints/deployments/create-a-new-deployment) |
| `session_affinity` | `string` | `"ip"` | see the zeit API [documentation](https://zeit.co/api#endpoints/deployments/create-a-new-deployment) |


## Example

```san
# .rocket.san
zeit_now = {
  directory = "dist"
  name = "my-app"
  engines =  { "node" = "^8.0.0" }
  env = { "NODE_ENV" = "production" }
  public = true
}
```
