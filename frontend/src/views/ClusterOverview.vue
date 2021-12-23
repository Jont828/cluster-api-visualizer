<template>
  <div id="overview">
    <AppBar
      title="Cluster Management Overview"
      :isStraight="this.isStraight"
      @togglePathStyle="linkHandler"
    />
    <OverviewTree :isStraight="this.isStraight" />
  </div>
</template>

<script>
import OverviewTree from "../components/OverviewTree.vue";
import AppBar from "../components/AppBar.vue";

export default {
  name: "Overview",
  components: {
    OverviewTree,
    AppBar,
  },
  mounted() {
    document.title = "Cluster Management Overview";
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
#overview {
  height: 100%;
}
</style>