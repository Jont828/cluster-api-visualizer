<template>
  <div class="treeContainer">
    <h1>Overview of Clusters</h1>
    <vue-tree
      id="overviewTree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="false"
    >
      <template v-slot:node="{ node, collapsed }">
        <div
          class="node"
          :style="{ border: collapsed ? '2px solid grey' : '' }"
        >
          <p>{{ node.name }}</p>
          <router-link
            :to="'/target-cluster/' + node.name"
            class="node-router-link"
            v-if="!node.isRoot"
          >
            <p>View Representation</p>
          </router-link>
          <!-- <router-link
            :to="'/'"
            class="node-router-link"
          >
            <p>View Cluster</p>
          </router-link> -->
        </div>
      </template>
    </vue-tree>
  </div>
</template>

<script>
import VueTree from "./VueTree.vue";

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
        provider: "Local",
        children: [
          {
            name: "default-1",
            provider: "Azure",
            children: [],
          },
          {
            name: "public-cluster",
            provider: "Azure",
            children: [
              {
                name: "private-cluster",
                provider: "Azure",
                children: [],
              },
            ],
          },
          {
            name: "default-2",
            provider: "Azure",
            children: [],
          },
        ],
      },
      treeConfig: { nodeWidth: 250, nodeHeight: 80, levelHeight: 200 },
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    };
  },
};
</script>

<style lang="less" scoped>
#overviewTree {
  width: 80%;
  height: 80%;
  border: 1px solid black;
}

.treeContainer {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.node-slot {
  cursor: default !important;
}

.node {
  cursor: default !important;
  width: 140px;
  height: 60px;
  /* padding: 8px; */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #a8c8ff;
  border-radius: 4px;
  box-shadow: 2px 3px 3px rgba(0, 0, 0, 0.3);

  p {
    font-size: 12px;
    margin: 2px;
    color: #2c3e50;
  }
}

.node-router-link {
  text-decoration: none;
  font-style: italic;
}
</style>
