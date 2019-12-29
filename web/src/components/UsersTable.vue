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

      <template v-slot:cell(role)="data">
        <b-badge
          :variant="['admin'].includes(data.item.role) ? 'danger' : 'primary'"
        >{{ data.item.role }}</b-badge>
      </template>

      <template v-slot:cell(active)="data">
        {{ data.item.active ? '✔️' : '❌' }}
      </template>
    </b-table>


    <b-form class="mb-3" @submit="onSubmit" v-show="data">
      <b-alert
        :show="isMasterUser"
        variant="warning">
      Cannot update master user settings.</b-alert>

      <b-form-group id="input-group-role" label="Role:" label-for="input-role">
        <b-form-select
          id="input-role"
          v-model="form.role"
          :options="roles"
          required
          :disabled="isMasterUser"
        ></b-form-select>
      </b-form-group>

      <b-form-group id="input-group-active">
        <b-form-checkbox
          v-model="form.active"
          :disabled="isMasterUser"
        >Active?</b-form-checkbox>
      </b-form-group>

      <b-button type="submit" variant="primary" :disabled="isMasterUser">Save</b-button>
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
import api from '../api';

export default {
  data() {
    return {
      fields: ['id', 'username', 'role', 'active'],
      items: [],
      data: null,
      form: {
        role: '',
        active: false,
      },
      roles: ['admin', 'operator', 'guest'],
      dismissSecs: 5,
      dismissCountDown: 0,
      showDismissibleAlert: false,
    };
  },

  computed: {
    pages() {
      return this.count === 0 ? 1 : Math.ceil(this.count / 10);
    },
    isMasterUser() {
      return this.data && this.data.id === 1;
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
      api.getUsers(this, (err, json) => {
        this.items = (json && json.users) || [];
      });
    },

    updateUser(data) {
      api.updateUser(this, data, (err, json) => {
        console.log('Submitted form: ', json);
        this.dismissCountDown = this.dismissSecs;
      });
    },

    onRowSelected(items) {
      if (items.length > 0) {
        [this.form] = items;

        this.data = this.form || {};
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
