provider "aws" {
  region = "us-east-1"
  alias  = "virginia"
}

provider "kubernetes" {
  host                   = module.eks_blueprints.eks_cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks_blueprints.eks_cluster_certificate_authority_data)
  token                  = data.aws_eks_cluster_auth.this.token
}

provider "helm" {
  kubernetes {
    host                   = module.eks_blueprints.eks_cluster_endpoint
    cluster_ca_certificate = base64decode(module.eks_blueprints.eks_cluster_certificate_authority_data)
    token                  = data.aws_eks_cluster_auth.this.token
  }
}

provider "kubectl" {
  apply_retry_count      = 10
  host                   = module.eks_blueprints.eks_cluster_endpoint
  cluster_ca_certificate = base64decode(module.eks_blueprints.eks_cluster_certificate_authority_data)
  load_config_file       = false
  token                  = data.aws_eks_cluster_auth.this.token
}

module "eks_blueprints" {
  source = "github.com/aws-ia/terraform-aws-eks-blueprints?ref=v4.27.0"

  cluster_name = local.name

  vpc_id             = module.vpc.vpc_id
  private_subnet_ids = module.vpc.private_subnets

  cluster_version = local.cluster_version

  managed_node_groups = {

    workload = {
      node_group_name = "managed-workload"
      instance_types  = ["m5.8xlarge"]
      subnet_ids      = module.vpc.private_subnets
      min_size        = 1
      max_size        = 5
      desired_size    = 1
      k8s_labels = {
        node_type = "workload"
      }
    }

    tooling = {
      node_group_name = "managed-tooling"
      instance_types  = ["m5.2xlarge"]
      subnet_ids      = module.vpc.private_subnets
      min_size        = 1
      max_size        = 5
      desired_size    = 1
      k8s_labels = {
        node_type = "tooling"
      }
    }
  }

  platform_teams = {
    admin = {
      users = [
        data.aws_caller_identity.current.arn
      ]
    }
  }

  tags = local.tags
}


module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "3.16.0"

  name = local.name
  cidr = local.vpc_cidr

  azs             = local.azs
  public_subnets  = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k)]
  private_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 10)]

  enable_nat_gateway   = true
  create_igw           = true
  enable_dns_hostnames = true
  single_nat_gateway   = true

  # Manage so we can name
  manage_default_network_acl    = true
  default_network_acl_tags      = { Name = "${local.name}-default" }
  manage_default_route_table    = true
  default_route_table_tags      = { Name = "${local.name}-default" }
  manage_default_security_group = true
  default_security_group_tags   = { Name = "${local.name}-default" }

  public_subnet_tags = {
    "kubernetes.io/cluster/${local.name}" = "shared"
    "kubernetes.io/role/elb"              = "1"
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${local.name}" = "shared"
    "kubernetes.io/role/internal-elb"     = "1"
  }

  tags = local.tags
}

module "kubernetes_addons" {
  source = "github.com/aws-ia/terraform-aws-eks-blueprints?ref=v4.27.0/modules/kubernetes-addons"

  eks_cluster_id = module.eks_blueprints.eks_cluster_id

  enable_argocd         = true
  argocd_manage_add_ons = true

  argocd_applications = {
    addons    = local.addon_application
    workloads = local.workload_application
  }

  argocd_helm_config = {
    set = [
      {
        name  = "server.service.type"
        value = "LoadBalancer"
      }
    ]
  }

  # EKS Addons (AWS Managed)
  enable_amazon_eks_aws_ebs_csi_driver = true
  enable_amazon_eks_coredns            = true
  enable_amazon_eks_kube_proxy         = true
  enable_amazon_eks_vpc_cni            = true

  enable_aws_load_balancer_controller = true
  enable_aws_for_fluentbit            = true
  enable_metrics_server               = true
  enable_cluster_autoscaler           = true

  enable_prometheus      = true
  enable_amazon_eks_adot = true
  amazon_eks_adot_config = {
    most_recent        = true
    kubernetes_version = module.eks_blueprints.eks_cluster_version
    resolve_conflicts  = "OVERWRITE"
  }
}

module "adot_collector_irsa_addon" {
  source = "github.com/aws-ia/terraform-aws-eks-blueprints?ref=v4.21.0/modules/irsa"

  create_kubernetes_namespace       = true
  create_kubernetes_service_account = true
  kubernetes_namespace              = "aws-otel"
  kubernetes_service_account        = "adot-collector"
  irsa_iam_policies                 = ["arn:aws:iam::aws:policy/AmazonPrometheusRemoteWriteAccess"]
  eks_cluster_id                    = module.eks_blueprints.eks_cluster_id
  eks_oidc_provider_arn             = module.eks_blueprints.eks_oidc_provider_arn
}