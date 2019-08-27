<template>
  <div id="queue-table">
    <b-table striped hover dark :items="items"></b-table>
    <h2 class="text-center" v-if="items.length === 0">No Items in Queue</h2>
  </div>
</template>

<script>
import auth from '../auth';

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
      const url = '/api/worker/queue';

      this.$http.get(url, {
        headers: auth.getAuthHeader(),
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
