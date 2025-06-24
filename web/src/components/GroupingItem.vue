<template>
  <v-card class="resource-card mx-auto">
    <link
      v-if="$vuetify.theme.dark"
      rel="stylesheet" 
      href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/styles/stackoverflow-dark.min.css"
    >
    <link
      v-else
      rel="stylesheet" 
      href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/10.7.2/styles/stackoverflow-light.min.css"
    >

    <v-sheet
      :color="($vuetify.theme.dark) ? '#272727' : color"
      class="resourceSheet pa-4"
    >
      <v-card-title
        class="text-h5"
        :style="{
          color: ($vuetify.theme.dark) ? color : 'white'
        }"
      >
        {{ name }}
        <StatusIcon
            :type="ready ? 'success' : severity.toLowerCase()"
            :spinnerWidth="2"
        >
        </StatusIcon>
        <v-spacer></v-spacer>
        <v-btn
          icon
          color="white"
          @click="() => { this.$emit('unselectNode', {}); }"
        >
          <v-icon>mdi-window-close</v-icon>
        </v-btn>

      </v-card-title>
    </v-sheet>
    <v-card-text>
      <v-list rounded>
        <v-list-item-group
          color="primary"
        >
          <v-list-item
            v-for="(item, i) in items"
            :key="i"
          >
            <!-- <v-list-item-icon> -->

            <!-- </v-list-item-icon> -->
            <v-list-item-content
            @click="$emit('selectGroupedItem', item)"
            >
              <v-list-item-title>
                {{ item.name }}<v-icon>mdi-chevron-right</v-icon>
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list-item-group>
      </v-list>
      <!-- <v-treeview
        hoverable
        :items="items"
        :search="search"
        :filter="filter"
        :open.sync="open"
        :open-all="true"
        :active.sync="active"
        activatable
        rounded
        class="text-wrap"
      >
        <template v-slot:label="{ item }">
          <highlightjs 
            language="yaml" 
            :ref="item.id"
            :code="item.name" 
            class="text-wrap yaml-code" 
          />
        </template>
      </v-treeview> -->
    </v-card-text>
  </v-card>
</template>

<script>
import yaml from "js-yaml";
import StatusIcon from "./StatusIcon.vue";
import { useSettingsStore } from "../stores/settings.js";
import colors from "vuetify/lib/util/colors";

export default {
  name: "GroupingItem",
  components: {
    StatusIcon,
  },
  props: {
    items: Array,
    name: String,
    color: String,
    ready: Boolean,
    severity: String,
  },
  setup() {
    const store = useSettingsStore();

    return { store };
  },
  data() {
    return {
      open: [],
      active: [], // for auto-highlighting statuses
      search: null,
      caseSensitive: false,
      conditions: [],
      url: "",
      scrollY: 0,
    };
  },
  mounted() {
    console.log("Got items", this.items);
    // this.setConditions(this.jsonItems?.status?.conditions);
    // window.addEventListener("scroll", this.onScroll);
    // console.log("JSON items are", this.jsonItems);
    // let kind = this.jsonItems.kind;
    // let name = this.jsonItems.metadata.name;
    // let namespace = this.jsonItems.metadata.namespace;
    // this.url = "/logs?kind=" + kind + "&name=" + name + "&namespace=" + namespace;
    // console.log("URL is", this.url);
  },
  methods: {
    getType(condition) {
      if (condition.status === "True") return "success";
      else if (condition.isError || !condition.severity || condition.status === "Unknown") return "error"; // if severity is undefined, we assume it's an error
      else return "warning";
    },
    onScroll(e) {
      this.scrollY = window.scrollY;
      // this.windowTop = window.top.scrollY /* or: e.target.documentElement.scrollTop */
    },
  },
  watch: {
    jsonItems: {
      handler(val, old) {
        console.log("Val is", val);
        this.setConditions(val?.status?.conditions);
      },
    },
    items: {
      handler(val, old) {
        let recurse = function(items, open = []) {
          items.forEach((item) => {
            if (item.children) {
              open = open.concat(recurse(item.children));
            }
            open.push(item.id);
          });
          return open;
        };

        this.open = recurse(val);
      },
    }
  },
  computed: {
    filter() {
      return this.caseSensitive
        ? (item, search, textKey) => {
            // console.log(item, search, textKey);
            return item["name"].indexOf(search) > -1;
          }
        : (item, search, textKey) => {
            // console.log(item, search, textKey);
            return (
              item["name"].toLowerCase().indexOf(search.toLowerCase()) > -1
            );
          };
    },
  },
};
</script>

<style lang="less" scoped>
.conditionChipListWrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
</style>

<style lang="less">
.v-treeview-node__label {
  padding: 10px 0;
}

.resource-card .yaml-code code {
  font-size: 100%;
}
</style>