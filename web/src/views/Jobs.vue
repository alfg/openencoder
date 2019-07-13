<template>
  <div id="jobs" class="container">
    <JobsTable />
  </div>
</template>

<script>
import JobsTable from '@/components/JobsTable.vue';

export default {
  name: 'Jobs',
  components: {
    JobsTable,
  },
};
</script>
