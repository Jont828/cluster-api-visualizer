<template>
  <div id="overview">
    <AppBar
      :title="'Logs: ' + kind + '/' + name"
      :showBack="true"
      @reload="fetchLogs(forceReload=true);"
      :scale="String(fontSize)"
      scaleIcon="format-size"
      @zoomIn="() => { fontSize += 1 }"
      @zoomOut="() => { 
        if (fontSize > 1) {
          fontSize -= 1;
        }
      }"
      @showSettings="() => { showSettingsOverlay = true }"
    />
    <!-- Wrap with v-sheet so we can force the dark prop to be true -->
    <v-sheet
      id="chartLoadWrapper"
      class="logs-card"
      dark
      v-show="logsReady"
    >
        <v-container dark class="log-container ma-0 px-0 py-2" fluid>
          <DynamicScroller
            :items="logData"
            :min-item-size="2"
            class="scroller"
          >
            <template v-slot="{ item, index, active }">
              <DynamicScrollerItem
                :item="item"
                :active="active"
                :size-dependencies="[
                  item.message,
                ]"
                :data-index="index"
              >
                <v-row 
                  rounded
                  no-gutters
                >
                  <v-col cols="12" rounded>
                    <highlightjs 
                      language="json" 
                      :code="item.message" 
                      class="wrap pa-0 ml-2 mr-2 my-0" 
                      :style="dynamicFontSize" 
                    />
                  </v-col>
                </v-row>
              </DynamicScrollerItem>
            </template>
          </DynamicScroller>
        </v-container>
    </v-sheet>
    <div
      id="resourceTree"
      class="spinner"
      v-if="!logsReady"
    >
      <v-progress-circular
        :size="50"
        :width="5"
        indeterminate
        color="primary"
      ></v-progress-circular>
    </div>
    
    <v-overlay
      absolute
      :value="showSettingsOverlay"
      z-index="99999"
    >
      <SettingsCard
        @close="() => { showSettingsOverlay = !showSettingsOverlay }"
        class="settingsCard"
        :version="gitVersion"
      />
    </v-overlay>
    <ScrollButton />
  </div>
</template>

<script>
import Vue from "vue";

import SettingsCard from "../components/SettingsCard.vue";
import AppBar from "../components/AppBar.vue";
import ScrollButton from "../components/ScrollButton.vue";

import { useSettingsStore } from "../stores/settings.js";
import { setVersion } from "../mixins/setVersion.js";

export default {
  name: "ResourceLogs",
  components: {
    SettingsCard,
    AppBar,
    ScrollButton,
  },
  mixins: [setVersion],
  data() {
    return {
      showSettingsOverlay: false,
      kind: "",
      name: "",
      namespace: "",
      logData: [],
      cachedLogString: [],
      logsReady: false,
      fontSize: 16,
    };
  },
  setup() {
    const store = useSettingsStore();
    return { store };
  },
  async beforeMount() {
    await this.fetchLogs();
  },
  computed: {
    theme() {
      return this.$vuetify.theme.dark ? "dark" : "light";
    },
    dynamicFontSize() {
      return {
        fontSize: this.fontSize + "px !important",
      };
    }
  },
  mounted() {
    this.kind = this.$route.query.kind;
    this.namespace = this.$route.query.namespace;
    this.name = this.$route.query.name;
    document.title = "Logs: " + this.kind + "/" + this.name;
    this.intervalHandler(this.store.selectedInterval);
  },
  beforeDestroy() {
    this.selected = {};
    clearInterval(this.polling);
  },
  watch: {
    "store.selectedInterval": function (val) {
      console.log("Overview store.selectedInterval: " + val);
      this.intervalHandler(val);
    },
  },
  methods: {
    intervalHandler(val) {
      console.log("Setting polling interval to " + val);
      clearInterval(this.polling);
      if (val === "Off") return;

      let totalSeconds = 0;

      let seconds = val.match(/(\d+)\s*s/);
      let minutes = val.match(/(\d+)\s*m/);

      if (seconds) {
        totalSeconds += parseInt(seconds[1]);
      }
      if (minutes) {
        totalSeconds += parseInt(minutes[1]) * 60;
      }

      console.log("Setting interval to " + totalSeconds + " seconds");
      this.polling = setInterval(
        function () {
          this.fetchLogs();
        }.bind(this),
        totalSeconds * 1000
      );
    },
    linkHandler(val) {
      this.isStraight = val;
    },
    async fetchLogs(forceReload = false) {
      try {
        const params = new URLSearchParams();

        let kind = this.$route.query.kind;
        let namespace = this.$route.query.namespace;
        let name = this.$route.query.name;
        params.append("kind", kind);
        params.append("name", name);
        params.append("namespace", namespace);

        console.log("Fetching log data");
        const response = await Vue.axios.get("/resource-logs", {
          params: params,
        });
        
        if (response.data == null) {
          this.errorMessage = "No log data found";
          return;
        }

        if (
          forceReload ||
          this.cachedLogString !== JSON.stringify(response.data)
        ) {
          this.cachedLogString = JSON.stringify(response.data);
          let data = [];
          for (let i = 0; i < response.data.length; i++) {
            data.push({
              id: i,
              message: response.data[i],
            });
          }
          // If the log array gets too large, the application will slow down.
          if (data.length > this.store.maxLogLines) {
            data = data.slice(0, this.store.maxLogLines);
          }
          this.logData = data;
          this.logsReady = true;
        }

      } catch (error) {
        console.log("Error:", error.toJSON());
        this.alert = true;
        if (error.response) {
          if (error.response.status == 404) {
            this.errorMessage =
              "Logs not found";
          } else {
            this.errorMessage =
              "Failed to get logs";
          }
        } else if (error.request) {
          this.errorMessage = "No server response received";
        } else {
          this.errorMessage = "Unable to create request";
        }
      }
    },
  },
};
</script>

<style lang="less">
code {
  background-color: transparent !important;
  word-break: normal; // TODO: look into breaking on commas.
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>

<style lang="less" scoped>
#overview {
  height: 100%;
}

#chartLoadWrapper {
  min-height: 100% !important;
}

.scroller {
  height: 100%;
}
</style>