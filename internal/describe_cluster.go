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
	Namespace   string                 `json:"namespace"`
	DisplayName string                 `json:"displayName"`
	Kind        string                 `json:"kind"`
	Group       string                 `json:"group"`
	Version     string                 `json:"version"`
	Provider    string                 `json:"provider"`
	UID         string                 `json:"uid"`
	Collapsible bool                   `json:"collapsible"`
	Collapsed   bool                   `json:"collapsed"`
	Ready       bool                   `json:"ready"`
	Severity    string                 `json:"severity"`
	HasReady    bool                   `json:"hasReady"`
	Children    []*ClusterResourceNode `json:"children"`
}

type ClusterResourceTreeOptions struct {
	GroupMachines                bool
	AddControlPlaneVirtualNode   bool
	KindsToCollapse              map[string]struct{}
	VNodesToInheritChildProvider map[string]struct{}
}

// Note: ObjectReferenceObjects do not have the virtual annotation so we can assume that all virtual objects are collapsible
func ConstructClusterResourceTree(defaultClient client.Client, dcOptions client.DescribeClusterOptions) (*ClusterResourceNode, *HTTPError) {
	objTree, err := defaultClient.DescribeCluster(dcOptions)
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			return nil, &HTTPError{Status: 404, Message: err.Error()}
		}

		return nil, NewInternalError(err)
	}

	treeOptions := ClusterResourceTreeOptions{
		GroupMachines:              true,
		AddControlPlaneVirtualNode: true,
		KindsToCollapse: map[string]struct{}{
			"TemplateGroup":           {},
			"ClusterResourceSetGroup": {},
		},
		VNodesToInheritChildProvider: map[string]struct{}{
			"ClusterResourceSetGroup": {},
			// "WorkerGroup":             {},
		},
	}
	resourceTree := objectTreeToResourceTree(objTree, objTree.GetRoot(), treeOptions)

	return resourceTree, nil
}

func objectTreeToResourceTree(objTree *tree.ObjectTree, object ctrlclient.Object, treeOptions ClusterResourceTreeOptions) *ClusterResourceNode {
	log := klogr.New()

	if object == nil {
		return nil
	}

	group := object.GetObjectKind().GroupVersionKind().Group
	kind := object.GetObjectKind().GroupVersionKind().Kind
	version := object.GetObjectKind().GroupVersionKind().Version

	_, collapsed := treeOptions.KindsToCollapse[kind]
	node := &ClusterResourceNode{
		Name:        object.GetName(),
		DisplayName: getDisplayName(object),
		Kind:        kind,
		Group:       group,
		Version:     version,
		Collapsible: tree.IsVirtualObject(object),
		Collapsed:   collapsed,
		Children:    []*ClusterResourceNode{},
		UID:         string(object.GetUID()),
	}
	if node.Namespace = object.GetNamespace(); node.Namespace == "" {
		node.Namespace = "default"
	}

	children := objTree.GetObjectsByParent(object.GetUID())
	provider, err := getProvider(object, children, treeOptions)
	if err != nil {
		log.Error(err, "failed to get provider for object", "kind", kind, "name", object.GetName())
	}
	node.Provider = provider

	setReadyFields(object, node)

	childTrees := []*ClusterResourceNode{}
	for _, child := range children {
		childTrees = append(childTrees, objectTreeToResourceTree(objTree, child, treeOptions))
	}

	log.V(4).Info("Node is", "node", node.Kind+"/"+node.Name)
	if treeOptions.GroupMachines {
		node.Children = createKindGroupNode(object.GetNamespace(), "Machine", "cluster", childTrees, false)
	} else {
		node.Children = childTrees
	}

	sort.Slice(node.Children, func(i, j int) bool {
		// TODO: make sure this is deterministic!
		if getSortKeys(node.Children[i])[0] == getSortKeys(node.Children[j])[0] {
			return getSortKeys(node.Children[i])[1] < getSortKeys(node.Children[j])[1]
		}
		return getSortKeys(node.Children[i])[0] < getSortKeys(node.Children[j])[0]
	})

	if treeOptions.AddControlPlaneVirtualNode && tree.GetMetaName(object) == "ControlPlane" {
		parent := &ClusterResourceNode{
			Name:        "control-plane-parent",
			Namespace:   object.GetNamespace(),
			DisplayName: "ControlPlane",
			Kind:        kind,
			Provider:    "virtual", // TODO: should this be provider=controlplane or provider=virtual?
			Group:       group,
			Version:     version,
			Collapsible: true,
			Collapsed:   false,
			Children:    []*ClusterResourceNode{node},
			UID:         "control-plane-parent",
		}

		return parent
	}

	return node
}

func getSortKeys(node *ClusterResourceNode) []string {
	if node.Group == "virtual.cluster.x-k8s.io" {
		return []string{node.DisplayName, ""}
	}
	return []string{node.Kind, node.DisplayName}
}

