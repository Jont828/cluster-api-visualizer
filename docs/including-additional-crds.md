# Including additional CRD types in the Cluster view

The Cluster view lists all the CRDs returned by the discovery process in `clusterctl describe`. To add additional CRDs, you must modify the CRD type (not the instances) with the `visualizer.cluster.x-k8s.io: ""` label as follows:

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  labels:
    visualizer.cluster.x-k8s.io: ""
```

The instances of the CRD must also have the `cluster.x-k8s.io/cluster-name: <cluster-name>` label in order to be affiliated with a Cluster. 

The CRD instance will then be displayed in the Cluster view such that it will be listed as a child of its controller reference. If the controller reference is not in the tree, it will be added as well, and if there is no controller reference, it will be a child of the Cluster object at the root of the tree.
