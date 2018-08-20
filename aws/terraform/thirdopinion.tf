resource "aws_security_group" "allow_ssh_alex" {
  name        = "allow_ssh_alex"
  description = "Allow ssh from alex house"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["5.150.231.78/32"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
      from_port = 0
      to_port = 0
      protocol = "-1"
      cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "allow_ssh_alex"
  }
}

resource "aws_instance" "thirdopinion" {
  ami           = "ami-0ebee003d8eab75a4"
  instance_type = "t2.micro"
  key_name = "id_rsa"
  security_groups = [
    "allow_ssh_alex"
  ]
}

