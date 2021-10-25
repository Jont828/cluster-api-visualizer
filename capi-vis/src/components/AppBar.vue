<template>
  <v-app-bar
    color="blue darken-2"
    app
    dense
    dark
  >
    <router-link
      :to="'/'"
      class="router-link"
      v-if="showBack"
    >
      <v-btn
        icon
        text
        class="ma-2"
      >
        <v-icon color="white">
          mdi-chevron-left
        </v-icon>
      </v-btn>
    </router-link>
    <v-app-bar-nav-icon v-else></v-app-bar-nav-icon>

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
    showBack: Boolean,
  },
};
</script>

<style lang="less" scoped>
.router-link {
  text-decoration: none;
}
</style>