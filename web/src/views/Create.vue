<template>
  <div id="create" class="container">
    <h1>Create Job</h1>
    <JobForm />
  </div>
</template>

<script>
import JobForm from '@/components/JobForm.vue';

export default {
  name: 'Jobs',
  components: {
    JobForm,
  },
};
</script>
