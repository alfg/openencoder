<template>
  <div id="queue-table">
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
    this.getQueue();
  },

  methods: {
    getQueue() {
      const url = '/api/worker/queue';

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
