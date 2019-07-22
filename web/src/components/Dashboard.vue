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
const STATS_INTERVAL = 5000;

export default {
  name: 'dashboard',

  data() {
    return {
      stats: {},
    };
  },

  mounted() {
    this.getStats();
    setInterval(() => {
      this.getStats();
    }, STATS_INTERVAL);
  },

  methods: {
    getStatusColor(o) {
      const statusMap = {
        queued: 'primary',
        completed: 'success',
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
