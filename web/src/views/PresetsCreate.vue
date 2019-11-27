<template>
  <div id="presets-create" class="container">
    <PresetForm />
  </div>
</template>

<script>
import PresetForm from '@/components/PresetForm.vue';

export default {
  name: 'PresetsCreate',
  components: {
    PresetForm,
  },
};
</script>
