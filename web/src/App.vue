<template>
  <div id="app">
    <b-navbar class="mb-4" toggleable="lg" type="dark" variant="dark">
      <b-navbar-brand href="#">Open Encoder</b-navbar-brand>
      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    </b-navbar>

    <div
      class="container mb-4"
      v-if="!$route.meta.hideNavigation">
      <b-nav tabs>
        <b-nav-item to="/">Dashboard</b-nav-item>
        <b-nav-item to="/create">Create</b-nav-item>
        <b-nav-item to="/jobs">Jobs</b-nav-item>
        <b-nav-item to="/queue">Queue</b-nav-item>
        <b-nav-item to="/workers">Workers</b-nav-item>
        <b-nav-item to="/machines">Machines</b-nav-item>
      </b-nav>
    </div>

    <router-view/>
  </div>
</template>

<script>
import store from './store';
import cookie from './cookie';

export default {
  created() {
    // Check if token exists from cookie and set the store.
    // If not, then redirect to the login page to get a new token.
    const token = cookie.get('token');
    if (!store.state.token && this.$route.name !== 'register') {
      if (token) {
        store.setTokenAction(token);
      } else {
        this.$router.push({ name: 'login' });
      }
    }
  },
};
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /* text-align: center; */
  color: #2c3e50;
}

#app a.router-link-exact-active {
  /* color: #ffffff; */
  color: #495057;
  background-color: #fff;
  border-color: #dee2e6 #dee2e6 #fff;
}
</style>
