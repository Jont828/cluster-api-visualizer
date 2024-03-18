<template>
  <div class="targetTreeWrapper">
    <vue-tree
      id="resourceTree"
      ref="tree"
      :dataset="treeData"
      :config="treeConfig"
      :collapse-enabled="true"
      :linkStyle="(store.straightLinks) ? 'straight' : 'curve'"
      @scale="(val) => $emit('scale', val)"
    >
      <template v-slot:node="{ node, collapsed }">
        <v-hover>
          <template v-slot:default="{ hover }">
            <div class="card-wrap shadow">

              <!-- Wrapper card for Dark theme transparent background -->
              <v-card elevation="0">
                <v-card
                  class="node mx-auto transition-swing"
                  dark
                  :elevation="hover ? 6 : 3"
                  v-on:click="selectNode(node)"
                  :style="($vuetify.theme.dark) ? {
                    //background: $vuetify.theme.themes[theme].legend[node.provider].background + (hover ? '55' : '44'),
                    //color: hover ? '#fff' : $vuetify.theme.themes[theme].legend[node.provider].text,
                    'background-color': hover ? '#383838' : '#272727',
                    'color': $vuetify.theme.themes[theme].legend[node.provider],
                    // 'border-color': $vuetify.theme.themes[theme].legend[node.provider],
                  } : {
                    'background-color': $vuetify.theme.themes[theme].legend[node.provider],
                  }"
                >

                  <p class="kind font-weight-medium text-truncate">{{ (node.collapsible) ? node.displayName : node.kind }}</p>

                  <p
                    class="name font-italic text-truncate"
                    v-if="!node.collapsible"
                  >{{ node.displayName }}</p>
                  <v-icon
                    class="chevron"
                    size="18"
                    :style="($vuetify.theme.dark) ? {
                      'color': $vuetify.theme.themes[theme].legend[node.provider],
                    } : null"
                    v-else-if="collapsed"
                  >mdi-chevron-down</v-icon>
                  <v-icon
                    class="chevron"
                    size="18"
                    :style="($vuetify.theme.dark) ? {
                      'color': $vuetify.theme.themes[theme].legend[node.provider],
                    } : null"
                    v-else
                  >mdi-chevron-up</v-icon>

                </v-card>
              </v-card>
              <StatusBadge
                v-if="node.hasReady"
                :type="(node.ready) ? 'success' : node.severity.toLowerCase()"
                :size="18"
              ></StatusBadge>

            </div>
          </template>
        </v-hover>
      </template>
    </vue-tree>
    <div class="legend">
      <v-card class="legend-card">
        <div
          class="legend-entry"
          v-for="(displayName, provider) in legend"
          :key="provider"
        >
          <div class="legend-entry-content">
            <v-icon
              class="legend-entry-icon"
              :color="$vuetify.theme.themes[theme].legend[provider]"
            >
              mdi-square-rounded
            </v-icon>
            <div class="legend-entry-text">{{ displayName }}</div>
          </div>

        </div>

      </v-card>
    </div>
  </div>
</template>

<script>
import VueTree from "../components/VueTree.vue";
import StatusBadge from "./StatusBadge.vue";

import { useSettingsStore } from "../stores/settings.js";

export default {
  name: "DescribeClusterTree",
  components: {
    VueTree,
    StatusBadge,
  },
  data() {
    return {
      legendGradient: "",
      index: 0,
    };
  },
  setup() {
    const store = useSettingsStore();
    return { store };
  },
  computed: {
    theme() {
      return this.$vuetify.theme.dark ? "dark" : "light";
    },
  },
  props: {
    treeConfig: Object,
    treeData: Object,
    selectedNode: Object,
    legend: Object,
  },
  methods: {
    selectNode(node) {
      if (!node.collapsible) {
        this.$emit("selectNode", node);
      } else {
        // TODO: store info about which nodes are open or closed and on a reload, preserve the collapse state of nodes that still exist
        // TODO: use a UID map and refactor GroupByKindNodes to have the same UID, i.e. make it dependent on the parent instead
        // this.$emit("toggleNodeCollapse", node);
      }
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
    // border: 1px solid black;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}

@-webkit-keyframes SlideRight {
  0% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}
@-moz-keyframes SlideRight {
  0% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}
@keyframes SlideRight {
  0% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.animated {
  background-size: 500% 100% !important;
  -webkit-animation: SlideRight 7.5s linear infinite !important;
  -moz-animation: SlideRight 7.5s linear infinite !important;
  animation: SlideRight 7.5s linear infinite !important;
}

.node-slot {
  cursor: default !important;
}

.card-wrap {
  // background-color: #f8f3f2;
  position: relative;
}

// .shadow {
//   box-shadow: 2px 2px 10px 0px rgba(255, 0, 0, 0.75), 2px -2px 10px 0px rgba(255, 0, 0, 0.75);
// }

.badge {
  padding: 0;
}

.v-badge {
  padding: 0 !important;
}

.node-hidden {
  width: 170px;
  height: 50px;
  position: absolute;
  background: red;
  z-index: -100;
  top: 0;
  left: 0;
  text-align: right;

  .node-hidden-inner {
    display: inline-block;
    width: 50px;
    height: 50px;
  }
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

  // display: flex;
  // flex-direction: row;
  // align-items: center;

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
    max-width: 145px;
  }

  .name {
    padding: 0 4px;
    // Note: this stops the end of the italicized text from getting cut off. The name will overflow but it doesn't affect max-width.
    max-width: 160px;
  }

  .name,
  .kind {
    // text-shadow: rgba(0, 0, 0, 1) 3px 0 5px;
    text-align: center;
    // white-space: nowrap;
    display: inline-block;
    // overflow: hidden;
    // text-overflow: ellipsis;
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
    padding: 10px 5px;
    display: inline-block;

    .legend-entry {
      display: inline-block;
      margin-right: 10px;
      line-height: 24px;
      height: 24px;

      .legend-entry-content {
        line-height: 24px;
        height: 24px !important;
        // display: flex;
        // position: relative;
        // align-items: center;
        // justify-content: center;

        .legend-entry-icon {
          display: inline-block;
          border-radius: 3px;
          margin: 0 6px;
          height: 24px;
          line-height: 24px;
          font-size: 24px;
          vertical-align: top;
        }
        .legend-entry-text {
          line-height: 24px;
          vertical-align: top;
          display: inline-block;
        }

        .overlapping-icon-wrapper {
          // vertical-align: middle;
          position: relative;
          height: 18px;
          width: 18px;
          margin: 0 8px;
          display: inline-block;

          .overlapping-icon {
            position: absolute;
            top: 0;
            left: 0;
            height: 18px;
            margin: 0 !important;
            transition: opacity 1s linear;
          }
          .overlapping-icon + .overlapping-icon {
            opacity: 0;
          }
        }
      }
    }
  }
}
</style>