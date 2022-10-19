# Development

This document describes the process for running this application on your local computer.

#### 1. Prerequisites

For local development, install [Go 1.17](https://go.dev/doc/install), [Node.js](https://nodejs.org/en/), [npm](https://www.npmjs.com/), [Vue CLI](https://cli.vuejs.org/guide/installation.html), [Cluster API](https://github.com/kubernetes-sigs/cluster-api), and its [prerequisites](https://cluster-api.sigs.k8s.io/user/quick-start.html#common-prerequisites).

#### 2. Create a management cluster and workload cluster with Cluster API

Create a local management cluster with kind and a workload cluster by following the [Cluster API quick start](https://cluster-api.sigs.k8s.io/user/quick-start.html). Alternatively, this can be used with a Cluster API provider repo if you have it set up, i.e. Cluster API Provider Azure, GCP, or AWS.

#### 3. Clone the repository

```
cd ${GOPATH}/src # Or your go directory if GOPATH is not set
git clone https://github.com/Jont828/cluster-api-visualizer.git
```

#### 4. Install Go packages

```
cd ${GOPATH}/src/cluster-api-visualizer/
go get ./...
```

#### 5. Run make

This will install node packages, build the web app, build the Go backend, and start the app.

```
make
```
Alternatively, you can run the steps individually if `node_modules` are already installed or the web app assets are already built.

```
make npm-install    # Install node packages
make build          # Build the web app into `web/dist` and the Go binary
make run            # Run the Go binary if the binary is built
```

### Hot reloading

For development and testing, the app can be run with hot reloading for both Go and Vue. 

To hot reload the Go server, install [Air](https://github.com/cosmtrek/air) and run `air init` in the root folder. After installing the Go packages, do not run `make` and instead and start the Go server from `${GOPATH}/src/cluster-api-visualizer/` with

```
air
```

To hot reload the front end server, open a separate terminal, enter `${GOPATH}/src/cluster-api-visualizer/web`, and install node packages if you haven't done so with


```
npm install
```

Then start the Vue app with

```
npm run serve
```