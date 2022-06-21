<template>
  <v-avatar
    :size="size"
    min-width="0"
    min-height="0"
    :color="getColor(type)"
    :left="left"
    v-if="circle"
  >
    <v-icon
      v-if="type==='success'"
      class="readyIcon"
      color="white"
      :size="size-2"
    > mdi-check</v-icon>
    <v-icon
      v-else-if="type==='error'"
      class="readyIcon"
      color="white"
      :size="size-2"
    > mdi-exclamation</v-icon>
    <v-progress-circular
      v-else-if="spinner && (type==='warning' || type==='info')"
      class="readySpinner"
      indeterminate
      :size="spinnerSize"
      :width="spinnerWidth"
      color="white"
    >
    </v-progress-circular>
    <v-icon
      v-else-if="type==='warning'"
      class="readyIcon"
      color="white"
      :size="size-2"
    > mdi-exclamation</v-icon>
    <v-icon
      v-else-if="type==='info'"
      class="readyIcon"
      color="white"
      :size="size-2"
    > mdi-information-variant</v-icon>
  </v-avatar>
  <v-avatar
    :size="size"
    min-width="0"
    min-height="0"
    :color="getColor(type)"
    :left="left"
    class="mr-1"
    v-else
  >
    <v-icon
      v-if="type==='success'"
      class="readyIcon"
      :color="getColor(type)"
      :size="size"
    > mdi-check-circle</v-icon>
    <v-icon
      v-else-if="type==='error'"
      class="readyIcon"
      :color="getColor(type)"
      :size="size"
    > mdi-alert-circle</v-icon>
    <v-progress-circular
      v-else-if="spinner && (type==='warning' || type==='info')"
      class="readySpinner"
      indeterminate
      :size="spinnerSize"
      :width="spinnerWidth"
      :color="getColor(type)"
    >
    </v-progress-circular>
    <v-icon
      v-else-if="type==='warning'"
      class="readyIcon"
      :color="getColor(type)"
      :size="size"
    > mdi-alert-circle</v-icon>
    <v-icon
      v-else-if="type==='info'"
      class="readyIcon"
      :color="getColor(type)"
      :size="size"
    > mdi-information</v-icon>
  </v-avatar>
</template>

<script>
export default {
  name: "StatusIcon",
  props: {
    type: {
      default: "success",
      type: String,
    },
    size: Number,
    spinnerWidth: {
      default: 1.5,
      type: Number,
    },
    circle: Boolean,
    spinner: Boolean,
    left: Boolean,
  },
  data() {
    return {
      spinnerSize: this.size ? this.size * 0.75 : 12,
    };
  },
  methods: {
    getColor(type) {
      if (!this.circle) {
        return "";
      }

      if (this.spinner && (type === "warning" || type === "info")) {
        return "warning";
      }
      return type;
    },
  },
};
</script>