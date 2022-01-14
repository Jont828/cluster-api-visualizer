package internal

import (
	"strings"

	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/tree"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterResourceNode struct {
	Name        string                 `json:"name"`
	Kind        string                 `json:"kind"`
	Group       string                 `json:"group"`
	Version     string                 `json:"version"`
	Provider    string                 `json:"provider"`
	UID         string                 `json:"uid"`
	IsVirtual   bool                   `json:"isVirtual"`
	Collapsable bool                   `json:"collapsable"`
	Children    []*ClusterResourceNode `json:"children"`
}

func ConstructClusterResourceTree(defaultClient client.Client, dcOptions client.DescribeClusterOptions) (*ClusterResourceNode, error) {
	objTree, err := defaultClient.DescribeCluster(dcOptions)
	if err != nil {
		return nil, err
	}

	resourceTree := objectTreeToResourceTree(objTree, objTree.GetRoot())

	return resourceTree, nil
}

func objectTreeToResourceTree(objTree *tree.ObjectTree, object ctrlclient.Object) *ClusterResourceNode {
	if object == nil {
		return nil
	}

	group := object.GetObjectKind().GroupVersionKind().Group
	kind := object.GetObjectKind().GroupVersionKind().Kind
	version := object.GetObjectKind().GroupVersionKind().Version

	// fmt.Printf("%s %s %s %s\n", group, kind, version, object.GetObjectKind().GroupVersionKind().String())

	node := &ClusterResourceNode{
		Name:        object.GetName(),
		Kind:        kind,
		Group:       group,
		Version:     version,
		Provider:    group[:strings.IndexByte(group, '.')],
		IsVirtual:   tree.IsVirtualObject(object),
		Collapsable: tree.IsVirtualObject(object),
		Children:    []*ClusterResourceNode{},
		UID:         string(object.GetUID()),
	}

	for _, child := range objTree.GetObjectsByParent(object.GetUID()) {
		node.Children = append(node.Children, objectTreeToResourceTree(objTree, child))
	}

	return node
}
