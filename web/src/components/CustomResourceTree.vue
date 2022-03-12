<template>
  <v-card class="resourceCard mx-auto">
    <v-sheet
      :color="color"
      class="resourceSheet pa-4"
      dark
    >
      <v-card-title class="text-h5">
        Resource: {{ name }}
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
        class="text-wrap"
      >
        <!-- <template v-slot:prepend="{ item }">
          <v-icon
            v-if="item.children"
            v-text="`mdi-${item.id === 1 ? 'home-variant' : 'folder-network'}`"
          ></v-icon>
        </template> -->
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
      open: [1, 2],
      search: null,
      caseSensitive: false,
    };
  },
  methods: {
    downloadYaml() {
      const yamlCRD = yaml.dump(this.jsonItems);
      const link = document.createElement("a");
      link.href = `data:text/plain;charset=utf-8,${yamlCRD}`;
      link.download = this.name + ".yaml";
      link.click();
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
.resourceSheet {
  // padding: 0 16px 16px 16px;
}
</style>