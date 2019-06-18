<template>
  <div class="jobs">
    <h1>Jobs</h1>
    <JobForm />
    <JobsTable />
  </div>
</template>

<script>
import JobsTable from '@/components/JobsTable.vue';
import JobForm from '@/components/JobForm.vue';

export default {
  name: 'Jobs',
  components: {
    JobsTable,
    JobForm,
  },
};
</script>
