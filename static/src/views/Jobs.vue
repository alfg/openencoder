<template>
  <div class="jobs">
    <h1>Jobs</h1>
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
