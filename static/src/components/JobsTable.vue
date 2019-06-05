<template>
  <div id="jobs-table">
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
    this.getJobs();
  },

  methods: {
    getJobs() {
      const url = '/api/jobs';

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
#jobs-table {
  width: 90%;
  margin: 0 auto;
}
</style>
