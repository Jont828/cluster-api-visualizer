package api

const (
	// VisualizeResourceLabel is used to include additional CRDs outside of clusterctl discovery
	// into the Cluster view. The label should be applied to the CRD type itself, and the resource
	// instances must have the `cluster.x-k8s.io/cluster-name` in order to be affiliated with a
	// Cluster.
	// The resource will be inserted into the tree based off its owner reference, and its owners
	// will be inserted as well until the owner is a Cluster, or it has no owner, in which case it
	// will be added as a child of the Cluster.
	VisualizeResourceLabel = "visualizer.cluster.x-k8s.io"
)
