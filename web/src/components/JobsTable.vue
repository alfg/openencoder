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
      :items="items">

      <template v-slot:cell(progress)="data">
        <b-progress
          :value="data.item.progress"
          :animated="data.value !== 100"
          :variant="data.value === 100 ? 'success' : 'primary'"
          show-progress></b-progress>
      </template>

      <template v-slot:cell(show_details)="row">
        <b-button size="sm" @click="row.toggleDetails" class="mr-2">
          {{ row.detailsShowing ? 'Hide' : 'Show'}} Details
        </b-button>
      </template>

      <template v-slot:row-details="row">
        <div class="code">
          <b-form-textarea
            rows="3"
            max-rows="6"
            :value="row.item.encode"
          ></b-form-textarea>
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

const UPDATE_INTERVAL = 5000;
let intervalId;

export default {
  data() {
    return {
      fields: ['id', 'guid', 'preset', 'created_date', 'status', 'progress', 'show_details'],
      items: [],
      count: 0,
      autoUpdate: true,
    };
  },

  computed: {
    pages() {
      return this.count === 0 ? 1 : Math.ceil(this.count / 10);
    },
  },

  mounted() {
    const page = this.$route.query.page || 0;

    this.getJobs(page);
    intervalId = setInterval(() => {
      if (this.autoUpdate) {
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

    getJobs(page) {
      const url = `/api/jobs?page=${page}`;

      this.$http.get(url, {
        headers: auth.getAuthHeader(),
      })
        .then(response => (
          response.json()
        ))
        .then((json) => {
          this.items = json && json.items;
          this.count = json && json.count;
        });
    },
  },
};
</script>

<style scoped>
.code {
  background-color: #f4f4f4;
  border: 1px solid #aaa;
  color: #000;
  font-family: monospace;
  margin-top: 10px;
  padding: 5px;
}
</style>
