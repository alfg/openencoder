<template>
  <div id="app">
    <b-navbar class="mb-4" toggleable="lg" type="dark" variant="dark">
      <b-navbar-brand href="#">Cloud Encoder</b-navbar-brand>
      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    </b-navbar>

    <div class="container mb-4">
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
