provider "aws" {
  alias                   = "localstack"
  region                  = "us-east-1"

  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    dynamodb = "http://localhost:4566"
  }
}

resource "aws_dynamodb_table" "jwk_keys" {
  provider = aws.localstack

  name           = "jwk_keys"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "kid"

  attribute {
    name = "kid"
    type = "S"
  }

  attribute {
    name = "status"
    type = "S"
  }

  global_secondary_index {
    name            = "status-index"
    hash_key        = "status"
    projection_type = "ALL"
  }

  ttl {
    attribute_name = "not_after_epoch"
    enabled        = true
  }

  tags = {
    Name = "jwk_keys"
  }
}
