# Cluster API Visualization App [WIP]

Cluster API developers and operators often need to quickly get insight multicluster configuration. This app provides that insight by making Cluster API significantly more accessible and easier to understand for both new and experienced users. It gives a bird’s eye view of a multicluster architecture, visualizes all the Cluster API custom resources for each cluster, and provides quick access to the specs of any resource.

**Note:** This app is a prototype meant to serve as a proof of concept for a Cluster API GUI. The back-end logic is still a work in progress. Currently, only Docker and Azure are supported as infrastructure providers until more can be added.

### Quick start

##### 1. Clone the repository

``` 
git clone git@github.com:Jont828/capi-visualization.git
```

##### 2. Install node packages for the front-end and back-end

```
cd capi-visualization/capi-vis
npm install
cd ../api
npm install
```

##### 3. Start the app

From `capi-vis/`, start the Vue app with

```
npm run serve -- --port 8081
```

And from `api/`, start the Node server with

```
npm run dev
```

##### 4. Create a management cluster and workload cluster

Create a local management cluster with kind and a workload cluster by following the [Cluster API Quickstart](https://cluster-api.sigs.k8s.io/user/quick-start.html). Currently, only Docker and Azure workload clusters are supported. Once the clusters are running, open the app in your browser at `localhost:8081`.

### Contributing:

All contributions are welcome. If you'd like to help out, feel free fork the repo and submit a pull request. 

### Acknowledgements:

- The cluster trees are drawn in D3 using a modified version of [ssthouse/vue-tree-chart](https://github.com/ssthouse/vue-tree-chart)
- The method for building the tree is inspired by [this article](https://typeofnan.dev/an-easy-way-to-build-a-tree-with-object-references/)
- The categorization of cluster resources into cluster infrastructure, control plane, and workers is based on `clusterctl describe`, which can be found in [Cluster API](https://github.com/kubernetes-sigs/cluster-api)
- Thanks to [Vuetify](https://vuetifyjs.com/en/) for providing a great UI component toolkit
