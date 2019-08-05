import Vue from 'vue';
import Router from 'vue-router';
import BootstrapVue from 'bootstrap-vue';
import Home from './views/Home.vue';
import Create from './views/Create.vue';
import Jobs from './views/Jobs.vue';
import Queue from './views/Queue.vue';
import Workers from './views/Workers.vue';
import Machines from './views/Machines.vue';

import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-vue/dist/bootstrap-vue.css';

Vue.use(Router);
Vue.use(BootstrapVue);

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
      path: '/create',
      name: 'create',
      component: Create,
    },
    {
      path: '/jobs',
      name: 'jobs',
      component: Jobs,
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
      path: '/status',
      name: 'status',
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "about" */ './views/Status.vue'),
    },
  ],
});
