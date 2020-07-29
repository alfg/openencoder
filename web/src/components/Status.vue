<template>
  <div class="status container">
    <h2 class="text-muted">Jobs</h2>
    <b-card-group deck>
      <div
        v-for="(o, i) in stats.jobs"
        v-bind:key="i"
        class="col-lg-4 col-md-4 p-1 mb-2"
      >
          <b-card
            :border-variant="getStatusColor(o.status)"
            :header="o.status.toUpperCase()"
            class="text-center"
            style=""
            >
            <b-card-text>{{o.count}}</b-card-text>
          </b-card>
      </div>
    </b-card-group>
    <hr />

    <h2 class="text-muted">Health</h2>
    <b-card-group deck>
      <div
        v-for="(o, i) in health"
        v-bind:key="i"
        class="col-lg-3 col-md-3 p-1 mb-2"
      >
          <b-card
            :header="i.toUpperCase()"
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
    <hr />

    <h2 class="text-muted">Machines</h2>
    <b-card-group deck>
      <div
        v-for="(o, i) in pricing"
        v-bind:key="i"
        class="col-lg-4 col-md-4 p-1 mb-2"
      >
          <b-card
            :header="i.toUpperCase()"
            class="text-center"
            style=""
            >
            <b-card-text>
              <b-alert
                show
              >{{ o }}</b-alert>
            </b-card-text>
          </b-card>
      </div>
    </b-card-group>
    <h2 class="text-center" v-if="Object.keys(pricing).length === 0">No Machines Running</h2>

    <h2 class="text-center" v-if="!stats.jobs">No Stats Found</h2>
  </div>
</template>

<script>
import api from '../api';

const UPDATE_INTERVAL = 5000;
const HEALTH_UPDATE_INTERVAL = 10000;
const PRICING_UPDATE_INTERVAL = 30000;

let intervalId;
let healthIntervalId;
let pricingIntervalId;

export default {
  name: 'status',

  data() {
    return {
      stats: {},
      health: {},
      pricing: {},
    };
  },

  mounted() {
    this.getStats();
    this.getHealth();
    this.getPricing();

    intervalId = setInterval(this.getStats, UPDATE_INTERVAL);
    healthIntervalId = setInterval(this.getHealth, HEALTH_UPDATE_INTERVAL);
    pricingIntervalId = setInterval(this.getPricing, PRICING_UPDATE_INTERVAL);
  },

  destroyed() {
    clearInterval(intervalId);
    clearInterval(healthIntervalId);
    clearInterval(pricingIntervalId);
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
      api.getStats(this, (err, json) => {
        if (json.stats) {
          this.stats = json.stats;
        }
      });
    },

    getHealth() {
      api.getHealth(this, (err, json) => {
        if (json) {
          this.health = json;
        }
      });
    },

    getPricing() {
      api.getPricing(this, (err, json) => {
        if (err) {
          console.log(err);
          return;
        }
        if (json.pricing) {
          this.pricing = json.pricing;
        }
      });
    },
  },
};
</script>
