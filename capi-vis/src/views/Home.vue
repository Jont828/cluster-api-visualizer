<template>
  <div id="home">
    <AppBar
      title="Cluster Management Overview"
      :isStraight="this.isStraight"
      @togglePathStyle="linkHandler"
    />
    <Tree :isStraight="this.isStraight" />
  </div>
</template>

<script>
import Tree from "../components/Tree.vue";
import AppBar from "../components/AppBar.vue";

export default {
  name: "Home",
  components: {
    Tree,
    AppBar,
  },
  data() {
    return {
      isStraight: false,
    };
  },
  methods: {
    linkHandler(val) {
      this.isStraight = val;
    },
  },
};
</script>

<style lang="less" scoped>
#home {
  height: 100%;
}
</style>