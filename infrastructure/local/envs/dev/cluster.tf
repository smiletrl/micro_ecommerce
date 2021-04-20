resource "kubernetes_namespace" "env" {
  metadata {
    labels = {
    }

    name = var.env
  }
}

// istio support
resource "kubectl_manifest" "istio_manifest_gateway" {
    yaml_body = templatefile("./gateway.yaml",
        {env = var.env}
    )
}

// file virtual_services.yaml includes many items to be created.
data "kubectl_path_documents" "istio_files_virtual_services" {
    pattern = "./virtual_services.yaml"
    vars = {env = var.env}
}

resource "kubectl_manifest" "istio_manifest_virtual_services" {
    count = length(data.kubectl_path_documents.istio_files_virtual_services.documents)
    yaml_body = element(data.kubectl_path_documents.istio_files_virtual_services.documents, count.index)
}

resource "kubernetes_service" "cart" {
    metadata {
        name = "cart"

        labels = {
            app = "cart"
        }

        namespace = var.env
    }

    spec {
        selector = {
            app = "cart"
        }

        port {
            port        = 1325
            target_port = 1325
        }

        type = "NodePort"
    }
}

resource "kubernetes_service" "customer" {
    metadata {
        name = "customer"

        labels = {
            app = "customer"
        }

        namespace = var.env
    }

    spec {
        selector = {
            app = "customer"
        }

        port {
            port        = 1325
            target_port = 1325
        }

        type = "NodePort"
    }
}

resource "kubernetes_service" "product" {
    metadata {
        name = "product"

        labels = {
            app = "product"
        }

        namespace = var.env
    }

    spec {
        selector = {
            app = "product"
        }

        port {
            name        = "rest-api"
            port        = 1325
            target_port = 1325
        }

        port {
            name        = "grpc"
            port        = 50051
            target_port = 50051
        }

        type = "NodePort"
    }
}

resource "kubernetes_service" "order" {
    metadata {
        name = "order"

        labels = {
            app = "order"
        }

        namespace = var.env
    }

    spec {
        selector = {
            app = "order"
        }

        port {
            name        = "rest-api"
            port        = 1325
            target_port = 1325
        }

        type = "NodePort"
    }
}

resource "kubernetes_service" "payment" {
    metadata {
        name = "payment"

        labels = {
            app = "payment"
        }

        namespace = var.env
    }

    spec {
        selector = {
            app = "payment"
        }

        port {
            name        = "rest-api"
            port        = 1325
            target_port = 1325
        }

        type = "NodePort"
    }
}

resource "kubernetes_deployment" "cart" {
    metadata {
        name = "cart"

        annotations = {
        }

        namespace = var.env
    }

    spec {
        progress_deadline_seconds = 6000

        replicas = 2

        selector {
            match_labels = {
                app = "cart"
            }
        }
    
        template {
            metadata {
                name = "cart"

                labels = {
                    app = "cart"
                }
            }
        
            spec {
                container {
                    name = "cart"
                    image = "${var.docker_registry}cart:${var.env}"
                    image_pull_policy = "Always"

                    liveness_probe {
                        http_get {
                            path = "/health"
                            port = "1325"
                        }

                        initial_delay_seconds = 1
                        period_seconds = 1
                    }

                    env {
                        name = "ENV"
                        value = var.env
                    }
                    env {
                        name = "STAGE"
                        value = var.stage
                    }
                }
            }
        }
    }

    timeouts {
        create = "10m"
        update = "10m"
        delete = "10m"
    }
}

resource "kubernetes_deployment" "customer" {
    metadata {
        name = "customer"

        annotations = {
        }

        namespace = var.env
    }

    spec {
        progress_deadline_seconds = 6000

        replicas = 2

        selector {
            match_labels = {
                app = "customer"
            }
        }
    
        template {
            metadata {
                name = "customer"

                labels = {
                    app = "customer"
                }
            }
        
            spec {
                container {
                    name = "customer"
                    image = "${var.docker_registry}customer:${var.env}"
                    image_pull_policy = "Always"

                    liveness_probe {
                        http_get {
                            path = "/health"
                            port = "1325"
                        }

                        initial_delay_seconds = 1
                        period_seconds = 1
                    }

                    env {
                        name = "ENV"
                        value = var.env
                    }
                    env {
                        name = "STAGE"
                        value = var.stage
                    }
                }
            }
        }
    }

    timeouts {
        create = "10m"
        update = "10m"
        delete = "10m"
    }
}

resource "kubernetes_deployment" "product" {
    metadata {
        name = "product"

        annotations = {
        }

        namespace = var.env
    }

    spec {
        progress_deadline_seconds = 6000

        replicas = 2

        selector {
            match_labels = {
                app = "product"
            }
        }
    
        template {
            metadata {
                name = "product"

                labels = {
                    app = "product"
                }
            }
        
            spec {
                container {
                    name = "product"
                    image = "${var.docker_registry}product:${var.env}"
                    image_pull_policy = "Always"

                    liveness_probe {
                        http_get {
                            path = "/health"
                            port = "1325"
                        }

                        initial_delay_seconds = 1
                        period_seconds = 1
                    }

                    env {
                        name = "ENV"
                        value = var.env
                    }
                    env {
                        name = "STAGE"
                        value = var.stage
                    }
                }
            }
        }
    }

    timeouts {
        create = "10m"
        update = "10m"
        delete = "10m"
    }
}

resource "kubernetes_deployment" "order" {
    metadata {
        name = "order"

        annotations = {
        }

        namespace = var.env
    }

    spec {
        progress_deadline_seconds = 6000

        replicas = 2

        selector {
            match_labels = {
                app = "order"
            }
        }
    
        template {
            metadata {
                name = "order"

                labels = {
                    app = "order"
                }
            }
        
            spec {
                container {
                    name = "order"
                    image = "${var.docker_registry}order:${var.env}"
                    image_pull_policy = "Always"

                    liveness_probe {
                        http_get {
                            path = "/health"
                            port = "1325"
                        }

                        initial_delay_seconds = 1
                        period_seconds = 1
                    }

                    env {
                        name = "ENV"
                        value = var.env
                    }
                    env {
                        name = "STAGE"
                        value = var.stage
                    }
                }
            }
        }
    }

    timeouts {
        create = "10m"
        update = "10m"
        delete = "10m"
    }
}

resource "kubernetes_deployment" "payment" {
    metadata {
        name = "payment"

        annotations = {
        }

        namespace = var.env
    }

    spec {
        progress_deadline_seconds = 6000

        replicas = 2

        selector {
            match_labels = {
                app = "payment"
            }
        }
    
        template {
            metadata {
                name = "payment"

                labels = {
                    app = "payment"
                }
            }
        
            spec {
                container {
                    name = "payment"
                    image = "${var.docker_registry}payment:${var.env}"
                    image_pull_policy = "Always"

                    liveness_probe {
                        http_get {
                            path = "/health"
                            port = "1325"
                        }

                        initial_delay_seconds = 1
                        period_seconds = 1
                    }

                    env {
                        name = "ENV"
                        value = var.env
                    }
                    env {
                        name = "STAGE"
                        value = var.stage
                    }
                }
            }
        }
    }

    timeouts {
        create = "10m"
        update = "10m"
        delete = "10m"
    }
}

