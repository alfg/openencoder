<template>
  <div id="s3-browser">
    <div>
      <ul>
      <!-- Folders -->
      <!-- <li v-for="(o, i) in tree.folders" v-bind:key="i">
        <a href="#" @click.prevent="onFolderSelect">{{ o }}</a>
      </li> -->

      <!-- Files -->
      <!-- <li v-for="o in tree.files" v-bind:key="o.name">
        <a href="#" @click.prevent="onFileSelect">{{ o.name }}</a>
      </li> -->

      <li v-for="o in tree" v-bind:key="o.label">
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
      tree: {},
    };
  },

  mounted() {
    this.getData();
  },

  methods: {
    onFileSelect(event) {
      const { text } = event.target;

      if (text[text.length - 1] !== '/') {
        this.$emit('file', event.target.text);
      } else {
        this.getData(text);
      }
    },

    onFolderSelect(event) {
      console.log(event.target.text);
      console.log(this.data);
    },

    getData(prefix = '') {
      const url = `/api/s3/list?prefix=${prefix}`;

      fetch(url)
        .then(response => (
          response.json()
        ))
        .then((json) => {
          console.log(json);
          // this.tree = json.data;
          this.updateTree(json.data);
        });
    },

    updateTree(data) {
      const items = [];

      if (data && data.folders) {
        data.folders.forEach((item) => {
          const o = {
            label: item,
            children: [],
          };
          items.push(o);
        });
      }

      if (data && data.files) {
        data.files.forEach((item) => {
          const o = {
            label: item.name,
            children: [],
          };
          items.push(o);
        });
      }

      console.log(items);

      this.tree = items;
    },
  },
};
</script>
