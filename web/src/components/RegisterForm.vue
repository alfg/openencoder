<template>
  <div id="login-form" class="container">
    <h2>Register</h2>
    <b-form @submit="onSubmit" v-if="show">
      <b-form-group
        label="Email:"
        label-for="input-email"
      >
        <b-form-input
          id="input-email"
          v-model="form.username"
          type="email"
          required
          placeholder="you@email.com"
        ></b-form-input>
      </b-form-group>

      <b-form-group label="Password:" label-for="input-password">
        <b-form-input
          id="input-password"
          v-model="form.password"
          type="password"
          required
        ></b-form-input>
      </b-form-group>

      <b-button type="submit" variant="primary">Submit</b-button>
    </b-form>
  </div>
</template>

<script>
import auth from '../auth';

export default {
  components: {
  },
  data() {
    return {
      form: {
        username: '',
        password: '',
      },
      show: true,
    };
  },

  methods: {
    onSubmit(event) {
      event.preventDefault();
      auth.register(this, this.form, 'login');
    },
  },
};
</script>
