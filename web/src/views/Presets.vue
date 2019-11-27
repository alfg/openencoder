<template>
  <div id="presets" class="container">
    <PresetsTable />
  </div>
</template>

<script>
import PresetsTable from '@/components/PresetsTable.vue';

export default {
  name: 'Presets',
  components: {
    PresetsTable,
  },
};
</script>
