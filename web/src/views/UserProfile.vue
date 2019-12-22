<template>
  <div id="user-profile" class="container">
    <UserProfile />
  </div>
</template>

<script>
import UserProfile from '@/components/UserProfile.vue';

export default {
  name: 'user-profile',
  components: {
    UserProfile,
  },
};
</script>
