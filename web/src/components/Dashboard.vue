<template>
  <div class="dashboard container">
    <div class="mt-1">
      <b-card-group deck>
        <b-card border-variant="success" header="Completed" class="text-center">
          <b-card-text>0</b-card-text>
        </b-card>

        <b-card border-variant="secondary" header="In-Progress" class="text-center">
          <b-card-text>0</b-card-text>
        </b-card>

        <b-card border-variant="danger" header="Failed" class="text-center">
          <b-card-text>0</b-card-text>
        </b-card>
      </b-card-group>
    </div>
  </div>
</template>

<script>
export default {
  name: 'dashboard',
  components: {
  },
};
</script>
