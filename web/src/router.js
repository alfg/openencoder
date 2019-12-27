import Vue from 'vue';
import Router from 'vue-router';
import VueResource from 'vue-resource';
import Moment from 'vue-moment';
import BootstrapVue from 'bootstrap-vue';
import Login from './views/Login.vue';
import Register from './views/Register.vue';
import UpdatePassword from './views/UpdatePassword.vue';

import cookie from './cookie';

import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';

Vue.use(Router);
Vue.use(BootstrapVue);
Vue.use(VueResource);
Vue.use(Moment);

Vue.http.headers.common.Authorization = `Bearer ${cookie.get('token')}`;

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      redirect: '/status',
    },
    {
      path: '/status',
      name: 'home',
      component: () => import(/* webpackChunkName: "status" */ './views/Status.vue'),
    },
    {
      path: '/jobs',
      name: 'jobs',
      component: () => import(/* webpackChunkName: "jobs" */ './views/Jobs.vue'),
    },
    {
      path: '/encode',
      name: 'encode',
      component: () => import(/* webpackChunkName: "encode" */ './views/Encode.vue'),
    },
    {
      path: '/queue',
      name: 'queue',
      component: () => import(/* webpackChunkName: "queue" */ './views/Queue.vue'),
    },
    {
      path: '/workers',
      name: 'workers',
      component: () => import(/* webpackChunkName: "workers" */ './views/Workers.vue'),
    },
    {
      path: '/machines',
      name: 'machines',
      component: () => import(/* webpackChunkName: "machines" */ './views/Machines.vue'),
    },
    {
      path: '/presets',
      name: 'presets',
      component: () => import(/* webpackChunkName: "presets" */ './views/Presets.vue'),
    },
    {
      path: '/presets/create',
      name: 'presets-create',
      component: () => import(/* webpackChunkName: "presets" */ './views/PresetsCreate.vue'),
    },
    {
      path: '/users',
      name: 'users',
      component: () => import(/* webpackChunkName: "users" */ './views/Users.vue'),
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import(/* webpackChunkName: "settings" */ './views/Settings.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import(/* webpackChunkName: "profile" */ './views/UserProfile.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: { hideNavigation: true },
    },
    {
      path: '/register',
      name: 'register',
      component: Register,
      meta: { hideNavigation: true },
    },
    {
      path: '/update-password',
      name: 'update-password',
      component: UpdatePassword,
      meta: { hideNavigation: true },
    },
  ],
});
