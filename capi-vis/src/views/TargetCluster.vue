<template>
  <div class="treeContainer">
    <AppBar
      :title="'Cluster Resources: ' + this.$route.params.id"
      :showBack="true"
    />

    <vue-tree
      id="resourceTree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="true"
    >
      <template v-slot:node="{ node }">
        <!-- <template v-slot:node="{ node, collapsed }"> -->
        <div
          class="machine"
          v-if="node.id == 'machine1'"
        >
          <span>x2</span>
        </div>
        <v-hover>
          <template v-slot:default="{ hover }">
            <v-card
              class="node mx-auto transition-swing"
              :elevation="hover ? 8  : 4"
              :style="{ 
                'background-color': legend[node.provider].color, 
              }"
              v-on:click="selectNode(node)"
            >
              <!-- border: collapsed ? '2px solid white' : '', -->
              <!-- <router-link
                :to="'/'"
                class="node-router-link"
              > -->
              <p class="kind font-weight-medium">{{ node.kind }}</p>
              <p
                class="font-italic"
                v-if="node.name"
              >{{ node.name }}</p>
              <v-icon
                class="chevron"
                size="18"
                color="white"
                v-else-if="collapsed"
              >mdi-chevron-down</v-icon>
              <v-icon
                class="chevron"
                size="18"
                color="white"
                v-else
              >mdi-chevron-up</v-icon>
              <!-- </router-link> -->
            </v-card>
          </template>
        </v-hover>
      </template>

    </vue-tree>
    <div class="legend">
      <v-card class="legend-card">
        <div
          class="legend-entry"
          v-for="(entry, provider) in legend"
          :key="provider"
        >
          <div :style="{
            'background-color': entry.color
          }" />
          <span>{{ entry.name }}</span>
        </div>
      </v-card>
    </div>
    <div
      class="left"
      v-if="selected.name"
    >
      <h1>Resource: {{ selected.kind }}/{{ selected.name }} </h1>
      <p>{{ selected.name }} is ...</p>
      <!-- <ul
        v-for="(value, key) in this.crd"
        :key="key"
      >
        <li>{{ key }}: {{ value }}</li>
      </ul> -->
    </div>
  </div>
</template>

<script>
/* eslint-disable */
import VueTree from "../components/VueTree.vue";
import AppBar from "../components/AppBar.vue";
// import AzureCluster from "../assets/yaml/default/azurecluster.yaml";

import colors from "vuetify/lib/util/colors";

import yaml from "js-yaml";
// import VueTree from '@ssthouse/vue-tree-chart';

export default {
  name: "TargetCluster",
  components: {
    VueTree,
    AppBar,
  },
  methods: {
    selectNode(node) {
      this.selected = node;
      // this.crd = yaml.load(AzureCluster);
      // console.log(this.crd);
    },
  },
  data() {
    return {
      crd: "",
      selected: {},
      legend: {
        bootstrap: {
          name: "Bootstrap Provider (Kubeadm)",
          color: colors.amber.darken2,
        },
        ctrlPlane: {
          name: "Control Plane (Kubeadm)",
          color: colors.purple.darken1,
        },
        infra: {
          name: "Infrastructure (Azure)",
          color: colors.green.base,
        },
        capi: {
          name: "Cluster API",
          color: colors.blue.darken1,
        },
        addons: {
          name: "Addons",
          color: colors.red.darken1,
        },
        "": {
          name: "None",
          color: colors.grey.darken1,
        },
      },
      treeData: {
        name: this.$route.params.id,
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
                    name: this.$route.params.id + "",
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
            provider: "",
            collapsable: true,
            children: [
              {
                name: this.$route.params.id + "-control-plane",
                kind: "AzureMachineTemplate",
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
      },
      treeConfig: { nodeWidth: 170, nodeHeight: 50, levelHeight: 120 },
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    };
  },
};
</script>

<style lang="less" scoped>
#resourceTree {
  width: 100%;
  height: 100%;
  background-color: #f8f3f2;
  // border: 1px solid black;
}

.treeContainer {
  height: 750px;
  // height: 100%;
  width: 100%;
  max-width: 100%;
  margin: 0 !important;
}

.node-slot {
  cursor: default !important;
}

.node {
  width: 160px;
  height: 50px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  // background-color: #dae8fc;
  // border-radius: 4px;
  // box-shadow: 2px 3px 3px rgba(0, 0, 0, 0.3);
  color: white;

  p {
    font-size: 11px;
    margin: 2px;
  }

  .chevron {
    margin: 0;
  }

  .node-router-link {
    text-decoration: none;
  }

  .kind {
    font-size: 13px;
  }
}

.legend {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  .legend-card {
    padding: 10px 10px;

    .legend-entry {
      display: inline-block;
      margin-right: 10px;

      div {
        display: inline-block;
        border-radius: 3px;
        // border: 1px solid black;
        margin: 0 5px;
        width: 12px;
        height: 12px;
      }
    }
  }

}

.machine {
  position: absolute;
  transform: translate(0, 65px);
  width: 375px;
  height: 230px;
  border: 3px solid #1e88e5;
  // border: 3px solid #a8c8ff;
  box-shadow: 3px 4px 3px rgba(0, 0, 0, 0.3);
  border-radius: 5px;
  z-index: -10000;

  span {
    position: absolute;
    bottom: 5px;
    right: 10px;
  }
}

.left {
  text-align: left;
}
</style>
