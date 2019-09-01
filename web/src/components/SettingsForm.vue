<template>
  <div id="settings-form">
    <h4>Settings</h4>

    <b-form class="" @submit="onSubmit">
      <b-form-group
        id="fieldset-horizontal"
        label-cols-sm="4"
        label-cols-lg="4"
        label="AWS ACCESS KEY"
        label-for="input-horizontal"
      >
        <b-form-input id="input-horizontal"></b-form-input>
      </b-form-group>

      <b-form-group
        id="fieldset-horizontal"
        label-cols-sm="4"
        label-cols-lg="4"
        label="AWS SECRET KEY"
        label-for="input-horizontal"
      >
        <b-form-input id="input-horizontal"></b-form-input>
      </b-form-group>

      <b-form-group
        id="fieldset-horizontal"
        label-cols-sm="4"
        label-cols-lg="4"
        label="DIGITAL OCEAN ACCESS KEY"
        label-for="input-horizontal"
      >
        <b-form-input id="input-horizontal"></b-form-input>
      </b-form-group>

      <b-form-group
        id="fieldset-horizontal"
        label-cols-sm="4"
        label-cols-lg="4"
        label="DIGITAL OCEAN SECRET KEY"
        label-for="input-horizontal"
      >
        <b-form-input id="input-horizontal"></b-form-input>
      </b-form-group>

      <b-button type="submit" variant="primary">Save</b-button>
      <b-button
        class="ml-2"
        type="submit" variant="primary">Show</b-button>
    </b-form>
  </div>
</template>

<script>
import auth from '../auth';

export default {
  data() {
    return {
      form: {
      },
    };
  },

  computed: {
  },

  mounted() {
  },

  methods: {
    onSubmit(evt) {
      evt.preventDefault();
      console.log(JSON.stringify(this.form));
    },
  },
};
</script>

<style scoped>
</style>
