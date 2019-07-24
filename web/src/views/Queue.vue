<template>
  <div id="queues" class="container">
    <QueueTable />
  </div>
</template>

<script>
import QueueTable from '@/components/QueueTable.vue';

export default {
  name: 'Queues',
  components: {
    QueueTable,
  },
};
</script>
