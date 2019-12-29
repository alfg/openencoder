<template>
  <div id="workers-table">
    <b-table striped hover dark :items="items"></b-table>
    <h2 class="text-center" v-if="items.length === 0">No Active Workers</h2>
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
    this.getWorkers();
  },

  methods: {
    getWorkers() {
      api.getWorkerPools(this, (err, json) => {
        this.items = json;
      });
    },
  },
};
</script>
