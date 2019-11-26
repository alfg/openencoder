<template>
  <div id="presets-form">
    <h4>Presets</h4>

    <b-form class="mb-3" @submit="onSubmit">
    </b-form>

    <b-alert
      :show="dismissCountDown"
      dismissible
      fade
      variant="success"
      @dismissed="dismissCountDown=0"
      @dismiss-count-down="countDownChanged"
    >
      Updated!
    </b-alert>
  </div>
</template>

<script>
import auth from '../auth';

export default {
  data() {
    return {
      form: {},
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

    onSubmit(evt) {
      evt.preventDefault();
    },
  },
};
</script>

<style scoped>
</style>
