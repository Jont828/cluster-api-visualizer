<template>
  <div id="app">
    <Tree/>
  </div>
</template>

<script>
import Tree from './components/Tree.vue'

export default {
  name: 'App',
  components: {
    Tree
  }
}
</script>

<style>
html, body {
  height: 100%;
  padding: 0;
  margin: 0;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  height: 100%;
}
</style>
