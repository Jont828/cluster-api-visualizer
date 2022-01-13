<template>
  <div class="targetTreeWrapper">
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
</template>

<script>
import VueTree from "../components/VueTree.vue";

export default {
  name: "TargetClusterTree",
  components: {
    VueTree,
  },
  props: {
    treeConfig: Object,
    treeData: Object,
    isStraight: Boolean,
    selectedNode: Object,
    legend: Object,
  },
  methods: {
    selectNode(node) {
      this.$emit("selectNode", node);
    },
  },
};
</script>

<style lang="less" scoped>
.targetTreeWrapper {
  height: 100%;
  position: relative;

  #resourceTree {
    width: 100%;
    // height: 100%;
    height: 100%;
    background-color: #f8f3f2;
    // border: 1px solid black;
  }
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

.legend {
  text-align: center;
  // display: inline-block;
  position: absolute;
  bottom: 30px;
  z-index: 1000;
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
</style>