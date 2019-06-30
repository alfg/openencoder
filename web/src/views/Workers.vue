<template>
  <div id="workers" class="container">
    <h1>Workers</h1>
    <WorkersTable />
  </div>
</template>

<script>
import WorkersTable from '@/components/WorkersTable.vue';

export default {
  name: 'Workers',
  components: {
    WorkersTable,
  },
};
</script>
