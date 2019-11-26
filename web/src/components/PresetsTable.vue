<template>
  <div id="presets-table">
    <b-table
      striped hover dark
      selectable
      select-mode="single"
      :fields="fields"
      :items="items"
      @row-selected="onRowSelected">
    </b-table>

    <div class="mb-4" v-show="data">
      <div ref="editor"></div>
      <div class="mt-2 text-right">
        <b-button>Save</b-button>
      </div>
    </div>

    <h2 class="text-center" v-if="items.length === 0">No Presets Found</h2>

    <b-pagination-nav
      @change="onChangePage"
      :link-gen="linkGen"
      :number-of-pages="pages"
      use-router></b-pagination-nav>
  </div>
</template>

<script>
import JSONEditor from 'jsoneditor';
import 'jsoneditor/dist/jsoneditor.min.css';
import auth from '../auth';

export default {
  data() {
    return {
      fields: ['id', 'name', 'description', 'active'],
      items: [],
      count: 0,
      editor: null,
      data: null,
    };
  },

  computed: {
    pages() {
      return this.count === 0 ? 1 : Math.ceil(this.count / 10);
    },
  },

  mounted() {
    const page = this.$route.query.page || 0;
    this.getPresets(page);

    // Load JSONEditor.
    const container = this.$refs.editor;
    const options = {
      mode: 'tree',
      modes: ['code', 'text', 'tree', 'preview'],
    };
    this.editor = new JSONEditor(container, options);
  },

  destroyed() {
    this.editor = null;
  },

  methods: {
    linkGen(pageNum) {
      return pageNum === 1 ? '?' : `?page=${pageNum}`;
    },

    onChangePage(event) {
      this.getPresets(event);
    },

    getPresets(page) {
      const url = `/api/presets?page=${page}`;

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

    onRowSelected(items) {
      console.log('row selected', items[0].options);
      this.data = items[0].options || {};
      this.editor.set(JSON.parse(this.data));
    },
  },
};
</script>
