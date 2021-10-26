const k8s = require('@kubernetes/client-node');
const { default: cluster } = require('cluster');

const kc = new k8s.KubeConfig();
kc.loadFromDefault();
const k8sCrd = kc.makeApiClient(k8s.CustomObjectsApi);

module.exports = function constructTargetClusterTree(clusterId) {
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