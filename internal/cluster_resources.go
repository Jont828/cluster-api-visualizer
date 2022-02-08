package internal

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
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
	provider, err := getProvider(object, group)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	node := &ClusterResourceNode{
		Name:        object.GetName(),
		Kind:        kind,
		Group:       group,
		Version:     version,
		Provider:    provider,
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

func getProvider(object ctrlclient.Object, group string) (string, error) {
	providerIndex := strings.IndexByte(group, '.')
	if tree.IsVirtualObject(object) {
		return "virtual", nil
	} else if providerIndex > -1 {
		return group[:providerIndex], nil
	} else {
		return "", errors.Errorf("No provider found for object %s of %s \n", object.GetName(), object.GetObjectKind().GroupVersionKind().String())
	}
}
