package internal

import (
	"context"
	"fmt"
	"sort"
	"strings"

	visualizerv1 "github.com/Jont828/cluster-api-visualizer/api/v1"
	"github.com/gobuffalo/flect"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/tree"
	"sigs.k8s.io/cluster-api/controllers/external"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	maxGroupSize = 3

	addonsAmountToCollapse = 4
)

// ClusterResourceNode represents a node in the Cluster API resource tree and is used to configure the frontend with additional
// options like collapsibility and provider.
type ClusterResourceNode struct {
	Name            string                 `json:"name"`
	Namespace       string                 `json:"namespace"`
	DisplayName     string                 `json:"displayName"`
	Kind            string                 `json:"kind"`
	Group           string                 `json:"group"`
	Version         string                 `json:"version"`
	Provider        string                 `json:"provider"`
	UID             string                 `json:"uid"`
	CollapseWithTab bool                   `json:"collapseWithTab"`
	CollapseOnClick bool                   `json:"collapseOnClick"`
	Collapsible     bool                   `json:"collapsible"`
	Collapsed       bool                   `json:"collapsed"`
	Ready           bool                   `json:"ready"`
	Status          string                 `json:"status"` // Used for calculating groups, `Ready` is used for displaying the status of the node.
	Severity        string                 `json:"severity"`
	HasReady        bool                   `json:"hasReady"`
	Reason          string                 `json:"reason"`
	Children        []*ClusterResourceNode `json:"children"`
	IsOverflowNode  bool                   `json:"isOverflowNode"` // Used to indicate that this node is an overflow node, e.g. when there are too many machines in a group.
}

type ClusterResourceTreeOptions struct {
	GroupMachines                bool
	AddControlPlaneVirtualNode   bool
	KindsToCollapse              map[string]struct{}
	VNodesToInheritChildProvider map[string]struct{}
	providerTypeOverrideMap      map[string]string
}

// ConstructClusterResourceTree returns a tree with nodes representing the Cluster API resources in the Cluster.
// Note: ObjectReferenceObjects do not have the virtual annotation so we can assume that all virtual objects are collapsible
func ConstructClusterResourceTree(ctx context.Context, defaultClient client.Client, runtimeClient ctrlclient.Client, dcOptions client.DescribeClusterOptions) (*ClusterResourceNode, *HTTPError) {
	objTree, err := defaultClient.DescribeCluster(ctx, dcOptions)
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			return nil, &HTTPError{Status: 404, Message: err.Error()}
		}

		return nil, NewInternalError(err)
	}

	treeOptions := ClusterResourceTreeOptions{
		GroupMachines: true,
		KindsToCollapse: map[string]struct{}{
			"TemplateGroup":           {},
			"ClusterResourceSetGroup": {},
			"Machine":                 {},
		},
		VNodesToInheritChildProvider: map[string]struct{}{
			"ClusterResourceSetGroup": {},
			// "WorkerGroup":             {},
		},
	}

	overrides, err := injectCustomResourcesToObjectTree(ctx, runtimeClient, dcOptions, objTree)
	if err != nil {
		return nil, NewInternalError(err)
	}
	treeOptions.providerTypeOverrideMap = overrides

	resourceTree, err := objectTreeToResourceTree(ctx, runtimeClient, objTree, objTree.GetRoot(), treeOptions)
	if err != nil {
		return nil, NewInternalError(err)
	}

	return resourceTree, nil
}

