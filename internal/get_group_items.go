package internal

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api/cmd/clusterctl/client/tree"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// GetGroupItems returns logs for a given resource from the CAPI controllers containing the provider name label.
func GetGroupItems(ctx context.Context, runtimeClient ctrlclient.Client, kind, apiVersion, namespace, status, severity, reason string) ([]unstructured.Unstructured, error) {
	log := ctrl.LoggerFrom(ctx)

	objects := make([]unstructured.Unstructured, 0)
	objectList := &unstructured.UnstructuredList{}
	objectList.SetAPIVersion(apiVersion)
	objectList.SetKind(kind)

	if err := runtimeClient.List(ctx, objectList); err != nil {
		log.Error(err, "Failed to list objects", "kind", kind, "apiVersion", apiVersion, "namespace", namespace)
		return nil, err
	}

	// TODO: filter so they match the machine deployment or control plane.
	for _, object := range objectList.Items {
		condition := tree.GetReadyCondition(&object)
		if status == string(condition.Status) &&
			severity == string(condition.Severity) &&
			reason == string(condition.Reason) {
			// If the object matches the status, reason, and severity, add it to the list
			objects = append(objects, object)
		}
	}

	return objects, nil
}

func GroupItemsToResourceNodes(ctx context.Context, objects []unstructured.Unstructured) []*ClusterResourceNode {
	_ = ctrl.LoggerFrom(ctx)

	nodes := make([]*ClusterResourceNode, 0, len(objects))

	for _, object := range objects {
		group := object.GetObjectKind().GroupVersionKind().Group
		kind := object.GetObjectKind().GroupVersionKind().Kind
		version := object.GetObjectKind().GroupVersionKind().Version

		node := &ClusterResourceNode{
			Name:        object.GetName(),
			DisplayName: getDisplayName(&object),
			Kind:        kind,
			Group:       group,
			Version:     version,
			Collapsed:   false,
			Collapsible: false,
			Children:    []*ClusterResourceNode{},
			UID:         string(object.GetUID()),
		}

		if node.Namespace = object.GetNamespace(); node.Namespace == "" {
			node.Namespace = "default"
		}
		setReadyFields(&object, node)
		// TODO: get conditions
		nodes = append(nodes, node)
	}

	return nodes
}
