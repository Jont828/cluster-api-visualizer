# Cluster API Visualizer (CAPI V8r)

Cluster API developers and operators often need to quickly get insight multicluster configuration. This app provides that insight by making Cluster API significantly more accessible and easier to understand for both new and experienced users. It gives a bird’s eye view of a multicluster architecture, visualizes all the Cluster API custom resources for each cluster, and provides quick access to the specs and status of any resource.

https://github.com/user-attachments/assets/f1d5c036-eaac-4dc4-a4c5-8de9b64ac3cd

### Lens Integration Setup

#### Prerequisites
1 - Install and setup Lens version 2024.9.300059-latest (or later)
Debian - https://api.k8slens.dev/binaries/Lens-2024.9.300059-latest.amd64.deb
Mac - https://api.k8slens.dev/binaries/Lens-2024.9.300059-latest-arm64.dmg

2 - Launch Lens

3 - Click the menu button and select File->Preferences

4 - In the preferences dialog click the Kubernetes menu item

5 - Click the Sync Files and Folders button

6 - Navigate to the current user's Downloads folder (where the browser to be used stores downloads) and click the Sync button

#### Note that it's important to keep the download folder free from large files for Lens performance.


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
