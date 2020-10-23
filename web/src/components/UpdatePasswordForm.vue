<template>
  <div id="update-password-form" class="container">
    <h2>Update Your Password</h2>
    <b-alert show variant="warning">Your account requires a password update.</b-alert>
    <b-form @submit="onSubmit">
      <b-form-group id="input-group-username" label="Username:" label-for="input-username">
          <b-form-input
            id="input-username"
            required
            v-model="form.username"
          ></b-form-input>
      </b-form-group>

      <b-form-group
        id="input-group-current-password"
        label="Current Password (required):"
        label-for="input-current-password">
          <b-form-input
            id="input-current-password"
            v-model="form.current_password"
            type="password"
            required
          ></b-form-input>
      </b-form-group>
      <b-form-group label="Password:" label-for="input-password">
        <b-form-input
          id="input-password"
          v-model="form.new_password"
          type="password"
          required
        ></b-form-input>
      </b-form-group>

      <b-form-group label="Verify Password:" label-for="input-verify-password">
        <b-form-input
          id="input-verify-password"
          v-model="form.verify_password"
          type="password"
          :state="form.verify_password !== '' && form.new_password === form.verify_password"
          required
        ></b-form-input>
      </b-form-group>

      <b-button type="submit" variant="primary">Submit</b-button>
    </b-form>

    <b-alert
      class="mt-4"
      :show="dismissCountDown"
      dismissible
      fade
      variant="danger"
      @dismissed="dismissCountDown=0"
      @dismiss-count-down="countDownChanged"
    >
      {{ errorMessage }}
    </b-alert>
  </div>
</template>

<script>
import auth from '../auth';

export default {
  data() {
    return {
      form: {
        username: '',
        current_password: '',
        new_password: '',
        verify_password: '',
      },
      dismissSecs: 5,
      dismissCountDown: 0,
      showDismissibleAlert: false,
      errorMessage: '',
    };
  },

  methods: {
    countDownChanged(dismissCountDown) {
      this.dismissCountDown = dismissCountDown;
    },
    onSubmit(event) {
      event.preventDefault();
      auth.updatePassword(this, this.form, 'login', (err) => {
        if (err) {
          this.errorMessage = err.body && err.body.message;
          this.dismissCountDown = this.dismissSecs;
        }
      });
    },
  },
};
</script>
