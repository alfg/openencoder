<template>
  <div id="presets-form">
    <div class="mb-2 text-right">
      <b-button v-b-modal.create-preset>Create</b-button>
      <b-modal id="create-preset" title="Create a Preset">
        <p class="my-4">Hello from modal!</p>
      </b-modal>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
    };
  },

  mounted() {
  },

  methods: {
  },
};
</script>
