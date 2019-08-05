<template>
  <div id="machines" class="container">
    <MachinesTable />
  </div>
</template>

<script>
import MachinesTable from '@/components/MachinesTable.vue';

export default {
  name: 'Machines',
  components: {
    MachinesTable,
  },
};
</script>
