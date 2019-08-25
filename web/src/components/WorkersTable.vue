<template>
  <div id="workers-table">
    <b-table striped hover dark :items="items"></b-table>
    <h2 class="text-center" v-if="items.length === 0">No Active Workers</h2>
  </div>
</template>

<script>
import store from '../store';

export default {
  data() {
    return {
      items: [],
      storeState: store.state,
    };
  },

  mounted() {
    this.getWorkers();
  },

  methods: {
    getWorkers() {
      const url = '/api/worker/pools';

      fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${this.storeState.token}`,
        },
      })
        .then(response => (
          response.json()
        ))
        .then((json) => {
          this.items = json;
        });
    },
  },
};
</script>
