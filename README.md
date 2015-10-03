# terraform-provider-amaxo-aws

Wraps the builtin `aws_s3_bucket_object` resource to permit the object source to be specified
as a URL. If the `source` attribute parses as a URL, this plugin will download from that URL to
a temporary file before uploading to the specified bucket.

Install by including the following in your ~/.terraformrc

    providers {
        aws = "terraform-provider-amaxo-aws"
    }
