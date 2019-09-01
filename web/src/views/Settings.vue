<template>
  <div id="settings" class="container">
    <SettingsForm />
  </div>
</template>

<script>
import SettingsForm from '@/components/SettingsForm.vue';

export default {
  name: 'Settings',
  components: {
    SettingsForm,
  },
};
</script>
