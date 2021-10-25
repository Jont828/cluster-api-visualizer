<template>
  <div class="treeContainer">
    <vue-tree
      id="overviewTree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="false"
    >
      <template v-slot:node="{ node, collapsed }">
        <router-link
          :to="'/target-cluster/' + node.name"
          class="node-router-link"
        >
          <v-card
            class="node"
            :style="{ 
            'background-color': '#fff', 
            border: collapsed ? '2px solid grey' : '',
          }"
          >
            <v-card-title>{{ node.name }}
              <v-spacer></v-spacer>
              <v-icon color="blue">
                mdi-{{node.icon}}
              </v-icon>
            </v-card-title>
            <v-card-actions v-if="!node.isRoot">
              <v-card-text class="card-bottom-text">Resources</v-card-text>
              <v-spacer></v-spacer>
              <v-icon>mdi-arrow-top-right</v-icon>

            </v-card-actions>

          </v-card>
        </router-link>
      </template>
    </vue-tree>
  </div>
</template>

<script>
import VueTree from "./VueTree.vue";

import colors from "vuetify/lib/util/colors";

export default {
  name: "Tree",
  components: {
    VueTree,
  },
  data() {
    return {
      treeData: {
        name: "kind-capz",
        isRoot: true,
        icon: "kubernetes",
        children: [
          {
            name: "default-1",
            icon: "microsoft-azure",
            children: [],
          },
          {
            name: "public-cluster",
            icon: "microsoft-azure",
            children: [
              {
                name: "private-cluster",
                icon: "microsoft-azure",
                children: [],
              },
            ],
          },
          {
            name: "default-2",
            icon: "microsoft-azure",
            children: [],
          },
        ],
      },
      treeConfig: { nodeWidth: 250, nodeHeight: 120, levelHeight: 200 },
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
  width: 200px;
  height: 120px;
  /* padding: 8px; */
  // display: flex;
  // flex-direction: column;
  // align-items: center;
  // justify-content: center;
  background-color: #a8c8ff;

  p {
    font-size: 12px;
    margin: 2px;
    // color: #2c3e50;
  }

  .card-bottom-text {
    padding-left: 8px;
  }
}

.node-router-link {
  text-decoration: none;
  // font-style: italic;
}
</style>
