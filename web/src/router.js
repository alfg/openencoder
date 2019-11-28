import Vue from 'vue';
import Router from 'vue-router';
import VueResource from 'vue-resource';
import BootstrapVue from 'bootstrap-vue';
import Home from './views/Home.vue';
import Jobs from './views/Jobs.vue';
import Encode from './views/Encode.vue';
import Queue from './views/Queue.vue';
import Workers from './views/Workers.vue';
import Machines from './views/Machines.vue';
import Presets from './views/Presets.vue';
import PresetsCreate from './views/PresetsCreate.vue';
import Settings from './views/Settings.vue';
import Login from './views/Login.vue';
import Register from './views/Register.vue';

import cookie from './cookie';

import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';

Vue.use(Router);
Vue.use(BootstrapVue);
Vue.use(VueResource);

Vue.http.headers.common.Authorization = `Bearer ${cookie.get('token')}`;

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/jobs',
      name: 'jobs',
      component: Jobs,
    },
    {
      path: '/encode',
      name: 'encode',
      component: Encode,
    },
    {
      path: '/queue',
      name: 'queue',
      component: Queue,
    },
    {
      path: '/workers',
      name: 'workers',
      component: Workers,
    },
    {
      path: '/machines',
      name: 'machines',
      component: Machines,
    },
    {
      path: '/presets',
      name: 'presets',
      component: Presets,
    },
    {
      path: '/presets/create',
      name: 'presets-create',
      component: PresetsCreate,
    },
    {
      path: '/settings',
      name: 'settings',
      component: Settings,
    },
    {
      path: '/status',
      name: 'status',
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "about" */ './views/Status.vue'),
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
  ],
});
