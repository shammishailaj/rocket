# AWS S3

## Description

The `aws_s3` provider ease the uploading of artifacts to AWS S3 buckets.

**Note**:  if `access_key_id` or `secret_access_key` is empty, even after environment expanded
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
| `bucket` | `string` | **$AWS_S3_BUCKET** | The S3 bucket to use |
| `local_directory` | `string` | `"."` | The base local directory to upload |
| `remote_directory` | `string` | `"/"` | The base remote directory to upload to |


## Example

```san
# .rocket.san
aws_s3 = {
  region = "$MY_AWS_S3_REGION"
  bucket = "my-bucket"
  remote_directory = "/my/app/directory"
}
```
