<template>
  <div class="wrapper">
    <AppBar
      :title="'Cluster Resources: ' + this.$route.params.id"
      :showBack="true"
      :isStraight="this.isStraight"
      @togglePathStyle="linkHandler"
    />
    <div
      id="chartLoadWrapper"
      v-if="treeIsReady"
    >
      <div
        id="treeChartWrapper"
        :style="{
          height: Object.keys(selected).length == 0 ? '100%' : 'calc(100% - 84px)' 
        }"
      >
        <vue-tree
          id="resourceTree"
          :dataset="treeData"
          :config="treeConfig"
          :collapse-enabled="true"
          :linkStyle="(isStraight) ? 'straight' : 'curve'"
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
                  :elevation="hover ? 6 : 3"
                  :style="{ 
                  'background-color': legend[node.provider].color, 
                // 'background-color': legend[node.provider][hover ? 'hoverColor' : 'color'], 
                border: collapsed ? '' : '',
              }"
                  v-on:click="selectNode(node)"
                >
                  <!-- <router-link
                :to="'/'"
                class="node-router-link"
              > -->
                  <p class="kind font-weight-medium">{{ node.kind }}</p>
                  <p
                    class="name font-italic"
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
      </div>

      <div
        class="resourceView"
        v-if="resourceIsReady && this.selected.name"
      >
        <CustomResourceTree
          :items="resource"
          :title="'Resource: ' + selected.kind + '/' + selected.name"
          :color="legend[selected.provider].color"
          :selectedNode="this.selected.name"
          @unselectNode="(val) => { this.selected=val; }"
        />

      </div>
    </div>
    <div
      id="resourceTree"
      class="spinner"
      v-else
    >
      <v-progress-circular
        :size="50"
        :width="5"
        indeterminate
        color="primary"
      ></v-progress-circular>
    </div>
    <AlertError :message="errorMessage" />
  </div>
</template>

<script>
/* eslint-disable */
import VueTree from "../components/VueTree.vue";
import AppBar from "../components/AppBar.vue";
import CustomResourceTree from "../components/CustomResourceTree.vue";
import AlertError from "../components/AlertError.vue";

import colors from "vuetify/lib/util/colors";

import { getCluster, getClusterResource } from "../services/Service.js";

export default {
  name: "TargetCluster",
  components: {
    VueTree,
    AppBar,
    CustomResourceTree,
    AlertError,
  },
  methods: {
    linkHandler(val) {
      this.isStraight = val;
    },
    async selectNode(node) {
      if (node.provider == "") return;
      this.selected = node;
      try {
        const response = await getClusterResource(
          this.selected.group,
          this.selected.plural.toLowerCase(),
          this.selected.name
        );
        // console.log(JSON.stringify(response));
        this.resource = response;
        this.resourceIsReady = true;
      } catch (error) {
        this.errorMessage =
          "Failed to fetch CRD for `" +
          this.selected.kind +
          this.selected.name +
          "`";
        console.log("Error fetching CRD");
        console.log(error);
      }
    },
    async fetchCluster() {
      try {
        const response = await getCluster(this.$route.params.id);
        this.treeData = response;
        this.treeIsReady = true;
      } catch (error) {
        this.errorMessage =
          "Failed to fetch resources for cluster `" +
          this.$route.params.id +
          "`";
        console.log(error);
      }
    },
  },
  async beforeMount() {
    await this.fetchCluster();
  },
  data() {
    return {
      errorMessage: "",
      treeIsReady: false,
      resourceIsReady: false,
      resource: [],
      selected: {},
      isStraight: false,
      legend: {
        bootstrap: {
          name: "Bootstrap Provider (Kubeadm)",
          color: colors.amber.darken2,
          hoverColor: colors.amber.darken3,
        },
        controlplane: {
          name: "Control Plane (Kubeadm)",
          color: colors.purple.darken1,
          hoverColor: colors.purple.darken2,
        },
        infrastructure: {
          name: "Infrastructure (Azure)",
          color: colors.green.base,
          hoverColor: colors.green.darken1,
        },
        cluster: {
          name: "Cluster API",
          color: colors.blue.darken1,
          hoverColor: colors.blue.darken2,
        },
        addons: {
          name: "Addons",
          color: colors.red.darken1,
          hoverColor: colors.red.darken2,
        },
        "": {
          name: "None",
          color: colors.grey.darken1,
          hoverColor: colors.grey.darken2,
        },
      },
      treeData: {},
      treeConfig: { nodeWidth: 180, nodeHeight: 50, levelHeight: 120 },
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    };
  },
};
</script>

<style lang="less" scoped>
.wrapper {
  height: 100%;
  width: 100%;
  max-width: 100%;
  margin: 0 !important;
}

.node-slot {
  cursor: default !important;
}

.node {
  width: 170px;
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
    font-size: 12.5px;
  }

  .name,
  .kind {
    max-width: 160px;
    text-align: center;
    white-space: nowrap;
    display: inline-block;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

#chartLoadWrapper {
  height: 100%;

  #treeChartWrapper {
    width: 100%;
    // height: 100%;
    position: relative;
    text-align: center;

    #resourceTree {
      width: 100%;
      // height: 100%;
      height: 100%;
      background-color: #f8f3f2;
      // border: 1px solid black;
    }
  }
}

.legend {
  text-align: center;
  // display: inline-block;
  position: absolute;
  bottom: 30px;
  z-index: 9999;
  width: 100%;

  .legend-card {
    padding: 10px 10px;
    display: inline-block;

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

.resourceView {
  margin: 0 30px;
  padding-bottom: 30px;
}
</style>

<style lang="less">
.spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
}
</style>
