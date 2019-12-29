<template>
  <div class="user-profile container">
    <h2>User Profile Settings</h2>

    <b-form class="mb-3" @submit="onSubmit">

      <b-form-group id="input-group-username" label="Username:" label-for="input-username">
          <b-form-input
            id="input-username"
            required
            aria-describedby="username-help-block"
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

      <b-form-group
        id="input-group-new-password"
        label="New Password:"
        label-for="input-new-password">
          <b-form-input
            id="input-new-password"
            v-model="form.new_password"
            type="password"
          ></b-form-input>
      </b-form-group>

      <b-form-group
        id="input-group-password-verify"
        label="Verify Password:"
        label-for="input-password-verify">
          <b-form-input
            id="input-password-verify"
            v-model="form.verify_password"
            type="password"
          ></b-form-input>
      </b-form-group>

      <b-form-group id="input-group-role" label="Role:" label-for="input-role">
          <b-form-input
            id="input-role"
            v-model="role"
            readonly
          ></b-form-input>
      </b-form-group>

      <b-button type="submit" variant="primary">Save</b-button>
    </b-form>

    <b-alert
      class="mt-4"
      :show="dismissCountDown"
      dismissible
      fade
      :variant="messageType"
      @dismissed="dismissCountDown=0"
      @dismiss-count-down="countDownChanged"
    >
      {{ message }}
    </b-alert>
  </div>
</template>

<script>
import api from '../api';
import auth from '../auth';

export default {
  name: 'user-profile',

  data() {
    return {
      form: {
        username: '',
        current_password: '',
        new_password: '',
        verify_password: '',
      },
      role: '',
      dismissSecs: 5,
      dismissCountDown: 0,
      showDismissibleAlert: false,
      message: '',
      messageType: '',
    };
  },

  mounted() {
    this.getUser();
  },

  methods: {
    countDownChanged(dismissCountDown) {
      this.dismissCountDown = dismissCountDown;
    },

    getUser() {
      api.getCurrentUser(this, (err, json) => {
        this.form.username = json.username;
        this.role = json.role;
      });
    },

    submitForm(data) {
      api.updateCurrentUser(this, data, (err, json) => {
        if (err) {
          this.message = err.body && err.body.message;
          this.dismissCountDown = this.dismissSecs;
        }

        console.log('Submitted form: ', json);
        this.messageType = 'success';
        this.message = json && json.message;
        this.dismissCountDown = this.dismissSecs;

        // Re-auth with new token.
        auth.login(this, {
          username: this.form.username,
          password: this.form.new_password || this.form.current_password,
        }, 'profile', (err2) => {
          if (err2) {
            console.log('err', err2);
          }
        });
      });
    },

    onSubmit(evt) {
      evt.preventDefault();
      console.log(JSON.stringify(this.form));

      if (this.form.new_password.length > 0
        && (this.form.new_password !== this.form.verify_password)) {
        this.messageType = 'danger';
        this.message = 'Passwords do not match';
        this.dismissCountDown = this.dismissSecs;
        return;
      }

      this.submitForm(this.form);
    },
  },
};
</script>
