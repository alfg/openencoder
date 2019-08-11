<template>
  <div id="machines-table">
    <b-table
      striped hover dark
      :fields="fields"
      :items="items">
      <template slot="action" slot-scope="data">
        <b-button variant="light" @click="onClickDelete(data.item.id)">❌</b-button>
      </template>
      </b-table>
    <h2 class="text-center" v-if="items.length === 0">No Active Machines</h2>
  </div>
</template>

<script>
const UPDATE_INTERVAL = 5000;
let intervalId;

export default {
  data() {
    return {
      fields: ['id', 'name', 'status', 'size_slug', 'created_at', 'region', 'tags', 'provider', 'action'],
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

    deleteMachine(id) {
      const url = `/api/machines/${id}`;

      fetch(url, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
        },
      }).then(response => (
        response.json()
      )).then((json) => {
        console.log('Deleting machine: ', json);
      });
    },

    onClickDelete(evt) {
      const id = evt;
      this.deleteMachine(id);
    },
  },
};
</script>
