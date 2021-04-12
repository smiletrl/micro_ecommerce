resource "null_resource" "dashboard_download" {

  provisioner "local-exec" {
    command = <<EOT
        echo "Downloading dashboard config"
        curl -L https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0/aio/deploy/recommended.yaml -o dashboard.yaml
    EOT
  }
}

data "kubectl_path_documents" "dashboard_manifests" {
    pattern = "./dashboard.yaml"
    depends_on = [null_resource.dashboard_download]
}

// This controller file has multiple configs.
// See https://github.com/gavinbunney/terraform-provider-kubectl/issues/52
resource "kubectl_manifest" "dashboard_install" {
    count     = length(data.kubectl_path_documents.dashboard_manifests.documents)
    yaml_body = element(data.kubectl_path_documents.dashboard_manifests.documents, count.index)
}

resource "kubernetes_service_account" "admin-user" {
    depends_on = [
        kubectl_manifest.dashboard_install
    ]

    metadata {
        name = "admin-user"
        namespace = "kubernetes-dashboard"
    }
}

resource "kubernetes_cluster_role_binding" "admin-user" {
    depends_on = [
        kubectl_manifest.dashboard_install
    ]

    metadata {
        name = "admin-user"
    }
    role_ref {
        api_group = "rbac.authorization.k8s.io"
        kind      = "ClusterRole"
        name      = "cluster-admin"
    }
    subject {
        kind      = "ServiceAccount"
        name      = "admin-user"
        namespace = "kubernetes-dashboard"
    }
}
