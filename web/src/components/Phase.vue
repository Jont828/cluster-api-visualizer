<template>
  <v-card-subtitle class="cardSubtitle">
    <!-- TODO: add loading spinner for "-ing" phases, x for failed, and ? for unknown -->
    <!-- <span>Target Cluster</span> -->
    <div class="wrap">
      <v-icon
        v-if="icon != ''"
        class="phase-icon"
        :color="color"
      > mdi-{{ icon }} </v-icon>
      <v-progress-circular
        v-else
        class="phase-spinner"
        indeterminate
        :size="14"
        :width="1"
        :color="color"
      ></v-progress-circular>
      <span :class="color + '--text'">
        {{ phase }}
      </span>
    </div>
  </v-card-subtitle>
</template>

<script>
export default {
  name: "Phase",
  props: {
    phase: String,
  },
  data() {
    return {
      color: "",
      icon: "",
    };
  },
  methods: {
    getColor(phase) {
      switch (phase) {
        case "Provisioned":
          this.color = "green";
          break;
        case "Pending":
        case "Provisioning":
        case "Deleting":
          this.color = "amber";
          break;
        case "Failed":
        case "Unknown":
          this.color = "red";
          break;
        default:
          this.color = "grey";
          break;
      }
    },
    setIcon(phase) {
      switch (phase) {
        case "Provisioned":
          this.icon = "check";
          break;
        case "Pending":
        case "Provisioning":
        case "Deleting":
          this.icon = "";
          break;
        case "Failed":
          this.icon = "close";
          break;
        case "Unknown":
          this.icon = "help";
          break;
        default:
          this.icon = "";
          break;
      }
    },
  },
  mounted() {
    this.getColor(this.phase);
    this.setIcon(this.phase);
  },
};
</script>

<style lang="less" scoped>
.cardSubtitle {
  padding-bottom: 0;

  .wrap {
    display: flex;
    flex-direction: row;
    align-items: center;

    .phase-icon {
      display: inline-block;
      font-size: 16px;
      margin-right: 2px;
    }

    .phase-spinner {
      display: inline-block;
      margin-right: 4px;
    }

    span {
      display: inline-block;
    }
  }
}
</style>