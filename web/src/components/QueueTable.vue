<template>
  <div id="queue-table">
    <b-table striped hover dark :items="items"></b-table>
    <h2 class="text-center" v-if="items.length === 0">No Items in Queue</h2>
  </div>
</template>

<script>
import api from '../api';

export default {
  data() {
    return {
      items: [],
    };
  },

  mounted() {
    this.getQueue();
  },

  methods: {
    getQueue() {
      api.getWorkerQueue(this, (err, json) => {
        this.items = json;
      });
    },
  },
};
</script>
