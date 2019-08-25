<template>
  <div id="queue-table">
    <b-table striped hover dark :items="items"></b-table>
    <h2 class="text-center" v-if="items.length === 0">No Items in Queue</h2>
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
    this.getQueue();
  },

  methods: {
    getQueue() {
      const url = '/api/worker/queue';

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
          console.log(json);
          this.items = json;
        });
    },
  },
};
</script>
