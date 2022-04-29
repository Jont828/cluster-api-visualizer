<template>
  <div id="overview">
    <AppBar
      title="Management Cluster Overview"
      :isStraight="this.isStraight"
      :scale="scale"
      @togglePathStyle="linkHandler"
      @reload="fetchOverview"
      @zoomIn="() => { $refs.overviewTree.$refs.tree.zoomIn() }"
      @zoomOut="() => { $refs.overviewTree.$refs.tree.zoomOut() }"
    />
    <OverviewTree
      ref="overviewTree"
      :isStraight="this.isStraight"
      :treeConfig="treeConfig"
      :treeData="treeData"
      :treeIsReady="treeIsReady"
      @scale="(val) => { scale = val }"
    />
  </div>
</template>

<script>
import Vue from "vue";

import OverviewTree from "../components/OverviewTree.vue";
import AppBar from "../components/AppBar.vue";

export default {
  name: "Overview",
  components: {
    OverviewTree,
    AppBar,
  },
  mounted() {
    document.title = "Management Cluster Overview";
  },
  data() {
    return {
      isStraight: false,
      treeConfig: { nodeWidth: 300, nodeHeight: 120, levelHeight: 200 },
      treeData: {},
      treeIsReady: false,
      scale: 1,
    };
  },
  methods: {
    linkHandler(val) {
      this.isStraight = val;
    },
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
        console.log("Error:", error.toJSON());
        this.alert = true;
        if (error.response) {
          if (error.response.status == 404) {
            this.errorMessage =
              "Management cluster not found, is the kubeconfig set?";
          } else {
            this.errorMessage =
              "Unable to load management cluster and workload clusters";
          }
        } else if (error.request) {
          this.errorMessage = "No server response received";
        } else {
          this.errorMessage = "Unable to create request";
        }
      }
    },
  },
  async beforeMount() {
    await this.fetchOverview();
  },
};
</script>

<style lang="less" scoped>
#overview {
  height: 100%;
}
</style>