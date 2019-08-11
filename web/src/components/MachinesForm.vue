<template>
  <div id="machines-form">
    <h4>Create Machines</h4>
    <b-form inline>
      <b-form-select
        id="inline-form-custom-select-provider"
        class="mb-2 mr-sm-2 mb-sm-0"
        v-model="form.provider"
        :value="null"
        :options="providers"
        required
      >
      </b-form-select>

      <b-form-select
        id="inline-form-custom-select-regions"
        class="mb-2 mr-sm-2 mb-sm-0"
        v-model="form.region"
        :value="null"
        :options="regions"
        @change="onRegionChange"
      >
        <option slot="first" :value="null">Regions</option>
      </b-form-select>

      <b-form-select
        id="inline-form-custom-select-sizes"
        class="mb-2 mr-sm-2 mb-sm-0"
        v-model="form.size"
        :value="null"
        :options="sizes"
      >
        <option slot="first" :value="null">Size</option>
      </b-form-select>

      <b-form-select
        class="mb-2 mr-sm-2 mb-sm-0"
        :value="null"
        :options="count"
        id="inline-form-custom-select-count"
      >
        <option slot="first" :value="null">Count</option>
      </b-form-select>

      <b-button variant="primary">Apply</b-button>
    </b-form>
    <hr />
  </div>
</template>

<script>
export default {
  data() {
    return {
      form: {
        provider: 'digitalocean',
        region: null,
        size: null,
        count: null,
      },
      providers: { digitalocean: 'Digital Ocean' },
      regionsData: [],
      sizesData: [],
      sizes: [],
      count: Object.assign({}, [...Array(10).keys()]), // 0..10.
    };
  },

  computed: {
    regions() {
      return this.regionsData.map(o => o.name).sort();
    },
  },

  mounted() {
    this.getRegions();
    this.getSizes();
  },

  methods: {
    getRegions() {
      const url = '/api/machines/regions';

      fetch(url)
        .then(response => (
          response.json()
        ))
        .then((json) => {
          this.regionsData = json && json.regions;
        });
    },

    getSizes() {
      const url = '/api/machines/sizes';

      fetch(url)
        .then(response => (
          response.json()
        ))
        .then((json) => {
          this.sizesData = json && json.sizes;
        });
    },

    onRegionChange() {
      const { region } = this.form;

      // Get pricing data with regions available.
      const sizesAvailable = this.regionsData.find(o => o.name === region).sizes;

      // Filter sizes data and sort by price to array.
      const filtered = this.sizesData
        .filter(o => sizesAvailable.includes(o.slug) && o.available)
        .sort((a, b) => a.price_monthly > b.price_monthly)
        .map(o => `${o.slug} -- $${o.price_monthly} /mo`);

      this.sizes = filtered;
    },
  },
};
</script>
