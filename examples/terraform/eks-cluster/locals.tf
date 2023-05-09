locals {

  name            = basename(path.cwd)
  region          = data.aws_region.current.name
  cluster_version = "1.24"

  vpc_cidr = "10.0.0.0/16"
  azs      = slice(data.aws_availability_zones.available.names, 0, 3)

  #---------------------------------------------------------------
  # ARGOCD ADD-ON APPLICATION
  #---------------------------------------------------------------

  addon_application = {
    path               = "chart"
    repo_url           = "https://github.com/aws-samples/eks-blueprints-add-ons.git"
    add_on_application = true
  }

  #---------------------------------------------------------------
  # ARGOCD WORKLOAD APPLICATION
  #---------------------------------------------------------------
  workload_repo = "https://github.com/aws-samples/amazon-eks-litmus-chaos"

  workload_application = {
    path               = "examples/helm/argocd"
    repo_url           = local.workload_repo
    target_revision    = "main"
    add_on_application = false
    values = {
      labels = {
        myapp = "argocd-apps"
      }
      spec = {
        source = {
          repoURL = local.workload_repo
        }

        blueprint   = "terraform"
        clusterName = local.name
      }
    }
  }
}


