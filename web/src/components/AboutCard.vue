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
        About
        <v-spacer></v-spacer>
        <v-btn
          icon
          @click="() => { this.$emit('close', {}); }"
        >
          <v-icon>mdi-window-close</v-icon>
        </v-btn>

      </v-card-title>
    </v-sheet>

    <!-- <v-card-text class="py-0">
      <div class="text-subtitle-2">Appearance</div>
    </v-card-text> -->
    <v-list>
      <v-list-item>
        <v-list-item-icon>
          <v-icon>mdi-information</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Version</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          v1.0.1
        </v-list-item-action>
      </v-list-item>

      <v-list-item>
        <v-list-item-icon>
          <v-icon>mdi-account-circle</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>Author</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          Jonathan Tong
        </v-list-item-action>
      </v-list-item>
      <v-list-item>
        <v-list-item-icon>
          <v-icon>mdi-github</v-icon>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title>GitHub Repo</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-icon>mdi-open-in-new</v-icon>
        </v-list-item-action>
      </v-list-item>
    </v-list>

    <!-- <v-card-text class="py-0">
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
            class="downloadFormat"
            :items="fileFormats"
            :value="fileFormats[0]"
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
            class="downloadFormat"
            :items="pollingPeriod"
            :value="pollingPeriod[0]"
            dense
            hide-details
          >
          </v-select>
        </v-list-item-action>
      </v-list-item>
    </v-list> -->
  </v-card>

</template>

<script>
import colors from "vuetify/lib/util/colors";

export default {
  name: "AboutCard",
  components: {},
  props: {},
  methods: {},
  data() {
    return {
      fileFormats: ["YAML", "JSON"],
      linkStyles: ["Curved", "Straight"],
      pollingPeriod: ["1s", "5s", "10s", "30s", "1m", "5m", "Off"],
    };
  },
};
</script>

<style lang="less" scoped>
.downloadFormat {
  width: 100px;
}
</style>