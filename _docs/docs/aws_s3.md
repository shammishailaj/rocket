# AWS S3

## Description

The `aws_s3` provider ease the uploading of artifacts to AWS S3 buckets.

## Fields

| Field | Type | Default Value | Description |
| ----- | -----| ------------- |------------ |
| `region` | `string` | **$AWS_REGION** | The AWS region to use |
| `access_key_id` | `string` | **AWS_ACCESS_KEY_ID** | The AWS access key ID |
| `secret_access_key` | `string` | **AWS_SECRET_ACCESS_KEY** | The AWS secret access key |
| `bucket` | `string` | **AWS_S3_BUCKET** | The S3 bucket to use |
| `local_dir` | `string` | `"."` | The base local directory to upload |
| `remote_dir` | `string` | `"/"` | The base remote directory to upload to |


## Example

```toml
# .rocket.toml
[aws_s3]
```
