<template>
  <div id="update-password" class="container col-lg-4 col-md-4">
    <UpdatePasswordForm />
  </div>
</template>

<script>
import UpdatePasswordForm from '@/components/UpdatePasswordForm.vue';

export default {
  name: 'update-password',
  components: {
    UpdatePasswordForm,
  },
};
</script>
