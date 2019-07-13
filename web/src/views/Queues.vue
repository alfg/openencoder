<template>
  <div id="queues" class="container">
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
