<template>
  <div class="alertWrapper">
    <v-alert
      class="alertError"
      type="error"
      dismissible
      border="left"
      elevation="2"
      colored-border
      v-if="message"
    >
      {{ message }}
    </v-alert>

  </div>
</template>

<script>
export default {
  name: "AlertError",
  props: {
    message: String,
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