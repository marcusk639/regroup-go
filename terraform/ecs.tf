resource "aws_ecs_service" "regroup_api" {
  name            = "regroup-api"
  task_definition = aws_ecs_task_definition.regroup_api.arn
  launch_type     = "FARGATE"
  cluster         = aws_ecs_cluster.app.id
  desired_count   = 1

  network_configuration {
    assign_public_ip = false

    security_groups = [
      aws_security_group.egress_all.id,
      aws_security_group.ingress_api.id,
    ]

    subnets = [
      aws_subnet.private_d.id,
      aws_subnet.private_e.id,
    ]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.regroup_api.arn
    container_name   = "regroup-api"
    container_port   = "8080"
  }
}

resource "aws_cloudwatch_log_group" "regroup_api" {
  name = "/ecs/regroup-api"
}

resource "aws_ecs_task_definition" "regroup_api" {
  family = "regroup-api"
  container_definitions = jsonencode(
    [
      {
        name  = "regroup-api"
        image = "799934209842.dkr.ecr.us-east-1.amazonaws.com/golang:regroup-service"
        portMappings = [
          {
            containerPort = 8080
          }
        ]
        logConfiguration = {
          logDriver = "awslogs"
          options = {
            awslogs-group         = aws_cloudwatch_log_group.regroup_api.name
            awslogs-region        = "us-west-1"
            awslogs-stream-prefix = "ecs"
          }
        }
      }
    ]
  )
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  network_mode             = "awsvpc"
  execution_role_arn       = aws_iam_role.regroup_api_task_execution_role.arn
  tags = {
    "Name" = "force deployment 1"
  }
}

# Pick up here
resource "aws_iam_role" "regroup_api_task_execution_role" {
  name               = "regroup-api-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume_role.json
}

data "aws_iam_policy_document" "ecs_task_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

data "aws_iam_policy" "ecs_task_execution_role" {
  arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role" {
  role       = aws_iam_role.regroup_api_task_execution_role.name
  policy_arn = data.aws_iam_policy.ecs_task_execution_role.arn
}

resource "aws_ecs_cluster" "app" {
  name = "app"
}

resource "aws_lb_target_group" "regroup_api" {
  name        = "regroup-api"
  port        = 8080
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.app_vpc.id

  health_check {
    enabled = true
    path    = "/health"
  }

  depends_on = [aws_alb.regroup_api]
}

resource "aws_alb" "regroup_api" {
  name               = "regroup-api-lb"
  internal           = false
  load_balancer_type = "application"

  subnets = [
    aws_subnet.public_d.id,
    aws_subnet.public_e.id,
  ]

  security_groups = [
    aws_security_group.http.id,
    aws_security_group.https.id,
    aws_security_group.egress_all.id,
  ]

  depends_on = [aws_internet_gateway.igw]
}

resource "aws_alb_listener" "regroup_api_http" {
  load_balancer_arn = aws_alb.regroup_api.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

output "alb_url" {
  value = "http://${aws_alb.regroup_api.dns_name}"
}

# ecs.tf
resource "aws_acm_certificate" "regroup_api" {
  domain_name       = "regroupgolang.com"
  validation_method = "DNS"
}

output "domain_validations" {
  value = aws_acm_certificate.regroup_api.domain_validation_options
}

resource "aws_alb_listener" "regroup_api_https" {
  load_balancer_arn = aws_alb.regroup_api.arn
  port              = "443"
  protocol          = "HTTPS"
  certificate_arn   = aws_acm_certificate.regroup_api.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.regroup_api.arn
  }
}
