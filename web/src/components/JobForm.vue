<template>
  <div id="job-form">
    <b-form class="mb-3" @submit="onSubmit" @reset="onReset" v-if="show">
      <b-form-group id="input-group-1" label="Select Encoding Preset:" label-for="input-1">
        <b-form-select
          id="input-1"
          v-model="form.preset"
          :options="presets"
          required
        ></b-form-select>
      </b-form-group>

      <b-form-group id="input-group-2" label="Select file:" label-for="input-2">
          <b-form-input
            id="input-2"
            v-model="form.source"
            placeholder=""
            @focus="onFileFocus"
          ></b-form-input>
      </b-form-group>

      <div v-if="showFileBrowser">
        <S3Browser v-on:file="onFileSelect" />
      </div>

      <b-form-group id="input-group-3" label="Destination:" label-for="input-3">
          <b-form-input
            id="input-3"
            v-model="form.dest"
            placeholder=""
            readonly
          ></b-form-input>
      </b-form-group>

      <b-button type="submit" variant="primary">Submit</b-button>
    </b-form>

    <b-alert
      :show="dismissCountDown"
      dismissible
      fade
      variant="success"
      @dismissed="dismissCountDown=0"
      @dismiss-count-down="countDownChanged"
    >
      Submitted Job!
    </b-alert>
  </div>
</template>

<script>
import auth from '../auth';
import S3Browser from '@/components/S3Browser.vue';

export default {
  components: {
    S3Browser,
  },
  data() {
    return {
      form: {
        preset: null,
        source: null,
        dest: null,
      },
      presetsData: [],
      show: true,
      showFileBrowser: false,
      dismissSecs: 5,
      dismissCountDown: 0,
      showDismissibleAlert: false,
    };
  },

  computed: {
    presets() {
      return this.presetsData.filter(x => x.active).map(x => x.name);
    },
  },

  mounted() {
    this.getPresets();
  },

  methods: {
    countDownChanged(dismissCountDown) {
      this.dismissCountDown = dismissCountDown;
    },

    getPresets() {
      const url = '/api/presets';

      this.$http.get(url, {
        headers: auth.getAuthHeader(),
      })
        .then(response => (
          response.json()
        ))
        .then((json) => {
          this.presetsData = json.presets;
        });
    },

    submitJob(data) {
      const url = '/api/jobs';

      this.$http.post(url, data, {
        headers: auth.getAuthHeader(),
      }).then(response => (
        response.json()
      )).then((json) => {
        console.log('Submitted job: ', json);
        this.dismissCountDown = this.dismissSecs;
      });
    },

    onFileSelect(file) {
      this.form.source = file;
      this.form.dest = file.replace('src', 'dst').replace('.mp4', '/');
      this.showFileBrowser = false;
    },

    onFileFocus() {
      this.showFileBrowser = true;
    },

    onSubmit(evt) {
      evt.preventDefault();
      console.log(JSON.stringify(this.form));
      this.submitJob(this.form);
    },

    onReset(evt) {
      evt.preventDefault();

      // Reset our form values
      this.form.preset = null;

      // Trick to reset/clear native browser form validation state
      this.show = false;
      this.$nextTick(() => {
        this.show = true;
      });
    },
  },
};
</script>
