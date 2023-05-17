#---------------------------------------------------------------
# ECR Resources
#---------------------------------------------------------------

resource "aws_ecr_repository" "like-service" {
  name                 = "like-service"
  image_tag_mutability = "IMMUTABLE"
  force_delete         = true
}

resource "aws_ecr_repository" "counter-service" {
  name                 = "counter-service"
  image_tag_mutability = "IMMUTABLE"
  force_delete         = true
}

resource "aws_ecr_repository" "ui-app" {
  name                 = "ui-app"
  image_tag_mutability = "IMMUTABLE"
  force_delete         = true
}