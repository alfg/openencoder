<template>
  <div id="settings-form">
    <h4>Settings</h4>

    <b-form class="mb-3" @submit="onSubmit">
      <div
        v-for="(o, i) in sortedSettings"
        v-bind:key="i"
      >
        <div v-show="!isHidden(o.name)">
          <b-form-group
            :id="`fieldset-horizontal-${i}`"
            label-cols-sm="4"
            label-cols-lg="4"
            :label-for="`input-horizontal-${i}`"
            :label="o.title"
            :description="o.description"
          >
            <div v-if="isSelectInput(o.name)">
              <b-form-select
                v-model="form[o.name]"
                :options="getSelectOptions(o.name)"
                @change="onSelectChange"
              >
              </b-form-select>
            </div>
            <div v-else>
              <b-form-input
                :id="`input-horizontal-${i}`"
                v-model="form[o.name]"
                autocomplete="off"
                :type="o.secure && hide ? 'password' : 'text'"
              ></b-form-input>
            </div>
          </b-form-group>
        </div>
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

    <div v-show="false">
      {{forceUpdate}}
    </div>
  </div>
</template>

<script>
import api from '../api';

const SORTED_OPTIONS = [
  'DIGITAL_OCEAN_ACCESS_TOKEN',
  'STORAGE_DRIVER',
  'S3_PROVIDER',
  'S3_ACCESS_KEY',
  'S3_SECRET_KEY',
  'S3_INBOUND_BUCKET',
  'S3_INBOUND_BUCKET_REGION',
  'S3_OUTBOUND_BUCKET',
  'S3_OUTBOUND_BUCKET_REGION',
  'S3_PROVIDER',
  'S3_STREAMING',
  'FTP_ADDR',
  'FTP_USERNAME',
  'FTP_PASSWORD',
  'SLACK_WEBHOOK',
];

export default {
  data() {
    return {
      form: {},
      settings: [],
      providers: [
        { value: '', text: 'Select an S3 Provider', disabled: true },
        { value: 'digitalocean', text: 'Digital Ocean' },
        { value: 'aws', text: 'Amazon AWS' },
      ],
      streamingOptions: [
        { value: '', text: 'Select an S3 Streaming Option', disabled: true },
        { value: 'enabled', text: 'Enabled' },
        { value: 'disabled', text: 'Disabled' },
      ],
      storageOptions: [
        { value: '', text: 'Select a Storage Option', disabled: true },
        { value: 's3', text: 'S3' },
        { value: 'ftp', text: 'FTP' },
      ],
      hide: true,
      dismissSecs: 5,
      dismissCountDown: 0,
      showDismissibleAlert: false,
      forceUpdate: 0,
    };
  },

  computed: {
    sortedSettings() {
      const newOrder = [];
      if (this.settings.length > 0) {
        SORTED_OPTIONS.forEach((item) => {
          const a = this.settings.find(o => o.name === item);
          newOrder.push(a);
        });
      }
      return newOrder;
    },
  },

  mounted() {
    this.getSettings();
  },

  methods: {
    onSelectChange() {
      this.forceUpdate = this.forceUpdate + 1;
    },

    isSelectInput(inputName) {
      return ['S3_PROVIDER', 'S3_STREAMING', 'STORAGE_DRIVER'].includes(inputName);
    },

    isHidden(inputName) {
      const options = ['FTP', 'S3'];
      const prefix = inputName.split('_')[0];
      const { STORAGE_DRIVER } = this.form;
      if (options.includes(prefix) && STORAGE_DRIVER.toUpperCase() !== prefix) {
        return true;
      }
      return false;
    },

    getSelectOptions(inputName) {
      switch (inputName) {
        case 'S3_PROVIDER':
          return this.providers;

        case 'S3_STREAMING':
          return this.streamingOptions;

        case 'STORAGE_DRIVER':
          return this.storageOptions;

        default:
          return [];
      }
    },

    countDownChanged(dismissCountDown) {
      this.dismissCountDown = dismissCountDown;
    },

    getSettings() {
      api.getSettings(this, (err, json) => {
        if (json.settings) {
          this.settings = json.settings;

          // Populate form items if availble.
          this.settings.forEach((item) => {
            this.form[item.name] = item.value;
          });
        }
      });
    },

    updateSettings(data) {
      api.updateSettings(this, data, (err, json) => {
        this.dismissCountDown = this.dismissSecs;
        console.log('Settings updated', json);
      });
    },

    onSubmit(evt) {
      evt.preventDefault();
      this.updateSettings(this.form);
    },

    onClickShow() {
      this.hide = !this.hide;
    },
  },
};
</script>

<style scoped>
</style>
