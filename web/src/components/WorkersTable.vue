<template>
  <div id="workers-table">
    <b-table striped hover dark :items="items"></b-table>
    <h2 class="text-center" v-if="items.length === 0">No Active Workers</h2>
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
    this.getWorkers();
  },

  methods: {
    getWorkers() {
      const url = '/api/worker/pools';

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
