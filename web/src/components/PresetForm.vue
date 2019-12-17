<template>
  <div id="presets-form">
    <b-form class="mb-3" @submit="onSubmit" @reset="onReset" v-if="show">

      <b-form-group id="input-group-name" label="Name:" label-for="input-name">
          <b-form-input
            id="input-name"
            v-model="form.name"
            placeholder="Enter a preset name..."
          ></b-form-input>
      </b-form-group>

      <b-form-group id="input-group-description" label="Description:" label-for="input-description">
        <b-form-textarea
          id="textarea"
          v-model="form.description"
          placeholder="Enter a description of the preset..."
          rows="2"
          max-rows="6"
        ></b-form-textarea>
      </b-form-group>

      <b-form-group id="input-group-output" label="Output:" label-for="input-output">
          <b-form-input
            id="input-output"
            v-model="form.output"
            placeholder="Enter an output filename..."
          ></b-form-input>
      </b-form-group>

      <b-form-group id="input-group-active">
        <b-form-checkbox v-model="form.active">Active?</b-form-checkbox>
      </b-form-group>

      <div class="mb-4">
        <b-alert show variant="info">FFmpeg presets follow the <a href="https://alfg.github.io/ffmpeg-commander">ffmpeg-commander</a> JSON format. See documentation for details.</b-alert>
        <div ref="editor" class="editor"></div>
      </div>

      <b-button type="submit" variant="primary">Create</b-button>
    </b-form>

    <b-alert
      :show="dismissCountDown"
      dismissible
      fade
      variant="success"
      @dismissed="dismissCountDown=0"
      @dismiss-count-down="countDownChanged"
    >
      Submitted Preset!
    </b-alert>
  </div>
</template>

<script>
import JSONEditor from 'jsoneditor';
import 'jsoneditor/dist/jsoneditor.min.css';
import auth from '../auth';

const tmpl = `
{
  "input": "tears-of-steel.mp4",
  "output": "tears-of-steel.out.mp4",
  "container": "mp4",
  "video": {
    "codec": "x264",
    "preset": "none",
    "hardware_acceleration_option": "off",
    "pass": "1",
    "crf": 23,
    "bitrate": "3000k",
    "minrate": "3000k",
    "maxrate": "3000k",
    "bufsize": "3000k",
    "pixel_format": "auto",
    "frame_rate": "auto",
    "speed": "auto",
    "tune": "none",
    "profile": "none",
    "level": "none"
  },
  "audio": {
    "codec": "copy"
  }
}
`;

export default {
  data() {
    return {
      form: {
        name: '',
        description: '',
        output: '',
        data: '',
        active: false,
      },
      show: true,
      dismissSecs: 5,
      dismissCountDown: 0,
      showDismissibleAlert: false,
      editor: null,
      data: null,
    };
  },

  mounted() {
    // Load JSONEditor.
    const container = this.$refs.editor;
    const options = {
      mode: 'code',
      modes: ['code', 'text', 'tree', 'preview'],
    };
    this.editor = new JSONEditor(container, options);
    this.editor.set(JSON.parse(tmpl));
  },

  destroyed() {
    this.editor = null;
  },

  methods: {
    countDownChanged(dismissCountDown) {
      this.dismissCountDown = dismissCountDown;
    },

    submitForm(data) {
      const url = '/api/presets';

      this.$http.post(url, data, {
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
      const presetData = JSON.stringify(this.editor.get());
      this.form.data = presetData;
      this.submitForm(this.form);
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

<style scoped>
.editor {
  height: 50vh;
}
</style>