// objectTreeToResourceTree converts an clusterctl ObjectTree to a ClusterResourceNode tree.
func objectTreeToResourceTree(ctx context.Context, runtimeClient ctrlclient.Client, objTree *tree.ObjectTree, object ctrlclient.Object, treeOptions ClusterResourceTreeOptions) (*ClusterResourceNode, error) {
	log := ctrl.LoggerFrom(ctx)

	if object == nil {
		return nil, nil
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
		Collapsed:   collapsed,
		Children:    []*ClusterResourceNode{},
		UID:         string(object.GetUID()),
	}

	children := objTree.GetObjectsByParent(object.GetUID())
	provider, err := getProvider(ctx, object, children, treeOptions)
	if err != nil {
		log.Error(err, "failed to get provider for object", "kind", kind, "name", object.GetName())
	}
	node.Provider = provider
	if tree.IsGroupObject(object) {
		node.Kind = strings.TrimSuffix(object.GetObjectKind().GroupVersionKind().Kind, "Group")
		// TODO: don't hard code this.
		if node.Kind == "Machine" {
			node.Provider = "cluster" // MachineGroup is a virtual node that represents a group of machines, so we set the provider to the same one as Machines.
			node.Group = "cluster.x-k8s.io"
		}
	}

	if node.Namespace = object.GetNamespace(); node.Namespace == "" {
		node.Namespace = "default"
	}

	setReadyFields(object, node)

	childTrees := []*ClusterResourceNode{}
	for _, child := range children {
		tree, err := objectTreeToResourceTree(ctx, runtimeClient, objTree, child, treeOptions)
		if err != nil {
			return nil, err
		}
		childTrees = append(childTrees, tree)
	}

	log.V(4).Info("Node is", "node", node.Kind+"/"+node.Name)
	// if treeOptions.GroupMachines {
	// 	node.Children = createKindGroupNode(ctx, object.GetNamespace(), "Machine", "cluster", childTrees, maxGroupSize)
	// } else {
	node.Children = childTrees
	// }

	if kind == "Cluster" {
		node.Children = addAddonsGroupNode(ctx, node.Children)
	}

	// If the resource represents a real CRD we want to collapse, and it has children, we can collapse it with tab.
	node.CollapseOnClick = tree.IsVirtualObject(object)
	node.CollapseWithTab = len(node.Children) > 0 && !node.CollapseOnClick
	node.Collapsible = node.CollapseWithTab || node.CollapseOnClick

	sort.Slice(node.Children, func(i, j int) bool {
		// TODO: make sure this is deterministic!
		if getSortKeys(node.Children[i])[0] == getSortKeys(node.Children[j])[0] {
			return getSortKeys(node.Children[i])[1] < getSortKeys(node.Children[j])[1]
		}
		return getSortKeys(node.Children[i])[0] < getSortKeys(node.Children[j])[0]
	})

	if tree.IsGroupObject(object) {
		node.Collapsed = true // Group objects are always collapsed by default.
		node.Collapsible = true
		node.CollapseOnClick = true
		node.CollapseWithTab = false

		if node.Kind == "Machine" {
			machineNames := strings.Split(tree.GetGroupItems(object), tree.GroupItemsSeparator)
			log.Info("Len machineNames", "machineNames", machineNames, "maxGroupSize", maxGroupSize)
			numTotalMachines := len(machineNames)
			if numTotalMachines > maxGroupSize {
				machineNames = machineNames[:maxGroupSize]
			}
			log.Info("Length of machineNames after slicing", "machineNames", machineNames, "maxGroupSize", maxGroupSize)

			for _, machineName := range machineNames {
				m := clusterv1.Machine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      machineName,
						Namespace: object.GetNamespace(),
					},
				}

				err := runtimeClient.Get(ctx, ctrlclient.ObjectKeyFromObject(&m), &m)
				if err != nil {
					log.Error(err, "failed to get machine object", "name", machineName, "namespace", object.GetNamespace())
					return nil, errors.Wrapf(err, "failed to get machine object %s/%s", object.GetNamespace(), machineName)
				}

				machineNode := &ClusterResourceNode{
					Name:            machineName,
					DisplayName:     machineName,
					Kind:            "Machine",
					Group:           "cluster.x-k8s.io",
					Version:         "v1beta1",
					Provider:        "cluster", // TODO: don't hardcode this
					UID:             string(m.GetUID()),
					Namespace:       m.GetNamespace(),
					Collapsed:       true,
					Collapsible:     true,
					CollapseWithTab: true,
					CollapseOnClick: false,
					Children:        []*ClusterResourceNode{},
				}
				setReadyFields(&m, machineNode)
				if (m.Spec.InfrastructureRef != corev1.ObjectReference{}) {
					if machineInfra, err := external.Get(ctx, runtimeClient, &m.Spec.InfrastructureRef, m.Namespace); err == nil {
						infraNode := &ClusterResourceNode{
							Name:            m.Spec.InfrastructureRef.Name,
							DisplayName:     m.Spec.InfrastructureRef.Name,
							Kind:            m.Spec.InfrastructureRef.Kind,
							Group:           m.Spec.InfrastructureRef.GroupVersionKind().Group,
							Version:         m.Spec.InfrastructureRef.GroupVersionKind().Version,
							Provider:        "infrastructure", // TODO: don't hardcode this
							UID:             string(m.Spec.InfrastructureRef.UID),
							Namespace:       m.Spec.InfrastructureRef.Namespace,
							Collapsed:       true,
							Collapsible:     false,
							CollapseWithTab: false,
							Children:        []*ClusterResourceNode{},
						}
						setReadyFields(machineInfra, infraNode)
						machineNode.Children = append(machineNode.Children, infraNode)
						// tree.Add(m, machineInfra, ObjectMetaName("MachineInfrastructure"), NoEcho(true))
					}
				}

				if m.Spec.Bootstrap.ConfigRef != nil {
					if machineBootstrap, err := external.Get(ctx, runtimeClient, m.Spec.Bootstrap.ConfigRef, m.Namespace); err == nil {
						bootstrapNode := &ClusterResourceNode{
							Name:            m.Spec.Bootstrap.ConfigRef.Name,
							DisplayName:     m.Spec.Bootstrap.ConfigRef.Name,
							Kind:            m.Spec.Bootstrap.ConfigRef.Kind,
							Group:           m.Spec.Bootstrap.ConfigRef.GroupVersionKind().Group,
							Version:         m.Spec.Bootstrap.ConfigRef.GroupVersionKind().Version,
							Provider:        "bootstrap", // TODO: don't hardcode this
							UID:             string(m.Spec.Bootstrap.ConfigRef.UID),
							Namespace:       m.Spec.Bootstrap.ConfigRef.Namespace,
							Collapsed:       false,
							Collapsible:     false,
							CollapseWithTab: false,
							Children:        []*ClusterResourceNode{},
						}
						setReadyFields(machineBootstrap, bootstrapNode)
						machineNode.Children = append(machineNode.Children, bootstrapNode)
					}
				}

				node.Children = append(node.Children, machineNode)
			}

			if numTotalMachines > maxGroupSize {
				node.Children = append(node.Children, &ClusterResourceNode{
					Name:            "",
					DisplayName:     fmt.Sprintf("%d %s...", numTotalMachines-maxGroupSize, flect.Pluralize(node.Kind)),
					Kind:            "Machine",
					Group:           "cluster.x-k8s.io",
					Version:         "v1beta1",
					Provider:        "cluster", // TODO: don't hardcode this
					UID:             node.UID + "more",
					Namespace:       node.Namespace,
					Collapsed:       true,
					Collapsible:     false,
					CollapseWithTab: false,
					Children:        []*ClusterResourceNode{},
					IsOverflowNode:  true, // This is an overflow node, meaning it contains more machines than the max group size.
				})
			}
		}
	}

	return node, nil
}

