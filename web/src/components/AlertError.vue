<template>
  <div class="alertWrapper">
    <v-alert
      :value="value"
      @input="update"
      class="alertError"
      type="error"
      dismissible
      border="left"
      elevation="2"
      colored-border
    >
      {{ message }}
    </v-alert>

  </div>
</template>

<script>
export default {
  name: "AlertError",
  props: {
    value: Boolean, // v-model
    message: String,
  },
  methods: {
    update(val) {
      // console.log("Emitting 'input' as ", val);
      this.$emit("input", val);
    },
  },
};
</script>

<style lang="less" scoped>
.alertError {
  min-width: 500px;
  text-align: left;
  display: inline-block;
}

.alertWrapper {
  text-align: center;
  width: 100%;
  position: fixed;
  bottom: 30px;
  z-index: 9999;
}
</style>