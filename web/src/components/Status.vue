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
    <hr />

    <b-card-group deck>
      <div
        v-for="(o, i) in health"
        v-bind:key="i"
        class="col-lg-3 col-md-3 p-1 mb-2"
      >
          <b-card
            :header="i"
            class="text-center"
            style=""
            >
            <b-card-text>
              <b-alert
                show
                :variant="['NOTOK', 0].includes(o) ? 'danger' : 'success'">{{ o }}</b-alert>
            </b-card-text>
          </b-card>
      </div>
    </b-card-group>

    <h2 class="text-center" v-if="!stats.jobs">No Stats Found</h2>
  </div>
</template>

<script>
import auth from '../auth';

const UPDATE_INTERVAL = 5000;
const HEALTH_UPDATE_INTERVAL = 5000;

let intervalId;
let healthIntervalId;

export default {
  name: 'status',

  data() {
    return {
      stats: {},
      health: {},
    };
  },

  mounted() {
    this.getStats();
    this.getHealth();

    intervalId = setInterval(this.getStats, UPDATE_INTERVAL);
    healthIntervalId = setInterval(this.getHealth, HEALTH_UPDATE_INTERVAL);
  },

  destroyed() {
    clearInterval(intervalId);
    clearInterval(healthIntervalId);
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

    getHealth() {
      const url = '/api/health';

      this.$http.get(url, {})
        .then(response => (
          response.json()
        ))
        .then((json) => {
          this.health = json;
        })
        .catch(() => {
          this.health = {
            API: 'NOTOK',
            DB: 'NOTOK',
            Redis: 'NOTOK',
            Workers: 0,
          };
        });
    },
  },
};
</script>
