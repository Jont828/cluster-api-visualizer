<template>
  <v-card
    class="py-1"
    min-width="500px"
    elevation="6"
    :light="!$vuetify.theme.dark"
  >

    <v-card-title class="ml-2">
      Settings
      <v-spacer></v-spacer>
      <v-btn
        icon
        @click="() => { this.$emit('close', {}); }"
      >
        <v-icon>mdi-window-close</v-icon>
      </v-btn>

    </v-card-title>

    <v-card-text class="py-0 ml-2">
      <div class="text-subtitle-2">About</div>
    </v-card-text>
    <v-list rounded>
      <v-list-item>
        <v-list-item-icon>
          <v-icon>mdi-information</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Version</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          v1.1.1
        </v-list-item-action>
      </v-list-item>
      <v-list-item
        href="https://github.com/Jont828/cluster-api-visualizer"
        target="_blank"
      >
        <v-list-item-icon>
          <v-icon>mdi-github</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Source code</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-icon>mdi-open-in-new</v-icon>
        </v-list-item-action>
      </v-list-item>
    </v-list>

    <v-card-text class="py-0 ml-2">
      <div class="text-subtitle-2">Appearance</div>
    </v-card-text>
    <v-list rounded>
      <v-list-item>
        <v-list-item-icon>
          <v-icon>mdi-brightness-6</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Dark theme</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-switch
            v-model="store.darkTheme"
            @change="toggleDarkTheme"
          >
          </v-switch>
        </v-list-item-action>
      </v-list-item>

      <v-list-item>
        <v-list-item-icon>
          <v-icon>mdi-sitemap-outline</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Straighten links</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-switch v-model="store.straightLinks">
          </v-switch>
        </v-list-item-action>
      </v-list-item>
    </v-list>

    <v-card-text class="py-0 ml-2">
      <div class="text-subtitle-2">General</div>
    </v-card-text>
    <v-list rounded>
      <v-list-item>
        <v-list-item-icon>
          <v-icon>mdi-file-download</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Download Format</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-select
            class="selectBox"
            v-model="store.selectedFileType"
            :items="fileTypes"
            dense
            hide-details
          >
          </v-select>
        </v-list-item-action>
      </v-list-item>

      <v-list-item>
        <v-list-item-icon>
          <v-icon>mdi-timer-sync</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Polling period</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-select
            class="selectBox"
            :items="pollingInterval"
            v-model="store.selectedInterval"
            dense
            hide-details
          >
          </v-select>
        </v-list-item-action>
      </v-list-item>
    </v-list>
  </v-card>

</template>

<script>
import { useSettingsStore } from "../stores/settings.js";

export default {
  name: "SettingsCard",
  components: {},
  methods: {
    toggleDarkTheme(val) {
      this.$vuetify.theme.dark = val;
    },
  },
  setup() {
    const store = useSettingsStore();

    return { store };
  },
  data() {
    return {
      fileTypes: ["YAML", "JSON"],
      pollingInterval: ["1s", "5s", "10s", "30s", "1m", "5m", "Off"],
    };
  },
};
</script>

<style lang="less" scoped>
.selectBox {
  width: 100px;
}
</style>