# Cluster API Visualizer (CAPI V8r)

Cluster API developers and operators often need to quickly get insight multicluster configuration. This app provides that insight by making Cluster API significantly more accessible and easier to understand for both new and experienced users. It gives a bird’s eye view of a multicluster architecture, visualizes all the Cluster API custom resources for each cluster, and provides quick access to the specs and status of any resource.

![Demo Recording](demo/demo.gif)


### Quick start for local deployment 

## 1. Prerequisites

Install and set up [kind](https://kind.sigs.k8s.io/), [Docker](https://www.docker.com/), [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/), and [helm](https://helm.sh/). In addition, install any additional prerequisites needed for [Cluster API](https://cluster-api.sigs.k8s.io/).
#### 2. Create a Cluster API management cluster

Create a local management cluster with [kind](https://kind.sigs.k8s.io/) and a workload cluster by following the [Cluster API quick start guide](https://cluster-api.sigs.k8s.io/user/quick-start.html).

##### 3. To deploy to a local kind cluster 
```
make build-and-deploy
```

### Quick start

#### 1. Prerequisites

Install and set up [kind](https://kind.sigs.k8s.io/), [Docker](https://www.docker.com/), [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/), and [helm](https://helm.sh/). In addition, install any additional prerequisites needed for [Cluster API](https://cluster-api.sigs.k8s.io/).
#### 2. Create a Cluster API management cluster

Create a local management cluster with [kind](https://kind.sigs.k8s.io/) and a workload cluster by following the [Cluster API quick start guide](https://cluster-api.sigs.k8s.io/user/quick-start.html).

#### 3. Deploy with Helm

Then, run the following command to start the app:
```
./hack/deploy-repo-to-kind.sh
```

This will run the app as a deployment on management clusters built with kind.

### Contributing:

All contributions are welcome. If you'd like to help out, feel free fork the repo and submit a pull request. 

### Acknowledgements:

- Thanks to [@fabriziopandini](https://github.com/fabriziopandini) for helping guide the backend development.
- The cluster trees are drawn in D3 using a modified version of [ssthouse/vue-tree-chart](https://github.com/ssthouse/vue-tree-chart)
- The tree is generated using the clusterctl client from [Cluster API](https://github.com/kubernetes-sigs/cluster-api)
- The Go server was developed from Trevor Taubitz's [tutorial](https://hackandsla.sh/posts/2021-06-18-embed-vuejs-in-go/) on embedding VueJS in Go and [tutorial](https://hackandsla.sh/posts/2021-11-06-serve-spa-from-go/) on serving single page apps from Go 
- Thanks to [Vuetify](https://vuetifyjs.com/en/) for providing a great UI component toolkit
