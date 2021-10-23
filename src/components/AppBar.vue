<template>
  <v-app-bar
    color="blue darken-2"
    app
    dense
    dark
  >
    <v-btn
      icon
      text
      class="ma-2"
    >
      <router-link
        :to="'/'"
        class="router-link"
      >
        <v-icon color="white">
          mdi-chevron-left
        </v-icon>
      </router-link>

    </v-btn>
    <!-- <v-spacer></v-spacer> -->

    <v-toolbar-title class="text-no-wrap">{{ title }}</v-toolbar-title>

    <!-- <v-spacer></v-spacer> -->

  </v-app-bar>
</template>

<script>
export default {
  name: "AppBar",
  props: {
    title: String,
  },
};
</script>

<style lang="less" scoped>
.router-link {
  text-decoration: none;
}
</style>