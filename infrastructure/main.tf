terraform {
  required_providers {
    kubectl = {
      source = "gavinbunney/kubectl"
      version = "~> 1.10.0"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "~> 1.13.3"
    }
  }
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
}

provider "kubectl" {
  load_config_file = true
}
