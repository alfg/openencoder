<template>
  <div id="machines" class="container">
    <MachinesForm />
    <MachinesTable />
  </div>
</template>

<script>
import MachinesTable from '@/components/MachinesTable.vue';
import MachinesForm from '@/components/MachinesForm.vue';

export default {
  name: 'Machines',
  components: {
    MachinesForm,
    MachinesTable,
  },
};
</script>
