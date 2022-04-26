package internal

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2/klogr"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/tree"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterResourceNode struct {
	Name        string                 `json:"name"`
	DisplayName string                 `json:"displayName"`
	Kind        string                 `json:"kind"`
	Group       string                 `json:"group"`
	Version     string                 `json:"version"`
	Provider    string                 `json:"provider"`
	UID         string                 `json:"uid"`
	IsVirtual   bool                   `json:"isVirtual"`
	Collapsable bool                   `json:"collapsable"`
	Collapsed   bool                   `json:"collapsed"`
	Ready       bool                   `json:"ready"`
	HasReady    bool                   `json:"hasReady"`
	Children    []*ClusterResourceNode `json:"children"`
}

func ConstructClusterResourceTree(defaultClient client.Client, dcOptions client.DescribeClusterOptions) (*ClusterResourceNode, *HTTPError) {
	objTree, err := defaultClient.DescribeCluster(dcOptions)
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			return nil, &HTTPError{Status: 404, Message: err.Error()}
		}

		return nil, NewInternalError(err)
	}

	resourceTree := objectTreeToResourceTree(objTree, objTree.GetRoot(), true)

	return resourceTree, nil
}

func objectTreeToResourceTree(objTree *tree.ObjectTree, object ctrlclient.Object, groupMachines bool) *ClusterResourceNode {
	log := klogr.New()

	if object == nil {
		return nil
	}

	group := object.GetObjectKind().GroupVersionKind().Group
	kind := object.GetObjectKind().GroupVersionKind().Kind
	version := object.GetObjectKind().GroupVersionKind().Version

	// fmt.Printf("%s %s %s %s\n", group, kind, version, object.GetObjectKind().GroupVersionKind().String())
	provider, err := getProvider(object, group)
	if err != nil {
		log.Error(err, "failed to get provider for object", "kind", kind, "name", object.GetName())
	}

	readyCondition := tree.GetReadyCondition(object)

	node := &ClusterResourceNode{
		Name:        object.GetName(),
		DisplayName: object.GetName(),
		Kind:        kind,
		Group:       group,
		Version:     version,
		Provider:    provider,
		IsVirtual:   tree.IsVirtualObject(object),
		HasReady:    readyCondition != nil,
		Ready:       readyCondition != nil && readyCondition.Status == corev1.ConditionTrue,
		Collapsable: tree.IsVirtualObject(object),
		Collapsed:   false,
		Children:    []*ClusterResourceNode{},
		UID:         string(object.GetUID()),
	}

	children := objTree.GetObjectsByParent(object.GetUID())
	sort.Slice(children, func(i, j int) bool {
		return children[i].GetObjectKind().GroupVersionKind().Kind < children[j].GetObjectKind().GroupVersionKind().Kind
	})

	childTrees := []*ClusterResourceNode{}
	for _, child := range children {
		childTrees = append(childTrees, objectTreeToResourceTree(objTree, child, true))
		// node.Children = append(node.Children, objectTreeToResourceTree(objTree, child, true))
	}

	log.V(3).Info("Node is", "node", node.Kind+"/"+node.Name)
	node.Children = createKindGroupNode(object.GetNamespace(), "Machine", "cluster", childTrees)

	return node
}

// TODO: create map of kinds to group by
// For each kind in map, get count of the kind in children
// If count > 1, create a group node and add children to group node
// Look into adding a striped background for nodes that aren't ready

func createKindGroupNode(namespace string, kind string, provider string, children []*ClusterResourceNode) []*ClusterResourceNode {
	log := klogr.New()

	log.V(2).Info("Starting children are ", "children", nodeArrayNames(children))

	resultChildren := []*ClusterResourceNode{}
	groupNode := &ClusterResourceNode{
		Name:        "",
		DisplayName: "",
		Kind:        kind,
		Group:       "virtual.cluster.x-k8s.io",
		Version:     "v1beta1",
		Provider:    "cluster",
		IsVirtual:   true,
		Collapsable: true,
		Collapsed:   true,
		Children:    []*ClusterResourceNode{},
		HasReady:    false,
		Ready:       false,
		UID:         kind + ": ",
	}

	for _, child := range children {
		if child.Kind == kind {
			groupNode.Children = append(groupNode.Children, child)
			groupNode.UID += child.UID + " "
			if child.HasReady {
				groupNode.HasReady = true
				groupNode.Ready = child.Ready && groupNode.Ready
			}
		} else {
			resultChildren = append(resultChildren, child)
		}
	}

	if len(groupNode.Children) > 1 {
		groupNode.DisplayName = fmt.Sprintf("%d %s", len(groupNode.Children), flect.Pluralize(kind))
		resultChildren = append(resultChildren, groupNode)
	} else {
		resultChildren = append(resultChildren, groupNode.Children...)
	}

	log.V(3).Info("Result children are ", "children", nodeArrayNames(resultChildren))

	return resultChildren
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

func nodeArrayNames(nodes []*ClusterResourceNode) string {
	result := ""
	for _, node := range nodes {
		result += node.Kind + "/" + node.Name + " "
	}

	return result
}
