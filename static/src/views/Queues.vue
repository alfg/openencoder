<template>
  <div id="queues" class="container">
    <h1>Queues</h1>
    <QueuesTable />
  </div>
</template>

<script>
import QueuesTable from '@/components/QueuesTable.vue';

export default {
  name: 'Queues',
  components: {
    QueuesTable,
  },
};
</script>
