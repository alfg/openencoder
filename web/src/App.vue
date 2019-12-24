<template>
  <div id="app">
    <b-navbar class="mb-4" toggleable="lg" type="dark" variant="dark">
      <b-navbar-brand href="#">Open Encoder <sup class="alpha">Alpha</sup></b-navbar-brand>
      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
      <b-navbar-nav class="ml-auto">
      <b-nav-item-dropdown right v-if="user.auth">
        <template slot="button-content">{{ user.username }}</template>
        <b-dropdown-item disabled>{{ user.role }}</b-dropdown-item>
        <b-dropdown-item to="/profile">Profile</b-dropdown-item>
        <b-dropdown-item href="#" @click="logout">Sign Out</b-dropdown-item>
      </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-navbar>

    <div
      class="container mb-4"
      v-if="!$route.meta.hideNavigation">
      <b-nav tabs>
        <b-nav-item to="/status">Status</b-nav-item>
        <b-nav-item to="/jobs">Jobs</b-nav-item>
        <b-nav-item v-if="isOperatorAdmin" to="/encode">Encode</b-nav-item>
        <b-nav-item v-if="isOperatorAdmin" to="/queue">Queue</b-nav-item>
        <b-nav-item v-if="isOperatorAdmin" to="/workers">Workers</b-nav-item>
        <b-nav-item v-if="isAdmin" to="/machines">Machines</b-nav-item>
        <b-nav-item v-if="isAdmin" to="/presets">Presets</b-nav-item>
        <b-nav-item v-if="isAdmin" to="/settings">Settings</b-nav-item>
      </b-nav>
    </div>

    <router-view />
  </div>
</template>

<script>
import auth from './auth';

export default {
  data() {
    return {
      user: auth.user,
      role: auth.role,
    };
  },

  computed: {
    isOperator() {
      return this.user.role === 'operator';
    },
    isAdmin() {
      return this.user.role === 'admin';
    },
    isOperatorAdmin() {
      return ['operator', 'admin'].includes(this.user.role);
    },
  },

  created() {
    auth.checkAuth(this);
  },

  methods: {
    logout() {
      auth.logout(this);
    },
  },
};
</script>

<style scoped>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

#app a.router-link-exact-active,
#app a.router-link-active {
  color: #495057;
  background-color: #fff;
  border-color: #dee2e6 #dee2e6 #fff;
}

.alpha {
  font-size: 12px;
  text-transform: uppercase;
}
</style>
