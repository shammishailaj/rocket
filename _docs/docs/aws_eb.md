# AWS Elastic Beanstalk

## Description

The `aws_eb` provider is a helper to deploy to the [AWS Elastic Beanstalk](https://aws.amazon.com/fr/elasticbeanstalk/)
service.

It does not aim to have all the `eb` CLI functionalities but to provide a standard and easy way
to create releases.


## Fields


| Field | Type | Default Value | Description |
| ----- | -----| ------------- |------------ |
| `access_key_id` | `string` | **$AWS_ACCESS_KEY_ID** | The AWS access key ID |
| `secret_access_key` | `string` | **$AWS_SECRET_ACCESS_KEY** | The AWS secret access key |
| `region` | `string` | **$AWS_REGION** | The AWS region to use |
| `application` | `string` | **$AWS_EB_APPLICATION** | The EB application to use |

## Example

```toml
# .rocket.toml
application = "myapp"
environment = "myapp-production"
directory = "."
```
