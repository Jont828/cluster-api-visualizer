<template>
  <v-app-bar
    id="appBar"
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
    <v-tooltip bottom>
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          icon
          text
          class="ma-2"
          @click="$emit('reload', true)"
          v-bind="attrs"
          v-on="on"
        >
          <v-icon color="white">
            {{"mdi-refresh"}}
          </v-icon>
        </v-btn>
      </template>
      <span>Reload resources</span>
    </v-tooltip>

    <v-spacer></v-spacer>
    <v-tooltip bottom>
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          icon
          text
          class="ma-2"
          @click="$emit('zoomOut', true)"
          v-bind="attrs"
          v-on="on"
        >
          <v-icon color="white">
            {{"mdi-minus"}}
          </v-icon>
        </v-btn>
      </template>
      <span>Zoom out</span>
    </v-tooltip>
    <span>{{ Math.round(scale*100) }}%</span>
    <v-tooltip bottom>
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          icon
          text
          class="ma-2"
          @click="$emit('zoomIn', true)"
          v-bind="attrs"
          v-on="on"
        >
          <v-icon color="white">
            {{"mdi-plus"}}
          </v-icon>
        </v-btn>
      </template>
      <span>Zoom in</span>
    </v-tooltip>
    <v-tooltip bottom>
      <template v-slot:activator="{ on, attrs }">
        <v-btn
          icon
          text
          class="ma-2"
          @click="$emit('togglePathStyle', !isStraight);"
          v-bind="attrs"
          v-on="on"
        >
          <v-icon color="white">
            {{ isStraight ? 'mdi-sine-wave' : 'mdi-square-wave' }}
          </v-icon>
        </v-btn>
      </template>
      <span>Toggle link style</span>
    </v-tooltip>

  </v-app-bar>
</template>

<script>
export default {
  name: "AppBar",
  props: {
    title: String,
    showBack: Boolean,
    isStraight: Boolean,
    scale: Number,
  },
};
</script>

<style lang="less" scoped>
#appBar {
  z-index: 2000;
}

.router-link {
  text-decoration: none;
}
</style>