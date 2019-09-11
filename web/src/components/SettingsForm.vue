<template>
  <div id="settings-form">
    <h4>Settings</h4>

    <b-form class="mb-3" @submit="onSubmit">
      <div
        v-for="(o, i) in settings"
        v-bind:key="i"
      >
        <b-form-group
          id="fieldset-horizontal"
          label-cols-sm="4"
          label-cols-lg="4"
          label-for="input-horizontal"
          :label="o.title"
          :description="o.description"
        >
          <b-form-input
            id="input-horizontal"
            v-model="form[o.name]"
            autocomplete="off"
            :type="o.secure && hide ? 'password' : 'text'"
          ></b-form-input>
        </b-form-group>
      </div>

      <b-button type="submit" variant="primary">Save</b-button>
      <b-button
        class="ml-2"
        variant="primary"
        @click="onClickShow">{{this.hide ? 'Show' : 'Hide'}}</b-button>
    </b-form>

    <b-alert
      :show="dismissCountDown"
      dismissible
      fade
      variant="success"
      @dismissed="dismissCountDown=0"
      @dismiss-count-down="countDownChanged"
    >
      Updated settings!
    </b-alert>
  </div>
</template>

<script>
import auth from '../auth';

export default {
  data() {
    return {
      form: {},
      settings: {},
      hide: true,
      dismissSecs: 5,
      dismissCountDown: 0,
      showDismissibleAlert: false,
    };
  },

  mounted() {
    this.getSettings();
  },

  methods: {
    countDownChanged(dismissCountDown) {
      this.dismissCountDown = dismissCountDown;
    },

    getSettings() {
      const url = '/api/settings';

      this.$http.get(url, {
        headers: auth.getAuthHeader(),
      })
        .then(response => (
          response.json()
        ))
        .then((json) => {
          if (json.settings) {
            this.settings = json.settings;

            // Populate form items if availble.
            this.settings.forEach(item => {
              this.form[item.name] = item.value;
            });
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },

    updateSettings(data) {
      const url = '/api/settings';

      this.$http.put(url, data, {
        headers: auth.getAuthHeader(),
      }).then(response => (
        response.json()
      )).then((json) => {
        this.dismissCountDown = this.dismissSecs;
      });
    },

    onSubmit(evt) {
      evt.preventDefault();
      this.updateSettings(this.form);
    },

    onClickShow() {
      console.log('show');
      this.hide = !this.hide;
    },
  },
};
</script>

<style scoped>
</style>