// addAddonsGroupNode finds all objects in children with `provider=addons` and create a parent node for them.
func addAddonsGroupNode(_ context.Context, children []*ClusterResourceNode) []*ClusterResourceNode {
	resultChildren := []*ClusterResourceNode{}

	addonsParent := &ClusterResourceNode{
		Name:            "",
		DisplayName:     "Add-ons",
		Kind:            "AddonsGroup",
		Provider:        "addons",
		CollapseWithTab: false,
		CollapseOnClick: true,
		Collapsible:     true,
		Collapsed:       false,
		HasReady:        false,
		Ready:           true,
		Severity:        "",
		Reason:          "",
		UID:             "addons",
	}

	for _, child := range children {
		if child.Provider == "addons" {
			addonsParent.Children = append(addonsParent.Children, child)
		} else {
			resultChildren = append(resultChildren, child)
		}
	}

	if len(addonsParent.Children) == 1 && addonsParent.Children[0].Kind == "ClusterResourceSetGroup" && addonsParent.Children[0].Name == "ClusterResourceSets" {
		// If the only add-ons are the CRS group node, just remove it and make the add-ons the new parent.
		addonsParent.Children = addonsParent.Children[0].Children
	}

	if len(addonsParent.Children) >= addonsAmountToCollapse {
		addonsParent.Collapsed = true
	}

	if len(addonsParent.Children) > 0 {
		resultChildren = append(resultChildren, addonsParent)
	}

	return resultChildren
}

