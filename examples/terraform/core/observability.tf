resource "aws_prometheus_workspace" "this" {
  alias = local.name
  tags  = var.tags
}

resource "aws_prometheus_alert_manager_definition" "this" {
  workspace_id = local.amp_ws_id

  definition = <<EOF
alertmanager_config: |
    route:
      receiver: 'default'
    receivers:
      - name: 'default'
EOF
}

module "managed_grafana" {
  source  = "terraform-aws-modules/managed-service-grafana/aws"
  version = "~> 1.3"

  name                      = local.name
  stack_set_name            = local.name
  data_sources              = ["PROMETHEUS"]
  associate_license         = false
  notification_destinations = ["SNS"]
  account_access_type       = "CURRENT_ACCOUNT"
  authentication_providers  = ["AWS_SSO"]
  permission_type           = "SERVICE_MANAGED"

  tags = var.tags
}


resource "aws_grafana_workspace_api_key" "this" {
  key_name        = "api-key-admin"
  key_role        = "ADMIN"
  seconds_to_live = 2592000
  workspace_id    = local.amg_ws_id
}

provider "grafana" {
  url  = local.amg_ws_endpoint
  auth = aws_grafana_workspace_api_key.this.key
}

resource "grafana_data_source" "amp" {
  type       = "prometheus"
  name       = "${local.name}-prometheus"
  is_default = true
  url        = local.amp_ws_endpoint
  uid        = "amp-ds"

  json_data {
    http_method     = "GET"
    sigv4_auth      = true
    sigv4_auth_type = "workspace-iam-role"
    sigv4_region    = local.amp_ws_region
  }
}

resource "grafana_folder" "this" {
  title = "Amazon EKS Chaos Demo Dashboards"
}

resource "grafana_dashboard" "health" {
  config_json = file("service-dashboard.json")
  overwrite   = true
  folder      = grafana_folder.this.id
}