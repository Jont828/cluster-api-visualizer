<template>
  <v-card
    class="pa-1"
    light
    min-width="500px"
    elevation="6"
  >
    <v-sheet>
      <!-- <v-sheet
      class="pa-0"
      color="primary"
      dark
    > -->

      <v-card-title>
        Settings
        <v-spacer></v-spacer>
        <v-btn
          icon
          @click="() => { this.$emit('close', {}); }"
        >
          <v-icon>mdi-window-close</v-icon>
        </v-btn>

      </v-card-title>
    </v-sheet>

    <v-card-text class="py-0">
      <div class="text-subtitle-2">Appearance</div>
    </v-card-text>
    <v-list>
      <v-list-item>
        <v-list-item-icon>
          <v-icon>mdi-brightness-6</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Dark theme</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-switch v-model="darkTheme">
          </v-switch>
        </v-list-item-action>
      </v-list-item>

      <v-list-item>
        <v-list-item-icon>
          <v-icon>mdi-sine-wave</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Curved link style</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-switch v-model="curvedLinks">
          </v-switch>
        </v-list-item-action>
      </v-list-item>
    </v-list>

    <v-card-text class="py-0">
      <div class="text-subtitle-2">General</div>
    </v-card-text>
    <v-list>
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
            v-model="selectedFileType"
            :items="fileType"
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
            v-model="selectedInterval"
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
import colors from "vuetify/lib/util/colors";

export default {
  name: "SettingsCard",
  components: {},
  props: {
    interval: String,
  },
  methods: {},
  mounted: function () {
    this.selectedInterval = this.interval;
  },
  data() {
    return {
      darkTheme: false,
      curvedLinks: true,
      fileType: ["YAML", "JSON"],
      selectedFileType: "YAML",
      pollingInterval: ["1s", "5s", "10s", "30s", "1m", "5m", "Off"],
      selectedInterval: null,
    };
  },
  watch: {
    darkTheme: function (val) {
      console.log("darkTheme: " + val);
      this.$emit("setDarkTheme", val);
    },
    curvedLinks: function (val) {
      console.log("curvedLinks: " + val);
      this.$emit("setStraightLinks", !val);
    },
    selectedFileType: function (val) {
      console.log("selectedFileType: " + val);
      this.$emit("setFileType", val);
    },
    selectedInterval: function (val) {
      console.log("selectedInterval: " + val);
      this.$emit("setInterval", val);
    },
  },
};
</script>

<style lang="less" scoped>
.selectBox {
  width: 100px;
}
</style>