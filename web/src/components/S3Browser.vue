<template>
  <div id="s3-browser">
    <div>
      <ul>
        <li v-if="prefix !== ''">
          <a href="#" @click.prevent="goBack">...</a>
        </li>
        <li v-for="o in filteredFiles" v-bind:key="o.label">
          <a href="#" @click.prevent="onFileSelect">{{ o.label }}</a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      prefix: '',
      data: [],
    };
  },

  computed: {
    filteredFiles() {
      const items = [];

      if (this.data.folders) {
        this.data.folders.forEach((item) => {
          const o = {
            label: item,
            children: [],
          };
          items.push(o);
        });
      }

      if (this.data.files) {
        this.data.files.forEach((item) => {
          const o = {
            label: item.name,
            children: [],
          };
          items.push(o);
        });
      }
      return items.filter(o => o.label !== this.prefix);
    },
  },

  mounted() {
    this.getData();
  },

  methods: {
    onFileSelect(event) {
      const { text } = event.target;

      if (text[text.length - 1] !== '/') {
        this.$emit('file', `s3:///${event.target.text}`);
      } else {
        this.getData(text);
      }
    },

    goBack() {
      const arr = this.prefix.split('/');
      arr.splice(-2, 1); // Remove last path, but keep leading slash.
      const newPrefix = arr.join('/');
      this.getData(newPrefix);
    },

    getData(prefix = '') {
      const url = `/api/s3/list?prefix=${prefix}`;

      fetch(url)
        .then(response => (
          response.json()
        ))
        .then((json) => {
          this.data = json.data;
          this.prefix = prefix;
        });
    },
  },
};
</script>
