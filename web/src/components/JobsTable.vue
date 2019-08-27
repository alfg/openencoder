<template>
  <div id="jobs-table">
    <b-table
      striped hover dark
      :fields="fields"
      :items="items">
      <template slot="progress" slot-scope="data">
        <b-progress
          :value="data.item.progress"
          :animated="data.value !== 100"
          :variant="data.value === 100 ? 'success' : 'primary'"
          show-progress></b-progress>
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
      fields: ['id', 'guid', 'profile', 'created_date', 'status', 'progress'],
      items: [],
      count: 0,
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
    intervalId = setInterval(() => { this.getJobs(page); }, UPDATE_INTERVAL);
  },

  destroyed() {
    clearInterval(intervalId);
  },

  methods: {
    linkGen(pageNum) {
      return pageNum === 1 ? '?' : `?page=${pageNum}`;
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
