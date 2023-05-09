data "aws_partition" "current" {}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

data "aws_availability_zones" "available" {
  state = "available"
}

locals {
  amp_ws_region   = coalesce(var.managed_prometheus_workspace_region, data.aws_region.current.name)
  amp_ws_id       = aws_prometheus_workspace.this.id
  amp_ws_endpoint = "https://aps-workspaces.${local.amp_ws_region}.amazonaws.com/workspaces/${local.amp_ws_id}/"

  amg_ws_endpoint = "https://${module.managed_grafana.workspace_endpoint}"
  amg_ws_id       = split(".", module.managed_grafana.workspace_endpoint)[0]

  name = "amazon-eks-chaos-demo"
}
