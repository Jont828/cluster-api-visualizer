<template>
  <div class="treeContainer">
    <vue-tree
      id="overviewTree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="false"
      :linkStyle="(isStraight) ? 'straight' : 'curve'"
      v-if="treeIsReady"
    >
      <template v-slot:node="{ node, collapsed }">
        <v-hover>
          <template v-slot:default="{ hover }">
            <router-link
              :to="'/target-cluster/' + node.name"
              class="node-router-link"
            >
              <v-card
                class="node mx-auto transition-swing"
                :elevation="hover ? 6 : 3"
                :style="{ 
                  border: collapsed ? '2px solid grey' : '',
                  // 'background-color': hover ? '#f0f0f0' : '#fff'
                }"
              >
                <v-card-title>
                  <span class="cardTitle">
                    {{ node.name }}
                  </span>
                  <v-spacer></v-spacer>
                  <v-icon color="blue">
                    mdi-{{node.icon}}
                  </v-icon>
                </v-card-title>
                <v-card-subtitle class="cardSubtitle">{{ (node.children.length) ? "Management Cluster" : "Target Cluster" }}</v-card-subtitle>
                <v-card-actions
                  class="cardActions"
                  v-if="!node.isRoot"
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
    <AlertError :message="errorMessage" />
  </div>
</template>

<script>
import VueTree from "./VueTree.vue";
import AlertError from "./AlertError.vue";
import { getClusterOverview } from "../services/Service.js";

export default {
  name: "Tree",
  components: {
    VueTree,
    AlertError,
  },
  props: {
    isStraight: Boolean,
  },
  methods: {
    async fetchOverview() {
      try {
        const response = await getClusterOverview();
        this.treeData = response;
        this.treeIsReady = true;
      } catch (error) {
        this.errorMessage = "Failed to construct cluster overview";
        console.log(error);
      }
    },
  },
  async beforeMount() {
    await this.fetchOverview();
  },
  data() {
    return {
      errorMessage: "",
      treeIsReady: false,
      treeData: {},
      treeConfig: { nodeWidth: 300, nodeHeight: 120, levelHeight: 200 },
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    };
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

.node-slot {
  // cursor: default !important;
}

.node {
  // cursor: default !important;
  width: 250px;
  height: 120px;
  /* padding: 8px; */
  // display: flex;
  // flex-direction: column;
  // align-items: center;
  // justify-content: center;
  background-color: #fff;

  p {
    font-size: 12px;
    margin: 2px;
    // color: #2c3e50;
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
      padding-left: 8px;
    }
  }
}

.node-router-link {
  text-decoration: none;
  // font-style: italic;
}
</style>
