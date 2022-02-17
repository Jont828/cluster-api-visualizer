package internal

import (
	"context"
	"log"
	"strings"

	"k8s.io/client-go/tools/clientcmd/api"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/cluster"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type MultiClusterTreeNode struct {
	Name                   string                  `json:"name"`
	Namespace              string                  `json:"namespace"`
	Icon                   string                  `json:"icon"`
	InfrastructureProvider string                  `json:"infrastructureProvider"`
	IsManagement           bool                    `json:"isManagement"`
	Children               []*MultiClusterTreeNode `json:"children"`
}

// ConstructMultiClusterTree returns a tree representing the workload cluster discovered in the management cluster.
func ConstructMultiClusterTree(clusterClient cluster.Client, k8sConfigClient *api.Config) (*MultiClusterTreeNode, *HTTPError) {
	currentContextName := k8sConfigClient.CurrentContext
	currentContext, ok := k8sConfigClient.Contexts[currentContextName]
	if !ok {
		return nil, &HTTPError{Status: 404, Message: "current context not found"}
	}
	name := currentContext.Cluster
	namespace, err := clusterClient.Proxy().CurrentNamespace()
	if err != nil {
		return nil, NewInternalError(err)
	}

	root := &MultiClusterTreeNode{
		Name:                   name,
		Namespace:              namespace,
		InfrastructureProvider: "",
		Icon:                   getIcon(""),
		Children:               []*MultiClusterTreeNode{},
		IsManagement:           true,
	}

	workloadClusters, err := clusterClient.Proxy().GetResourceNames("cluster.x-k8s.io/v1beta1", "Cluster", []ctrlclient.ListOption{}, "")
	if err != nil {
		if strings.Contains(err.Error(), "no matches for kind") {
			log.Println(err)
			return root, nil
		}
		return nil, NewInternalError(err)
	}

	ctrlClient, err := clusterClient.Proxy().NewClient()
	if err != nil {
		return nil, NewInternalError(err)
	}

	for _, clusterName := range workloadClusters {
		cluster := &clusterv1.Cluster{}
		clusterKey := ctrlclient.ObjectKey{
			Namespace: namespace,
			Name:      clusterName,
		}

		if err := ctrlClient.Get(context.TODO(), clusterKey, cluster); err != nil {
			// TODO: do we want to return a 404 if a workload cluster is not found?
			return nil, NewInternalError(err)
		}
		// Don't get the kubeconfig for now until we use it to find additional clusters.
		// kubeconfig, err := pkgClient.GetKubeconfig(client.GetKubeconfigOptions{
		// 	WorkloadClusterName: clusterName,
		// })
		infraProvider := cluster.Spec.InfrastructureRef.Kind

		workloadCluster := MultiClusterTreeNode{
			Name:                   clusterName,
			Namespace:              namespace,
			InfrastructureProvider: infraProvider,
			Icon:                   getIcon(infraProvider),
			Children:               []*MultiClusterTreeNode{},
			IsManagement:           false,
		}

		root.Children = append(root.Children, &workloadCluster)
	}

	return root, nil
}

func getIcon(provider string) string {
	switch provider {
	case "AzureCluster":
		return "microsoft-azure"
	case "DockerCluster":
		return "docker"
	case "GCPCluster":
		return "google-cloud"
	case "AWSCluster":
		return "aws"
	default:
		return "kubernetes"
	}
}
