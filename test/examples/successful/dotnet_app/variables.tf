variable "key_path" {
  type = "string"
}

variable "key_name" {
  type = "string"
}

variable "access_key"{
  type = "string"
}

variable "region" {
  default = "us-east-1"
}

variable "win_amis" {
  type = "map"
  default = {
    us-east-1 = "ami-027a14492d667b8f5"
    us-west-2 = "ami-9f5efbff"
    eu-west-1 = "ami-7ac78809"
  }
}

variable "win_username" {
  default = "Terraform"
}
variable "win_password" { }