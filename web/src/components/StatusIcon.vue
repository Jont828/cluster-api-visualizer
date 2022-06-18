<template>
  <v-avatar
    :size="size"
    min-width="0"
    min-height="0"
    :color="type"
    :left="left"
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
      v-else-if="type==='warning'"
      class="readySpinner"
      indeterminate
      :size="size*0.8"
      :width="spinnerWidth"
      color="white"
    >
    </v-progress-circular>
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
    size: {
      default: 16,
      type: Number,
    },
    spinnerWidth: {
      default: 1.5,
      type: Number,
    },
    left: Boolean,
  },
};
</script>