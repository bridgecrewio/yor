resource "aws_instance" "no_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"
}

resource "aws_instance" "simple_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"

  tags = {
    Name = "my-instance"
  }
}

resource "aws_instance" "rendered_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"

  tags = var.tags
}

resource "aws_instance" "merge_tags" {
  ami           = "some-ami"
  instance_type = "t3.micro"

  tags = merge(var.tags,
  {
    Name = "merged-tags-instance",
    Env  = var.env
  })
}

variable "tags" {
  default = {}
  type = map(string)
}

variable "env" {
  default = "dev"
  type = string
}