package internal

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/cluster-api/controllers/external"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// GetGroupItems returns logs for a given resource from the CAPI controllers containing the provider name label.
func GetGroupItems(ctx context.Context, runtimeClient ctrlclient.Client, kind string, apiVersion string, namespace string, names []string) ([]unstructured.Unstructured, error) {
	log := ctrl.LoggerFrom(ctx)

	objects := make([]unstructured.Unstructured, 0)

	for _, name := range names {
		ref := &corev1.ObjectReference{
			APIVersion: apiVersion,
			Kind:       kind,
			Namespace:  namespace,
			Name:       name,
		}

		if object, err := external.Get(ctx, runtimeClient, ref, namespace); err != nil {
			log.Error(err, "Failed to get object", "kind", kind, "name", name, "namespace", namespace)
			return nil, err
		} else {
			objects = append(objects, *object)
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
