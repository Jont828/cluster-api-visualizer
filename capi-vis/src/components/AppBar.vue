<template>
  <v-app-bar
    color="blue darken-2"
    app
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
    <v-app-bar-nav-icon
      class="ma-2"
      v-else
    ></v-app-bar-nav-icon>

    <v-toolbar-title class="text-no-wrap pa-0">{{ title }}</v-toolbar-title>

    <v-spacer></v-spacer>
    <v-btn
      icon
      text
      class="ma-2"
      @click="linkHandler"
    >
      <v-icon color="white">
        {{ isStraight ? 'mdi-sine-wave' : 'mdi-square-wave' }}
      </v-icon>
    </v-btn>

  </v-app-bar>
</template>

<script>
export default {
  name: "AppBar",
  props: {
    title: String,
    showBack: Boolean,
    isStraight: Boolean,
  },
  methods: {
    linkHandler(value) {
      console.log(this.isStraight);
      this.$emit("togglePathStyle", !this.isStraight);
    },
  },
};
</script>

<style lang="less" scoped>
.router-link {
  text-decoration: none;
}
</style>