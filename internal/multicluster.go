package internal

import (
	"github.com/Azure/go-autorest/autorest/to"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
)

type MultiClusterTreeNode struct {
	Name                   string `json:"name"`
	Namespace              string
	Icon                   string `json:"icon"`
	InfrastructureProvider string
	IsVirtual              bool
	Children               []*MultiClusterTreeNode `json:"children"`
	Kubeconfig             string
}

func ConstructMultiClusterTree(tree *client.MultiClusterTree) *MultiClusterTreeNode {
	if tree == nil {
		return nil
	}

	node := &MultiClusterTreeNode{
		Name:      tree.Name,
		Namespace: tree.Namespace,
		Icon:      getIcon(to.String(tree.InfrastructureProvider)),
		Children:  []*MultiClusterTreeNode{},
		IsVirtual: false,
	}
	if tree.Kubeconfig != nil {
		node.Kubeconfig = *tree.Kubeconfig
	} else {
		node.Kubeconfig = ""
	}
	if tree.InfrastructureProvider != nil {
		node.InfrastructureProvider = *tree.InfrastructureProvider
	} else {
		node.InfrastructureProvider = ""
	}

	for _, child := range tree.WorkloadClusters {
		node.Children = append(node.Children, ConstructMultiClusterTree(child))
	}
	return node
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
