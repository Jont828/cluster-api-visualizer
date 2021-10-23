<template>
  <div class="treeContainer">
    <AppBar :title="'Cluster Resources: ' + this.$route.params.id" />

    <vue-tree
      id="resourceTree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="true"
    >
      <template v-slot:node="{ node, collapsed }">
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
              :elevation="hover ? 12 : 3"
              :style="{ 
                'background-color': legend[node.provider].color, 
                border: collapsed ? '2px solid grey' : '',
              }"
              v-on:click="selectNode(node)"
            >
              <!-- <router-link
                :to="'/'"
                class="node-router-link"
              > -->
              <p class="kind font-weight-medium">{{ node.kind }}</p>
              <p
                class="font-italic"
                v-if="node.name"
              >{{ node.name }}</p>
              <p
                class="chevron"
                v-else-if="collapsed"
              >&#9660;</p>
              <p
                class="chevron"
                v-else
              >&#9650;</p>
              <!-- </router-link> -->
            </v-card>
          </template>
        </v-hover>
      </template>
    </vue-tree>
    <div class="legend">
      <div
        class="legend-entry"
        v-for="(entry, provider) in this.legend"
        :key="provider"
      >
        <div :style="{
          'background-color': entry.color
        }" />
        <span>{{ entry.name }}</span>
      </div>
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
          color: colors.yellow.lighten2,
        },
        ctrlPlane: {
          name: "Control Plane (Kubeadm)",
          color: colors.purple.lighten3,
        },
        infra: {
          name: "Infrastructure (Azure)",
          color: colors.green.lighten2,
        },
        capi: {
          name: "Cluster API",
          color: colors.blue.lighten3,
        },
        addons: {
          name: "Addons",
          color: colors.red.lighten2,
        },
        "": {
          name: "None",
          color: colors.grey.lighten2,
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
      treeConfig: { nodeWidth: 160, nodeHeight: 50, levelHeight: 120 },
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    };
  },
};
</script>

<style lang="less" scoped>
#resourceTree {
  width: 100%;
  height: 100%;
  // border: 1px solid black;
}

.treeContainer {
  height: 750px;
  width: 100%;
  max-width: 100%;
  margin: 0 !important;
}

.node-slot {
  cursor: default !important;
}

.node {
  width: 150px;
  height: 50px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  // background-color: #dae8fc;
  // border-radius: 4px;
  // box-shadow: 2px 3px 3px rgba(0, 0, 0, 0.3);

  p {
    font-size: 10px;
    margin: 2px;
  }

  p.chevron {
    margin: 0;
  }

  .node-router-link {
    text-decoration: none;
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

.left {
  text-align: left;
}
</style>
