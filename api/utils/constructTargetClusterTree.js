const k8s = require('@kubernetes/client-node');
const { default: cluster } = require('cluster');
const { assert } = require('console');

const kc = new k8s.KubeConfig();
kc.loadFromDefault();
const k8sCrd = kc.makeApiClient(k8s.CustomObjectsApi);

const customResources = [
  { group: "addons.cluster.x-k8s.io", plural: "clusterresourcesetbindings", category: "clusterInfra" },
  { group: "addons.cluster.x-k8s.io", plural: "clusterresourcesets", category: "clusterInfra" },
  // { group: "cluster.x-k8s.io", plural: "clusterclasses", category: "clusterInfra" },
  { group: "cluster.x-k8s.io", plural: "clusters", category: "" },
  { group: "cluster.x-k8s.io", plural: "machinedeployments", category: "" },
  // { group: "cluster.x-k8s.io", plural: "machinehealthchecks", category: "clusterInfra" },
  // { group: "cluster.x-k8s.io", plural: "machinepools", category: "clusterInfra" },
  { group: "cluster.x-k8s.io", plural: "machinesets", category: "workers" },
  { group: "cluster.x-k8s.io", plural: "machines", category: "" },
  { group: "infrastructure.cluster.x-k8s.io", plural: "azureclusteridentities", category: "clusterInfra" },
  { group: "infrastructure.cluster.x-k8s.io", plural: "azureclusters", category: "clusterInfra" },
  // { group: "infrastructure.cluster.x-k8s.io", plural: "azuremachinepoolmachines", category: "clusterInfra" },
  // { group: "infrastructure.cluster.x-k8s.io", plural: "azuremachinepools", category: "clusterInfra" },
  { group: "infrastructure.cluster.x-k8s.io", plural: "azuremachines", category: "" },
  { group: "infrastructure.cluster.x-k8s.io", plural: "azuremachinetemplates", category: "" },
  // { group: "infrastructure.cluster.x-k8s.io", plural: "azuremanagedclusters", category: "clusterInfra" },
  // { group: "infrastructure.cluster.x-k8s.io", plural: "azuremanagedcontrolplanes", category: "clusterInfra" },
  // { group: "infrastructure.cluster.x-k8s.io", plural: "azuremanagedmachinepools", category: "clusterInfra" },
  // { group: "infrastructure.cluster.x-k8s.io", plural: "azureserviceprincipals", category: "clusterInfra" },
  // { group: "infrastructure.cluster.x-k8s.io", plural: "azuresystemassignedidentites", category: "clusterInfra" },
  // { group: "infrastructure.cluster.x-k8s.io", plural: "azureuserassignedidentites", category: "clusterInfra" },
  { group: "bootstrap.cluster.x-k8s.io", plural: "kubeadmconfigs", category: "" },
  { group: "bootstrap.cluster.x-k8s.io", plural: "kubeadmconfigtemplates", category: "clusterInfra" },
  { group: "controlplane.cluster.x-k8s.io", plural: "kubeadmcontrolplanes", category: "controlPlane" },
  // { group: "controlplane.cluster.x-k8s.io", plural: "kubeadmcontrolplanetemplates", category: "controlPlane" }
]

async function getCRDInstances(group, plural, category, clusterName) {

  // Hack since getClusterCustomObject is getting a 404
  const response = await k8sCrd.listClusterCustomObject('cluster.x-k8s.io', 'v1beta1', 'clusters');
  let clusters = response.body.items.filter(e => e.metadata.name == clusterName);
  assert(clusters.length == 1)
  let clusterUid = clusters[0].metadata.uid;
  // End hack

  const res = await k8sCrd.listClusterCustomObject(group, 'v1beta1', plural);
  let crds = [];
  res.body.items.forEach((e, i) => {
    let crd = {
      id: e.metadata.uid,
      name: e.metadata.name,
      kind: e.kind,
      group: group,
      provider: group.substr(0, group.indexOf('.')),
    }

    if (!category) {
      if (crd.name.indexOf(clusterName + '-control-plane') == 0) {
        crd.category = 'controlPlane'
      } else if (crd.name.indexOf(clusterName + '-md') == 0) {
        crd.category = 'workers'
      }
    } else {
      crd.category = category
    }

    console.log(crd.kind + '/' + crd.name);
    let owners = e.metadata.ownerReferences;
    let owner;

    if (owners === undefined) {
      owner = clusterUid;
    } else if (owners.length > 1) {
      owners = owners.filter(elt => elt.kind != 'Cluster');
      assert(owners.length == 1);
      owner = owners[0].uid;
    }

    crds.parent = owner
    crds.push(crd)

  })

  return crds;
}

module.exports = async function constructTargetClusterTree(clusterId) {
  let allCrds = [];

  for (e of customResources) {
    const instances = await getCRDInstances(e.group, e.plural, e.category, clusterId);
    allCrds = allCrds.concat(instances);
  }

  const whitelist = ['crs-calico', 'crs-calico-ipv6', 'flannel-windows', 'cluster-identity'];

  let crds = allCrds.filter((crd) => (crd.name.indexOf(clusterId) == 0 || whitelist.includes(crd.name)));

  console.log('Printing categories', crds.length);
  crds.forEach((e, i) => {
    console.log(e);
  })
  console.log('Started tree for', clusterId);
  let tree = {
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
        children: [],
      },
      {
        name: "",
        kind: "ControlPlane",
        id: "controlPlane",
        provider: "",
        collapsable: true,
        children: [],
      },
      {
        name: "",
        kind: "Workers",
        id: "workers",
        provider: "",
        collapsable: true,
        children: [],
      },
    ],
  }

  let sampleTree = {
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
  };

  return sampleTree;
}

