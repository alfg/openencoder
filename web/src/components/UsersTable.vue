<template>
  <div id="users-table">
    <div class="mb-2 text-right">
      <b-button to="register">Create User</b-button>
    </div>

    <b-table
      striped hover dark
      selectable
      select-mode="single"
      :fields="fields"
      :items="items"
      @row-selected="onRowSelected">
    </b-table>

    <b-form class="mb-3" @submit="onSubmit" v-show="data">
      <b-form-group id="input-group-active">
        <b-form-checkbox v-model="form.active">Active?</b-form-checkbox>
      </b-form-group>

      <b-button type="submit" variant="primary">Save</b-button>
    </b-form>

    <b-alert
      :show="dismissCountDown"
      dismissible
      fade
      variant="success"
      @dismissed="dismissCountDown=0"
      @dismiss-count-down="countDownChanged"
    >
      Updated User!
    </b-alert>

    <h2 class="text-center" v-if="items.length === 0">No Users Found</h2>
  </div>
</template>

<script>
import auth from '../auth';

export default {
  data() {
    return {
      fields: ['id', 'username', 'role', 'active'],
      items: [],
      data: null,
      form: {
        active: false,
      },
      dismissSecs: 5,
      dismissCountDown: 0,
      showDismissibleAlert: false,
    };
  },

  computed: {
    pages() {
      return this.count === 0 ? 1 : Math.ceil(this.count / 10);
    },
  },

  mounted() {
    this.getUsers();
  },

  methods: {
    countDownChanged(dismissCountDown) {
      this.dismissCountDown = dismissCountDown;
    },

    linkGen(pageNum) {
      return pageNum === 1 ? '?' : `?page=${pageNum}`;
    },

    onChangePage(event) {
      this.getUsers(event);
    },

    getUsers() {
      const url = '/api/users';

      this.$http.get(url, {
        headers: auth.getAuthHeader(),
      })
        .then(response => (
          response.json()
        ))
        .then((json) => {
          this.items = (json && json.users) || [];
        });
    },

    updateUser(data) {
      const url = `/api/users/${data.id}`;

      this.$http.put(url, data, {
        headers: auth.getAuthHeader(),
      }).then(response => (
        response.json()
      )).then((json) => {
        console.log('Submitted form: ', json);
        this.dismissCountDown = this.dismissSecs;
      });
    },


    onRowSelected(items) {
      if (items.length > 0) {
        [this.form] = items;

        this.data = this.form.data || {};
      } else {
        this.data = null;
      }
    },

    onSubmit(evt) {
      evt.preventDefault();
      console.log(JSON.stringify(this.form));
      this.updateUser(this.form);
    },
  },
};
</script>
