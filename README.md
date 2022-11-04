# Kube Glass

Kube Glass is a project that provides Kubernetes-as-a-Service by taking advantage of the ClusterAPI and VCluster.
It provides a single pane of glass and at-a-glance information that allows users to focus on developing applications rather than managing Kubernetes clusters.

## Components

* APIServer - GraphQL API that allows users to interact with the system.
* Operator - Watches for user changes to Kubernetes resources and manages the creation and updating of managed clusters.
* Web - Web frontend for the APIServer.
