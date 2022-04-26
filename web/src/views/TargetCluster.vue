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
      <TargetClusterTree
        id="targetTree"
        :treeConfig="treeConfig"
        :treeData="treeData"
        :isStraight="isStraight"
        :legend="legend"
        @selectNode="selectNodeHandler"
        :style="{
          height: Object.keys(selected).length == 0 ? '100%' : 'calc(100% - 84px)' 
        }"
      />

      <div
        class="resourceView"
        v-if="resourceIsReady && this.selected.name"
      >
        <CustomResourceTree
          :items="treeviewResource"
          :jsonItems="resource"
          :name="selected.kind + '/' + selected.name"
          :color="legend[selected.provider].color"
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
    <AlertError
      v-model="alert"
      :message="errorMessage"
    />
  </div>
</template>

<script>
/* eslint-disable */
import Vue from "vue";
import TargetClusterTree from "../components/TargetClusterTree.vue";
import AppBar from "../components/AppBar.vue";
import CustomResourceTree from "../components/CustomResourceTree.vue";
import AlertError from "../components/AlertError.vue";

import colors from "vuetify/lib/util/colors";

export default {
  name: "TargetCluster",
  components: {
    TargetClusterTree,
    AppBar,
    CustomResourceTree,
    AlertError,
  },
  mounted() {
    document.title = "Cluster Resources: " + this.$route.params.id;
  },
  methods: {
    linkHandler(val) {
      this.isStraight = val;
    },
    async selectNodeHandler(node) {
      try {
        const params = new URLSearchParams();
        params.append("kind", node.kind);
        params.append("apiVersion", node.group + "/" + node.version);
        params.append("name", node.name);

        const response = await Vue.axios.get("/custom-resource", {
          params: params,
        });
        console.log(response.data);
        this.resource = response.data;
        this.treeviewResource = this.formatToTreeview(response.data);
        this.selected = node; // Don't select until an error won't pop up
        this.resourceIsReady = true;
      } catch (error) {
        console.log("Error:", error.toJSON());
        this.alert = true;
        if (error.response) {
          if (error.response.status == 404) {
            this.errorMessage =
              "Custom Resource Definition `" +
              node.kind +
              "/" +
              node.name +
              "` not found";
          } else {
            this.errorMessage =
              "Unable to load Custom Resource Definition `" +
              node.kind +
              "/" +
              node.name +
              "`";
          }
        } else if (error.request) {
          this.errorMessage = "No server response received";
        } else {
          this.errorMessage = "Unable to create request";
        }
      }
    },
    async fetchCluster() {
      try {
        // const response = await getCluster(this.$route.params.id);
        const response = await Vue.axios.get(
          "/cluster-resources/" + this.$route.params.id
        );
        this.treeData = response.data;
        console.log(this.treeData);
        this.treeIsReady = true;
      } catch (error) {
        console.log("Error:", error.toJSON());
        this.alert = true;
        if (error.response) {
          if (error.response.status == 404) {
            this.errorMessage =
              "Cluster `" + this.$route.params.id + "` not found";
          } else {
            this.errorMessage =
              "Failed to construct object tree for cluster `" +
              this.$route.params.id +
              "`";
          }
        } else if (error.request) {
          this.errorMessage = "No server response received";
        } else {
          this.errorMessage = "Unable to create request";
        }
      }
    },
    formatToTreeview(resource, id = 0) {
      let result = [];
      if (typeof resource == "string") {
        return [{ name: resource }];
      } else if (Array.isArray(resource)) {
        let children = [];
        resource.forEach((e, i) => {
          result.push({
            id: id++,
            name: i.toString(),
            children: this.formatToTreeview(e, id),
          });
        });
      } else {
        // isObject
        Object.entries(resource).forEach(([key, value]) => {
          let name = "";
          let children = [];
          if (typeof value == "string") {
            name = key + ": " + value;
          } else {
            name = key;
            children = this.formatToTreeview(value, id);
          }
          result.push({
            id: id++,
            name: name,
            children: children,
          });
        });
      }

      return result;
    },
  },
  async beforeMount() {
    await this.fetchCluster();
  },
  data() {
    return {
      alert: false,
      errorMessage: "",
      treeIsReady: false,
      resourceIsReady: false,
      resource: [],
      selected: {},
      isStraight: false,
      treeData: {},
      treeConfig: { nodeWidth: 180, nodeHeight: 50, levelHeight: 120 },
      legend: {
        bootstrap: {
          name: "Bootstrap Provider",
          color: colors.amber.darken2,
          // altColor: colors.amber.darken1,
        },
        controlplane: {
          name: "Control Plane Provider",
          color: colors.purple.darken1,
          // altColor: colors.purple.lighten1,
        },
        infrastructure: {
          name: "Infrastructure Provider",
          color: colors.green.base,
          // altColor: colors.green.lighten1,
        },
        cluster: {
          name: "Cluster API",
          color: colors.blue.darken1,
          // altColor: colors.blue.lighten1,
        },
        addons: {
          name: "Addons",
          color: colors.red.darken1,
          // altColor: colors.red.lighten2,
        },
        virtual: {
          name: "None",
          color: colors.grey.darken1,
          // altColor: colors.grey.base,
        },
      },
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

#chartLoadWrapper {
  height: 100%;

  #treeChartWrapper {
    width: 100%;
    height: 100%;
    position: relative;
    text-align: center;
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
