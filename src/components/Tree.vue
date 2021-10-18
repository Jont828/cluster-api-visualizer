<template>
  <div class="container">
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
          <router-link :to="'/cluster/' + node.name">
            <p>{{ node.name }}</p>
          </router-link>
          <p>{{ node.provider }}</p>
          <p>{{ node.status }}</p>
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
        provider: 'Local',
        status: 'Provisioned',
        children: [
          {
            name: 'default-1',
            provider: 'Azure',
            status: 'Provisioned',
            children: []
          },
          {
            name: 'public-cluster',
            provider: 'Azure',
            status: 'Provisioned',
            children: [
              {
                name: 'private-cluster',
                provider: 'Azure',
                status: 'Provisioning',
              }
            ]
          },
          {
            name: 'default-2',
            provider: 'Azure',
            status: 'Provisioned',
          }
        ]
      },
      treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
      // treeConfig: { nodeWidth: 250, nodeHeight: 150, levelHeight: 250 }
    }
  }
}
</script>

<style>
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

.node-slot {
  cursor: default !important;
}

.node {
  width: 150px;
  height: 150px;
  /* padding: 8px; */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #2c3e50;
  background-color: #DAE8FC;
  border-radius: 4px;
}

.node p {
  margin: 6px;
}

</style>
