<template>
  <div class="treeContainer">
    <vue-tree
      id="overviewTree"
      ref="tree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="false"
      :linkStyle="(isStraight) ? 'straight' : 'curve'"
      @scale="(val) => $emit('scale', val)"
      v-if="treeIsReady"
    >
      <template v-slot:node="{ node, collapsed }">
        <v-hover>
          <template v-slot:default="{ hover }">
            <!-- :to="{ path: 'clusters', params: { name: node.name, namespace: node.namespace }}" -->
            <router-link
              :to="'/cluster?name=' + node.name + '&namespace=' + node.namespace"
              :event="node.isManagement ? '' : 'click' /* disable link on management cluster */"
              class="node-router-link"
            >
              <v-card
                class="node mx-auto transition-swing"
                :elevation="hover ? 6 : 3"
                :style="{ 
                  border: collapsed ? '2px solid grey' : '',
                  // height: (node.isManagement) ? '120px' : '140px',
                  // 'background-color': hover ? '#f0f0f0' : '#fff'
                }"
              >
                <v-card-title>
                  <span class="cardTitle">
                    {{ node.name }}
                  </span>
                  <v-spacer></v-spacer>
                  <v-icon color="blue">
                    mdi-{{ getIcon(node.infrastructureProvider) }}
                  </v-icon>
                </v-card-title>
                <!-- <v-card-subtitle class="cardSubtitle">{{ (node.isManagement) ? "Management Cluster" : "Target Cluster" }}</v-card-subtitle> -->
                <v-card-subtitle class="pb-1">{{ node.namespace }}</v-card-subtitle>

                <!-- <v-card-subtitle v-if="node.isManagement">Management Cluster</v-card-subtitle> -->
                <Phase
                  v-if="!node.isManagement"
                  :phase="node.phase"
                />
                <v-card-actions
                  class="cardActions"
                  v-if="!node.isManagement"
                >
                  <v-card-text class="card-bottom-text">Resources</v-card-text>
                  <v-spacer></v-spacer>
                  <v-icon>mdi-arrow-top-right</v-icon>

                </v-card-actions>

              </v-card>
            </router-link>
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
    <AlertError
      v-model=alert
      :message="errorMessage"
    />
  </div>
</template>

<script>
import VueTree from "./VueTree.vue";
import AlertError from "./AlertError.vue";
import Phase from "./Phase.vue";

export default {
  name: "ManagementClusterTree",
  components: {
    VueTree,
    AlertError,
    Phase,
  },
  props: {
    isStraight: Boolean,
    treeData: Object,
    treeConfig: Object,
    treeIsReady: Boolean,
  },
  data() {
    return {
      alert: false,
      errorMessage: "",
    };
  },
  methods: {
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
  background-color: #f8f3f2;
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
  background-color: #fff;

  p {
    font-size: 12px;
    margin: 2px;
  }

  .cardTitle {
    max-width: 194px;
    white-space: nowrap;
    display: inline-block;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .cardSubtitle {
    padding-bottom: 0;
  }

  .cardActions {
    padding-top: 0;
    padding-right: 12px;

    .card-bottom-text {
      padding-top: 8px;
      padding-left: 8px;
    }
  }
}

.node-router-link {
  text-decoration: none;
  // font-style: italic;
}
</style>
