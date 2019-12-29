<template>
  <div id="jobs-table">
    <div class="text-right mb-2">
      <b-button
        @click="toggleAutoUpdate"
      >{{ autoUpdate ? '&#10227; Auto Update' : '&#10074;&#10074; Paused' }}</b-button>
    </div>

    <b-table
      striped hover dark
      :fields="fields"
      :busy="!items"
      :items="items">

      <template v-slot:cell(created_date)="data">
        <span
          v-b-tooltip="data.item.created_date">
          {{ data.item.created_date |  moment("from", "now") }}
        </span>
      </template>

      <template v-slot:cell(status)="data">
        <b-badge
          :variant="['error', 'cancelled'].includes(data.item.status) ? 'danger' : 'primary'"
        >{{ data.item.status }}</b-badge>
      </template>

      <template v-slot:cell(progress)="data">
        <b-progress
          v-if="!['error', 'cancelled'].includes(data.item.status)"
          :value="data.item.progress"
          :animated="data.value !== 100"
          :variant="data.value === 100 ? 'success' : 'primary'"
          show-progress></b-progress>
        <p
          class="text-monospace text-center"
          style="font-size: 0.7em; margin: 0;"
          v-if="(data.item.speed && data.item.fps) && data.item.status === 'encoding'"
        >{{ data.item.speed }} @ {{ data.item.fps }} FPS</p>
      </template>

      <template v-slot:cell(details)="row">
        <b-button size="sm" @click="row.toggleDetails" class="mr-2">
          {{ row.detailsShowing ? 'Hide' : 'Show'}}
        </b-button>
      </template>

      <template v-slot:cell(action)="data" v-if="isOperatorAdmin">
        <b-button-group size="sm">
          <b-button
            variant="light"
            v-if="!['error', 'cancelled', 'completed'].includes(data.item.status)"
            @click="onClickCancel(data.item.id)">❌</b-button>
          <b-button
            variant="light"
            v-if="['error', 'cancelled'].includes(data.item.status)"
            @click="onClickRestart(data.item.id)">&#10227;</b-button>
        </b-button-group>
      </template>

      <template v-slot:row-details="row">
          <b-row class="mb-2">
            <b-col sm="2" class="text-sm-right"><b>Guid:</b></b-col>
            <b-col>{{ row.item.guid }}</b-col>
          </b-row>

          <b-row class="mb-2">
            <b-col sm="2" class="text-sm-right"><b>Source:</b></b-col>
            <b-col>{{ row.item.source }}</b-col>
          </b-row>

          <b-row class="mb-2">
            <b-col sm="2" class="text-sm-right"><b>Destination:</b></b-col>
            <b-col>{{ row.item.destination }}</b-col>
          </b-row>

          <b-row class="mb-2">
            <b-col sm="2" class="text-sm-right"><b>Probe Data:</b></b-col>
            <b-col>
              <div class="code">
                <b-form-textarea
                  rows="3"
                  max-rows="6"
                  :value="row.item.probe"
                ></b-form-textarea>
              </div>
            </b-col>
          </b-row>

          <b-row class="mb-2">
            <b-col sm="2" class="text-sm-right"><b>Encode Options:</b></b-col>
            <b-col>
              <div class="code">
                <b-form-textarea
                  rows="3"
                  max-rows="6"
                  :value="row.item.options"
                ></b-form-textarea>
              </div>
            </b-col>
          </b-row>
      </template>

      <template v-slot:table-busy>
        <div class="text-center text-danger my-2">
          <b-spinner class="align-middle"></b-spinner>
          <strong>Loading...</strong>
        </div>
      </template>

    </b-table>
    <h2 class="text-center" v-if="items.length === 0">No Jobs Found</h2>

    <b-pagination-nav
      @change="onChangePage"
      :link-gen="linkGen"
      :number-of-pages="pages"
      use-router></b-pagination-nav>
  </div>
</template>

<script>
import auth from '../auth';
import api from '../api';

const UPDATE_INTERVAL = 5000;
let intervalId;

export default {
  data() {
    return {
      fields: ['id', 'source', 'preset', 'created_date', 'status', 'progress', 'details', 'action'],
      items: [],
      count: 0,
      autoUpdate: true,
      user: auth.user,
      role: auth.role,
    };
  },

  computed: {
    pages() {
      return this.count === 0 ? 1 : Math.ceil(this.count / 10);
    },
    isOperatorAdmin() {
      return ['operator', 'admin'].includes(this.user.role);
    },
  },

  mounted() {
    let page = this.$route.query.page || 0;
    this.getJobs(page);

    // Start auto-update interval.
    intervalId = setInterval(() => {
      if (this.autoUpdate) {
        page = this.$route.query.page || 0;
        this.getJobs(page);
      }
    }, UPDATE_INTERVAL);
  },

  destroyed() {
    clearInterval(intervalId);
  },

  methods: {
    linkGen(pageNum) {
      return pageNum === 1 ? '?' : `?page=${pageNum}`;
    },

    toggleAutoUpdate() {
      this.autoUpdate = !this.autoUpdate;
    },

    onChangePage(event) {
      this.getJobs(event);
    },

    onClickCancel(evt) {
      const id = evt;
      this.cancelJob(id);
    },

    onClickRestart(evt) {
      const id = evt;
      this.restartJob(id);
    },

    cancelJob(id) {
      api.cancelJob(this, id, (err, json) => {
        console.log('Cancel job: ', json);
      });
    },

    restartJob(id) {
      api.restartJob(this, id, (err, json) => {
        console.log('Restart job: ', json);
      });
    },

    getJobs(page) {
      api.getJobs(this, page, (err, json) => {
        this.items = json && json.items;
        this.count = json && json.count;
      });
    },
  },
};
</script>

<style>
.code {
  background-color: #f4f4f4;
  border: 1px solid #aaa;
  color: #000;
  font-family: monospace;
  margin-top: 10px;
  padding: 5px;
}
#jobs-table .table td {
  vertical-align: middle;
}
</style>
