<template>
  <div class="status container">
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
    <h2 class="text-center" v-if="!stats.jobs">No Stats Found</h2>
  </div>
</template>

<script>
import auth from '../auth';

const UPDATE_INTERVAL = 5000;
let intervalId;

export default {
  name: 'status',

  data() {
    return {
      stats: {},
    };
  },

  mounted() {
    this.getStats();
    intervalId = setInterval(this.getStats, UPDATE_INTERVAL);
  },

  destroyed() {
    clearInterval(intervalId);
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

      this.$http.get(url, {
        headers: auth.getAuthHeader(),
      })
        .then(response => (
          response.json()
        ))
        .then((json) => {
          if (json.stats) {
            this.stats = json.stats;
          }
        })
        .catch((err) => {
          console.log(err);
        });
    },
  },
};
</script>
