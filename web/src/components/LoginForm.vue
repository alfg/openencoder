<template>
  <div id="login-form" class="container">
    <h2>Login</h2>
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

    <b-card class="mt-3" header="Form Data Result">
      <pre class="m-0">{{ form }}</pre>
    </b-card>
  </div>
</template>

<script>
import cookie from '../cookie';
import store from '../store';

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

  computed: {
  },

  mounted() {
  },

  methods: {
    async getJWTToken(data) {
      const url = '/api/login';

      const res = await fetch(url, {
        method: 'POST',
        body: JSON.stringify(data),
        headers: {
          'Content-Type': 'application/json',
        },
      });

      const json = await res.json();
      if (json.code === 200) {
        console.log(json);

        // Set the cookie.
        cookie.set('token', json.token);
        store.setTokenAction(json.token);

        // Redirect to home.
        this.$router.push({ name: 'home' });
      }
    },

    onSubmit(event) {
      event.preventDefault();
      alert(JSON.stringify(this.form));
      this.getJWTToken(this.form);
    },
  },
};
</script>