// createKindGroupNode finds all objects in children with `kind` and create a parent node for them.
func createKindGroupNode(ctx context.Context, namespace string, kind string, provider string, children []*ClusterResourceNode, maxGroupSize int) []*ClusterResourceNode {
	log := ctrl.LoggerFrom(ctx)

	log.V(4).Info("Starting children are ", "children", nodeArrayNames(children))

	// sort.Slice(children, func(i, j int) bool {
	// 	// TODO: make sure this is deterministic!
	// 	if getSortKeys(children[i])[0] == getSortKeys(children[j])[0] {
	// 		return getSortKeys(children[i])[1] < getSortKeys(children[j])[1]
	// 	}
	// 	return getSortKeys(children[i])[0] < getSortKeys(children[j])[0]
	// })

	resultChildren := []*ClusterResourceNode{}

	// TODO: maybe in the future, we can group based on severity/error, but we'd still need a way to make sure the groups aren't too large.
	// Init a parent node, if the child groups need to be broken up. For example, if we have 100 machines, it would be
	// [MachineSet] -> [30 Machines] -> [10 Machines, 10 Machines, 10 Machines]

	conditionToGroupNode := map[string]*ClusterResourceNode{}
	for _, child := range children {
		if child.Kind == kind {
			conditionMapKey := fmt.Sprintf("HasReady: %s, Ready: %s, Severity: %s, Reason: %s", child.HasReady, child.Ready, child.Severity, child.Reason)
			if existingGroupNode, ok := conditionToGroupNode[conditionMapKey]; ok {
				// If we already have a group node for this condition, just add the child to it.
				existingGroupNode.Children = append(existingGroupNode.Children, child)
				existingGroupNode.UID += child.UID + " "
			} else {
				// If we don't have a group node for this condition, create a new one.
				groupNode := &ClusterResourceNode{
					Name:            "",
					Namespace:       namespace,
					DisplayName:     "",
					Kind:            kind,
					Group:           child.Group,
					Version:         child.Version,
					Provider:        provider, // TODO: don't hardcode this
					CollapseWithTab: false,
					CollapseOnClick: true,
					Collapsible:     true,
					Collapsed:       true,
					Children:        []*ClusterResourceNode{child},
					HasReady:        child.HasReady,
					Ready:           child.Ready,
					Severity:        child.Severity,
					Reason:          child.Reason,
					UID:             kind + ": " + child.UID + " ",
				}

				conditionToGroupNode[conditionMapKey] = groupNode
			}
		} else {
			resultChildren = append(resultChildren, child)
		}
	}

	for _, groupNode := range conditionToGroupNode {
		// If the group node has one child, we can just add it to the result children.
		if len(groupNode.Children) == 1 {
			resultChildren = append(resultChildren, groupNode.Children[0])
		} else {
			// If the group node has more than one child, we need to add it to the result children.
			groupNode.DisplayName = fmt.Sprintf("%d %s", len(groupNode.Children), flect.Pluralize(groupNode.Kind))
			resultChildren = append(resultChildren, groupNode)
		}
	}

	log.V(4).Info("Result children are", "children", nodeArrayNames(resultChildren))

	return resultChildren
}

// func hasSameReadyStatusSeverityAndReason(a, b *clusterv1.Condition) bool {
// 	if a == nil && b == nil {
// 		return true
// 	}
// 	if (a == nil) != (b == nil) {
// 		return false
// 	}

// 	return a.Status == b.Status &&
// 		a.Severity == b.Severity &&
// 		a.Reason == b.Reason
// }

