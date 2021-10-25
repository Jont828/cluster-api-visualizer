const express = require('express');
const path = require('path');
const app = express(),
  bodyParser = require("body-parser");
port = 3080;

const yaml = require('js-yaml');
const fs = require('fs');

// place holder for the data
const clusters = [];

app.use(bodyParser.json());
app.use(express.static(path.join(__dirname, '../capi-vis/build')));

app.get('/api/cluster', (req, res) => {
  console.log('api/clusters called!')
  // console.log(req.query);
  let id = req.query.ID;
  // console.log('Got cluster ID ' + id);
  res.json(getTree(id));
});

app.get('/api/cluster-resource', (req, res) => {
  console.log('api/cluster-resource called!')
  console.log(req.query);
  let id = req.query.ID;

  try {
    const file = yaml.load(fs.readFileSync('./temp-assets/azureclusters.infrastructure.cluster.x-k8s.io-default-1495.yaml', 'utf8'));
    // console.log(JSON.stringify(file, null, 2));
    let result = formatToTreeview(file);
    res.json(result);
  } catch (e) {
    console.log(e);
    res.status(500);
  }
});


app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, '../capi-vis/build/index.html'));
});

app.listen(port, () => {
  console.log(`Server listening on the port::${port}`);
});

function formatToTreeview(resource, id = 0) {
  let result = [];
  if (Array.isArray(resource)) {
    let children = [];
    resource.forEach((e, i) => {
      result.push({
        id: id++,
        name: i,
        children: formatToTreeview(e, id)
      });
    });

  } else { // isObject
    Object.entries(resource).forEach(([key, value]) => {
      let name = '';
      let children = [];
      if (typeof (value) == 'string') {
        name = key + ': ' + value
      } else {
        name = key;
        children = formatToTreeview(value, id);
      }
      result.push({
        id: id++,
        name: name,
        children: children
      });
    });

  }

  return result;
}

function getTree(clusterId) {
  return {
    name: clusterId,
    kind: "Cluster",
    id: "cluster",
    provider: "capi",
    children: [
      {
        name: "",
        kind: "ClusterInfrastructure",
        id: "clusterInfra",
        provider: "",
        collapsable: true,
        children: [
          {
            name: "crs-calico",
            kind: "ClusterResourceSets",
            id: "crsCalico",
            provider: "addons",
            children: [
              {
                name: clusterId + "",
                kind: "ClusterResourceSetBinding",
                id: "clusterResourceSetBinding",
                provider: "addons",
                children: [],
              },
            ],
          },
          {
            name: "crs-calico-ipv6",
            kind: "ClusterResourceSets",
            id: "crsCalicoIpv6",
            provider: "addons",
            children: [],
          },
          {
            name: "flannel-windows",
            kind: "ClusterResourceSet",
            id: "flannelWindows",
            provider: "addons",
            children: [],
          },
          {
            name: clusterId + "",
            kind: "AzureCluster",
            id: "azureCluster",
            provider: "infra",
            children: [],
          },
          {
            name: clusterId + "-md",
            kind: "KubeadmConfigTemplate",
            id: "kubeadmConfigTemp",
            provider: "bootstrap",
            children: [],
          },
          {
            name: "cluster-identity",
            kind: "AzureClusterIdentity",
            id: "clusterIdentity",
            provider: "infra",
            children: [],
          },
        ],
      },
      {
        name: "",
        kind: "ControlPlane",
        id: "controlPlane",
        provider: "",
        collapsable: true,
        children: [
          {
            name: clusterId + "-control-plane",
            kind: "KubeadmControlPlane",
            id: "kubeadmCtrlPlane",
            provider: "ctrlPlane",
            children: [
              {
                name: clusterId + "-control-plane",
                kind: "Machine",
                id: "machineCtrlPlane",
                provider: "capi",
                children: [
                  {
                    name: clusterId + "-control-plane",
                    kind: "AzureMachine",
                    id: "azureMachineCtrl",
                    provider: "infra",
                    children: [],
                  },
                  {
                    name: clusterId + "-control-plane",
                    kind: "KubeadmConfig",
                    id: "kubeadmConfigCtrl",
                    provider: "bootstrap",
                    children: [],
                  },
                ],
              },
            ],
          },
          {
            name: clusterId + "-control-plane",
            kind: "AzureMachineTemplate",
            id: "azureMachineTemplateCtrl",
            provider: "infra",
            children: [],
          },
        ],
      },
      {
        name: "",
        kind: "Workers",
        id: "workers",
        provider: "",
        collapsable: true,
        children: [
          {
            name: clusterId + "-control-plane",
            kind: "AzureMachineTemplate",
            id: "azureMachineTempMd",
            provider: "infra",
            children: [],
          },
          {
            name: clusterId + "-md",
            kind: "MachineDeployment",
            id: "machineDeployment",
            provider: "capi",
            children: [
              {
                name: clusterId + "",
                kind: "MachineSet",
                id: "machineSet",
                provider: "capi",
                children: [
                  {
                    name: clusterId + "-md-1",
                    kind: "Machine",
                    id: "machine1",
                    provider: "capi",
                    children: [
                      {
                        name: clusterId + "-md-1",
                        kind: "AzureMachine",
                        id: "azureMachine1",
                        provider: "infra",
                        children: [],
                      },
                      {
                        name: clusterId + "-control-plane",
                        kind: "KubeadmConfig",
                        id: "kubeadmConfig1",
                        provider: "bootstrap",
                        children: [],
                      },
                    ],
                  },
                ],
              },
            ],
          },
        ],
      },
    ],
    links: [
      // {
      //   parent: "machine1",
      //   child: "azureMachine1",
      //   styles: {
      //     "stroke-width": "4px",
      //     stroke: "#555",
      //   },
      // },
      // {
      //   parent: "machineCtrlPlane",
      //   child: "azureMachineCtrl",
      //   styles: {
      //     "stroke-width": "4px",
      //     stroke: "#555",
      //   },
      // },
      // {
      //   parent: "cluster",
      //   child: "clusterInfra",
      //   styles: {
      //     "stroke-width": "4px",
      //     stroke: "#555",
      //   },
      // },
      // {
      //   parent: "clusterInfra",
      //   child: "azureCluster",
      //   styles: {
      //     "stroke-width": "4px",
      //     stroke: "#555",
      //   },
      // },
    ],
    identifier: "id",
  }
}