package internal

import (
	"fmt"
	"strings"

	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/tree"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterResourceNode struct {
	Name      string                 `json:"name"`
	Kind      string                 `json:"kind"`
	Group     string                 `json:"group"`
	Version   string                 `json:"version"`
	Provider  string                 `json:"provider"`
	UID       string                 `json:"uid"`
	IsVirtual bool                   `json:"isVirtual"`
	Children  []*ClusterResourceNode `json:"children"`
}

func ConstructClusterResourceTree(clusterName string) (*ClusterResourceNode, error) {
	objTree, err := getObjectTree(clusterName)
	if err != nil {
		return nil, err
	}

	resourceTree := objectTreeToResourceTree(objTree, objTree.GetRoot())

	return resourceTree, nil
}

func getObjectTree(name string) (*tree.ObjectTree, error) {
	kubeconfig := "/home/jonathan/.kube/config"
	kubeconfigContext := ""
	namespace := "default"
	showOtherConditions := ""
	showAllResources := true
	showMachineSets := false
	c, err := client.New("")
	if err != nil {
		return nil, err
	}

	return c.DescribeCluster(client.DescribeClusterOptions{
		Kubeconfig:          client.Kubeconfig{Path: kubeconfig, Context: kubeconfigContext},
		Namespace:           namespace,
		ClusterName:         name,
		ShowOtherConditions: showOtherConditions,
		ShowAllResources:    showAllResources,
		ShowMachineSets:     showMachineSets,
	})
}

func objectTreeToResourceTree(objTree *tree.ObjectTree, object ctrlclient.Object) *ClusterResourceNode {
	if object == nil {
		return nil
	}

	group := object.GetObjectKind().GroupVersionKind().Group
	kind := object.GetObjectKind().GroupVersionKind().Kind
	version := object.GetObjectKind().GroupVersionKind().Version
	fmt.Printf("Name: %s/%s\n", kind, object.GetName())
	fmt.Printf("Group: %s\n", group)
	fmt.Printf("Version: %s\n\n", version)
	node := &ClusterResourceNode{
		Name:      object.GetName(),
		Kind:      kind,
		Group:     group,
		Version:   version,
		Provider:  group[:strings.IndexByte(group, '.')],
		IsVirtual: tree.IsVirtualObject(object),
		Children:  []*ClusterResourceNode{},
		UID:       string(object.GetUID()),
	}
	for _, child := range objTree.GetObjectsByParent(object.GetUID()) {
		node.Children = append(node.Children, objectTreeToResourceTree(objTree, child))
	}
	return node
}
