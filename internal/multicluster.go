package internal

import (
	"context"
	"fmt"

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
	Kubeconfig             string                  `json:"kubeconfig"`
}

// ConstructMultiClusterTree returns a tree representing the workload cluster discovered in the management cluster.
func ConstructMultiClusterTree(clusterClient cluster.Client) (*MultiClusterTreeNode, error) {
	name, err := clusterClient.Proxy().CurrentClusterName()
	if err != nil {
		return nil, err
	}
	namespace, err := clusterClient.Proxy().CurrentNamespace()
	if err != nil {
		return nil, err
	}

	contexts, err := clusterClient.Proxy().GetContexts("")
	if err != nil {
		fmt.Printf("Error is %s\n", err)
	} else {
		fmt.Printf("Contexts are %+v\n", contexts)
	}

	root := &MultiClusterTreeNode{
		Name:                   name,
		Namespace:              namespace,
		InfrastructureProvider: "",
		Icon:                   getIcon(""),
		Children:               []*MultiClusterTreeNode{},
		IsManagement:           true,
		Kubeconfig:             "",
	}

	// pkgClient, err := client.New("")
	// if err != nil {
	// 	return nil, err
	// }

	workloadClusters, err := clusterClient.Proxy().GetResourceNames("cluster.x-k8s.io/v1beta1", "Cluster", []ctrlclient.ListOption{}, "")
	if err != nil {
		return nil, err
	}
	ctrlClient, err := clusterClient.Proxy().NewClient()
	if err != nil {
		return nil, err
	}

	for _, clusterName := range workloadClusters {
		cluster := &clusterv1.Cluster{}
		clusterKey := ctrlclient.ObjectKey{
			Namespace: namespace,
			Name:      clusterName,
		}

		if err := ctrlClient.Get(context.TODO(), clusterKey, cluster); err != nil {
			return nil, err
		}
		// Don't get the kubeconfig for now until we use it to find additional clusters.
		// kubeconfig, err := pkgClient.GetKubeconfig(client.GetKubeconfigOptions{
		// 	WorkloadClusterName: clusterName,
		// })
		// if err != nil {
		// 	// Don't return an error if we can't get the kubeconfig since it won't be found until the cluster is ready.
		// 	kubeconfig = ""
		// }
		infraProvider := cluster.Spec.InfrastructureRef.Kind

		workloadCluster := MultiClusterTreeNode{
			Name:                   clusterName,
			Namespace:              namespace,
			InfrastructureProvider: infraProvider,
			Icon:                   getIcon(infraProvider),
			Children:               []*MultiClusterTreeNode{},
			IsManagement:           false,
			// Kubeconfig:             kubeconfig,
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
