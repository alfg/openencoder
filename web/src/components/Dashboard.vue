<template>
  <div class="dashboard container">
    <b-card-group deck>
      <div
        v-for="(o, i) in stats.jobs"
        v-bind:key="i"
        class="col-lg-4 col-md-4 p-1 mb-2"
      >
          <b-card
            :border-variant="getStatusColor(o.status)"
            :header="o.status"
            class="text-center"
            style=""
            >
            <b-card-text>{{o.count}}</b-card-text>
          </b-card>
      </div>
    </b-card-group>
  </div>
</template>

<script>
export default {
  name: 'dashboard',

  data() {
    return {
      stats: {},
    };
  },

  mounted() {
    this.getStats();
  },

  methods: {
    getStatusColor(o) {
      const statusMap = {
        created: 'primary',
        completed: 'success',
        pending: 'warning',
        downloading: 'warning',
        encoding: 'info',
        uploading: 'warning',
        error: 'danger',
      };
      return statusMap[o] || 'secondary';
    },
    getStats() {
      const url = '/api/stats';

      fetch(url)
        .then(response => (
          response.json()
        ))
        .then((json) => {
          this.stats = json.stats;
        });
    },
  },
};
</script>
