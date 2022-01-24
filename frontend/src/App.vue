<template>
  <div id="app">
    <VotingPage msg="Welcome to Your Vue.js App"/>
  </div>
</template>

<script>
// import HelloWorld from './components/HelloWorld.vue'
import VotingPage from './components/VotingPage.vue'

export default {
  name: 'App',
  components: {
    VotingPage
  }
}
</script>

<style>

html, body {
  height: 100vh;
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
