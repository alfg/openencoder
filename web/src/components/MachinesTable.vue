<template>
  <div id="machines-table">
    <b-table striped hover dark :items="items"></b-table>
    <h2 class="text-center" v-if="items.length === 0">No Active Machines</h2>
  </div>
</template>

<script>
const UPDATE_INTERVAL = 5000;
let intervalId;

export default {
  data() {
    return {
      items: [],
    };
  },

  mounted() {
    this.getMachines();
    intervalId = setInterval(() => { this.getMachines(); }, UPDATE_INTERVAL);
  },

  destroyed() {
    clearInterval(intervalId);
  },

  methods: {
    getMachines() {
      const url = '/api/machines';

      fetch(url)
        .then(response => (
          response.json()
        ))
        .then((json) => {
          this.items = json && json.machines;
        });
    },
  },
};
</script>
