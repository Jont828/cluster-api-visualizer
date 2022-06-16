<template>
  <v-card class="resourceCard mx-auto">
    <v-sheet
      :color="color"
      class="resourceSheet pa-4"
      dark
    >
      <v-card-title class="text-h5">
        {{ name }}
        <v-spacer></v-spacer>
        <v-btn
          icon
          color="white"
          @click="this.downloadYaml"
        >
          <v-icon>mdi-download</v-icon>
        </v-btn>
        <v-btn
          icon
          color="white"
          @click="() => { this.$emit('unselectNode', {}); }"
        >
          <v-icon>mdi-window-close</v-icon>
        </v-btn>

      </v-card-title>
      <v-card-subtitle>
        <v-chip-group
          column
          @change="selectConditions"
          multiple
        >
          <v-chip
            active
            v-for="condition in conditions"
            :key="condition.type"
            color="white"
            active-class="primary--text"
            :text-color="(condition.status) ? 'success' : 'error'"
          >
            <v-icon
              left
              v-if="condition.status"
            >
              mdi-check-circle
            </v-icon>
            <v-icon
              left
              v-else
            >
              mdi-alert-circle
            </v-icon>
            {{ condition.type }}
          </v-chip>
        </v-chip-group>
      </v-card-subtitle>
      <v-text-field
        v-model="search"
        label="Search Custom Resource Fields"
        dark
        flat
        solo-inverted
        hide-details
        clearable
        clear-icon="mdi-close-circle-outline"
        :color="color"
      ></v-text-field>
      <v-checkbox
        v-model="caseSensitive"
        dark
        hide-details
        label="Case sensitive search"
      ></v-checkbox>
    </v-sheet>
    <v-card-text>
      <v-treeview
        hoverable
        :items="items"
        :search="search"
        :filter="filter"
        :open.sync="open"
        :active.sync="active"
        rounded
        class="text-wrap"
      >
      </v-treeview>
    </v-card-text>
  </v-card>
</template>

<script>
import yaml from "js-yaml";

export default {
  name: "CustomResourceTree",
  props: {
    items: Array,
    jsonItems: Object,
    name: String,
    color: String,
  },
  data() {
    return {
      open: [],
      active: [], // for auto-highlighting statuses
      search: null,
      caseSensitive: false,
      conditions: [],
    };
  },
  mounted() {
    // Open all top level elements
    console.log("Items are", this.items);
    // this.items.forEach((e, i) => {
    //   if (e.children.length > 0) {
    //     this.open.push(i);
    //   }
    // });
    // this.open = [".status", ".status.conditions", ".status.conditions[0]"];
    this.setConditions(this.jsonItems?.status?.conditions);
  },
  methods: {
    downloadYaml() {
      const yamlCRD = yaml.dump(this.jsonItems);
      const link = document.createElement("a");
      link.href = `data:text/plain;charset=utf-8,${yamlCRD}`;
      link.download = this.name + ".yaml";
      link.click();
    },
    setConditions(conditions) {
      this.conditions = [];
      if (conditions !== undefined) {
        conditions.forEach((e, i) => {
          this.conditions.push({
            type: e.type,
            status: e.status === "True",
          });
        });
        // console.log("Conditions are", this.conditions);
      }
    },
    selectConditions(indexArr) {
      console.log(indexArr);
      if (indexArr.length > 0) {
        this.open = [".status", ".status.conditions"];
      } else {
        this.open = [];
      }
      this.active = []; // TODO: it looks like this array only highlights the last index
      indexArr.forEach((index) => {
        this.open.push(".status.conditions[" + index + "]");
        this.active.push(".status.conditions[" + index + "].type");
      });
      // console.log("Open is", this.open);
    },
  },
  watch: {
    jsonItems: {
      handler(val, old) {
        this.setConditions(val?.status?.conditions);
      },
    },
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
// .resourceSheet {
//   padding: 0 16px 16px 16px;
// }
</style>