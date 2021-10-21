<template>
  <div class="treeContainer">
    <h1>Cluster Resource Ownership: {{ this.$route.params.id }}</h1>
    <vue-tree
      id="resourceTree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="false"
    >
      <template v-slot:node="{ node, collapsed }">
        <div
          class="machine"
          v-if="node.id == 'machine1'"
        >
          <span>3x</span>
        </div>
        <div
          class="node"
          :style="{ 
            'background-color': colors[node.provider], 
            border: collapsed ? '2px solid grey' : '',
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
import VueTree from "./VueTree.vue";
// import VueTree from "@ssthouse/vue-tree-chart";

export default {
  name: "TargetCluster",
  components: {
    VueTree,
  },
  methods: {},
  data() {
    return {
      colors: {
        bootstrap: "#ffdf7d",
        capi: "#a8c8ff",
        ctrlPlane: "#daadf0",
        infra: "#bbf895",
        addons: "#ffb786",
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
                provider: "capi",
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
                kind: "KubeadmControlPlane",
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
      treeConfig: { nodeWidth: 160, nodeHeight: 40, levelHeight: 100 },
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    };
  },
};
</script>

<style lang="less" scoped>
#resourceTree {
  width: 100%;
  height: 800px;
  border: 1px solid black;
}

.treeContainer {
  height: 100%;
  width: 100%;
  max-width: 100%;
  margin: 0 !important;
}

.node-slot {
  cursor: default !important;
}

.node {
  width: 150px;
  height: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #dae8fc;
  border-radius: 4px;
  box-shadow: 2px 3px 3px rgba(0, 0, 0, 0.3);

  p {
    font-size: 10px;
    margin: 2px;
    color: #2c3e50;
  }

  .node-router-link {
    text-decoration: none;
  }

  .name {
    font-style: italic;
  }

  .kind {
    font-size: 12px;
  }
}

.legend-entry {
  display: inline-block;
  margin-right: 10px;

  div {
    display: inline-block;
    border: 1px solid black;
    margin: 0 5px;
    width: 12px;
    height: 12px;
  }
}

.machine {
  position: absolute;
  transform: translate(0, 55px);
  width: 375px;
  height: 220px;
  border: 3px solid #a8c8ff;
  box-shadow: 3px 4px 3px rgba(0, 0, 0, 0.3);
  border-radius: 5px;

  span {
    position: absolute;
    bottom: 5px;
    right: 10px;
  }
}
</style>
