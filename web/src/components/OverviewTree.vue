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
              :to="node.isManagement ? '#' : ('/target-cluster/' + node.name)"
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
                <v-card-subtitle class="cardSubtitle">{{ (node.isManagement) ? "Management Cluster" : "Target Cluster" }}</v-card-subtitle>
                <v-card-actions class="cardActions">
                  <!-- v-if="!node.isRoot" -->
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
import Vue from "vue";
import VueTree from "./VueTree.vue";
import AlertError from "./AlertError.vue";
// import { getClusterOverview } from "../services/Service.js";

export default {
  name: "OverviewTree",
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
        const response = await Vue.axios.get("/multicluster");
        this.treeData = response.data;
        console.log(this.treeData);
        if (this.treeData == null) {
          this.errorMessage =
            "Couldn't find a management cluster from default kubeconfig";
          return;
        }
        this.treeIsReady = true;
      } catch (error) {
        console.log(error);
        if (error.response) {
          // The request was made and the server responded with a status code that falls out of the range of 2xx
          // this.errorMessage = error.response.data;
          console.log("Error Status:", error.response.status);
          console.log("Error Data:", error.response.data);
          console.log("Error Headers:", error.response.headers);
          this.errorMessage = error.response.data;
        } else if (error.request) {
          // The request was made but no response was received
          // `error.request` is an instance of XMLHttpRequest in the browser and an instance of
          // http.ClientRequest in node.js
          console.log("Error Request:", error.request);
          this.errorMessage = "No server response received";
        } else {
          // Something happened in setting up the request that triggered an Error
          console.log("Error message:", error.message);
          this.errorMessage = error.message;
        }
        console.log("Error Config:", error.config);
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

.node {
  width: 250px;
  height: 120px;
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
      padding-left: 8px;
    }
  }
}

.node-router-link {
  text-decoration: none;
  // font-style: italic;
}
</style>
