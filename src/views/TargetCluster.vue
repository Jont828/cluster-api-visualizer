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
            'background-color': colors[node.category], 
            border: collapsed ? '2px solid grey' : '' 
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
        v-for="(color, category) in this.colors"
        :key="category"
      >
        <div :style="{
          'background-color': color
        }" />
        <span>{{ category }}</span>
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
        root: "#FFF2CC",
        infra: "#FBE5D6",
        ctrlPlane: "#DAE3F3",
        worker: "#E2F0D9",
        none: "#D0CECE",
      },
      treeData: {
        name: "",
        kind: "All Resources",
        id: "root",
        category: "root",
        children: [
          {
            name: "crs-calico",
            kind: "ClusterResourceSets",
            id: "crsCalico",
            category: "none",
            children: [],
          },
          {
            name: "crs-calico-ipv6",
            kind: "ClusterResourceSets",
            id: "crsCalicoIpv6",
            category: "none",
            children: [],
          },
          {
            name: "flannel-windows",
            kind: "ClusterResourceSet",
            id: "flannelWindows",
            category: "none",
            children: [],
          },
          {
            name: this.$route.params.id,
            kind: "Cluster",
            id: "cluster",
            category: "infra",
            children: [
              {
                name: this.$route.params.id + "",
                kind: "ClusterResourceSetBinding",
                id: "clusterResourceSetBinding",
                category: "none",
                children: [],
              },
              {
                name: this.$route.params.id + "-control-plane",
                kind: "KubeAdmCtrlPlane",
                id: "kubeAdmCtrlPlane",
                category: "ctrlPlane",
                children: [
                  {
                    name: this.$route.params.id + "-control-plane",
                    kind: "Machine",
                    id: "machineCtrlPlane",
                    category: "ctrlPlane",
                    children: [
                      {
                        name: this.$route.params.id + "-control-plane",
                        kind: "AzureMachine",
                        id: "azureMachineCtrl",
                        category: "ctrlPlane",
                        children: [],
                      },
                      {
                        name: this.$route.params.id + "-control-plane",
                        kind: "KubeAdmConfig",
                        id: "kubeAdmConfigCtrl",
                        category: "ctrlPlane",
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
                category: "ctrlPlane",
                children: [],
              },
              {
                name: this.$route.params.id + "",
                kind: "AzureCluster",
                id: "azureCluster",
                category: "infra",
                children: [],
              },
              {
                name: this.$route.params.id + "-md",
                kind: "KubeAdmConfigTemplate",
                id: "kubeAdmConfigTemp",
                category: "none",
                children: [],
              },
              {
                name: this.$route.params.id + "-control-plane",
                kind: "AzureMachineTemp",
                id: "azureMachineTempMd",
                category: "worker",
                children: [],
              },
              {
                name: this.$route.params.id + "-md",
                kind: "MachineDeployment",
                id: "machineDeployment",
                category: "worker",
                children: [
                  {
                    name: this.$route.params.id + "",
                    kind: "MachineSet",
                    id: "machineSet",
                    category: "worker",
                    children: [
                      {
                        name: this.$route.params.id + "-md-1",
                        kind: "Machine",
                        id: "machine1",
                        category: "worker",
                        children: [
                          {
                            name: this.$route.params.id + "-md-1",
                            kind: "AzureMachine",
                            id: "azureMachine1",
                            category: "worker",
                            children: [],
                          },
                          {
                            name: this.$route.params.id + "-control-plane",
                            kind: "KubeAdmConfig",
                            id: "kubeAdmConfig1",
                            category: "worker",
                            children: [],
                          },
                        ],
                      },
                      {
                        name: "... (3x)",
                        kind: "",
                        id: "machine2",
                        category: "worker",
                        children: [
                          // {
                          //   name: "...",
                          // kind: "",
                          //   id: "azureMachine2",
                          //   category: "worker",
                          //   children: [],
                          // },
                          // {
                          //   name: "...",
                          // kind: "",
                          //   id: "kubeAdmConfig2",
                          //   category: "worker",
                          //   children: [],
                          // },
                        ],
                      },
                      {
                        name: this.$route.params.id + "-md-3",
                        kind: "Machine",
                        id: "machine3",
                        category: "worker",
                        children: [
                          {
                            name: this.$route.params.id + "-md-3",
                            kind: "AzureMachine",
                            id: "azureMachine3",
                            category: "worker",
                            children: [],
                          },
                          {
                            name: "default-3-control-plane",
                            kind: "KubeAdmConfig",
                            id: "kubeAdmConfig3",
                            category: "worker",
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
          {
            name: "cluster-identity",
            kind: "AzureClusterIdentity",
            id: "clusterIdentity",
            category: "none",
            children: [],
          },
        ],
        links: [
          {
            parent: "kubeAdmCtrlPlane",
            child: "azureMachineCtrl",
          },
          {
            parent: "kubeAdmCtrlPlane",
            child: "kubeAdmConfigCtrl",
          },
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
