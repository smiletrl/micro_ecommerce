// Minikube host access(https://minikube.sigs.k8s.io/docs/handbook/host-access/) is working inside minikube
// VM, but not working inside k8s pods.
// Use the workaround(https://github.com/kubernetes/minikube/issues/8439#issuecomment-799801736).
// Then pods will be able to connect to service in developer's local machine using host `minikube-host.default.svc.cluster.local`.
// For example, a postgres service is created via local machine(a mac) docker container, then k8s pods can connect to local machine's
// postgres service using above host.
resource "kubernetes_service" "host" {
    metadata {
        name = "minikube-host"

        labels = {
            app = "minikube-host"
        }

        namespace = "default"
    }

    spec {
        port {
            name = "app"
            port = 8082
        }
        cluster_ip = "None"
    }
}

resource "kubernetes_endpoints" "host" {
    metadata {
        name = "minikube-host"
        namespace = "default"
    }
    subset {
        address {
            // This ip comes from command: minikube ssh 'grep host.minikube.internal /etc/hosts | cut -f1'
            // This ip can also be used as service host at k8s container directly, instead of `minikube-host.default.svc.cluster.local`.
            ip = "192.168.65.2"
        }

        port {
            name     = "app"
            port     = 8082
        }
    }
}
