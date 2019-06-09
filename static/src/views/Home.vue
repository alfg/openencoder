<template>
  <div class="home">
    <h1>Dashboard</h1>
    <Dashboard />
  </div>
</template>

<script>
import Dashboard from '@/components/Dashboard.vue';

export default {
  name: 'home',
  components: {
    Dashboard,
  },
};
</script>
