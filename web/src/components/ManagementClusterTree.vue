<template>
  <div class="treeContainer">
    <vue-tree
      id="overviewTree"
      ref="tree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="false"
      :linkStyle="(store.straightLinks) ? 'straight' : 'curve'"
      @scale="(val) => $emit('scale', val)"
      v-if="treeIsReady"
    >
      <template v-slot:node="{ node }">
        <v-hover>

            <template v-slot:default="{ hover }">
              <div style="padding-left: 10px; padding-right: 10px;">
              <!-- :to="{ path: 'clusters', params: { name: node.name, namespace: node.namespace }}" -->
              <v-card
                  class="mx-auto transition-swing"
                  :elevation="hover ? 6 : 3"
                  :style="($vuetify.theme.dark) ? {
                  'background-color': hover ? '#383838' : '#272727',
                } : null"
              >
                <v-card-title>
                  <span class="cardTitle text-truncate">
                    {{ node.name }}
                  </span>
                  <v-spacer></v-spacer>
                  <v-icon color="primary">
                    mdi-{{ getIcon(node.infrastructureProvider) }}
                  </v-icon>
                </v-card-title>
                <!-- <v-card-subtitle class="cardSubtitle">{{ (node.isManagement) ? "Management Cluster" : "Target Cluster" }}</v-card-subtitle> -->
                <v-card-subtitle class="pb-1 text-truncate">{{ (node.namespace == "") ? "default" : node.namespace }}</v-card-subtitle>

                <!-- <v-card-subtitle v-if="node.isManagement">Management Cluster</v-card-subtitle> -->
                <ClusterPhase
                    v-if="!node.isManagement"
                    :phase="node.phase"
                />
                <v-card-actions :class="[ 'cardActions', (node.isManagement) ? 'pt-8' : 'pt-2' ]">
                  <v-col>
                    <router-link
                        :to="'/cluster?name=' + node.name + '&namespace=' + node.namespace"
                        :is="node.isManagement ? 'span' : 'router-link'"
                        :event="node.isManagement ? '' : 'click' /* disable link on management cluster */"
                        class="node-router-link"
                    >
                      <v-card-text class="card-bottom-text">{{ (node.isManagement) ? 'Management Cluster' : 'View Workload Cluster' }}
                        <span v-if="!node.isManagement">
                    <v-icon>mdi-arrow-top-right</v-icon>
                  </span>
                      </v-card-text>
                    </router-link>
                  </v-col>
                  <v-col v-if="!node.isManagement">
                    <div v-if="showLens">
                      <a v-on:click="openLens(node)" v-bind:href="'api/v1/cluster-kubeconfig/?name=' + node.name + '&namespace=' + node.namespace"
                          v-bind:download="node.namespace + '.' + node.name + '.kubeconfig' ">
                        <v-img src="../assets/lens.png" max-width="120" />
                      </a>
                    </div>
                    <div v-if="!showLens">
                      <a v-bind:href="'api/v1/cluster-kubeconfig/?name=' + node.name + '&namespace=' + node.namespace" v-bind:download="node.namespace + '.' + node.name + '.kubeconfig' ">
                        <img src="../assets/download.png" width="19px"  style="float: left; padding-right: 5px" />
                        Kubeconfig</a>
                    </div>
                  </v-col>
                </v-card-actions>

              </v-card>
              </div>
            </template>

        </v-hover>
      </template>
    </vue-tree>
    <div
      id="overviewTree"
      class="spinner"
      v-else
    >
      <v-progress-circular
        indeterminate
        color="primary"
      ></v-progress-circular>
    </div>
  </div>
</template>

<script>
import VueTree from "./VueTree.vue";
import ClusterPhase from "./ClusterPhase.vue";

import { useSettingsStore } from "../stores/settings.js";

export default {
  name: "ManagementClusterTree",
  computed: {
    console: () => console,
    window: () => window,
  },
  components: {
    VueTree,
    ClusterPhase,
  },
  props: {
    treeData: Object,
    treeConfig: Object,
    treeIsReady: Boolean,
    showLens: Boolean,
  },
  data() {
    return {
      alert: false,
      errorMessage: "",
    };
  },
  setup() {
    const store = useSettingsStore();
    return { store };
  },
  methods: {
    openLens(node) {
      window.open('lens://app/open/direct/' + node.clusterUrl + '/cluster/pods')
    },
    getIcon(provider) {
      switch (provider) {
        case "AzureCluster":
          return "microsoft-azure";
        case "DockerCluster":
          return "docker";
        case "GCPCluster":
          return "google-cloud";
        case "AWSCluster":
          return "aws";
        default:
          return "kubernetes";
      }
    },
  },
};
</script>

<style lang="less" scoped>
#overviewTree {
  width: 100%;
  height: 100%;
}

.treeContainer {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.node {
  width: 250px;
  height: 140px;
  background: rgba(50, 52, 59, 1);

  p {
    font-size: 12px;
    margin: 2px;
  }

  .cardTitle {
    max-width: 194px;
    display: inline-block;
  }

  .cardSubtitle {
    padding-bottom: 0;
  }

  .cardActions {
    padding-right: 12px;

    .card-bottom-text {
      padding-top: 0px;
      padding-bottom: 0;
      padding-left: 8px;
    }
  }
}

.node-router-link {
  text-decoration: none;
  // font-style: italic;
}
</style>
