<template>
  <div id="queues-table">
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
    this.getQueues();
  },

  methods: {
    getQueues() {
      const url = '/api/worker/queues';

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
#queues-table {
  width: 90%;
  margin: 0 auto;
}
</style>
