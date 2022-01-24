<template>
  <div style="height: 100%;" >
  <b-navbar toggleable="lg" type="dark" variant="info">
    <b-navbar-brand href="#">Voting</b-navbar-brand>

    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

    <b-collapse id="nav-collapse" is-nav>

      <!-- Right aligned nav items -->
      <b-navbar-nav class="ml-auto">

        <b-nav-item-dropdown right>
          <!-- Using 'button-content' slot -->
          <template #button-content>
            <em>User</em>
          </template>
          <b-dropdown-item href="#">Profile</b-dropdown-item>
          <b-dropdown-item href="#">Sign Out</b-dropdown-item>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
  <b-container class="main" fluid=True>
  <b-row class="justify-content-md-center">
    <b-col cols="4">
      <SubPoll 
      v-for="subpoll in this.$store.state.subpolls" 
      :key="subpoll.id"
      v-bind:id="subpoll.id" 
      v-bind:title="subpoll.title" 
      v-bind:description="subpoll.description" 
      v-bind:options="subpoll.options"
      v-bind:voted="subpoll.voted"
      v-bind:isOpen="subpoll.isOpen"></SubPoll>
    </b-col>
  </b-row>

    
  </b-container>
</div>
</template>

<script>

import SubPoll from './SubPoll.vue'

export default {
  name: 'VotingPage',
  props: {
    msg: String
  },
  components: {
    SubPoll
  }
}
</script>

<style scoped>

.main {
  background-color: rgb(243, 243, 243);
  min-height: 100%;
  height: 100%;
}

</style>