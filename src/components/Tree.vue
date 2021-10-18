<template>
  <div class="container">
    <vue-tree
      id="tree"
      :dataset="treeData"
      :config="treeConfig"
    >
      <template v-slot:node="{ node, collapsed }">
        <div
          class="rich-media-node"
          :style="{ border: collapsed ? '2px solid grey' : '' }"
        >
          <!-- <img
            :src="node.avatar"
            style="width: 48px; height: 48px; border-raduis: 4px;"
          /> -->
          <span>{{ node.name }}</span>
        </div>
      </template>
    </vue-tree>
  </div>  
</template>

<script>
import VueTree from '@ssthouse/vue-tree-chart'

export default {
  name: 'Tree',
  components: {
    VueTree
  },
  data() {
    return {
      treeData: {
        name: 'kind-capz',
        value: 800,
        children: [
          {
            name: 'default-27983',
            value: 400,
            children: []
          },
          {
            name: 'Smith',
            value: 200,
            children: [
              {
                name: 'S1',
                value: 230,
              }
            ]
          },
          {
            name: 'Jackson',
            value: 300,
          }
        ]
      },
      treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    }
  }
}
</script>

<style scoped >
#tree {
  width: 80%;
  height: 80%;
  border: 1px solid black;
}

.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
}

.rich-media-node {
  width: 150px;
  height: 150px;
  /* padding: 8px; */
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  color: #2c3e50;
  background-color: #DAE8FC;
  border-radius: 4px;
}
</style>
