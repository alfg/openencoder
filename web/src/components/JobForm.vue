<template>
  <div id="job-form">
    <b-form @submit="onSubmit" @reset="onReset" v-if="show">
      <b-form-group id="input-group-1" label="Profile:" label-for="input-1">
        <b-form-select
          id="input-1"
          v-model="form.profile"
          :options="profiles"
          required
        ></b-form-select>
      </b-form-group>

      <b-form-group id="input-group-2" label="Select file:" label-for="input-2">
          <b-form-input
            id="input-2"
            v-model="form.file"
            placeholder=""
            @focus="onFileFocus"
            @blur="onFileBlur"
          ></b-form-input>
      </b-form-group>

      <div v-if="showFileBrowser">
        <S3Browser v-on:file="onFileSelect" />
      </div>

      <b-button type="submit" variant="primary">Submit</b-button>
    </b-form>
    <b-card class="mt-1" header="Form Data Result">
      <pre class="m-0">{{ form }}</pre>
    </b-card>
  </div>
</template>

<script>
import S3Browser from '@/components/S3Browser.vue';

export default {
  components: {
    S3Browser,
  },
  data() {
    return {
      form: {
        profile: null,
        file: null,
      },
      profiles: [{ text: 'Select One', value: null }, 'Test 1', 'Test 2', 'Test 3'],
      show: true,
      showFileBrowser: false,
    };
  },
  methods: {
    onFileSelect(file) {
      this.form.file = file;
      this.showFileBrowser = false;
    },
    onFileFocus() {
      this.showFileBrowser = true;
    },
    onFileBlur() {
      // this.showFileBrowser = false;
    },
    onSubmit(evt) {
      evt.preventDefault();
      console.log(JSON.stringify(this.form));
    },
    onReset(evt) {
      evt.preventDefault();

      // Reset our form values
      this.form.profile = null;

      // Trick to reset/clear native browser form validation state
      this.show = false;
      this.$nextTick(() => {
        this.show = true;
      });
    },
  },
};
</script>