// injectCustomResourcesToObjectTree amends the clusterctl ObjectTree with custom CRDs that are not included in the clusterctl resource discovery.
// It queries all CRD types and their instances containing the visualizer label and the cluster name label.
func injectCustomResourcesToObjectTree(ctx context.Context, c ctrlclient.Client, dcOptions client.DescribeClusterOptions, objTree *tree.ObjectTree) (map[string]string, error) {
	log := ctrl.LoggerFrom(ctx)

	log.V(4).Info("Adding user specified custom resources to object tree", "namespace", dcOptions.Namespace, "clusterName", dcOptions.ClusterName)

	crds, err := getCRDList(ctx, c, ctrlclient.MatchingLabels{visualizerv1.VisualizeResourceLabel: ""})
	if err != nil {
		return nil, err
	}

	namespace := dcOptions.Namespace
	clusterName := dcOptions.ClusterName

	clusterObjSelector := []ctrlclient.ListOption{
		ctrlclient.InNamespace(namespace),
		ctrlclient.MatchingLabels{clusterv1.ClusterNameLabel: clusterName},
	}

	providerTypeOverrideMap := make(map[string]string)
	clusterObjects := []ctrlclient.Object{}
	for _, crd := range crds {
		crdLabels := crd.GetLabels()
		if crdLabels != nil {
			if provider, ok := crdLabels[visualizerv1.ProviderTypeLabel]; ok {
				switch provider {
				case "cluster":
					fallthrough
				case "bootstrap":
					fallthrough
				case "controlplane":
					fallthrough
				case "infrastructure":
					fallthrough
				case "addons":
					fallthrough
				case "virtual":
					providerTypeOverrideMap[crd.Spec.Names.Kind] = provider
				default:
					return nil, errors.Errorf("Invalid provider type %s for CRD type %s \n", provider, crd.GetName())
				}
			}
		}

		for _, version := range crd.Spec.Versions {
			typeMeta := metav1.TypeMeta{
				Kind: crd.Spec.Names.Kind,
				APIVersion: metav1.GroupVersion{
					Group:   crd.Spec.Group,
					Version: version.Name,
				}.String(),
			}

			clusterObjList, err := getObjList(ctx, c, typeMeta, clusterObjSelector)
			if err != nil {
				return nil, err
			}

			for i := range clusterObjList.Items {
				clusterObj := &clusterObjList.Items[i]
				clusterObjects = append(clusterObjects, clusterObj)

				// Add the CRD to the object tree
			}

		}
	}

	for i := range clusterObjects {
		object := clusterObjects[i]
		// Make sure not to implicitly reference loop variable!
		if err := ensureObjConnectedTotree(ctx, c, objTree, object); err != nil {
			return nil, err
		}
	}

	return providerTypeOverrideMap, nil
}

// ensureObjConnectedTotree ensures that the object is connected to the tree by adding it and its parents until a parent is owned by the Cluster (root node).
// If a parent has no owner, it is set as a child of the Cluster.
// Note: At the moment, this only supports a use case where an object has only one owner which is also set the controller.
func ensureObjConnectedTotree(ctx context.Context, c ctrlclient.Client, objTree *tree.ObjectTree, object ctrlclient.Object) error {
	log := ctrl.LoggerFrom(ctx)

	if objTree.GetObject(object.GetUID()) != nil || objTree.GetRoot().GetUID() == object.GetUID() {
		log.V(4).Info("Object already in tree", "kind", object.GetObjectKind().GroupVersionKind().Kind, "name", object.GetName(), "namespace", object.GetNamespace())
		return nil
	}

	log.V(4).Info("Adding object to tree", "kind", object.GetObjectKind().GroupVersionKind().Kind, "name", object.GetName(), "namespace", object.GetNamespace())
	var parent ctrlclient.Object
	// TODO: handle case where there is no controllerRef or how to resolve multiple owners.
	ref := pickOwner(ctx, c, object)
	if ref != nil {
		if p, err := external.Get(ctx, c, ref, object.GetNamespace()); err != nil {
			return err
		} else {
			parent = p
		}
	} else {
		// If no ownerRef, set to root.
		parent = objTree.GetRoot()
		// TODO: look into creating an add-ons virtual node.
	}

	ensureObjConnectedTotree(ctx, c, objTree, parent)

	added, _ := objTree.Add(parent, object)
	if !added {
		return fmt.Errorf("failed to add object %s to tree", object.GetName())
	}

	return nil
}
