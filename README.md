<a href="https://cluster-api.sigs.k8s.io"><img alt="capi" src="./docs/images/cluster-api-logo.svg" width="160x" /></a>

[![Go Reference](https://pkg.go.dev/badge/github.com/Jont828/cluster-api-visualizer.svg)](https://pkg.go.dev/github.com/Jont828/cluster-api-visualizer)
[![Go Report Card](https://goreportcard.com/badge/github.com/Jont828/cluster-api-visualizer)](https://goreportcard.com/report/github.com/Jont828/cluster-api-visualizer)
[![Slack](https://img.shields.io/badge/join%20slack-%23cluster--api-brightgreen)](http://slack.k8s.io/)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Jont828/cluster-api-visualizer)
[![Latest Release](https://img.shields.io/github/v/release/Jont828/cluster-api-visualizer)](https://github.com/Jont828/cluster-api-visualizer/releases/latest)

# Cluster API Visualizer (CAPI V8r)

Cluster API developers and operators often need to quickly get insight multicluster configuration. This app provides that insight by making Cluster API significantly more accessible and easier to understand for both new and experienced users. It gives a bird’s eye view of a multicluster architecture, visualizes all the Cluster API custom resources for each cluster, and provides quick access to the specs and status of any resource.

https://github.com/user-attachments/assets/f1d5c036-eaac-4dc4-a4c5-8de9b64ac3cd

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