// Find all objects in children with `kind` and create a parent node for them
func createKindGroupNode(namespace string, kind string, provider string, children []*ClusterResourceNode, groupForOne bool) []*ClusterResourceNode {
	log := klogr.New()

	log.V(4).Info("Starting children are ", "children", nodeArrayNames(children))

	resultChildren := []*ClusterResourceNode{}
	groupNode := &ClusterResourceNode{
		Name:        "",
		Namespace:   namespace,
		DisplayName: "",
		Kind:        kind,
		Provider:    provider, // TODO: don't hardcode this
		Collapsible: true,
		Collapsed:   true,
		Children:    []*ClusterResourceNode{},
		HasReady:    false,
		Ready:       true,
		Severity:    "",
		UID:         kind + ": ",
	}

	for _, child := range children {
		if child.Kind == kind {
			groupNode.Group = child.Group
			groupNode.Version = child.Version
			groupNode.Children = append(groupNode.Children, child)
			groupNode.UID += child.UID + " "
			if child.HasReady {
				groupNode.HasReady = true
				groupNode.Ready = child.Ready && groupNode.Ready
				groupNode.Severity = updateSeverityIfMoreSevere(groupNode.Severity, child.Severity)
				// Set severity based on most severe child, i.e. Error > Warning > Info > Success
			}
		} else {
			resultChildren = append(resultChildren, child)
		}
	}

	if len(groupNode.Children) > 1 {
		groupNode.DisplayName = fmt.Sprintf("%d %s", len(groupNode.Children), flect.Pluralize(kind))
		resultChildren = append(resultChildren, groupNode)
	} else if len(groupNode.Children) == 1 && groupForOne {
		groupNode.DisplayName = fmt.Sprintf("1 %s", kind)
		resultChildren = append(resultChildren, groupNode)
	} else {
		resultChildren = append(resultChildren, groupNode.Children...)
	}

	log.V(4).Info("Result children are ", "children", nodeArrayNames(resultChildren))

	return resultChildren
}

func updateSeverityIfMoreSevere(existingSev string, newSev string) string {
	switch {
	case existingSev == "":
		return newSev
	case existingSev == "Info":
		if newSev == "Error" || newSev == "Warning" {
			return newSev
		}
		return existingSev
	case existingSev == "Warning":
		if newSev == "Error" {
			return newSev
		}
		return existingSev
	case existingSev == "Error":
		return existingSev
	}

	return existingSev
}

func getProvider(object ctrlclient.Object, children []ctrlclient.Object, treeOptions ClusterResourceTreeOptions) (string, error) {
	log := klogr.New()

	if tree.IsVirtualObject(object) {
		_, inherit := treeOptions.VNodesToInheritChildProvider[object.GetObjectKind().GroupVersionKind().Kind]
		if !inherit {
			return "virtual", nil
		}
		log.V(4).Info("Aggregating object w/ kind, name, and metaName", "kind", object.GetObjectKind().GroupVersionKind().Kind, "name", object.GetName(), "metaName", tree.GetMetaName(object))

		prev := ""
		for i, child := range children {
			provider, err := lookUpProvider(child)
			if err != nil {
				return "", err
			}
			log.V(4).Info("Child object w/ kind, name, and provider", "kind", object.GetObjectKind().GroupVersionKind().Kind, "name", object.GetName(), "metaName", tree.GetMetaName(object))

			if provider == "virtual" { // Do not inherit virtual provider
				return "virtual", nil
			}
			if i > 0 && provider != prev { // If two children have different providers, don't inherit
				return "virtual", nil
			}
			prev = provider
		}

		return prev, nil
	} else {
		return lookUpProvider(object)
	}
}

func lookUpProvider(object ctrlclient.Object) (string, error) {
	group := object.GetObjectKind().GroupVersionKind().Group
	providerIndex := strings.IndexByte(group, '.')
	if tree.IsVirtualObject(object) {
		return "virtual", nil
	} else if providerIndex > -1 {
		return group[:providerIndex], nil
	} else {
		return "", errors.Errorf("No provider found for object %s of %s \n", object.GetName(), object.GetObjectKind().GroupVersionKind().String())
	}
}

func getDisplayName(object ctrlclient.Object) string {
	metaName := tree.GetMetaName(object)
	displayName := object.GetName()
	if metaName != "" {
		if object.GetName() == "" || tree.IsVirtualObject(object) {
			displayName = metaName
		}
	}

	return displayName
}

func setReadyFields(object ctrlclient.Object, node *ClusterResourceNode) {
	if readyCondition := tree.GetReadyCondition(object); readyCondition != nil {
		node.HasReady = true
		node.Ready = readyCondition.Status == corev1.ConditionTrue
		node.Severity = string(readyCondition.Severity)
	}
}

func nodeArrayNames(nodes []*ClusterResourceNode) string {
	result := ""
	for _, node := range nodes {
		result += node.Kind + "/" + node.Name + " "
	}

	return result
}
