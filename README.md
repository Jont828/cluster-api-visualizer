# Cluster API Visualization App

Cluster API developers and operators often need to quickly get insight multicluster configuration. This app provides that insight by making Cluster API significantly more accessible and easier to understand for both new and experienced users. It gives a bird’s eye view of a multicluster architecture, visualizes all the Cluster API custom resources for each cluster, and provides quick access to the specs of any resource.

**Note:** This app is a prototype meant to serve as a proof of concept for a Cluster API GUI. It is built using the clusterctl client and VueJS.

![Demo Recording](demo/demo.gif)

### Quick start

#### 1. Prerequisites

This app requires [Go 1.17](https://go.dev/doc/install), [Cluster API](https://github.com/kubernetes-sigs/cluster-api), [Node.js](https://nodejs.org/en/), and the [Vue CLI](https://cli.vuejs.org/guide/installation.html).

#### 2. Clone the repository

```
cd ${GOPATH}/src # Or your go directory if GOPATH is not set
git clone git@github.com:Jont828/capi-visualization.git
```

#### 3. Install Go packages

```
cd ${GOPATH}/src/capi-visualization/
go get ./...
```

Clone a copy of [Cluster API](https://github.com/kubernetes-sigs/cluster-api), with the following commands. This will allow the app to access local changes to the repo.

```
mkdir ${GOPATH}/src/sigs.k8s.io
cd ${GOPATH}/src/sigs.k8s.io
git clone git@github.com:kubernetes-sigs/cluster-api.git
go get ./...
```
**Optional:** Install [Air](https://github.com/cosmtrek/air) to run the Go server with hot reloading for development. It can be installed using the instructions on the repo or your preferred package manager. Once it's it's installed, set it up with the following commands.

```
cd ${GOPATH}/src/capi-visualization/
air init
```

#### 4. Install node packages

```
cd ${GOPATH}/src/capi-visualization/web
npm install
```

#### 5. Start the app

From `${GOPATH}/src/capi-visualization/`, start the Go server with

```
go run main.go
```

In a separate terminal, enter `${GOPATH}/src/capi-visualization/web` and start the Vue app with

```
npm run serve
```

**Optional:** If you are using Air, you can start the Go server instead using
```
air
```

#### 6. Create a management cluster and workload cluster

Create a local management cluster with kind and a workload cluster by following the [Cluster API Quickstart](https://cluster-api.sigs.k8s.io/user/quick-start.html). Currently, only Docker and Azure workload clusters are supported. Once the clusters are running, open the app in your browser at `localhost:8080`.

### Contributing:

All contributions are welcome. If you'd like to help out, feel free fork the repo and submit a pull request. 

### Acknowledgements:

- The cluster trees are drawn in D3 using a modified version of [ssthouse/vue-tree-chart](https://github.com/ssthouse/vue-tree-chart)
- The method for building the tree is inspired by [this article](https://typeofnan.dev/an-easy-way-to-build-a-tree-with-object-references/)
- The tree is generated using the clusterctl client from [Cluster API](https://github.com/kubernetes-sigs/cluster-api)
- Thanks to [Vuetify](https://vuetifyjs.com/en/) for providing a great UI component toolkit
