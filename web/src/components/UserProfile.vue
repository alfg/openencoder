<template>
  <div class="user-profile container">
    <h2>User Profile Settings</h2>

    <b-form class="mb-3" @submit="onSubmit">

      <b-form-group id="input-group-username" label="Username:" label-for="input-username">
          <b-form-input
            id="input-username"
            v-model="form.username"
          ></b-form-input>
      </b-form-group>

      <b-form-group id="input-group-password" label="Password:" label-for="input-password">
          <b-form-input
            id="input-password"
            v-model="form.password"
            type="password"
          ></b-form-input>
      </b-form-group>

      <b-form-group
        id="input-group-password-confirm"
        label="Password Confirm:"
        label-for="input-password-confirm">
          <b-form-input
            id="input-password-confirm"
            v-model="form.password_confirm"
            type="password"
          ></b-form-input>
      </b-form-group>

      <b-button type="submit" variant="primary">Save</b-button>
    </b-form>
  </div>
</template>

<script>
import auth from '../auth';

export default {
  name: 'user-profile',

  data() {
    return {
      form: {
        username: '',
        password: '',
        password_confirm: '',
      },
    };
  },

  mounted() {
    this.getUser();
  },

  methods: {
    getUser() {
      const url = '/api/user';

      this.$http.get(url, {
        headers: auth.getAuthHeader(),
      })
        .then(response => (
          response.json()
        ))
        .then((json) => {
          console.log(json);
        });
    },

    submitForm(data) {
      const url = '/api/user';

      this.$http.put(url, data, {
        headers: auth.getAuthHeader(),
      }).then(response => (
        response.json()
      )).then((json) => {
        console.log('Submitted form: ', json);
        this.dismissCountDown = this.dismissSecs;
      });
    },

    onSubmit(evt) {
      evt.preventDefault();
      console.log(JSON.stringify(this.form));
      this.submitForm(this.form);
    },
  },
};
</script>
