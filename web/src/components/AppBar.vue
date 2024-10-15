<template>
  <v-app-bar
    id="appBar"
    app
    :color="($vuetify.theme.dark ? null: 'primary')"
    style="background: linear-gradient(104.44deg, #000000 25%, #32343B 67.88%);"
    dark
  >
    <v-img src="../assets/mirantis-logo-inverted-horizontal-one-color.png" max-width="160"/>
    <v-btn
      icon
      text
      class="ma-2"
      @click="() => { router.back() }"
      v-if="showBack"
    >
      <v-icon color="white">
        mdi-chevron-left
      </v-icon>
    </v-btn>
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
    <div>
        <v-switch label="Lens Integration" style="height: 25px;" v-model="lens" @change="lensChanged"></v-switch>
    </div>
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
    <v-icon
      v-if="scaleIcon"
      color="white"
    > mdi-{{ scaleIcon }}</v-icon>
    <span v-else>{{ scale }}</span>
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
          @click="$emit('showSettings', true)"
          v-bind="attrs"
          v-on="on"
        >
          <!-- TODO: should it be showSettings or toggleSettings, i.e. should clicking again close the overlay -->
          <v-icon color="white">
            mdi-cog
          </v-icon>
        </v-btn>
      </template>
      <span>Show settings</span>
    </v-tooltip>

  </v-app-bar>
</template>

<script>
import router from '../router';

export default {
  name: "AppBar",
  props: {
    title: String,
    showBack: Boolean,
    isStraight: Boolean,
    scale: String,
    scaleIcon: String,
    backURL: String,
  },
  methods: {
    lensChanged(e) {
      console.log("lens=", e)
      this.$emit('updateLens', e)
    },
  },
  data() {
    return {
      router: router,
      lens: true,
    };
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