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
  async beforeMount() {
    await this.fetchOverview();
  },
  mounted() {
    document.title = "Management Cluster Overview";
    setInterval(
      function () {
        console.log("Polling...");
        this.fetchOverview();
      }.bind(this),
      1000 * 3
    );
  },
  data() {
    return {
      isStraight: false,
      treeConfig: { nodeWidth: 300, nodeHeight: 120, levelHeight: 200 },
      treeData: {},
      cachedTreeString: "",
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

        if (response.data == null) {
          this.errorMessage =
            "Couldn't find a management cluster from default kubeconfig";
          return;
        }

        console.log("Cluster overview data:", response.data);
        if (this.cachedTreeString !== JSON.stringify(response.data)) {
          this.treeData = response.data;
          this.cachedTreeString = JSON.stringify(response.data);
          this.treeIsReady = true;
        }
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
};
</script>

<style lang="less" scoped>
#overview {
  height: 100%;
}
</style>