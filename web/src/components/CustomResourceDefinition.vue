<template>
  <v-card class="resourceCard mx-auto">
    <v-sheet
      :color="color"
      class="resourceSheet pa-4"
      dark
    >
      <v-card-title class="text-h5">
        {{ name }}
        <v-spacer></v-spacer>
        <v-btn
          icon
          color="white"
          @click="this.downloadYaml"
        >
          <v-icon>mdi-download</v-icon>
        </v-btn>
        <v-btn
          icon
          color="white"
          @click="() => { this.$emit('unselectNode', {}); }"
        >
          <v-icon>mdi-window-close</v-icon>
        </v-btn>

      </v-card-title>
      <v-card-subtitle>
        <div class="conditionChipListWrapper my-2">
          <div
            v-for="(condition, index) in conditions"
            :key="index"
          >
            <v-tooltip
              :top="scrollY <= 20"
              :bottom="scrollY > 20"
              :disabled="condition.status"
            >
              <template v-slot:activator="{ on, attrs }">
                <v-chip
                  active
                  link
                  :class="{
                    'conditionChip': true,
                  }"
                  color="white"
                  :text-color="(condition.status) ? 'success' : ((condition.isError) ? 'error' : 'warning')"
                  @click="selectCondition(index)"
                  v-bind="attrs"
                  v-on="on"
                >
                  <StatusIcon
                    :type="(condition.status) ? 'success' : condition.severity.toLowerCase()"
                    :spinnerWidth="2"
                    left
                  >
                  </StatusIcon>
                  {{ condition.type }}
                </v-chip>
              </template>
              <span>{{ condition.severity }}: {{ condition.reason }}</span>
            </v-tooltip>
          </div>
        </div>

        <div class="mt-4">
          <v-text-field
            v-model="search"
            label="Search Custom Resource Fields"
            dark
            flat
            solo-inverted
            hide-details
            clearable
            clear-icon="mdi-close-circle-outline"
            :color="color"
          ></v-text-field>
          <v-checkbox
            v-model="caseSensitive"
            dark
            hide-details
            label="Case sensitive search"
          ></v-checkbox>
        </div>

      </v-card-subtitle>
    </v-sheet>
    <v-card-text>
      <v-treeview
        hoverable
        :items="items"
        :search="search"
        :filter="filter"
        :open.sync="open"
        :active.sync="active"
        activatable
        rounded
        class="text-wrap"
      >
        <template v-slot:label="{ item }">
          <span
            :ref="item.id"
            class="text-wrap"
          >{{ item.name }}</span>
        </template>
      </v-treeview>
    </v-card-text>
  </v-card>
</template>

<script>
import yaml from "js-yaml";
import StatusIcon from "./StatusIcon.vue";

export default {
  name: "CustomResourceDefinition",
  components: {
    StatusIcon,
  },
  props: {
    items: Array,
    jsonItems: Object,
    name: String,
    color: String,
  },
  data() {
    return {
      open: [],
      active: [], // for auto-highlighting statuses
      search: null,
      caseSensitive: false,
      conditions: [],
      scrollY: 0,
    };
  },
  mounted() {
    this.setConditions(this.jsonItems?.status?.conditions);
    console.log("Scroll is", window.scrollY);
    window.addEventListener("scroll", this.onScroll);
  },
  methods: {
    onScroll(e) {
      console.log("Scroll is", window.scrollY);
      this.scrollY = window.scrollY;
      // this.windowTop = window.top.scrollY /* or: e.target.documentElement.scrollTop */
    },
    downloadYaml() {
      const yamlCRD = yaml.dump(this.jsonItems);
      const link = document.createElement("a");
      link.href = `data:text/plain;charset=utf-8,${yamlCRD}`;
      link.download = this.name + ".yaml";
      link.click();
    },
    setConditions(conditions) {
      this.conditions = [];
      if (conditions !== undefined) {
        conditions.forEach((e, i) => {
          this.conditions.push({
            type: e.type,
            status: e.status === "True",
            isError: e.severity === "Error",
            severity: e.severity,
            reason: e.reason,
          });
        });
      }
      console.log(this.conditions);
    },
    selectCondition(index) {
      this.open.push(".status");
      this.open.push(".status.conditions");
      this.open.push(".status.conditions[" + index + "]");
      this.active.push(".status.conditions[" + index + "].type");
      console.log(this.open);

      let refName = ".status.conditions[" + index + "].type";
      this.$nextTick(() =>
        this.$vuetify.goTo(this.$refs[refName], {
          easing: "easeInOutQuint",
          duration: 1000,
        })
      );
    },
  },
  watch: {
    jsonItems: {
      handler(val, old) {
        this.setConditions(val?.status?.conditions);
      },
    },
  },
  computed: {
    filter() {
      return this.caseSensitive
        ? (item, search, textKey) => {
            // console.log(item, search, textKey);
            return item["name"].indexOf(search) > -1;
          }
        : (item, search, textKey) => {
            // console.log(item, search, textKey);
            return (
              item["name"].toLowerCase().indexOf(search.toLowerCase()) > -1
            );
          };
    },
  },
};
</script>

<style lang="less" scoped>
.conditionChipListWrapper {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
// .conditionChip {
//   .loading {
//     height: 20px !important;
//     width: 20px !important;
//     min-height: 20px !important;
//     min-width: 20px !important;
//   }
// }
</style>

<style lang="less">
.v-treeview-node__label {
  padding: 10px 0;
}
</style>