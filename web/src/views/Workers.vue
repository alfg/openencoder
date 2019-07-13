<template>
  <div id="workers" class="container">
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
