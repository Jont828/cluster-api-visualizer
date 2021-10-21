<template>
  <div class="container">
    <h1>Cluster Resource Ownership: {{ this.$route.params.id }}</h1>
    <vue-tree
      id="tree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="false"
    >
      <!-- <template v-slot:node="{ node }"> -->
      <template v-slot:node="{ node, collapsed }">
        <div
          class="node"
          :style="{ 
            'background-color': colors[node.provider], 
            border: collapsed ? '2px solid grey' : '',
            stroke: 'black'
          }"
        >
          <router-link
            :to="'/'"
            class="node-router-link"
          >
            <p class="kind">{{ node.kind }}</p>
            <p class="name">{{ node.name }}</p>
          </router-link>
        </div>
      </template>
    </vue-tree>
    <div class="legend">
      <div
        class="legend-entry"
        v-for="(color, provider) in this.colors"
        :key="provider"
      >
        <div :style="{
          'background-color': color
        }" />
        <span>{{ provider }}</span>
      </div>
    </div>
  </div>
</template>

<script>
import VueTree from "@ssthouse/vue-tree-chart";

export default {
  name: "TargetCluster",
  components: {
    VueTree,
  },
  methods: {},
  data() {
    return {
      colors: {
        bootstrap: "#FFF2CC",
        capi: "#DAE3F3",
        infra: "#E2F0D9",
        addons: "#FBE5D6",
        none: "#D0CECE",
      },
      treeData: {
        name: "",
        kind: "All Resources",
        id: "root",
        provider: "none",
        children: [
          {
            name: "",
            kind: "ClusterInfrastructure",
            id: "clusterInfra",
            provider: "none",
            children: [
              {
                name: "crs-calico",
                kind: "ClusterResourceSets",
                id: "crsCalico",
                provider: "addons",
                children: [],
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
                name: this.$route.params.id,
                kind: "Cluster",
                id: "capi",
                provider: "cluster",
                children: [
                  {
                    name: this.$route.params.id + "",
                    kind: "ClusterResourceSetBinding",
                    id: "clusterResourceSetBinding",
                    provider: "addons",
                    children: [],
                  },

                  {
                    name: this.$route.params.id + "",
                    kind: "AzureCluster",
                    id: "azureCluster",
                    provider: "infra",
                    children: [],
                  },
                  {
                    name: this.$route.params.id + "-md",
                    kind: "KubeadmConfigTemplate",
                    id: "kubeadmConfigTemp",
                    provider: "bootstrap",
                    children: [],
                  },
                ],
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
            provider: "none",
            children: [
              {
                name: this.$route.params.id + "-control-plane",
                kind: "KubeadmCtrlPlane",
                id: "kubeadmCtrlPlane",
                provider: "ctrlPlane",
                children: [
                  {
                    name: this.$route.params.id + "-control-plane",
                    kind: "Machine",
                    id: "machineCtrlPlane",
                    provider: "capi",
                    children: [
                      {
                        name: this.$route.params.id + "-control-plane",
                        kind: "AzureMachine",
                        id: "azureMachineCtrl",
                        provider: "infra",
                        children: [],
                      },
                      {
                        name: this.$route.params.id + "-control-plane",
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
                name: this.$route.params.id + "-control-plane",
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
            provider: "none",
            children: [
              {
                name: this.$route.params.id + "-control-plane",
                kind: "AzureMachineTemp",
                id: "azureMachineTempMd",
                provider: "infra",
                children: [],
              },
              {
                name: this.$route.params.id + "-md",
                kind: "MachineDeployment",
                id: "machineDeployment",
                provider: "capi",
                children: [
                  {
                    name: this.$route.params.id + "",
                    kind: "MachineSet",
                    id: "machineSet",
                    provider: "capi",
                    children: [
                      {
                        name: this.$route.params.id + "-md-1",
                        kind: "Machine",
                        id: "machine1",
                        provider: "capi",
                        children: [
                          {
                            name: this.$route.params.id + "-md-1",
                            kind: "AzureMachine",
                            id: "azureMachine1",
                            provider: "infra",
                            children: [],
                          },
                          {
                            name: this.$route.params.id + "-control-plane",
                            kind: "KubeadmConfig",
                            id: "kubeadmConfig1",
                            provider: "bootstrap",
                            children: [],
                          },
                        ],
                      },
                      {
                        name: "... (3x)",
                        kind: "",
                        id: "machine2",
                        provider: "capi",
                        children: [
                          // {
                          //   name: "...",
                          // kind: "",
                          //   id: "azureMachine2",
                          //   provider: "worker",
                          //   children: [],
                          // },
                          // {
                          //   name: "...",
                          // kind: "",
                          //   id: "kubeadmConfig2",
                          //   provider: "worker",
                          //   children: [],
                          // },
                        ],
                      },
                      {
                        name: this.$route.params.id + "-md-3",
                        kind: "Machine",
                        id: "machine3",
                        provider: "capi",
                        children: [
                          {
                            name: this.$route.params.id + "-md-3",
                            kind: "AzureMachine",
                            id: "azureMachine3",
                            provider: "infra",
                            children: [],
                          },
                          {
                            name: "default-3-control-plane",
                            kind: "KubeadmConfig",
                            id: "kubeadmConfig3",
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
          //   parent: "kubeadmCtrlPlane",
          //   child: "azureMachineCtrl",
          // },
          // {
          //   parent: "kubeadmCtrlPlane",
          //   child: "kubeadmConfigCtrl",
          // },
          {
            parent: "crsCalico",
            child: "clusterResourceSetBinding",
          },
        ],
        identifier: "id",
      },
      treeConfig: { nodeWidth: 150, nodeHeight: 40, levelHeight: 100 },
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    };
  },
};
</script>

<style scoped>
#tree {
  width: 100%;
  height: 800px;
  border: 1px solid black;
}

.container {
  height: 100%;
  width: 100%;
  max-width: 100%;
  margin: 0 !important;
}

.node-slot {
  cursor: default !important;
}

.node {
  width: 130px;
  height: 40px;
  /* padding: 8px; */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #dae8fc;
  border-radius: 4px;
}

.node p {
  font-size: 10px;
  margin: 2px;
  color: #2c3e50;
}

.node .node-router-link {
  text-decoration: none;
}

.name {
  font-style: italic;
}

.legend-entry {
  display: inline-block;
  margin-right: 10px;
}

.legend-entry div {
  display: inline-block;
  border: 1px solid black;
  margin: 0 5px;
  width: 12px;
  height: 12px;
}
</style>
