package internal

import (
	"context"
	"sort"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2/klogr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type MultiClusterTreeNode struct {
	Name                   string                  `json:"name"`
	Namespace              string                  `json:"namespace"`
	InfrastructureProvider string                  `json:"infrastructureProvider"`
	IsManagement           bool                    `json:"isManagement"`
	Phase                  string                  `json:"phase"`
	Ready                  bool                    `json:"ready"`
	Children               []*MultiClusterTreeNode `json:"children"`
}

// ConstructMultiClusterTree returns a tree representing the workload cluster discovered in the management cluster.
func ConstructMultiClusterTree(ctrlClient ctrlclient.Client, k8sConfigClient *api.Config) (*MultiClusterTreeNode, *HTTPError) {
	log := klogr.New()

	currentContextName := k8sConfigClient.CurrentContext
	currentContext, ok := k8sConfigClient.Contexts[currentContextName]
	if !ok {
		return nil, &HTTPError{Status: 404, Message: "current context not found"}
	}
	name := currentContext.Cluster
	// namespace, err := clusterClient.Proxy().CurrentNamespace()
	// if err != nil {
	// 	return nil, NewInternalError(err)
	// }

	root := &MultiClusterTreeNode{
		Name:                   name,
		Namespace:              "",
		InfrastructureProvider: "",
		Children:               []*MultiClusterTreeNode{},
		IsManagement:           true,
	}

	clusterList := &clusterv1.ClusterList{}

	// TODO: should we use ctrlClient.MatchingLabels or try to use the labelSelector itself?
	if err := ctrlClient.List(context.TODO(), clusterList); err != nil {
		return nil, NewInternalError(err)
	}

	if clusterList == nil || len(clusterList.Items) == 0 {
		log.V(4).Info("No workload clusters found")
		return root, nil
	}
	sort.Slice(clusterList.Items, func(i, j int) bool {
		// This must be deterministic, otherwise the tree will be different between runs.
		// In this case, we can't have two clusters with the same name.
		return clusterList.Items[i].GetName() < clusterList.Items[j].GetName()
	})

	for _, cluster := range clusterList.Items {
		// Don't get the kubeconfig for now until we use it to find additional clusters.
		// kubeconfig, err := pkgClient.GetKubeconfig(client.GetKubeconfigOptions{
		// 	WorkloadClusterName: clusterName,
		// })
		infraProvider := cluster.Spec.InfrastructureRef.Kind

		readyCondition := conditions.Get(&cluster, clusterv1.ReadyCondition)

		workloadCluster := MultiClusterTreeNode{
			Name:                   cluster.GetName(),
			Namespace:              cluster.GetNamespace(),
			InfrastructureProvider: infraProvider,
			IsManagement:           false,
			Phase:                  cluster.Status.Phase,
			Ready:                  readyCondition != nil && readyCondition.Status == corev1.ConditionTrue,
			Children:               []*MultiClusterTreeNode{},
		}

		root.Children = append(root.Children, &workloadCluster)
	}

	return root, nil
}
