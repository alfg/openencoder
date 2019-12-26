<template>
  <div id="users" class="container">
    <UsersTable />
  </div>
</template>

<script>
import UsersTable from '@/components/UsersTable.vue';

export default {
  name: 'Users',
  components: {
    UsersTable,
  },
};
</script>
