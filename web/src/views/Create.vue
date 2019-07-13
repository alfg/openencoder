<template>
  <div id="create" class="container">
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
