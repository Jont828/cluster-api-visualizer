<template>
  <div id="home">
    <AppBar title="Cluster Management Overview" />
    <Tree />
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
};
</script>

<style lang="less" scoped>
#home {
  height: 100%;
}
</style>