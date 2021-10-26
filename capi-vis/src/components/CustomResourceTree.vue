<template>
  <v-card class="mx-auto">
    <v-sheet
      :color="color"
      class="pa-4"
      dark
    >
      <v-card-title class="text-h5">{{ title }}</v-card-title>
      <v-text-field
        v-model="search"
        label="Search Company Directory"
        dark
        flat
        solo-inverted
        hide-details
        clearable
        clear-icon="mdi-close-circle-outline"
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
export default {
  name: "CustomResourceTree",
  props: {
    items: Array,
    title: String,
    color: String,
  },
  data() {
    return {
      open: [1, 2],
      search: null,
      caseSensitive: false,
    };
  },
  computed: {
    filter() {
      return this.caseSensitive
        ? (item, search, textKey) => {
            console.log(item, search, textKey);
            return item["name"].indexOf(search) > -1;
          }
        : undefined;
    },
  },
};
</script>