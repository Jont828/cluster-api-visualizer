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
              :class="[ 'node', 'mx-auto', 'transition-swing', { animated: (node.hasReady && !node.ready) } ]"
              :elevation="hover ? 6 : 3"
              :style="{ 
                background: (node.hasReady && !node.ready) ? computeNotReadyGradient(legend[node.provider].color, 20) : legend[node.provider].color,
                border: collapsed ? '' : '',
              }"
              v-on:click="selectNode(node)"
            >
              <!-- <router-link
          :to="'/'"
          class="node-router-link"
        > -->
              <p
                class="kind font-weight-medium"
                v-if="!node.isVirtual"
              >{{ node.kind }}</p>
              <p
                class="kind font-weight-medium"
                v-else
              >{{ node.displayName }}</p>

              <p
                class="name font-italic"
                v-if="!node.isVirtual"
              >{{ node.displayName }}</p>
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
          <div class="legend-entry-content">
            <div
              class="legend-entry-icon"
              :style="{
                'background-color': entry.color
              }"
            />
            <span class="legend-entry-text">{{ entry.name }}</span>
          </div>

        </div>

        <div class="legend-entry">
          <div class="legend-entry-content">
            <div class="overlapping-icon-wrapper">
              <div
                :key="i"
                v-for="(entry, i) in Object.values(this.legend)"
                class="legend-entry-icon animated overlapping-icon"
                :style="{
                    background: computeNotReadyGradient(entry.color, 3),
                    opacity: (i == index) ? 1 : 0,
                  }"
              />
            </div>
            <span class="legend-entry-text">Not Ready</span>
          </div>
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
  data() {
    return {
      legendGradient: "",
      index: 0,
    };
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
      if (!node.isVirtual) {
        this.$emit("selectNode", node);
      }
    },
    computeNotReadyGradient(color, width) {
      return this.computeLinearGradient(
        [color, this.adjustColor(color, 20)],
        width
      );
    },
    computeLinearGradient(colors, width) {
      // console.log("colors", colors);
      let result = "repeating-linear-gradient(135deg";
      colors.forEach((color, i) => {
        result += ", " + color + " " + width * (2 * i) + "px,";
        result += color + "  " + width * (2 * i + 1) + "px";
        // Alternate effect is (i) and (i+1)
      });
      result += ")";

      return result;
    },
    iterateGradientColors() {
      let legendColors = Object.values(this.legend);
      // console.log("Color:", legendColors[this.index].color);
      this.legendGradient = this.computeNotReadyGradient(
        legendColors[this.index++].color,
        4
      );
      this.index = this.index % legendColors.length;
      // console.log("Gradient is", this.legendGradient);
    },
  },
  mounted() {
    setInterval(this.iterateGradientColors, 2000);
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
  -webkit-animation: SlideRight 5s linear infinite !important;
  -moz-animation: SlideRight 5s linear infinite !important;
  animation: SlideRight 5s linear infinite !important;
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
    // text-shadow: rgba(0, 0, 0, 1) 3px 0 5px;
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

      .legend-entry-content {
        display: flex;
        position: relative;
        align-items: center;
        justify-content: center;

        .overlapping-icon-wrapper {
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
        .legend-entry-icon {
          display: inline-block;
          border-radius: 3px;
          margin: 0 8px;
          width: 18px;
          height: 18px;
        }
        .legend-entry-text {
          display: inline-block;
        }
      }
    }
  }
}
</style>