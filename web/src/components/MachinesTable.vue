<template>
  <div id="machines-table">
    <b-table
      striped hover dark
      :fields="fields"
      :items="items">
      <template v-slot:cell(created_at)="data">
        <span
          v-b-tooltip="data.item.created_at">
          {{ data.item.created_at |  moment("from", "now") }}
        </span>
      </template>
      <template v-slot:cell(action)="data">
        <b-button variant="light" @click="onClickDelete(data.item.id)">❌</b-button>
      </template>
      </b-table>
    <h2 class="text-center" v-if="items.length === 0">No Active Machines</h2>
  </div>
</template>

<script>
import api from '../api';

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
      api.getMachines(this, (err, json) => {
        this.items = (json && json.machines) || [];
      });
    },

    deleteMachine(id) {
      api.deleteMachine(this, id, (err, json) => {
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
