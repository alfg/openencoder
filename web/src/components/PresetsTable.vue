<template>
  <div id="presets-table">
    <div class="mb-2 text-right">
      <b-button to="presets/create">Create Preset</b-button>
    </div>

    <b-table
      striped hover dark
      selectable
      select-mode="single"
      :fields="fields"
      :items="items"
      @row-selected="onRowSelected">

      <template v-slot:cell(active)="data">
        {{ data.item.active ? '✔️' : '❌' }}
      </template>
    </b-table>

    <b-form class="mb-3" @submit="onSubmit" v-show="data">
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
        <b-alert
          show
          variant="info"
        >FFmpeg presets follow the <a href="https://alfg.github.io/ffmpeg-commander">ffmpeg-commander</a> JSON
        format. See <a href="https://github.com/alfg/openencoder/wiki/Creating-Presets">wiki</a> for details.
        </b-alert>
        <div ref="editor"></div>
      </div>

      <b-button type="submit" variant="primary">Save</b-button>
    </b-form>

    <b-alert
      :show="dismissCountDown"
      dismissible
      fade
      variant="success"
      @dismissed="dismissCountDown=0"
      @dismiss-count-down="countDownChanged"
    >
      Updated Preset!
    </b-alert>

    <h2 class="text-center" v-if="items.length === 0">No Presets Found</h2>

    <b-pagination-nav
      @change="onChangePage"
      :link-gen="linkGen"
      :number-of-pages="pages"
      use-router></b-pagination-nav>
  </div>
</template>

<script>
import JSONEditor from 'jsoneditor';
import 'jsoneditor/dist/jsoneditor.min.css';
import api from '../api';

export default {
  data() {
    return {
      fields: ['id', 'name', 'description', 'active'],
      items: [],
      count: 0,
      editor: null,
      data: null,
      form: {
        name: '',
        description: '',
        output: '',
        data: '',
        active: false,
      },
      dismissSecs: 5,
      dismissCountDown: 0,
      showDismissibleAlert: false,
    };
  },

  computed: {
    pages() {
      return this.count === 0 ? 1 : Math.ceil(this.count / 10);
    },
  },

  mounted() {
    const page = this.$route.query.page || 0;
    this.getPresets(page);

    // Load JSONEditor.
    const container = this.$refs.editor;
    const options = {
      mode: 'tree',
      modes: ['code', 'text', 'tree', 'preview'],
    };
    this.editor = new JSONEditor(container, options);
  },

  destroyed() {
    this.editor = null;
  },

  methods: {
    countDownChanged(dismissCountDown) {
      this.dismissCountDown = dismissCountDown;
    },

    linkGen(pageNum) {
      return pageNum === 1 ? '?' : `?page=${pageNum}`;
    },

    onChangePage(event) {
      this.getPresets(event);
    },

    getPresets(page) {
      api.getPresets(this, page, (err, json) => {
        this.items = json && json.presets;
        this.count = json && json.count;
      });
    },

    updatePreset(data) {
      api.updatePreset(this, data, (err, json) => {
        console.log('Submitted form: ', json);
        this.dismissCountDown = this.dismissSecs;
      });
    },

    onSubmit(evt) {
      evt.preventDefault();
      console.log(JSON.stringify(this.form));
      this.form.data = JSON.stringify(this.editor.get());
      this.updatePreset(this.form);
    },

    onRowSelected(items) {
      if (items.length > 0) {
        [this.form] = items;

        this.data = this.form.data || {};
        this.editor.set(JSON.parse(this.data));
        this.editor.expandAll();
      } else {
        this.data = null;
      }
    },
  },
};
</script>
