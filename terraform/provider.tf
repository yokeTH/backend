# export AWS_ACCESS_KEY_ID=
# export AWS_SECRET_ACCESS_KEY=
# export AWS_EC2_METADATA_DISABLED=false
provider "aws" {
  region                      = "ap-southeast-7"
  skip_credentials_validation = var.env == "dev"
  skip_metadata_api_check     = var.env == "dev"
  skip_requesting_account_id  = var.env == "dev"
}

# export AWS_ACCESS_KEY_ID=test
# export AWS_SECRET_ACCESS_KEY=test
# export AWS_EC2_METADATA_DISABLED=true
provider "aws" {
  alias  = "localstack"
  region = "us-east-1"

  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    dynamodb = "http://localhost:4566"
  }
}
