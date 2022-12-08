<template>
  <div id="overview">
    <AppBar
      title="Management Cluster Overview"
      :isStraight="this.isStraight"
      :scale="scale"
      @togglePathStyle="linkHandler"
      @reload="fetchOverview(forceRedraw=true)"
      @zoomIn="() => { $refs.overviewTree.$refs.tree.zoomIn() }"
      @zoomOut="() => { $refs.overviewTree.$refs.tree.zoomOut() }"
      @toggleDrawer="() => { drawer = !drawer }"
    />
    <ManagementClusterTree
      ref="overviewTree"
      :isStraight="this.isStraight"
      :treeConfig="treeConfig"
      :treeData="treeData"
      :treeIsReady="treeIsReady"
      @scale="(val) => { scale = val }"
    />
    <v-navigation-drawer
      v-model="drawer"
      absolute
      temporary
    >
      <v-list nav>
        <v-list-item link>
          <v-list-item-icon>
            <v-icon>mdi-information</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>About</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item
          link
          @click="() => { showSettingsOverlay = !showSettingsOverlay; drawer = false }"
        >
          <v-list-item-icon>
            <v-icon>mdi-cog</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>Settings</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
      <!-- <v-list nav>
        <v-list-item
          v-for="([icon, text], i) in items"
          :key="i"
          link
        >
          <v-list-item-icon>
            <v-icon>mdi-{{ icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ text }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list> -->
    </v-navigation-drawer>

    <v-overlay
      absolute
      :value="showSettingsOverlay"
      z-index="99999"
      light
    >
      <SettingsCard
        @close="() => { showSettingsOverlay = !showSettingsOverlay }"
        class="settingsCard"
      />
    </v-overlay>
  </div>
</template>

<script>
import Vue from "vue";

import ManagementClusterTree from "../components/ManagementClusterTree.vue";
import SettingsCard from "../components/SettingsCard.vue";
import AppBar from "../components/AppBar.vue";

export default {
  name: "ManagementCluster",
  components: {
    ManagementClusterTree,
    SettingsCard,
    AppBar,
  },
  async beforeMount() {
    await this.fetchOverview();
  },
  mounted() {
    document.title = "Management Cluster Overview";
    const reloadTime = 60 * 1000;
    this.polling = setInterval(
      function () {
        this.fetchOverview();
      }.bind(this),
      reloadTime
    );
  },
  beforeDestroy() {
    this.selected = {};
    clearInterval(this.polling);
  },
  data() {
    return {
      showSettingsOverlay: false,
      // items: [
      //   ["information", "About"],
      //   ["cog", "Settings"],
      // ],
      drawer: false,
      polling: null,
      isStraight: false,
      treeConfig: { nodeWidth: 300, nodeHeight: 140, levelHeight: 275 },
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
    async fetchOverview(forceRedraw = false) {
      try {
        const response = await Vue.axios.get("/management-cluster");

        if (response.data == null) {
          this.errorMessage =
            "Couldn't find a management cluster from default kubeconfig";
          return;
        }

        console.log("Cluster overview data:", response.data);
        if (
          forceRedraw ||
          this.cachedTreeString !== JSON.stringify(response.data)
        ) {
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

.settingsCard {
  // width: 500px;
}
</style>