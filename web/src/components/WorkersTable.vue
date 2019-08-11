<template>
  <div id="workers-table">
    <b-table striped hover dark :items="items"></b-table>
    <h2 class="text-center" v-if="items.length === 0">No Active Workers</h2>
  </div>
</template>

<script>
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

      fetch(url)
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
