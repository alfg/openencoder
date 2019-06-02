<template>
  <div id="jobs-table">
    <b-table striped hover dark :items="items"></b-table>
  </div>
</template>

<script>
export default {
  data() {
    return {
      items: [
        { id: 1, profile: 'test', created_date: 'now' },
        { id: 2, profile: 'test', created_date: 'now' },
        { id: 3, profile: 'test', created_date: 'now' },
      ],
    };
  },
};
</script>

<style scoped>
#jobs-table {
  width: 90%;
  margin: 0 auto;
}
</style>
