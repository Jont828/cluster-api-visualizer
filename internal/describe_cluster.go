package internal

import (
	"context"
	"fmt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sort"
	"strings"

	visualizerv1 "github.com/Jont828/cluster-api-visualizer/api/v1"
	"github.com/gobuffalo/flect"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/tree"
	"sigs.k8s.io/cluster-api/controllers/external"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterResourceNode represents a node in the Cluster API resource tree and is used to configure the frontend with additional
// options like collapsibility and provider.
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
	providerTypeOverrideMap      map[string]string
}

func injectClusterTemplates(ctx context.Context, tree *ClusterResourceNode) error {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to get in cluster rest config %w", err)
	}
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to get in dynamic client %w", err)
	}

	resourceID := schema.GroupVersionResource{
		Group:    "hmc.mirantis.com",
		Version:  "v1alpha1",
		Resource: "clustertemplates",
	}

	list, err := dc.Resource(resourceID).Namespace(tree.Namespace).List(ctx, metav1.ListOptions{})

	if apierrors.IsNotFound(err) || len(list.Items) == 0 {
		return nil
	}

	serviceTemplateNode := &ClusterResourceNode{
		Name:        "ClusterTemplates",
		Namespace:   tree.Namespace,
		DisplayName: "ClusterTemplates",
		Kind:        "ClusterTemplates",
		Provider:    "",
		Collapsible: true,
		Collapsed:   false,
		Ready:       true,
		Severity:    "",
		HasReady:    false,
		Children:    nil,
	}

	for _, template := range list.Items {
		serviceTemplateNode.Children = append(serviceTemplateNode.Children, &ClusterResourceNode{
			Name:        template.GetName(),
			Namespace:   template.GetNamespace(),
			DisplayName: template.GetName(),
			Kind:        template.GetKind(),
			Group:       template.GroupVersionKind().Group,
			Version:     template.GroupVersionKind().Version,
			Provider:    "",
			UID:         string(template.GetUID()),
			Collapsible: false,
			Collapsed:   false,
			Ready:       true,
			Severity:    "",
			HasReady:    false,
			Children:    nil,
		})
	}

	tree.Children = append(tree.Children, serviceTemplateNode)
	return nil
}

func injectHmcResources(ctx context.Context, resourceName string, displayName string, tree *ClusterResourceNode) error {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to get in cluster rest config %w", err)
	}
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("failed to get in dynamic client %w", err)
	}

	resourceID := schema.GroupVersionResource{
		Group:    "hmc.mirantis.com",
		Version:  "v1alpha1",
		Resource: resourceName,
	}

	list, err := dc.Resource(resourceID).Namespace(tree.Namespace).List(ctx, metav1.ListOptions{})

	if apierrors.IsNotFound(err) || len(list.Items) == 0 {
		return nil
	}

	serviceTemplateNode := &ClusterResourceNode{
		Name:        displayName,
		Namespace:   tree.Namespace,
		DisplayName: displayName,
		Kind:        displayName,
		Provider:    "",
		Collapsible: true,
		Collapsed:   true,
		Ready:       true,
		Severity:    "",
		HasReady:    false,
		Children:    nil,
	}

	for _, template := range list.Items {
		serviceTemplateNode.Children = append(serviceTemplateNode.Children, &ClusterResourceNode{
			Name:        template.GetName(),
			Namespace:   template.GetNamespace(),
			DisplayName: template.GetName(),
			Kind:        template.GetKind(),
			Group:       template.GroupVersionKind().Group,
			Version:     template.GroupVersionKind().Version,
			Provider:    "",
			UID:         string(template.GetUID()),
			Collapsible: false,
			Collapsed:   false,
			Ready:       true, // TODO
			Severity:    "",
			HasReady:    false,
			Children:    nil,
		})
	}

	tree.Children = append(tree.Children, serviceTemplateNode)
	return nil
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

	overrides, err := injectCustomResourcesToObjectTree(ctx, runtimeClient, dcOptions, objTree)
	if err != nil {
		return nil, NewInternalError(err)
	}
	treeOptions.providerTypeOverrideMap = overrides

	resourceTree := objectTreeToResourceTree(ctx, objTree, objTree.GetRoot(), treeOptions)
	injectHmcResources(ctx, "clustertemplates", "ClusterTemplates", resourceTree)
	injectHmcResources(ctx, "servicetemplates", "ServiceTemplates", resourceTree)
	injectHmcResources(ctx, "credentials", "Credentials", resourceTree)
	return resourceTree, nil
}

// objectTreeToResourceTree converts an clusterctl ObjectTree to a ClusterResourceNode tree.
func objectTreeToResourceTree(ctx context.Context, objTree *tree.ObjectTree, object ctrlclient.Object, treeOptions ClusterResourceTreeOptions) *ClusterResourceNode {
	log := ctrl.LoggerFrom(ctx)

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
	provider, err := getProvider(ctx, object, children, treeOptions)
	if err != nil {
		log.Error(err, "failed to get provider for object", "kind", kind, "name", object.GetName())
	}
	node.Provider = provider

	setReadyFields(object, node)

	childTrees := []*ClusterResourceNode{}
	for _, child := range children {
		// log.Info("Child UID is ", "UID", child.GetUID())
		// obj := objTree.GetObject(child.GetUID())
		// log.Info("Obj is", "obj", obj)
		childTrees = append(childTrees, objectTreeToResourceTree(ctx, objTree, child, treeOptions))
	}

	log.V(4).Info("Node is", "node", node.Kind+"/"+node.Name)
	if treeOptions.GroupMachines {
		node.Children = createKindGroupNode(ctx, object.GetNamespace(), "Machine", "cluster", childTrees, false)
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

// createKindGroupNode finds all objects in children with `kind` and create a parent node for them.
func createKindGroupNode(ctx context.Context, namespace string, kind string, provider string, children []*ClusterResourceNode, groupForOne bool) []*ClusterResourceNode {
	log := ctrl.LoggerFrom(ctx)

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
