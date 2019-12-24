<template>
  <div id="login" class="container col-lg-4 col-md-4">
    <LoginForm />
  </div>
</template>

<script>
import LoginForm from '@/components/LoginForm.vue';

export default {
  name: 'login',
  components: {
    LoginForm,
  },
};
</script>
