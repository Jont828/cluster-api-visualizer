<template>

  <div class="container">
    <h1>Management Cluster: {{ this.$route.params.id }}</h1>
    <vue-tree
      id="tree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="false"
    >
      <template v-slot:node="{ node, collapsed }">
        <div
          class="node"
          :style="{ border: collapsed ? '2px solid grey' : '' }"
        >
          <router-link :to="'/management-cluster/' + node.name">
            <span>{{ node.name }}</span>
          </router-link>
          <router-link
            :to="'/target-cluster/' + node.name"
            v-if="!node.children.length"
          >
            <span>Managed resources</span>
          </router-link>
        </div>
      </template>
    </vue-tree>
    <b-container>
      <!-- <h1>Managed resources:</h1> -->
    </b-container>
  </div>

</template>

<script>
import VueTree from "@ssthouse/vue-tree-chart";

export default {
  name: "ManagementCluster",
  components: {
    VueTree,
  },
  data() {
    return {
      treeData: {
        name: "kind-capz",
        provider: "Local",
        status: "Provisioned",
        children: [
          {
            name: "default-1",
            provider: "Azure",
            status: "Provisioned",
            children: [],
          },
          {
            name: "public-cluster",
            provider: "Azure",
            status: "Provisioned",
            children: [],
          },
          {
            name: "default-2",
            provider: "Azure",
            status: "Provisioned",
            children: [],
          },
        ],
      },
      treeConfig: { nodeWidth: 250, nodeHeight: 100, levelHeight: 250 },
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    };
  },
};
</script>

<style scoped>
#tree {
  width: 80%;
  height: 1000px;
  border: 1px solid black;
}

.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.node-slot {
  cursor: default !important;
}

.node {
  width: 200px;
  height: 120px;
  /* padding: 8px; */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #2c3e50;
  background-color: #dae8fc;
  border-radius: 4px;
}

.node p {
  margin: 2px;
}
</style>