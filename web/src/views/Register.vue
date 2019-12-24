<template>
  <div id="register" class="container col-lg-4 col-md-4">
    <RegisterForm />
  </div>
</template>

<script>
import RegisterForm from '@/components/RegisterForm.vue';

export default {
  name: 'register',
  components: {
    RegisterForm,
  },
};
</script>
