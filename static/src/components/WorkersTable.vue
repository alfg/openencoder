<template>
  <div id="workers-table">
    <b-table striped hover dark :items="items"></b-table>
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
          console.log(json);
          this.items = json;
        });
    },
  },
};
</script>

<style scoped>
#workers-table {
  width: 90%;
  margin: 0 auto;
}
</style>
