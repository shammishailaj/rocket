# AWS Elastic Beanstalk

## Description

The `aws_eb` provider is a helper to deploy to the [AWS Elastic Beanstalk](https://aws.amazon.com/fr/elasticbeanstalk/)
service.

It does not aim to have all the `eb` CLI functionalities but to provide a standard and easy way
to create releases.

It will proceed the following steps:

1. ZIP the given directory
2. upload the `bundle.zip` to the given S3 bucket
3. create a new Application version

**Note**:  if `access_key_id` or `secret_access_key` is emtpy, even after environment expanded
and default values filled, the `aws_s3` provider will use the *shared credentials file* (`~/.aws/credentials`)
and if empty the *EC2 instance role credentials*.

The recommend way is by setting the **$AWS_ACCESS_KEY_ID** and **$AWS_SECRET_ACCESS_KEY** environment
variables.



## Fields


| Field | Type | Default Value | Description |
| ----- | -----| ------------- |------------ |
| `access_key_id` | `string` | **$AWS_ACCESS_KEY_ID** | The AWS access key ID |
| `secret_access_key` | `string` | **$AWS_SECRET_ACCESS_KEY** | The AWS secret access key |
| `region` | `string` | **$AWS_REGION** | The AWS region to use |
| `application` | `string` | **$AWS_EB_APPLICATION** | The EB application to use |
| `environment` | `string` | **$AWS_EB_ENVIRONMENT** | The EB environment to use |
| `s3_bucket` | `string` | **$AWS_S3_BUCKET** | The S3 bucket to upload the bundle to (MUST be the same region as the `eb` application) |
| `version` | `string` | **$ROCKET_COMMIT_HASH** | The version of the applicaiton to release |
| `directory` | `string` | `"."` | The directory of your project (files will be zipped and uploaded) |
| `s3_key` | `string` | /**${AWS_EB_APPLICATION}**\_**${AWS_EB_ENVIRONMENT}**\_**${ROCKET_COMMIT_HASH}**.zip | The S3 key to upload the bundle to |

## Example

```san
# .rocket.san
application = "myapp"
environment = "myapp-production"
s3_bucket = "mybucket"
```
